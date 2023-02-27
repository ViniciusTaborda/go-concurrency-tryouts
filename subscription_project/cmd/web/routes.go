package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) GetRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)

	mux.Get("/", app.HomePage)
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.Login)
	mux.Get("/logout", app.Logout)
	mux.Get("/register", app.RegisterPage)
	mux.Post("/register", app.Register)
	mux.Get("/activate-account", app.ActivateAccount)

	mux.Get("/test-email", func(w http.ResponseWriter, r *http.Request) {
		mail := Mail{
			Domain:       "localhost",
			Host:         "localhost",
			Port:         1025,
			Encryption:   "none",
			FromAddress:  "info@mycompany.com",
			FromName:     "Info",
			ErrorChannel: make(chan error),
		}

		message := Message{
			To:      "me@here.com",
			Subject: "Test mail",
			Data:    "Hello world!",
		}

		mail.sendMail(message, make(chan error))
	})

	return mux
}
