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

// SendMail takes a userAuth, subject and content to send an email
func SendMail(user *models.UserAuth, link string) error {
	from := mail.NewEmail(hostUser, hostAddress)
	to := mail.NewEmail(user.Username, user.Email)
	message := mail.NewSingleEmail(from, "Email Verification for TnyIE", to, makeTextContent(link), makeHTMLContent(link))
	client := sendgrid.NewSendClient(viper.GetString("SENDGRID_CREDENTIAL"))

	response, err := client.Send(message)
	log.Println(response)
	return err
}

func makeHTMLContent(link string) string {
	link = "https://tny.ie/api/verify/" + link
	return fmt.Sprintf(htmlTemplate, link, link)
}

func makeTextContent(link string) string {
	return fmt.Sprintf(textTemplate, "https://tny.ie/api/verify/"+link)
}
