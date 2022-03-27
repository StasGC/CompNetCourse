package main

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
)

const address = "smtp.gmail.com:587"
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

func SendMail(addr, from, subject, body string, to []string) error {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Mail(r.Replace(from)); err != nil {
		fmt.Println("Error during 'MAIL FROM' connection:", err)
		return err
	}

	for i := range to {
		to[i] = r.Replace(to[i])
		if err = client.Rcpt(to[i]); err != nil {
			fmt.Printf("Error during 'RCPT TO' connection with recipient %v.\n"+
				"Error: %s", to[i], err)
			return err
		}
	}

	writeCloser, err := client.Data()
	if err != nil {
		fmt.Println("Error during 'DATA' connection:", err)
		return err
	}

	msg := "From: " + from + "\r\n" +
		"To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	_, err = writeCloser.Write([]byte(msg))
	if err != nil {
		fmt.Println("Error during 'Write':", err)
		return err
	}

	if err = writeCloser.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func main() {
	var username, password, from, recipients, subject, body string
	fmt.Scanf("Username:%s", username)
	fmt.Scanf("password: %s", password)
	fmt.Scanf("Massage from: %s", from)
	fmt.Scanf("Recipients (replaced by space): %s", recipients)
	fmt.Scanf("Subject of message: %s", subject)
	fmt.Scanf("Email body: %s", body)

	to := strings.Split(recipients, " ")

	err := SendMail(address, from, subject, body, to)
	if err != nil {
		fmt.Println("Error in Send mail:", err)
		return
	}
}
