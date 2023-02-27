package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func Initialize() *sql.DB {
	conn := mustGetDatabaseConnection()

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println("Could not complete Database initialization!")
			fmt.Println("Exiting...")
		}
	}()

	if conn == nil {
		panic("Could not connect to the database!")
	}

	return conn
}

func mustGetDatabaseConnection() *sql.DB {
	connectionAttempts := 3

	defaultConnectionString := os.Getenv("DSN")

	for i := 0; i <= connectionAttempts; i++ {
		connection, err := openDatabase(defaultConnectionString)
		if err != nil {
			log.Println("Postgres not ready yet...")
			log.Println(err)
		} else {
			return connection
		}

		log.Println("Awaiting one second before trying again...")
		time.Sleep(1 * time.Second)
	}

	return nil
}

func openDatabase(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectionString)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
