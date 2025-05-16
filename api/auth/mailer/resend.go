package mailer

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
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

func SendConfirmation(conn *pgx.Conn, email string) (*Confirmation, error) {
	var code string
	err := conn.
		QueryRow(context.Background(), "select MakeVerificationCode($1)", email).
		Scan(&code)

	if err != nil {
		return nil, err
	}

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
