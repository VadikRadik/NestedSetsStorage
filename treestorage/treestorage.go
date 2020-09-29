package treestorage

import (
	"database/sql"
	"errors"
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
	DO $$ 
	BEGIN;
	DECLARE
		node RECORD;
		parent RECORD;
		delta_right INT;
		delta_left INT;
	BEGIN
		SELECT node_left, node_right 
		INTO node
		FROM departmens
		WHERE 
			name = $1;

		SELECT node_left, node_right 
		INTO parent
		FROM departmens
		WHERE 
			name = $2;

		IF node <> NULL and parent <> NULL
			delta_right := ABS(parent.node_right - node.node_rigt);
			delta_left := ABS(parent.node_left - node.node_left);

			IF (delta_left > delta_right)

				UPDATE departments 
				SET node_right = node_right - 1 
				WHERE node_right > node.node_left AND node_right < node.node_right;

				UPDATE departments 
				SET node_left = node_left - 1 
				WHERE node_left > node.node_left AND node_left < node.node_right;

				UPDATE departments 
				SET node_right = node_right - 2 
				WHERE 
					node_right >= node.node_right AND node_right <= parent.node_right
				OR
					node_right >= parent.node_right AND node_right <= node.node_right;

				UPDATE departments 
				SET node_left = node_left - 2 
				WHERE 
					node_left >= node.node_right AND node_left <= parent.node_right
				OR 
					node_left >= parent.node_right AND node_left <= node.node_right;
				
			ELSE

				UPDATE departments 
				SET node_right = node_right + 1 
				WHERE node_right > node.node_left AND node_right < node.node_right;

				UPDATE departments 
				SET node_left = node_left + 1 
				WHERE node_left > node.node_left AND node_left < node.node_right;

				UPDATE departments 
				SET node_right = node_right + 2 
				WHERE 
					node_right >= node.node_left AND node_right <= parent.node_left
				OR
					node_right >= parent.node_left AND node_right <= node.node_left;

				UPDATE departments 
				SET node_left = node_left + 2 
				WHERE 
					node_left >= node.node_left AND node_left <= parent.node_left
				OR
					node_left >= parent.node_left AND node_left <= node.node_left;

			END IF;
		END IF;
	COMMIT;
	END $$;
	
	`
	_, err = db.Exec(query, name, newParent)
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
