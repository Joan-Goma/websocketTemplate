package controller

import (
	"bufio"
	"fmt"
	"github.com/Joan-Goma/websocketTemplate/models"
	engine "github.com/JoanGTSQ/api"
	"os"
)

// ReadInput read in every moment the console and change maintenance and debug mode
func ReadInput(debug bool) {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		switch input.Text() {
		case "debug":
			debug = !debug
			engine.EnableDebug(debug)
			engine.Info.Println("debug mode changed ->", debug)
		}
	}
}

type DbConnection struct {
	URL      string `json:"dbDirection"`
	User     string `json:"dbUser"`
	Name     string `json:"dbName"`
	Password string `json:"dbPsswd"`
	SslMode  string `json:"sslMode"`
}

var database DbConnection

func InitDatabase(settings DbConnection, debugDB bool) {

	database = settings
	// Create connection with DB
	engine.Debug.Println("Creating connection with DB")
	var err error
	err = models.NewServices(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		database.URL,
		5432,
		database.User,
		database.Password,
		database.Name,
		database.SslMode),
		debugDB)
	if err != nil {
		engine.Error.Fatalln("Could not start db", err)
	}

	// Auto generate new tables or modifications in every start | Use DestructiveReset() to delete all data
	if err := models.AutoMigrate(); err != nil {
		engine.Error.Fatalln("Can not AutoMigrate the database", err)
	}

}
