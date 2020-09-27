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

	log.Println("storage started")
}
