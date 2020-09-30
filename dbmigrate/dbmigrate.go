package dbmigrate

import (
	"NestedSetsStorage/configs"
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const _ATTTEMPTS = 10
const _ATTTEMPT_INTERVAL = 1000

// Migrate updates data base tables structure
func Migrate(config *configs.Config) error {
	db, err := tryToConnect(config)
	if err != nil {
		return err
	}
	defer db.Close()

	return createDb(db, queriesForCreatingDb())
}

func queriesForCreatingDb() []string {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS departments
		(
			id SERIAL,
			name VARCHAR(100) NOT NULL UNIQUE,
			node_left INT NOT NULL,
			node_right INT NOT NULL,
			PRIMARY KEY (id)
		);`,

		`CREATE OR REPLACE FUNCTION increase_nodes_left ( range_start INT, range_finish INT, value INT) 
		RETURNS VOID AS $$
		BEGIN

			UPDATE departments 
			SET node_left = node_left + value
			WHERE range_start < node_left AND node_left < range_finish;

		END;
		$$  LANGUAGE plpgsql`,

		`CREATE OR REPLACE FUNCTION increase_nodes_right ( range_start INT, range_finish INT, value INT) 
		RETURNS VOID AS $$
		BEGIN

			UPDATE departments 
			SET node_right = node_right + value
			WHERE range_start < node_right AND node_right < range_finish;

		END;
		$$  LANGUAGE plpgsql`,

		`CREATE OR REPLACE PROCEDURE remove_node (node_name varchar(100)) 
		AS $$
		DECLARE
			node RECORD;
		BEGIN	

			SELECT node_left, node_right, name
			INTO node
			FROM departments
			WHERE 
				name = node_name;

			IF node IS NOT NULL THEN

				PERFORM increase_nodes_left(node.node_left, node.node_right, -1);
				PERFORM increase_nodes_right(node.node_left, node.node_right, -1);

				UPDATE departments 
				SET node_left = node_left -2
				WHERE node_left > node.node_right;

				UPDATE departments 
				SET node_right = node_right -2
				WHERE node_right > node.node_right;

				DELETE FROM departments
				WHERE name = node.name;

			END IF;

		END;
		$$  LANGUAGE plpgsql`,

		`CREATE OR REPLACE PROCEDURE add_node (node_name varchar(100), parent_name varchar(100)) 
		AS $$
		DECLARE
			parent RECORD;
			node RECORD;
		BEGIN	

			SELECT node_left, node_right
			INTO parent
			FROM departments
			WHERE 
				name = parent_name;

			SELECT name
			INTO node
			FROM departments
			WHERE 
				name = node_name;

			IF parent IS NOT NULL AND node IS NULL THEN

				UPDATE departments 
				SET node_left = node_left +2
				WHERE node_left >= parent.node_right;

				UPDATE departments 
				SET node_right = node_right +2
				WHERE node_right >= parent.node_right;

				INSERT INTO departments
				(name, node_left, node_right) 
				VALUES (node_name, parent.node_right, parent.node_right + 1);

			END IF;

		END;
		$$  LANGUAGE plpgsql`,

		`CREATE OR REPLACE PROCEDURE move_node (node_name varchar(100), parent_name varchar(100)) 
		AS $$
		DECLARE
			node RECORD;
			parent RECORD;
		BEGIN
			SELECT node_left, node_right, name
			INTO node
			FROM departments
			WHERE 
				name = node_name;

			SELECT node_left, node_right, name 
			INTO parent
			FROM departments
			WHERE 
				name = parent_name;

			IF node IS NOT NULL and parent IS NOT NULL THEN

				IF node.node_right < parent.node_left THEN

					PERFORM increase_nodes_left(node.node_left, node.node_right, -1);
					PERFORM increase_nodes_right(node.node_left, node.node_right, -1);

					PERFORM increase_nodes_left(node.node_right, parent.node_left + 1, -2);
					PERFORM increase_nodes_right(node.node_right, parent.node_left + 1, -2);

					UPDATE departments 
					SET node_left = parent.node_left - 1,
					node_right = parent.node_left
					WHERE name = node.name;

				ELSEIF node.node_left > parent.node_right THEN

					PERFORM increase_nodes_left(node.node_left, node.node_right, 1);
					PERFORM increase_nodes_right(node.node_left, node.node_right, 1);

					PERFORM increase_nodes_left(parent.node_right - 1, node.node_left, 2);
					PERFORM increase_nodes_right(parent.node_right - 1, node.node_left, 2);

					UPDATE departments 
					SET node_left = parent.node_right,
					node_right = parent.node_right + 1
					WHERE name = node.name;

				ELSEIF node.node_right < parent.node_right AND node.node_left > parent.node_left THEN

					IF  parent.node_right - node.node_right < node.node_left - parent.node_left THEN
						PERFORM increase_nodes_left(node.node_left, node.node_right, -1);
						PERFORM increase_nodes_right(node.node_left, node.node_right, -1);

						PERFORM increase_nodes_left(node.node_right, parent.node_right, -2);
						PERFORM increase_nodes_right(node.node_right, parent.node_right, -2);

						UPDATE departments 
						SET node_left = parent.node_right - 2,
						node_right = parent.node_right - 1
						WHERE name = node.name;
					ELSE
						PERFORM increase_nodes_left(node.node_left, node.node_right, 1);
						PERFORM increase_nodes_right(node.node_left, node.node_right, 1);

						PERFORM increase_nodes_left(parent.node_left, node.node_left, 2);
						PERFORM increase_nodes_right(parent.node_left, node.node_left, 2);

						UPDATE departments 
						SET node_left = parent.node_left + 1,
						node_right = parent.node_left + 2
						WHERE name = node.name;
					END IF;

				ELSEIF parent.node_right < node.node_right AND parent.node_left > node.node_left THEN
				
					PERFORM increase_nodes_left(node.node_left, parent.node_left + 1, -1);
					PERFORM increase_nodes_right(node.node_left, parent.node_left + 1, -1);

					PERFORM increase_nodes_left(parent.node_left, node.node_right, 1);
					PERFORM increase_nodes_right(parent.node_left, node.node_right, 1);

					UPDATE departments 
					SET node_left = parent.node_left,
					node_right = parent.node_left + 1
					WHERE name = node.name;

				END IF;

			END IF;

		END;
		$$  LANGUAGE plpgsql`,
	}
	return queries
}

func tryToConnect(config *configs.Config) (*sql.DB, error) {
	var db *sql.DB
	var err error = nil
	for i := 0; i < _ATTTEMPTS; i++ {
		db, err = sql.Open(config.DbDriver, config.DbConnectionSting)
		wait, err := checkErrorForWaitingDb(err)
		if err == nil {
			break
		} else if !wait {
			return db, err
		}
	}
	return db, err
}

func createDb(db *sql.DB, queries []string) error {
	for _, query := range queries {
		err := tryQueryExec(db, query)
		if err != nil {
			return err
		}
	}
	return nil
}

func tryQueryExec(db *sql.DB, query string) error {
	var err error
	for i := 0; i < _ATTTEMPTS; i++ {
		_, err = db.Exec(query)
		wait, err := checkErrorForWaitingDb(err)
		if err == nil {
			break
		} else if !wait {
			return err
		}
	}
	return err
}

func checkErrorForWaitingDb(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	isWaitingError := strings.Contains(err.Error(), "the database system is starting up") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "EOF")
	if isWaitingError {
		log.Println("waiting for db")
		time.Sleep(time.Duration(_ATTTEMPT_INTERVAL) * time.Millisecond)
	}
	return isWaitingError, err
}
