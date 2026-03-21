package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

func getSMTPServer(email string) (host, port string) {
	if strings.Contains(email, "gmail.com") {
		return "smtp.gmail.com", "587"
	}
	if strings.Contains(email, "qq.com") {
		return "smtp.qq.com", "587"
	}
	if strings.Contains(email, "163.com") {
		return "smtp.163.com", "25"
	}
	return "smtp." + strings.Split(email, "@")[1], "587"
}

func SendEmail(msg string, from string, password string, to string) error {
	host, port := getSMTPServer(from)

	auth := smtp.PlainAuth("", from, password, host)

	message := []byte(
		fmt.Sprintf("From: MyDaemon <%s>\r\n", from) +
			fmt.Sprintf("To: %s\r\n", to) +
			"Subject: Process Result\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			msg,
	)

	addr := fmt.Sprintf("%s:%s", host, port)

	return smtp.SendMail(addr, auth, from, []string{to}, message)
}
