package mail

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"

	"github.com/tnyie/tny-api/models"
)

const (
	hostUser    = "Tny.ie"
	hostAddress = "noreply@tny.ie"
)

// SendMailVerification takes a userAuth, subject and content to send an email
func SendMailVerification(user *models.UserAuth, link string) error {
	from := mail.NewEmail(hostUser, hostAddress)
	to := mail.NewEmail(user.Username, user.Email)
	message := mail.NewSingleEmail(from, "Email Verification for TnyIE", to, makeMailVerificationTextContent(link), makeMailVerificationHTMLContent(link))
	client := sendgrid.NewSendClient(viper.GetString("SENDGRID_CREDENTIAL"))

	response, err := client.Send(message)
	log.Println(response)
	return err
}

func makeMailVerificationHTMLContent(link string) string {
	link = "https://tny.ie/api/verify/" + link
	return fmt.Sprintf(htmlMailVerificationTemplate, link, link)
}

func makeMailVerificationTextContent(link string) string {
	return fmt.Sprintf(textMailVerificationTemplate, "https://tny.ie/api/verify/"+link)
}

func SendPasswordVerification(user *models.UserAuth, link string) error {
	from := mail.NewEmail(hostUser, hostAddress)
	to := mail.NewEmail(user.Username, user.Email)
	message := mail.NewSingleEmail(from, "Password Reset for TnyIE", to, makePasswordResetTextContent(link), makePasswordResetHTMLContent(link))
	client := sendgrid.NewSendClient(viper.GetString("SENDGRID_CREDENTIAL"))

	response, err := client.Send(message)
	log.Println(response)
	return err
}

func makePasswordResetHTMLContent(link string) string {
	link = "https://" + viper.GetString("tny.ui.url") + "/changepassword/" + link
	return fmt.Sprintf(htmlPasswordVerificationTemplate, link, link)
}

func makePasswordResetTextContent(link string) string {
	link = "https://" + viper.GetString("tny.ui.url") + "/changepassword/" + link
	return fmt.Sprintf(textPasswordVerificationTemplate, link)
}
