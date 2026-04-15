package main

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func GetSMTPServer(email string) (host, port string) {
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

func SendEmail(email EmailConfig, fromTitle string, subjectTitle string, msg string) error {
	host, port := GetSMTPServer(email.From)
	auth := smtp.PlainAuth("", email.From, email.Pwd, host)
	message := []byte(
		fmt.Sprintf("From: %s <%s>\r\n", fromTitle, email.From) +
			fmt.Sprintf("To: %s\r\n", email.To) +
			fmt.Sprintf("Subject: %s\r\n", subjectTitle) +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			msg,
	)
	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, auth, email.From, []string{email.To}, message)
}

func ListEmail(config *Config, extra bool) {
	size := len(config.Emails)
	if extra {
		fmt.Println("[0] Quit and save")
	}
	for i, email := range config.Emails {
		fmt.Printf("[%d] %s -> %s\n", i+1, email.From, email.To)
	}
	if extra && size > 0 {
		fmt.Printf("[%d] Remove all\n", size+1)
	}
}

func HandleEmailAdd(args *Args, config *Config) {
	newEmail := EmailConfig{
		From: args.EmailAddFrom,
		Pwd:  args.EmailAddPwd,
		To:   args.EmailAddTo,
	}
	config.Emails = append(config.Emails, newEmail)
	SaveConfig(config)
}

func HandleEmailLs(args *Args, config *Config) {
	ListEmail(config, false)
}

func HandleEmailRm(args *Args, config *Config) {
	var index int
	var size int
	var changed = false
	for {
		ListEmail(config, true)
		size = len(config.Emails)
		fmt.Print(PromptRemovingEmail)
		fmt.Scanf("%d", &index)
		if index == 0 {
			if changed {
				err := SaveConfig(config)
				if err != nil {
					fmt.Println(ErrorWritingFile)
					os.Exit(1)
				}
			}
			break
		} else if (index == len(config.Emails)+1) && (size > 0) {
			config.Emails = []EmailConfig{}
			changed = true
		} else if (index > 0) && (index <= len(config.Emails)) {
			config.Emails = append(config.Emails[:index-1], config.Emails[index:]...)
			changed = true
		} else {
			fmt.Println(ErrorInvalidIndex)
		}
	}
}

func HandleEmail(args *Args) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println(ErrorReadingFile)
		os.Exit(1)
	}
	switch {
	case args.EmailAdd:
		HandleEmailAdd(args, config)
	case args.EmailLs:
		HandleEmailLs(args, config)
	case args.EmailRm:
		HandleEmailRm(args, config)
	}
}
