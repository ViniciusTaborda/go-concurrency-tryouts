package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var pathToTemplates = "./cmd/web/templates"

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	Data          map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	// User *data.user
}

func (app *Config) render(
	w http.ResponseWriter,
	r *http.Request,
	templateName string,
	templateData *TemplateData) {

	templateNames := []string{
		fmt.Sprintf("%s/%s", pathToTemplates, templateName),
		fmt.Sprintf("%s/base.layout.gohtml", pathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", pathToTemplates),
	}

	if templateData == nil {
		templateData = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templateNames...)
	if err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, app.AddDefaultData(templateData, r)); err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Config) AddDefaultData(templateData *TemplateData, r *http.Request) *TemplateData {

	templateData.Flash = app.Session.PopString(r.Context(), "flash")
	templateData.Warning = app.Session.PopString(r.Context(), "warning")
	templateData.Error = app.Session.PopString(r.Context(), "error")
	if app.IsAuthenticated(r) {
		templateData.Authenticated = true
		// TODO - Get more user info
	}

	templateData.Now = time.Now()

	return templateData

}

func (app *Config) IsAuthenticated(r *http.Request) bool {

	return app.Session.Exists(r.Context(), "userID")

}
