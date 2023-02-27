package main

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	simplemail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain        string
	Host          string
	Port          int
	Username      string
	Password      string
	Encryption    string
	FromAddress   string
	FromName      string
	WaitGroup     *sync.WaitGroup
	MailerChannel (chan Message)
	ErrorChannel  (chan error)
	DoneChannel   (chan bool)
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
	Template    string
}

func (app *Config) listenForMail() {
	for {
		select {
		case message := <-app.Mailer.MailerChannel:
			go app.Mailer.sendMail(message, app.Mailer.ErrorChannel)
		case err := <-app.Mailer.ErrorChannel:
			// Just logging but actually should retry or do something else
			app.ErrorLog.Println(err)
		case <-app.Mailer.DoneChannel:
			//Stop listening for emails...
			return
		}
	}
}

func (mail *Mail) sendMail(message Message, errorChannel chan error) {

	defer mail.WaitGroup.Done()

	if message.Template == "" {
		message.Template = "mail"
	}

	if message.From == "" {
		message.From = mail.FromAddress
	}

	if message.FromName == "" {
		message.FromName = mail.FromName
	}

	data := map[string]any{
		"message": message.Data,
	}

	message.DataMap = data

	formattedMessage, err := mail.buildHTMLMessage(message)
	if err != nil {
		errorChannel <- err
	}

	plainMessage, err := mail.buildPlainTextMessage(message)
	if err != nil {
		errorChannel <- err
	}

	server := simplemail.NewSMTPClient()

	server.Host = mail.Host
	server.Port = mail.Port
	server.Username = mail.Username
	server.Password = mail.Password
	server.Encryption = mail.getEncryption(mail.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		errorChannel <- err
	}

	email := simplemail.NewMSG()
	email.SetFrom(message.From).AddTo(message.To).SetSubject(message.Subject)
	email.SetBody(simplemail.TextPlain, plainMessage)
	email.AddAlternative(simplemail.TextHTML, formattedMessage)

	if len(message.Attachments) > 0 {
		for _, attch := range message.Attachments {
			email.AddAttachment(attch)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		errorChannel <- err
	}
}

func (mail *Mail) buildHTMLMessage(message Message) (string, error) {
	templateToRender := fmt.Sprintf(
		"./cmd/web/templates/%s.html.gohtml",
		message.Template,
	)

	tmplt, err := template.New("email-html").ParseFiles(templateToRender)

	if err != nil {
		return "", err
	}

	var templateBuffer bytes.Buffer

	if err = tmplt.ExecuteTemplate(&templateBuffer, "body", message.DataMap); err != nil {
		return "", err
	}

	formattedMessage := templateBuffer.String()
	formattedMessage, err = mail.inlineCSS(formattedMessage)

	if err != nil {
		return "", err
	}

	return formattedMessage, nil

}

func (mail *Mail) inlineCSS(text string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	premailer, err := premailer.NewPremailerFromString(text, &options)
	if err != nil {
		return "", err
	}

	htmlText, err := premailer.Transform()
	if err != nil {
		return "", err
	}

	return htmlText, nil

}

func (mail *Mail) buildPlainTextMessage(message Message) (string, error) {
	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.plain.gohtml", message.Template)

	tmplt, err := template.New("email-plain").ParseFiles(templateToRender)

	if err != nil {
		return "", err
	}

	var templateBuffer bytes.Buffer

	if err = tmplt.ExecuteTemplate(&templateBuffer, "body", message.DataMap); err != nil {
		return "", err
	}

	plainMessage := templateBuffer.String()

	return plainMessage, nil

}

func (mail *Mail) getEncryption(encryptionName string) simplemail.Encryption {

	switch encryptionName {
	case "tlw":
		return simplemail.EncryptionSTARTTLS
	case "ssl":
		return simplemail.EncryptionSSL
	case "none":
		return simplemail.EncryptionNone
	default:
		return simplemail.EncryptionSTARTTLS
	}

}
