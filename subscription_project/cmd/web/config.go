package main

import (
	"database/sql"
	"goroutines/subscription_project/data"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
)

type Config struct {
	Session   *scs.SessionManager
	Db        *sql.DB
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
	WaitGroup *sync.WaitGroup
	Models    data.Models
	Mailer    Mail
}
