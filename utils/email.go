package utils

import (
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"

	"userstyles.world/modules/config"
	"userstyles.world/modules/errors"
)

var (
	auth = sasl.NewPlainClient("", config.EMAIL_ADDRESS, config.EMAIL_PWD)
	CLRF = "\r\n"
)

type EmailBuilder struct {
	To       string
	From     string
	Subject  string
	boundary string
	Parts    []MimePart
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

func (eb *EmailBuilder) parseMultiPart() (string, error) {
	output := ""
	boundary := "--" + eb.boundary
	partsLen := len(eb.Parts)

	if partsLen == 0 {
		return "", errors.ErrNoParts
	}

	if partsLen > 1 {
		output += "Content-Type: multipart/alternative; boundary=\"" + eb.boundary + "\"\n\n"
	} else {
		part0 := eb.Parts[0]
		if part0.ContentTransferEncoding == "" {
			part0.ContentTransferEncoding = "8bit"
		}
		if part0.ContentType == "" {
			part0.ContentType = "text/plain"
		}
		output += "Content-Type: " + part0.ContentType + ";charset=\"utf-8\"\n" +
			"Content-Transfer-Encoding: " + part0.ContentTransferEncoding + "\n\n"
	}

	for i := 0; i < partsLen; i++ {
		part := eb.Parts[i]

		if part.Body == "" {
			return "", errors.ErrNoPartBody
		}
		if part.ContentTransferEncoding == "" {
			part.ContentTransferEncoding = "8bit"
		}
		if part.ContentType == "" {
			part.ContentType = "text/plain"
		}
		if partsLen > 1 {
			output += boundary + "\n" +
				"Content-Type: " + part.ContentType + "; charset=\"utf-8\"\n" +
				"Content-Transfer-Encoding: " + part.ContentTransferEncoding + "\n" +
				"\n"
		}
		output += part.Body + "\n\n"
	}

	return output, nil
}

func correctLineBreak(message string) string {
	return string(strings.ReplaceAll(message, "\\n", CLRF))
}

func (eb *EmailBuilder) SendEmail() error {
	eb.boundary = UnsafeString(RandStringBytesMaskImprSrcUnsafe(30))

	if eb.From == "" {
		eb.From = config.EMAIL_ADDRESS
	}

	if eb.To == "" {
		return errors.ErrNoToParameter
	}

	if eb.Subject == "" {
		return errors.ErrNoSubject
	}

	bodyMessage, err := eb.parseMultiPart()
	if err != nil {
		return err
	}

	r := strings.NewReader(correctLineBreak("To: " + eb.To + "\n" +
		"From:" + eb.From + "\n" +
		"Subject:" + eb.Subject + "\n" +
		"MIME-Version: 1.0\n" +
		bodyMessage))

	return smtp.SendMail("mail.userstyles.world:587", auth, config.EMAIL_ADDRESS, []string{eb.To}, r)
}
