package main

import (
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
)

const htmlBody = `
<html>
<head>
   <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
   <title>Hello, World</title>
</head>
<body>
   <p>This is an email using Go</p>
</body>
`

func main() {
	var username, password, from, to, subject, body, attachment_name string
	fmt.Scanf("Username:%s", username)
	fmt.Scanf("password: %s", password)
	fmt.Scanf("Massage from: %s", from)
	fmt.Scanf("Recipient: %s", to)
	fmt.Scanf("Subject of message: %s", subject)
	fmt.Scanf("Email body: %s", body)
	fmt.Scanf("Email attachment (can be scipped): %s", attachment_name)

	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 587
	server.Username = username
	server.Password = password
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Create email
	email := mail.NewMSG()
	email.SetFrom(from)
	email.AddTo(to)
	email.SetSubject(subject)

	email.SetBody(mail.TextHTML, body)
	email.AddAttachment(attachment_name)

	// Send email
	err = email.Send(smtpClient)
	if err != nil {
		log.Fatal(err)
	}
}
