package util

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

//
// SendInviteEmail sends an invite email when a user wishes to invite another user to their account
//
func SendInviteEmail(address, name, inviteCode string) error {

	// Create a new session in the us-east-1 region.
	// Replace us-east-1 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return err
	}
	// Create an SES session.
	svc := ses.New(sess)

	getTemplate := ses.GetTemplateInput{}
	getTemplate.TemplateName = aws.String(os.Getenv("INVITATION_TEMPLATE_NAME"))

	destinations := ses.Destination{}
	sendTemplateInput := ses.SendTemplatedEmailInput{}
	toAddresses := []*string{&address}

	var url = os.Getenv("VIJNANA_CNAME_URL")
	url += "/invitation?invitation=" + inviteCode
	destinations.ToAddresses = toAddresses
	sendTemplateInput.Destination = &destinations
	sendTemplateInput.Source = aws.String(os.Getenv("TESPO_EMAIL"))
	sendTemplateInput.Template = aws.String(os.Getenv("INVITATION_TEMPLATE_NAME"))
	sendTemplateInput.TemplateData = aws.String("{\"name\":\"" + name + "\",\"invitationLandingPage\":\"" + url + "\"}")
	sendTemplateInput.SourceArn = aws.String(os.Getenv("SOURCE_ARN"))

	_, err = svc.SendTemplatedEmail(&sendTemplateInput)
	if err != nil {
		return err
	}
	return nil
}
