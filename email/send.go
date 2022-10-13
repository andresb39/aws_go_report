package aws_go_cost

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"gopkg.in/gomail.v2"
	"log"
	"time"
)

func SendMail(region, from, recipient string) {

	// create new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		log.Println("Error occurred while creating aws session", err)
		return
	}

	svc := ses.New(sess)

	// create raw message
	msg := gomail.NewMessage()

	date := time.Now().Format("2006-1")
	body := "Hi Everyone! <br />" + "This is the expense report for the month of " + date + "! <br />" + "Greetings! <br />"

	// Set to emails
	msg.SetAddressHeader("From", from, "Cost Report")
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", "Expense report!")
	msg.SetBody("text/html", body)
	msg.Attach("/tmp/report.csv")

	// create a new buffer to add raw data
	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	// create new raw message
	message := ses.RawMessage{Data: emailRaw.Bytes()}

	input := &ses.SendRawEmailInput{RawMessage: &message}

	// send raw email
	_, err = svc.SendRawEmail(input)
	if err != nil {
		log.Println("Error sending mail - ", err)
		return
	}
	fmt.Println("Email sent successfully to: %s", recipient)

}
