package utils

import (
	"crypto/sha1"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"net/mail"
	"net/smtp"
)

func SendEmail(subject, to, body string) (err error) {
	from := viper.GetString("mail_username")
	pass := viper.GetString("mail_password")
	e := email.NewEmail()
	// TODO config from prefix
	e.From = "Robbo <" + from + ">"
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(body)

	auth := smtp.PlainAuth("", from, pass, viper.GetString("smtp_server_host"))
	err = e.Send(viper.GetString("smtp_server_address"), auth)
	return
}

func Hash(s string) (hash string) {
	// TODO use bcrypt
	pwd := sha1.New()
	pwd.Write([]byte(s))
	pwd.Write([]byte(viper.GetString("auth_hash_salt")))
	hash = fmt.Sprintf("%x", pwd.Sum(nil))
	return
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GetOffsetAndLimit(page, pageSize *int) (offset, limit int) {
	if page == nil || pageSize == nil {
		limit = -1
		offset = 0
	} else {
		offset = (*page - 1) * *pageSize
		limit = *pageSize
	}
	return
}
