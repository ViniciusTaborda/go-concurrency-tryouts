package main

import (
	"net/http"
)

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) Login(w http.ResponseWriter, r *http.Request) {

	_ = app.Session.RenewToken(r.Context())

	err := r.ParseForm()

	if err != nil {
		app.ErrorLog.Println(err)
	}

	// Get email and password from form
	email := r.Form.Get("email")

	user, err := app.Models.User.GetByEmail(email)

	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials...")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	password := r.Form.Get("password")

	isPasswordValid, err := user.PasswordMatches(password)

	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials...")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !isPasswordValid {

		message := Message{
			To:      email,
			Subject: "Failed logged in attempt",
			Data:    "Invalid login...",
		}

		app.sendEmail(message)

		app.Session.Put(r.Context(), "error", "Invalid credentials...")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//Looks alright, log user in
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)

	app.Session.Put(r.Context(), "flash", "Successful login!")

	// Redirect user
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {

	// Clean up session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Config) Register(w http.ResponseWriter, r *http.Request) {

}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {

}
