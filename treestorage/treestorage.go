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
func (s *NestedSetsStorage) GetParents(name string) ([]string, error) {
	if name == "" {
		return []string{}, errors.New("invalid node name")
	}

	db, err := sql.Open(s.DbDriver, s.DbConnectionString)
	if err != nil {
		return []string{}, err
	}
	defer db.Close()

	query :=
		`WITH child AS (SELECT ch.node_left, ch.node_right 
						FROM nodes AS ch WHERE ch.name = $1)
		SELECT n.name
		FROM nodes AS n, child 
		WHERE n.node_left < child.node_left AND n.node_right > child.node_right;`
	rows, err := db.Query(query, name)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var nodeName string
		err := rows.Scan(&nodeName)
		if err != nil {
			return []string{}, err
		}
		result = append(result, nodeName)
	}

	return result, nil
}

// GetChildren returns children for the node name
func (s *NestedSetsStorage) GetChildren(name string) ([]string, error) {
	if name == "" {
		return []string{}, errors.New("invalid node name")
	}

	db, err := sql.Open(s.DbDriver, s.DbConnectionString)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}
	defer db.Close()

	query :=
		`WITH parent AS (SELECT p.node_left, p.node_right 
						FROM nodes AS p WHERE p.name = $1)
		SELECT n.name 
		FROM nodes AS n, parent 
		WHERE n.node_left > parent.node_left AND n.node_right < parent.node_right;`
	rows, err := db.Query(query, name)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var nodeName string
		err := rows.Scan(&nodeName)
		if err != nil {
			return []string{}, err
		}
		result = append(result, nodeName)
	}

	return result, nil
}

// GetWholeTree returns all nodes
func (s *NestedSetsStorage) GetWholeTree() ([]NestedSetsNode, error) {
	db, err := sql.Open(s.DbDriver, s.DbConnectionString)
	if err != nil {
		log.Println(err)
		return []NestedSetsNode{}, err
	}
	defer db.Close()

	query := `SELECT name, node_left, node_right
			  FROM nodes;`
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return []NestedSetsNode{}, err
	}
	defer rows.Close()

	var result []NestedSetsNode
	for rows.Next() {
		var node NestedSetsNode
		err := rows.Scan(&node.Name, &node.Left, &node.Right)
		if err != nil {
			return []NestedSetsNode{}, err
		}
		result = append(result, node)
	}

	return result, nil
}

// AddNode adds new child node with name name for parent node with name parent
func (s *NestedSetsStorage) AddNode(name string, parent string) error {
	if name == "" || parent == "" {
		return errors.New("invalid node name")
	}

	db, err := sql.Open(s.DbDriver, s.DbConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`CALL add_node($1, $2);`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(name, parent); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// RemoveNode removes node with name name
func (s *NestedSetsStorage) RemoveNode(name string) error {
	if name == "" {
		return errors.New("invalid node name")
	}

	db, err := sql.Open(s.DbDriver, s.DbConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`CALL remove_node($1);`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(name); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("CALL move_node($1,$2);")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(name, newParent); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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

	renameQuery := `UPDATE nodes 
					SET name = $1 
					WHERE name = $2;`
	_, err = db.Exec(renameQuery, newName, name)
	return err
}
