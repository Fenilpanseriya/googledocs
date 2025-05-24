package helpers

import (
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, resetLink string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "fenilpanseriya2004@gmail.com")
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", "Click the link to reset your password: <a href='"+resetLink+"'>Reset Password</a>")

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "fenilpanseriya2004@gmail.com", os.Getenv("PASSWORD"))

	if err := dialer.DialAndSend(mailer); err != nil {
		log.Fatal("Failed to send email:", err.Error())
	} else {
		log.Println("Reset password email sent successfully!")
	}
	return nil
}
