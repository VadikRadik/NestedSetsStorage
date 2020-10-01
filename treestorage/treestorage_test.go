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
	"github.com/stretchr/testify/assert"
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

func TestNestedSetsStorage_GetParents(t *testing.T) {
	refillTestData()

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "getting parents for not existing node",
			args: args{"Заместитель директора"},
			want: []string{},
		},
		{
			name: "getting parents for invalid name node",
			args: args{""},
			want: []string{},
		},
		{
			name: "getting parents for root",
			args: args{"Директор"},
			want: []string{},
		},
		{
			name: "getting parents for node",
			args: args{"Ученики"},
			want: []string{"Директор", "Совет лицея", "Ученическое самоуправление"},
		},
		{
			name: "getting parents for node case 2",
			args: args{"Служба сопровождения"},
			want: []string{"Директор", "Заместитель директора по ВР"},
		},
	}

	s := &treestorage.NestedSetsStorage{
		DbConnectionString: dbConnectionString,
		DbDriver:           dbDriver,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := s.GetParents(tt.args.name)
			assert.ElementsMatch(t, tt.want, got)
		})
	}

	clearTestDataFromDb()
}

func TestNestedSetsStorage_GetChildren(t *testing.T) {
	refillTestData()

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "getting children for not existing node",
			args: args{"Заместитель директора"},
			want: []string{},
		},
		{
			name: "getting children for invalid name node",
			args: args{""},
			want: []string{},
		},
		{
			name: "getting children for node without children",
			args: args{"Бухгалтерия"},
			want: []string{},
		},
		{
			name: "getting children for root",
			args: args{"Директор"},
			want: getChildrenCase1(),
		},
		{
			name: "getting children for node",
			args: args{"Совет лицея"},
			want: []string{"Благотворительный фонд \"Развитие школы\"", "Ученическое самоуправление", "Ученики"},
		},
	}

	s := &treestorage.NestedSetsStorage{
		DbConnectionString: dbConnectionString,
		DbDriver:           dbDriver,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := s.GetChildren(tt.args.name)
			assert.ElementsMatch(t, tt.want, got)
		})
	}

	clearTestDataFromDb()
}

func TestNestedSetsStorage_GetWholeTree(t *testing.T) {
	refillTestData()
	defaultNodes := createTestNodes()

	tests := []struct {
		name string
		want []treestorage.NestedSetsNode
	}{
		{
			name: "getting whole tree",
			want: defaultNodes,
		},
	}

	s := &treestorage.NestedSetsStorage{
		DbConnectionString: dbConnectionString,
		DbDriver:           dbDriver,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := s.GetWholeTree()
			assert.ElementsMatch(t, tt.want, got)
		})
	}

	clearTestDataFromDb()
}

func TestNestedSetsStorage_AddNode(t *testing.T) {
	refillTestData()
	defaultNodes := createTestNodes()

	type args struct {
		name   string
		parent string
	}
	tests := []struct {
		name string
		args args
		want []treestorage.NestedSetsNode
	}{
		{
			name: "adding existing nodes",
			args: args{"Совет лицея", "Заместитель директора по ВР"},
			want: defaultNodes,
		},
		{
			name: "adding invalid parent nodes case 1",
			args: args{"Совет лицея", "Заместитель директора"},
			want: defaultNodes,
		},
		{
			name: "adding invalid parent nodes case 2",
			args: args{"Совет лицея", ""},
			want: defaultNodes,
		},
		{
			name: "adding empty name nodes",
			args: args{"", "Совет лицея"},
			want: defaultNodes,
		},
		{
			name: "addNodeCase1",
			args: args{"Общешкольный родительский комитет", "Совет лицея"},
			want: addNodeCase1(),
		},
		{
			name: "addNodeCase2",
			args: args{"Психолог", "Заместитель директора по ВР"},
			want: addNodeCase2(),
		},
		{
			name: "addNodeCase3",
			args: args{"Общее собрание трудового коллектива", "Директор"},
			want: addNodeCase3(),
		},
	}

	s := &treestorage.NestedSetsStorage{
		DbConnectionString: dbConnectionString,
		DbDriver:           dbDriver,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.AddNode(tt.args.name, tt.args.parent)
			got, _ := s.GetWholeTree()
			assert.ElementsMatch(t, tt.want, got)
		})
	}

	clearTestDataFromDb()
}

func TestNestedSetsStorage_RemoveNode(t *testing.T) {
	refillTestData()
	defaultNodes := createTestNodes()

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want []treestorage.NestedSetsNode
	}{
		{
			name: "removing not existing nodes",
			args: args{"Психолог"},
			want: defaultNodes,
		},
		{
			name: "removing invalid name nodes",
			args: args{""},
			want: defaultNodes,
		},
		{
			name: "removeNodeCase1",
			args: args{"Служба сопровождения"},
			want: removeNodeCase1(),
		},
		{
			name: "removeNodeCase2",
			args: args{"Совет лицея"},
			want: removeNodeCase2(),
		},
		{
			name: "removeNodeCase3",
			args: args{"Директор"},
			want: removeNodeCase3(),
		},
	}

	s := &treestorage.NestedSetsStorage{
		DbConnectionString: dbConnectionString,
		DbDriver:           dbDriver,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.RemoveNode(tt.args.name)
			got, _ := s.GetWholeTree()
			assert.ElementsMatch(t, tt.want, got)
		})
	}

	clearTestDataFromDb()
}

func TestNestedSetsStorage_MoveNode(t *testing.T) {
	refillTestData()
	defaultNodes := createTestNodes()

	type args struct {
		name      string
		newParent string
	}
	tests := []struct {
		name string
		args args
		want []treestorage.NestedSetsNode
	}{
		{
			name: "moving invalid node",
			args: args{"", "Заместитель директора по ВР"},
			want: defaultNodes,
		},
		{
			name: "moving not existing node",
			args: args{"Психолог", "Заместитель директора по ВР"},
			want: defaultNodes,
		},
		{
			name: "moving to invalid parent",
			args: args{"Заместитель директора по ВР", ""},
			want: defaultNodes,
		},
		{
			name: "moving to not existing node",
			args: args{"Заместитель директора по ВР", "Психолог"},
			want: defaultNodes,
		},
		{
			name: "not modifying moving",
			args: args{"Благотворительный фонд \"Развитие школы\"", "Совет лицея"},
			want: defaultNodes,
		},
		{
			name: "moving node case 1: left direction move to the right parent node",
			args: args{"Педагогический совет", "Заместитель директора по ВР"},
			want: moveNodeCase1(),
		},
		{
			name: "moving node case 2: right direction move to the left parent node",
			args: args{"Совет лицея", "Заместитель директора по ВР"},
			want: moveNodeCase2(),
		},
		{
			name: "moving node case 3: right direction move to the left parent node",
			args: args{"Методическое объединение педагогов дополнительного образования", "Методическое объединение классных руководителей"},
			want: moveNodeCase3(),
		},
		{
			name: "moving node case 4: left direction move to hte right parent node",
			args: args{"Педагогический совет", "Заместитель директора по ВР"},
			want: moveNodeCase4(),
		},
		{
			name: "moving node case 5: moving down along branch",
			args: args{"Ученическое самоуправление", "Ученики"},
			want: moveNodeCase5(),
		},
		{
			name: "moving node case 6: moving up along branch to the right parent node",
			args: args{"Ученики", "Совет лицея"},
			want: moveNodeCase6(),
		},
		{
			name: "moving node case 7: moving up along branch to the left parent node",
			args: args{"Совет лицея", "Директор"},
			want: moveNodeCase7(),
		},
		{
			name: "moving node case 8: moving down along branch",
			args: args{"Совет лицея", "Ученики"},
			want: moveNodeCase8(),
		},
	}

	s := &treestorage.NestedSetsStorage{
		DbConnectionString: dbConnectionString,
		DbDriver:           dbDriver,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refillTestData()
			err := s.MoveNode(tt.args.name, tt.args.newParent)
			if err != nil {
				log.Println(err)
			}
			got, _ := s.GetWholeTree()
			assert.ElementsMatch(t, tt.want, got)
		})
	}

	clearTestDataFromDb()
}

func TestNestedSetsStorage_RenameNode(t *testing.T) {
	refillTestData()
	defaultNodes := createTestNodes()

	type args struct {
		name    string
		newName string
	}
	tests := []struct {
		name string
		args args
		want []treestorage.NestedSetsNode
	}{
		{
			name: "renaming invalid node",
			args: args{"", "Заместитель директора"},
			want: defaultNodes,
		},
		{
			name: "renaming not existing node",
			args: args{"Психолог", "Заместитель директора"},
			want: defaultNodes,
		},
		{
			name: "renaming node",
			args: args{"Заместитель директора по ВР", "Заместитель директора по воспитательной работе"},
			want: renameNodeCase(),
		},
	}

	s := &treestorage.NestedSetsStorage{
		DbConnectionString: dbConnectionString,
		DbDriver:           dbDriver,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.RenameNode(tt.args.name, tt.args.newName)
			got, _ := s.GetWholeTree()
			assert.ElementsMatch(t, tt.want, got)
		})
	}

	clearTestDataFromDb()
}

func TestNestedSetsStorage_AddRoot(t *testing.T) {
	refillTestData()
	defaultNodes := createTestNodes()

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want []treestorage.NestedSetsNode
	}{
		{
			name: "adding an invalid node",
			args: args{""},
			want: defaultNodes,
		},
		{
			name: "adding an existing node",
			args: args{"Заместитель директора по ВР"},
			want: defaultNodes,
		},
		{
			name: "adding a node",
			args: args{"Директор колледжа"},
			want: addingRootCase(),
		},
	}

	s := &treestorage.NestedSetsStorage{
		DbConnectionString: dbConnectionString,
		DbDriver:           dbDriver,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.AddRoot(tt.args.name)
			got, _ := s.GetWholeTree()
			assert.ElementsMatch(t, tt.want, got)
		})
	}

	clearTestDataFromDb()

	tests = []struct {
		name string
		args args
		want []treestorage.NestedSetsNode
	}{
		{
			name: "adding to empty tree",
			args: args{"Директор колледжа"},
			want: []treestorage.NestedSetsNode{{"Директор колледжа", 0, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.AddRoot(tt.args.name)
			got, _ := s.GetWholeTree()
			assert.ElementsMatch(t, tt.want, got)
		})
	}

	clearTestDataFromDb()
}

func loadTestDataToDb() {
	nodes := createTestNodes()

	db, err := sql.Open(dbDriver, dbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fullQuery := `INSERT INTO 
					nodes (name, node_left, node_right) 
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

	query := "DELETE FROM nodes;"
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

// added "Общешкольный родительский комитет" to "Совет лицея"
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

// added "Психолог" to "Заместитель директора по ВР"
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

// added "Общее собрание трудового коллектива" to "Директор"
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

// Removing a node without children // removed "Служба сопровождения"
func removeNodeCase1() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 33},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 5, 12},
		{"Благотворительный фонд \"Развитие школы\"", 6, 7},
		{"Ученическое самоуправление", 8, 11},
		{"Ученики", 9, 10},
		{"Заместитель директора по информатизации", 13, 16},
		{"Инженегр по ВТ", 14, 15},
		{"Заместитель директора по ВР", 17, 22},
		{"Методическое объединение педагогов дополнительного образования", 18, 19},
		{"Методическое объединение классных руководителей", 20, 21},
		{"Бухгалтерия", 23, 24},
		{"Педагогический совет", 25, 26},
		{"Заместитель директора по УВР", 27, 30},
		{"Кафедры профильного образования", 28, 29},
		{"Научно-методический совет", 31, 32},
	}
	return nodes
}

// Removing a node with children // removed "Совет лицея"
func removeNodeCase2() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 31},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Благотворительный фонд \"Развитие школы\"", 5, 6},
		{"Ученическое самоуправление", 7, 10},
		{"Ученики", 8, 9},
		{"Заместитель директора по информатизации", 11, 14},
		{"Инженегр по ВТ", 12, 13},
		{"Заместитель директора по ВР", 15, 20},
		{"Методическое объединение педагогов дополнительного образования", 16, 17},
		{"Методическое объединение классных руководителей", 18, 19},
		{"Бухгалтерия", 21, 22},
		{"Педагогический совет", 23, 24},
		{"Заместитель директора по УВР", 25, 28},
		{"Кафедры профильного образования", 26, 27},
		{"Научно-методический совет", 29, 30},
	}
	return nodes
}

// Removing a tree root // removed "Директор"
func removeNodeCase3() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Заместитель директора по АХЧ", 0, 3},
		{"Обслуживающий персонал", 1, 2},
		{"Благотворительный фонд \"Развитие школы\"", 4, 5},
		{"Ученическое самоуправление", 6, 9},
		{"Ученики", 7, 8},
		{"Заместитель директора по информатизации", 10, 13},
		{"Инженегр по ВТ", 11, 12},
		{"Заместитель директора по ВР", 14, 19},
		{"Методическое объединение педагогов дополнительного образования", 15, 16},
		{"Методическое объединение классных руководителей", 17, 18},
		{"Бухгалтерия", 20, 21},
		{"Педагогический совет", 22, 23},
		{"Заместитель директора по УВР", 24, 27},
		{"Кафедры профильного образования", 25, 26},
		{"Научно-методический совет", 28, 29},
	}
	return nodes
}

// left direction move to the right parent node
func moveNodeCase1() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 35},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 5, 12},
		{"Благотворительный фонд \"Развитие школы\"", 6, 7},
		{"Ученическое самоуправление", 8, 11},
		{"Ученики", 9, 10},
		{"Заместитель директора по информатизации", 13, 16},
		{"Инженегр по ВТ", 14, 15},
		{"Заместитель директора по ВР", 17, 26},
		{"Служба сопровождения", 18, 19},
		{"Методическое объединение педагогов дополнительного образования", 20, 21},
		{"Методическое объединение классных руководителей", 22, 23},
		{"Бухгалтерия", 27, 28},
		{"Педагогический совет", 24, 25},
		{"Заместитель директора по УВР", 29, 32},
		{"Кафедры профильного образования", 30, 31},
		{"Научно-методический совет", 33, 34},
	}
	return nodes
}

// right direction move to the left parent node
func moveNodeCase2() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 35},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 16, 17},
		{"Благотворительный фонд \"Развитие школы\"", 5, 6},
		{"Ученическое самоуправление", 7, 10},
		{"Ученики", 8, 9},
		{"Заместитель директора по информатизации", 11, 14},
		{"Инженегр по ВТ", 12, 13},
		{"Заместитель директора по ВР", 15, 24},
		{"Служба сопровождения", 18, 19},
		{"Методическое объединение педагогов дополнительного образования", 20, 21},
		{"Методическое объединение классных руководителей", 22, 23},
		{"Бухгалтерия", 25, 26},
		{"Педагогический совет", 27, 28},
		{"Заместитель директора по УВР", 29, 32},
		{"Кафедры профильного образования", 30, 31},
		{"Научно-методический совет", 33, 34},
	}
	return nodes
}

// move node in the same branch
func moveNodeCase3() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 35},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 5, 12},
		{"Благотворительный фонд \"Развитие школы\"", 6, 7},
		{"Ученическое самоуправление", 8, 11},
		{"Ученики", 9, 10},
		{"Заместитель директора по информатизации", 13, 16},
		{"Инженегр по ВТ", 14, 15},
		{"Заместитель директора по ВР", 17, 24},
		{"Служба сопровождения", 18, 19},
		{"Методическое объединение педагогов дополнительного образования", 21, 22},
		{"Методическое объединение классных руководителей", 20, 23},
		{"Бухгалтерия", 25, 26},
		{"Педагогический совет", 27, 28},
		{"Заместитель директора по УВР", 29, 32},
		{"Кафедры профильного образования", 30, 31},
		{"Научно-методический совет", 33, 34},
	}
	return nodes
}

// left moving to the right parent node
func moveNodeCase4() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 35},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 5, 12},
		{"Благотворительный фонд \"Развитие школы\"", 6, 7},
		{"Ученическое самоуправление", 8, 11},
		{"Ученики", 9, 10},
		{"Заместитель директора по информатизации", 13, 16},
		{"Инженегр по ВТ", 14, 15},
		{"Заместитель директора по ВР", 17, 26},
		{"Служба сопровождения", 18, 19},
		{"Методическое объединение педагогов дополнительного образования", 20, 21},
		{"Методическое объединение классных руководителей", 22, 23},
		{"Бухгалтерия", 27, 28},
		{"Педагогический совет", 24, 25},
		{"Заместитель директора по УВР", 29, 32},
		{"Кафедры профильного образования", 30, 31},
		{"Научно-методический совет", 33, 34},
	}
	return nodes
}

// moving down alang branch
func moveNodeCase5() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 35},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 5, 12},
		{"Благотворительный фонд \"Развитие школы\"", 6, 7},
		{"Ученическое самоуправление", 9, 10},
		{"Ученики", 8, 11},
		{"Заместитель директора по информатизации", 13, 16},
		{"Инженегр по ВТ", 14, 15},
		{"Заместитель директора по ВР", 17, 24},
		{"Служба сопровождения", 18, 19},
		{"Методическое объединение педагогов дополнительного образования", 20, 21},
		{"Методическое объединение классных руководителей", 22, 23},
		{"Бухгалтерия", 25, 26},
		{"Педагогический совет", 27, 28},
		{"Заместитель директора по УВР", 29, 32},
		{"Кафедры профильного образования", 30, 31},
		{"Научно-методический совет", 33, 34},
	}
	return nodes
}

// moving up along branch to the right parent node
func moveNodeCase6() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 35},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 5, 12},
		{"Благотворительный фонд \"Развитие школы\"", 6, 7},
		{"Ученическое самоуправление", 8, 9},
		{"Ученики", 10, 11},
		{"Заместитель директора по информатизации", 13, 16},
		{"Инженегр по ВТ", 14, 15},
		{"Заместитель директора по ВР", 17, 24},
		{"Служба сопровождения", 18, 19},
		{"Методическое объединение педагогов дополнительного образования", 20, 21},
		{"Методическое объединение классных руководителей", 22, 23},
		{"Бухгалтерия", 25, 26},
		{"Педагогический совет", 27, 28},
		{"Заместитель директора по УВР", 29, 32},
		{"Кафедры профильного образования", 30, 31},
		{"Научно-методический совет", 33, 34},
	}
	return nodes
}

// moving up along branch to the left parent node
func moveNodeCase7() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 35},
		{"Заместитель директора по АХЧ", 3, 6},
		{"Обслуживающий персонал", 4, 5},
		{"Совет лицея", 1, 2},
		{"Благотворительный фонд \"Развитие школы\"", 7, 8},
		{"Ученическое самоуправление", 9, 12},
		{"Ученики", 10, 11},
		{"Заместитель директора по информатизации", 13, 16},
		{"Инженегр по ВТ", 14, 15},
		{"Заместитель директора по ВР", 17, 24},
		{"Служба сопровождения", 18, 19},
		{"Методическое объединение педагогов дополнительного образования", 20, 21},
		{"Методическое объединение классных руководителей", 22, 23},
		{"Бухгалтерия", 25, 26},
		{"Педагогический совет", 27, 28},
		{"Заместитель директора по УВР", 29, 32},
		{"Кафедры профильного образования", 30, 31},
		{"Научно-методический совет", 33, 34},
	}
	return nodes
}

// moving down alnog pranch
func moveNodeCase8() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 35},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 9, 10},
		{"Благотворительный фонд \"Развитие школы\"", 5, 6},
		{"Ученическое самоуправление", 7, 12},
		{"Ученики", 8, 11},
		{"Заместитель директора по информатизации", 13, 16},
		{"Инженегр по ВТ", 14, 15},
		{"Заместитель директора по ВР", 17, 24},
		{"Служба сопровождения", 18, 19},
		{"Методическое объединение педагогов дополнительного образования", 20, 21},
		{"Методическое объединение классных руководителей", 22, 23},
		{"Бухгалтерия", 25, 26},
		{"Педагогический совет", 27, 28},
		{"Заместитель директора по УВР", 29, 32},
		{"Кафедры профильного образования", 30, 31},
		{"Научно-методический совет", 33, 34},
	}
	return nodes
}

func renameNodeCase() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 35},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 5, 12},
		{"Благотворительный фонд \"Развитие школы\"", 6, 7},
		{"Ученическое самоуправление", 8, 11},
		{"Ученики", 9, 10},
		{"Заместитель директора по информатизации", 13, 16},
		{"Инженегр по ВТ", 14, 15},
		{"Заместитель директора по воспитательной работе", 17, 24},
		{"Служба сопровождения", 18, 19},
		{"Методическое объединение педагогов дополнительного образования", 20, 21},
		{"Методическое объединение классных руководителей", 22, 23},
		{"Бухгалтерия", 25, 26},
		{"Педагогический совет", 27, 28},
		{"Заместитель директора по УВР", 29, 32},
		{"Кафедры профильного образования", 30, 31},
		{"Научно-методический совет", 33, 34},
	}
	return nodes
}

// get children for root
func getChildrenCase1() []string {
	nodes := []string{
		"Заместитель директора по АХЧ",
		"Обслуживающий персонал",
		"Совет лицея",
		"Благотворительный фонд \"Развитие школы\"",
		"Ученическое самоуправление",
		"Ученики",
		"Заместитель директора по информатизации",
		"Инженегр по ВТ",
		"Заместитель директора по ВР",
		"Служба сопровождения",
		"Методическое объединение педагогов дополнительного образования",
		"Методическое объединение классных руководителей",
		"Бухгалтерия",
		"Педагогический совет",
		"Заместитель директора по УВР",
		"Кафедры профильного образования",
		"Научно-методический совет",
	}
	return nodes
}

func addingRootCase() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 35},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 5, 12},
		{"Благотворительный фонд \"Развитие школы\"", 6, 7},
		{"Ученическое самоуправление", 8, 11},
		{"Ученики", 9, 10},
		{"Заместитель директора по информатизации", 13, 16},
		{"Инженегр по ВТ", 14, 15},
		{"Заместитель директора по ВР", 17, 24},
		{"Служба сопровождения", 18, 19},
		{"Методическое объединение педагогов дополнительного образования", 20, 21},
		{"Методическое объединение классных руководителей", 22, 23},
		{"Бухгалтерия", 25, 26},
		{"Педагогический совет", 27, 28},
		{"Заместитель директора по УВР", 29, 32},
		{"Кафедры профильного образования", 30, 31},
		{"Научно-методический совет", 33, 34},
		{"Директор колледжа", 36, 37},
	}
	return nodes
}

func createTestNodes() []treestorage.NestedSetsNode {
	nodes := []treestorage.NestedSetsNode{
		{"Директор", 0, 35},
		{"Заместитель директора по АХЧ", 1, 4},
		{"Обслуживающий персонал", 2, 3},
		{"Совет лицея", 5, 12},
		{"Благотворительный фонд \"Развитие школы\"", 6, 7},
		{"Ученическое самоуправление", 8, 11},
		{"Ученики", 9, 10},
		{"Заместитель директора по информатизации", 13, 16},
		{"Инженегр по ВТ", 14, 15},
		{"Заместитель директора по ВР", 17, 24},
		{"Служба сопровождения", 18, 19},
		{"Методическое объединение педагогов дополнительного образования", 20, 21},
		{"Методическое объединение классных руководителей", 22, 23},
		{"Бухгалтерия", 25, 26},
		{"Педагогический совет", 27, 28},
		{"Заместитель директора по УВР", 29, 32},
		{"Кафедры профильного образования", 30, 31},
		{"Научно-методический совет", 33, 34},
	}
	return nodes
}
