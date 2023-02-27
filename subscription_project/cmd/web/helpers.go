package main

func (app *Config) sendEmail(message Message) {
	app.WaitGroup.Add(1)
	app.Mailer.MailerChannel <- message
}
