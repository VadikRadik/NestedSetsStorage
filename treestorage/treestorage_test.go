package treestorage_test

import (
	"NestedSetsStorage/configs"
	"NestedSetsStorage/treestorage"
	"log"
	"testing"

	"github.com/BurntSushi/toml"
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
	s := &treestorage.NestedSetsStorage{
		DbConnectionString: dbConnectionString,
		DbDriver:           dbDriver,
	}
	s.AddNode("", "")
}

func TestRemoveNode(t *testing.T) {

}

func TestMoveNode(t *testing.T) {

}

func TestRenameNode(t *testing.T) {

}
