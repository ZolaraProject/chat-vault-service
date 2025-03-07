package main

import (
	"log"
	"os"

	server "github.com/ZolaraProject/chat-vault-service/chatvault"
)

const (
	defaultDbPort     string = "5432"
	defaultDbHostname string = "postgres.database.svc.cluster.local"
	defaultDbUser     string = "postgres"
)

func main() {
	var ok bool
	server.DbPort, ok = os.LookupEnv("DB_PORT")
	if !ok {
		log.Printf("Warning: could not read $DB_PORT, starting server with default DB port (%s)", defaultDbPort)
		server.DbPort = defaultDbPort
	}
	server.DbHostname, ok = os.LookupEnv("DB_HOSTNAME")
	if !ok {
		log.Printf("Warning: could not read $DB_HOSTNAME, starting server with default DB hostname (%s)", defaultDbHostname)
		server.DbHostname = defaultDbHostname
	}
	server.DbUser, ok = os.LookupEnv("DB_USERNAME")
	if !ok {
		log.Printf("Warning: could not read $DB_USERNAME, starting server with default DB username (%s)", defaultDbUser)
		server.DbUser = defaultDbUser
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if len(dbPassword) == 0 {
		log.Fatal("Error: could not read $DB_PASSWORD")
	}
	server.DbPassword = dbPassword
	server.Run()
}
