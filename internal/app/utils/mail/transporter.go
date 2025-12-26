package mailing

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	brevo "github.com/sendinblue/APIv3-go-library/v2/lib"
)

var (
	_       = godotenv.Load()
	API_KEY = os.Getenv("TOKEN_BREVO")
)

func BrevoMail(
	targetMail string,
	templateId int,
	params map[string]interface{},
) (bool, error) {
	fmt.Printf("API KEY %s", API_KEY)
	var ctx context.Context
	brevoConfiguration := brevo.NewConfiguration()
	brevoConfiguration.AddDefaultHeader("api-key", API_KEY)
	brevoConfiguration.AddDefaultHeader("partner-key", API_KEY)

	br := brevo.NewAPIClient(brevoConfiguration)

	to := brevo.SendSmtpEmailTo{
		Email: targetMail,
	}

	smptEmail := brevo.SendSmtpEmail{
		To:         []brevo.SendSmtpEmailTo{to},
		TemplateId: int64(templateId),
		Params:     params,
	}

	_, _, err := br.TransactionalEmailsApi.SendTransacEmail(ctx, smptEmail)
	if err != nil {
		return false, err
	}

	return true, nil
}
