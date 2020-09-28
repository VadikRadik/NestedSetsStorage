package treestorage_test

import (
	"NestedSetsStorage/configs"
	"NestedSetsStorage/treestorage"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

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

	defaultNodes := createTestNodes()

	testCases := []struct {
		name       string
		nodeName   string
		parentName string
		expected   []treestorage.NestedSetsNode
	}{
		{
			name:       "adding of existing nodes",
			nodeName:   "Совет лицея",
			parentName: "Заместитель директора по ВР",
			expected:   defaultNodes,
		},
		{
			name:       "adding of invalid parent nodes case 1",
			nodeName:   "Совет лицея",
			parentName: "Заместитель директора",
			expected:   defaultNodes,
		},
		{
			name:       "adding of invalid parent nodes case 2",
			nodeName:   "Совет лицея",
			parentName: "",
			expected:   defaultNodes,
		},
		{
			name:       "adding of empty name nodes",
			nodeName:   "",
			parentName: "Совет лицея",
			expected:   defaultNodes,
		},
		{
			name:       "addNodeCase1",
			nodeName:   "Общешкольный родительский комитет",
			parentName: "Совет лицея",
			expected:   addNodeCase1(),
		},
		{
			name:       "addNodeCase2",
			nodeName:   "Психолог",
			parentName: "Заместитель директора по ВР",
			expected:   addNodeCase2(),
		},
		{
			name:       "addNodeCase3",
			nodeName:   "Общее собрание трудового коллектива",
			parentName: "Директор",
			expected:   addNodeCase3(),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			s.AddNode(test.nodeName, test.parentName)
			assert.ElementsMatch(t, test.expected, s.GetWholeTree())
		})
	}

	s.AddNode("", "")
	clearTestDataFromDb()
}

func TestRemoveNode(t *testing.T) {

}

func TestMoveNode(t *testing.T) {

}

func TestRenameNode(t *testing.T) {

}

func TestGetParents(t *testing.T) {

}

func TestGetChildren(t *testing.T) {

}

func TestGetWholeTree(t *testing.T) {

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

func addNodeCase1() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 37},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 5, 14},
		{"Благотворительный фонд \"Развитие школы\"", 6, 7},
		{"Ученическое самоуправление", 8, 11},
		{"Ученики", 9, 10},
		{"Общешкольный родительский комитет", 12, 13},
		{"Заместитель директора по информатизации", 15, 18},
		{"Инженегр по ВТ", 16, 17},
		{"Заместитель директора по ВР", 19, 26},
		{"Служба сопровождения", 20, 21},
		{"Методическое объединение педагогов дополнительного образования", 22, 23},
		{"Методическое объединение классных руководителей", 24, 25},
		{"Бухгалтерия", 27, 28},
		{"Педагогический совет", 29, 30},
		{"Заместитель директора по УВР", 31, 34},
		{"Кафедры профильного образования", 32, 33},
		{"Научно-методический совет", 35, 36},
	}
	return nodes
}

func addNodeCase2() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 39},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 5, 14},
		{"Благотворительный фонд \"Развитие школы\"", 6, 7},
		{"Ученическое самоуправление", 8, 11},
		{"Ученики", 9, 10},
		{"Общешкольный родительский комитет", 12, 13},
		{"Заместитель директора по информатизации", 15, 18},
		{"Инженегр по ВТ", 16, 17},
		{"Заместитель директора по ВР", 19, 28},
		{"Служба сопровождения", 20, 21},
		{"Методическое объединение педагогов дополнительного образования", 22, 23},
		{"Методическое объединение классных руководителей", 24, 25},
		{"Психолог", 26, 27},
		{"Бухгалтерия", 29, 30},
		{"Педагогический совет", 31, 32},
		{"Заместитель директора по УВР", 33, 36},
		{"Кафедры профильного образования", 34, 35},
		{"Научно-методический совет", 37, 38},
	}
	return nodes
}

func addNodeCase3() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 41},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 5, 14},
		{"Благотворительный фонд \"Развитие школы\"", 6, 7},
		{"Ученическое самоуправление", 8, 11},
		{"Ученики", 9, 10},
		{"Общешкольный родительский комитет", 12, 13},
		{"Заместитель директора по информатизации", 15, 18},
		{"Инженегр по ВТ", 16, 17},
		{"Заместитель директора по ВР", 19, 28},
		{"Служба сопровождения", 20, 21},
		{"Методическое объединение педагогов дополнительного образования", 22, 23},
		{"Методическое объединение классных руководителей", 24, 25},
		{"Психолог", 26, 27},
		{"Бухгалтерия", 29, 30},
		{"Педагогический совет", 31, 32},
		{"Заместитель директора по УВР", 33, 36},
		{"Кафедры профильного образования", 34, 35},
		{"Научно-методический совет", 37, 38},
		{"Общее собрание трудового коллектива", 39, 40},
	}
	return nodes
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
