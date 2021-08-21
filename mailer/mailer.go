package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"mime/quotedprintable"
	"net/smtp"
)

type Sender struct {
	Email    string
	Password string
}

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func NewSender(email, password string) Sender {
	return Sender{email, password}
}

func (sender Sender) Mail(subject, message, dest string) (string, error) {
	host := "smtp.gmail.com"
	port := "smtp.gmail.com:587"

	auth := smtp.PlainAuth("", sender.Email, sender.Password, host)

	err := smtp.SendMail(port, auth, sender.Email, []string{dest}, []byte(message))

	if err != nil {
		return "failed to send", errors.New(err.Error())
	}
	return "Mail sent", nil
}

func (sender Sender) WriteEmail(dest, subject, msg string, Cc []string) string {
	header := make(map[string]string)

	if len(Cc) != 0 {
		var list string
		for _, c := range Cc {
			list += fmt.Sprint(c + ",")
		}
		header["Cc"] = list
	}

	header["From"] = sender.Email
	header["To"] = dest
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", "text/html")

	message := ""
	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMsg bytes.Buffer
	finalMsg := quotedprintable.NewWriter(&encodedMsg)
	finalMsg.Write([]byte(msg))
	finalMsg.Close()

	message += "\r\n" + encodedMsg.String()
	return message
}

func (sender Sender) WriteMessage(senderName, path string) (string, error) {

	t, err := template.ParseFiles(path)
	Check(err)

	data := struct{ Name string }{senderName}

	var tmp bytes.Buffer
	err = t.Execute(&tmp, data)
	Check(err)
	result := tmp.String()

	return result, nil
}
