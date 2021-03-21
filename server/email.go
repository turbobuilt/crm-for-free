package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"
)

func SendEmail(from string, to []string, templatePath string, subject string, data interface{}) {
	auth := smtp.PlainAuth("", SMTP_USERNAME, SMTP_PASSWORD, SMTP_SERVER)

	templateRaw, err := ioutil.ReadFile("templates/" + templatePath)
	tmpl, err := template.New(templatePath).Parse(string(templateRaw))

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	tmpl.Execute(writer, data)
	writer.Flush()

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	// fmt.Println("message is ", buf.String())
	// msg := []byte("To: dane@turbobuilt.com\r\n" +
	// 	"Subject: This is a test email.\r\n" +
	// 	"\r\n" +
	// 	"Hello Dane, \r\n I just wanted you to see this message\r\n")
	// err = smtp.SendMail(SMTP_SERVER+":2525", auth, from, to, msg)

	e := &email.Email{
		To:      to,
		From:    "TurboBuilt <" + FROM_EMAIL + ">",
		Subject: subject,
		HTML:    buf.Bytes(),
		Headers: textproto.MIMEHeader{},
	}
	err = e.Send(SMTP_SERVER+":2525", auth)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sent")
}
