package utils

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"regexp"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"userstyles.world/config"
)

var (
	auth      = sasl.NewPlainClient("", config.EMAIL_ADDRESS, config.EMAIL_PWD)
	emailLine = regexp.MustCompile(`(\\r)?\\n`)
	CLRF      = "\r\n"
)

type EmailBuilder struct {
	To        string
	From      string
	Subject   string
	PlainBody string
	HTMLBody  string
	boundary  string
}

func NewEmail() *EmailBuilder {
	return &EmailBuilder{}
}

func (eb *EmailBuilder) SetTo(to string) *EmailBuilder {
	eb.To = to
	return eb
}

func (eb *EmailBuilder) SetFrom(from string) *EmailBuilder {
	eb.From = from
	return eb
}

func (eb *EmailBuilder) SetDefaultFrom() *EmailBuilder {
	eb.From = config.EMAIL_ADDRESS
	return eb
}

func (eb *EmailBuilder) SetSubject(subject string) *EmailBuilder {
	eb.Subject = subject
	return eb
}

func (eb *EmailBuilder) SetPlainBody(message string) *EmailBuilder {
	eb.PlainBody = message
	return eb
}

func (eb *EmailBuilder) SetHTMLBody(htmlMessage string) *EmailBuilder {
	eb.HTMLBody = htmlMessage
	return eb
}

func randomBoundary() (string, error) {
	buf := make([]byte, 30)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func (eb *EmailBuilder) parseMultiPart() string {
	boundary := "--" + eb.boundary
	output := boundary + "\n" +
		"Content-Type: text/plain; charset=\"utf-8\"; format=\"fixed\"\n" +
		"Content-Transfer-Encoding: quoted-printable\n" +
		"\n" +
		eb.PlainBody +
		"\n\n" +
		boundary + "\n" +
		"Content-Type: text/html; charset=UTF-8\n" +
		"Content-Transfer-Encoding: quoted-printable\n" +
		"\n" +
		"<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n" +
		eb.HTMLBody
	return output
}

func correctLineBreak(message string) string {
	return string(emailLine.ReplaceAll([]byte(message), []byte(CLRF)))
}

func (eb *EmailBuilder) SendEmail() error {
	boundary, err := randomBoundary()
	if err != nil {
		return err
	}
	eb.boundary = boundary
	r := strings.NewReader(correctLineBreak("To: " + eb.To + "\n" +
		"From:" + eb.From + "\n" +
		"Subject:" + eb.Subject + "\n" +
		"Content-Type: multipart/alternative; boundary=\"" + eb.boundary + "\"\n" +
		"MIME-Version: 1.0\n" +
		"\n" +
		eb.parseMultiPart()))
	return smtp.SendMail("mail.userstyles.world:587", auth, config.EMAIL_ADDRESS, []string{eb.To}, r)
}
