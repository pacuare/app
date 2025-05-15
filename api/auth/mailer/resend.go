package mailer

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/resend/resend-go/v2"
)

var (
	apiKey string
	client *resend.Client
)

type Confirmation struct {
	Code     string
	Response *resend.SendEmailResponse
}

func init() {
	apiKey = os.Getenv("RESEND_API_KEY")
	client = resend.NewClient(apiKey)
}

func SendConfirmation(email string) (*Confirmation, error) {
	code := strconv.FormatInt(rand.Int63n(199999), 10)

	sent, err := client.Emails.Send(&resend.SendEmailRequest{
		From:    "Pacuare Reserve <support@farthergate.com>",
		To:      []string{email},
		Text:    fmt.Sprintf("Please use the code %s to log in to Pacuare.", code),
		Subject: "Login Confirmation",
	})

	if err != nil {
		return nil, err
	}

	return &Confirmation{Code: code, Response: sent}, nil
}
