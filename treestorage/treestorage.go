package treestorage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// NestedSetsNode is a tree node
type NestedSetsNode struct {
	Name  string
	Left  int
	Right int
}

// NestedSetsStorage is an interface for data base table
type NestedSetsStorage struct {
	DbConnectionString string
	DbDriver           string
}

// GetParents returns parents for the node name
func (s *NestedSetsStorage) GetParents(name string) []NestedSetsNode {
	return []NestedSetsNode{}
}

// GetChildren returns children for the node name
func (s *NestedSetsStorage) GetChildren(name string) []NestedSetsNode {
	return []NestedSetsNode{}
}

// GetWholeTree returns all nodes
func (s *NestedSetsStorage) GetWholeTree() []NestedSetsNode {
	db, err := sql.Open(s.DbDriver, s.DbConnectionString)
	if err != nil {
		log.Println(err)
		return []NestedSetsNode{}
	}
	defer db.Close()

	query := `SELECT name, node_left, node_right
			  FROM departments;`
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return []NestedSetsNode{}
	}
	defer rows.Close()

	var result []NestedSetsNode
	for rows.Next() {
		var node NestedSetsNode
		err := rows.Scan(&node.Name, &node.Left, &node.Right)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, node)
	}

	return result
}

// AddNode adds new child node with name name for parent node with name parent
func (s *NestedSetsStorage) AddNode(name string, parent string) {

}

// RemoveNode removes node with name name
func (s *NestedSetsStorage) RemoveNode(name string) {

}

// MoveNode moves node with name name
func (s *NestedSetsStorage) MoveNode(name string, newParent string) error {
	if name == "" || newParent == "" {
		return errors.New("invalid node name")
	}

	db, err := sql.Open(s.DbDriver, s.DbConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
	BEGIN;
	DO $$ 
	DECLARE
		node RECORD;
		parent RECORD;
	BEGIN
		SELECT node_left, node_right, name
		INTO node
		FROM departments
		WHERE 
			name = '%s';

		SELECT node_left, node_right, name 
		INTO parent
		FROM departments
		WHERE 
			name = '%s';

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
	END $$ LANGUAGE plpgsql;
	COMMIT;
	`
	_, err = db.Exec(fmt.Sprintf(query, name, newParent))
	return err
}

// RenameNode renames node with name name
func (s *NestedSetsStorage) RenameNode(name string, newName string) error {
	if name == "" || newName == "" {
		return errors.New("invalid node name")
	}

	db, err := sql.Open(s.DbDriver, s.DbConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	renameQuery := `UPDATE departments 
					SET name = $1 
					WHERE name = $2;`
	_, err = db.Exec(renameQuery, newName, name)
	return err
}
