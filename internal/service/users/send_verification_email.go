package users

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/pkg/errors"
	"gopkg.in/mail.v2"
)

const verificationEmailBody = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Email Verification</title>
	</head>
	<body style="margin:0; padding:0; background-color:#ffffff;">
		<table width="100%" cellpadding="0" cellspacing="0">
			<tr>
				<td align="center" style="padding:20px;">
					<table width="600" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4;">
						<tr>
							<td style="padding:20px; font-family:Arial, sans-serif; color:#333; line-height:1.6;">
								<h2 style="color:#2c3e50; text-align:center; margin-top:0;">
									Email Verification
								</h2>

								<p>Hello {{.FirstName}},</p>

								<p>Please take a moment to verify your email address by clicking the button below:</p>

								<table align="center" cellpadding="0" cellspacing="0" style="margin:30px auto;">
									<tr>
										<td bgcolor="#3498db" style="padding:12px 30px;">
											<a href="{{.Link}}"
											style="color:#ffffff; text-decoration:none; font-weight:bold; font-family:Arial, sans-serif;">
												Verify Email Address
											</a>
										</td>
									</tr>
								</table>

								<p style="font-size:14px; color:#666;">
									This link will expire in 1 hour.
								</p>

								<p>
									If you're unable to click the button above, copy and paste the link below:
								</p>

								<p style="word-break:break-all; background-color:#ecf0f1; padding:10px; font-size:12px;">
									{{.Link}}
								</p>

								<hr style="border:none; border-top:1px solid #ddd; margin:30px 0;">

								<p style="font-size:14px; color:#666; text-align:center;">
									Thank you!<br>
									FightPicker Team
								</p>
							</td>
						</tr>
					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>
`

type VerificationEmailSender interface {
	SendVerificationEmail(ctx context.Context, userID id.UUID7, hashedToken []byte, msg *mail.Message) error
}

func (s *Service) SendVerificationEmail(ctx context.Context, userID id.UUID7, firstName, email string) error {
	if userID == id.UUID7Nil {
		s.logger.ErrorContext(ctx, "invalid user id: nil UUID")
		return errors.Wrap(ErrInvalidParameter, "user_id")
	}

	if firstName == "" {
		s.logger.ErrorContext(ctx, "invalid first name: empty string")
		// no need to fail the req due to missing first name, use a reasonable default
		firstName = "User"
	}

	if !s.regexEmail.MatchString(email) {
		s.logger.ErrorContext(ctx, "invalid email format", "email", email)
		return errors.Wrap(ErrInvalidParameter, "email")
	}

	// generate a random verification token
	token, err := s.authTool.GenerateVerificationToken()
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to generate verification token", "error", err)
		return ErrInternalError
	}

	hashedToken, err := s.authTool.HashVerificationToken(token)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to hash verification token", "error", err)
		return ErrInternalError
	}

	verificationLink := fmt.Sprintf("%s/api/v1/users/verify?token=%s", s.domain, token)

	tmpl, err := template.New("email").Parse(verificationEmailBody)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, struct {
		FirstName string
		Link      string
	}{
		FirstName: firstName,
		Link:      verificationLink,
	})

	body := buf.String()

	// create the email message to be sent
	msg := mail.NewMessage()
	msg.SetHeader("From", s.emailAddressNoReply)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Please verify your email address")
	msg.SetDateHeader("Date", s.dateTimeTool.Now().UTC())
	msg.SetBody("text/html", body)
	msg.AddAlternative("text/plain", fmt.Sprintf("Hello %s,\n\nPlease verify your email address by clicking the link below:\n\n%s\n\nThis link will expire in 1 hour.\n\nThank you!\nFightPicker Team", firstName, verificationLink))

	return s.repo.SendVerificationEmail(ctx, userID, hashedToken, msg)
}
