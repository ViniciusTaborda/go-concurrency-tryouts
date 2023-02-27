package main

import (
	"fmt"
	"goroutines/subscription_project/data"
	"goroutines/subscription_project/database"
	"goroutines/subscription_project/session"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const defaultWebPort = "80"

func main() {

	//loggers
	infoLog := log.New(os.Stdout, " - INFO \t", log.Ldate|log.Ltime)
	errorLog := log.New(
		os.Stdout,
		" - ERROR \t",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	db := database.Initialize()
	session := session.Initialize()

	waitGroup := sync.WaitGroup{}

	app := Config{
		Session:   session,
		Db:        db,
		InfoLog:   infoLog,
		ErrorLog:  errorLog,
		WaitGroup: &waitGroup,
		Models:    data.New(db),
	}

	app.Mailer = app.createMail()

	go app.listenForMail()

	// Running in background
	go app.listerForShutdown()

	app.runServer()

}

func (app *Config) runServer() {
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", defaultWebPort),
		Handler: app.GetRouter(),
	}

	app.InfoLog.Println("Starting up Web server...")

	err := server.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) listerForShutdown() {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	//Perform clean up tasks
	app.InfoLog.Println("Cleaning up before shutting down...")

	// Block until Wait group is empty
	app.WaitGroup.Wait()

	//Done after all goroutines are done
	app.Mailer.DoneChannel <- true

	close(app.Mailer.MailerChannel)
	close(app.Mailer.DoneChannel)
	close(app.Mailer.ErrorChannel)

	app.InfoLog.Println("Closing all channels...")

}

func (app *Config) createMail() Mail {
	errorChannel := make(chan error)

	//Buffering
	mailerChannel := make(chan Message, 100)
	mailerDoneChannel := make(chan bool)

	mail := Mail{
		Domain:        "localhost",
		Host:          "localhost",
		Port:          1025,
		Encryption:    "none",
		FromAddress:   "info@mycompany.com",
		FromName:      "Info",
		ErrorChannel:  errorChannel,
		MailerChannel: mailerChannel,
		DoneChannel:   mailerDoneChannel,
	}

	return mail

}
