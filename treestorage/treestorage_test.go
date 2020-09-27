package treestorage_test

import (
	"NestedSetsStorage/configs"
	"NestedSetsStorage/treestorage"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
)

var dbConnectionString string
var dbDriver string

func init() {
	var config configs.Config
	_, err := toml.DecodeFile("../configs/config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}
	dbConnectionString = config.DbConnectionSting
	dbDriver = config.DbDriver
}

func TestAddNode(t *testing.T) {
	refillTestData()
	s := &treestorage.NestedSetsStorage{
		DbConnectionString: dbConnectionString,
		DbDriver:           dbDriver,
	}
	s.AddNode("", "")
	//clearTestDataFromDb()
}

func TestRemoveNode(t *testing.T) {

}

func TestMoveNode(t *testing.T) {

}

func TestRenameNode(t *testing.T) {

}

func loadTestDataToDb() {
	nodes := createTestNodes()

	db, err := sql.Open(dbDriver, dbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fullQuery := `INSERT INTO 
					departments (name, node_left, node_right) 
				VALUES 
					%s;`
	nodeFields := "('%s', %d, %d)"

	nodesValues := make([]string, len(nodes))
	for i, node := range nodes {
		nodesValues[i] = fmt.Sprintf(nodeFields, node.Name, node.Left, node.Right)
	}

	fullValues := strings.Join(nodesValues, ",\n")
	fullQuery = fmt.Sprintf(fullQuery, fullValues)
	_, err = db.Exec(fullQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func createTestNodes() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{
			Name:  "Директор",
			Left:  0,
			Right: 35,
		},
		{
			Name:  "Заместитель директора по АХЧ",
			Left:  1,
			Right: 4,
		},
		{
			Name:  "Обслуживающий персонал",
			Left:  2,
			Right: 3,
		},
		{
			Name:  "Совет лицея",
			Left:  5,
			Right: 12,
		},
		{
			Name:  "Благотворительный фонд \"Развитие школы\"",
			Left:  6,
			Right: 7,
		},
		{
			Name:  "Ученическое самоуправление",
			Left:  8,
			Right: 11,
		},
		{
			Name:  "Ученики",
			Left:  9,
			Right: 10,
		},
		{
			Name:  "Заместитель директора по информатизации",
			Left:  13,
			Right: 16,
		},
		{
			Name:  "Инженегр по ВТ",
			Left:  14,
			Right: 15,
		},
		{
			Name:  "Заместитель директора по ВР",
			Left:  17,
			Right: 24,
		},
		{
			Name:  "Служба сопровождения",
			Left:  18,
			Right: 19,
		},
		{
			Name:  "Методическое объединение педагогов дополнительного образования",
			Left:  20,
			Right: 21,
		},
		{
			Name:  "Методическое объединение классных руководителей",
			Left:  22,
			Right: 23,
		},
		{
			Name:  "Бухгалтерия",
			Left:  25,
			Right: 26,
		},
		{
			Name:  "Педагогический совет",
			Left:  27,
			Right: 28,
		},
		{
			Name:  "Заместитель директора по УВР",
			Left:  29,
			Right: 32,
		},
		{
			Name:  "Кафедры профильного образования",
			Left:  30,
			Right: 31,
		},
		{
			Name:  "Научно-методический совет",
			Left:  33,
			Right: 34,
		},
	}
	return nodes
}

func refillTestData() {
	clearTestDataFromDb()
	loadTestDataToDb()
}

func clearTestDataFromDb() {
	db, err := sql.Open(dbDriver, dbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "DELETE FROM departments;"
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
