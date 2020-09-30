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
	}

	return createTables(db, queries)
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

func createTables(db *sql.DB, queries []string) error {
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
