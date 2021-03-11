package mail

import (
	"fmt"

	"github.com/sendgrid/rest"
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
func SendMail(user *models.UserAuth, subject, link string) (*rest.Response, error) {
	from := mail.NewEmail(hostUser, hostAddress)
	to := mail.NewEmail(user.Username, user.Email)
	message := mail.NewSingleEmail(from, subject, to, makeTextContent(link), makeHTMLContent(link))
	client := sendgrid.NewSendClient(viper.GetString("SENDGRID_CREDENTIAL"))

	response, err := client.Send(message)
	return response, err
}

func makeTextContent(link string) string {
	return fmt.Sprintf(textTemplate, "https://tny.ie/api/verify/"+link)
}

func makeHTMLContent(link string) string {
	return fmt.Sprintf(htmlTemplate, "https://tny.ie/api/verify/"+link)
}
