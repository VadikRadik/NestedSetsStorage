package main

import (
	"NestedSetsStorage/configs"
	"NestedSetsStorage/dbmigrate"
	"flag"
	"log"

	"github.com/BurntSushi/toml"
)

func main() {
	var config = new(configs.Config)
	_, err := toml.DecodeFile("configs/config.toml", config)
	if err != nil {
		log.Fatal(err)
	}

	isMigrate := flag.Bool("dbmigrate", false, "runs version migration for data base")
	flag.Parse()
	if *isMigrate == true {
		log.Println("db version migration started")
		err := dbmigrate.Migrate(config)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	/*refillTestData(config.DbDriver, config.DbConnectionSting)
	s := treestorage.NestedSetsStorage{
		DbConnectionString: config.DbConnectionSting,
		DbDriver:           config.DbDriver}
	log.Println("storage started")
	for _, node := range s.GetWholeTree() {
		log.Println(node)
	}

	db, err := sql.Open(config.DbDriver, config.DbConnectionSting)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()*/
}

/*func refillTestData(dbDriver, dbConnectionString string) {
	clearTestDataFromDb(dbDriver, dbConnectionString)
	loadTestDataToDb(dbDriver, dbConnectionString)
}

func clearTestDataFromDb(dbDriver, dbConnectionString string) {
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

func loadTestDataToDb(dbDriver, dbConnectionString string) {
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
}*/
