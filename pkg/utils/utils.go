package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/jordan-wright/email"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
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

func HashPassword(s string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(hashed)
}

func ComparePassword(hashed string, normal string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(normal))
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

func DoesHaveRole(clientRole models.Role, roles []*models.Role) bool {
	for _, role := range roles {
		if role.String() == clientRole.String() {
			return true
		}
	}
	return false
}

func GetHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
