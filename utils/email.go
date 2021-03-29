package utils

import (
	"crypto/rand"
	"fmt"
	"log"
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
	To       string
	From     string
	Subject  string
	Parts    []MimePart
	boundary string
}

type MimePart struct {
	ContentType             string
	ContentTransferEncoding string
	Body                    string
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

func (eb *EmailBuilder) SetSubject(subject string) *EmailBuilder {
	eb.Subject = subject
	return eb
}

func NewPart() *MimePart {
	return &MimePart{}
}

func (mp *MimePart) SetContentType(contentType string) *MimePart {
	mp.ContentType = contentType
	return mp
}

func (mp *MimePart) SetContentTransferEncoding(contentTransferEncoding string) *MimePart {
	mp.ContentTransferEncoding = contentTransferEncoding
	return mp
}
func (mp *MimePart) SetBody(body string) *MimePart {
	mp.Body = body
	return mp
}

func (eb *EmailBuilder) AddPart(part MimePart) *EmailBuilder {
	eb.Parts = append(eb.Parts, part)
	return eb
}

func (eb *EmailBuilder) parseMultiPart() string {
	output := ""
	boundary := "--" + eb.boundary
	partsLen := len(eb.Parts)

	if partsLen > 1 {
		output += "Content-Type: multipart/alternative; boundary=\"" + eb.boundary + "\"\n" +
			"MIME-Version: 1.0\n\n"
	} else if partsLen == 0 {
		output += "\n"
	} else {
		panic("Wanted to send Email but no parts were detected.")
	}

	for i := 0; i < partsLen; i++ {
		part := eb.Parts[i]

		if part.Body == "" {
			panic("Wanted to send Email part, but it doesn't contain a body.")
		}
		if part.ContentTransferEncoding == "" {
			part.ContentTransferEncoding = "7bit"
		}
		if part.ContentType == "" {
			part.ContentType = "text/plain; charset=\"utf-8\""
		}

		output += boundary + "\n" +
			"Content-Type:" + part.ContentType + "\n" +
			"Content-Transfer-Encoding:" + part.ContentTransferEncoding + "\n" +
			"\n" +
			part.Body + "\n\n"

	}

	return output
}

func correctLineBreak(message string) string {
	return string(emailLine.ReplaceAll([]byte(message), []byte(CLRF)))
}

func RandomString(size int) string {
	b := make([]byte, size)

	if _, err := rand.Read(b); err != nil {
		log.Fatalln("Failed to generate RandomString, err:", err)
	}

	return fmt.Sprintf("%X", b[0:size])
}

func (eb *EmailBuilder) SendEmail() error {
	eb.boundary = RandomString(30)

	if eb.From == "" {
		eb.From = config.EMAIL_ADDRESS
	}

	if eb.To == "" {
		panic("Wanted to send Email but TO: is missing")
	}

	if eb.Subject == "" {
		panic("Wanted to send Email but SUBJECT: is missing")
	}

	r := strings.NewReader(correctLineBreak("To: " + eb.To + "\n" +
		"From:" + eb.From + "\n" +
		"Subject:" + eb.Subject + "\n" +
		eb.parseMultiPart()))
	return smtp.SendMail("mail.userstyles.world:587", auth, config.EMAIL_ADDRESS, []string{eb.To}, r)
}
