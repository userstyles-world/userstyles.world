package utils

import (
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"

	"userstyles.world/modules/config"
	"userstyles.world/modules/errors"
)

var (
	auth = sasl.NewPlainClient("", config.EmailAddress, config.EmailPassword)
	clrf = "\r\n"
)

type EmailBuilder struct {
	to       string
	from     string
	subject  string
	boundary string
	parts    []MimePart
}

type MimePart struct {
	contentType             string
	contentTransferEncoding string
	body                    string
}

func NewEmail() *EmailBuilder {
	return &EmailBuilder{}
}

func (eb *EmailBuilder) SetTo(to string) *EmailBuilder {
	eb.to = to
	return eb
}

func (eb *EmailBuilder) SetFrom(from string) *EmailBuilder {
	eb.from = from
	return eb
}

func (eb *EmailBuilder) SetSubject(subject string) *EmailBuilder {
	eb.subject = subject
	return eb
}

func NewPart() *MimePart {
	return &MimePart{}
}

func (mp *MimePart) HTML() *MimePart {
	mp.contentType = "text/html"
	return mp
}

func (mp *MimePart) SetContentType(contentType string) *MimePart {
	mp.contentType = contentType
	return mp
}

func (mp *MimePart) SetContentTransferEncoding(contentTransferEncoding string) *MimePart {
	mp.contentTransferEncoding = contentTransferEncoding
	return mp
}

func (mp *MimePart) SetBody(body string) *MimePart {
	mp.body = body
	return mp
}

func (eb *EmailBuilder) AddPart(part MimePart) *EmailBuilder {
	eb.parts = append(eb.parts, part)
	return eb
}

func (eb *EmailBuilder) parseMultiPart() (string, error) {
	output := ""
	boundary := "--" + eb.boundary
	partsLen := len(eb.parts)

	if partsLen == 0 {
		return "", errors.ErrNoParts
	}

	if partsLen > 1 {
		output += "Content-Type: multipart/alternative; boundary=\"" + eb.boundary + "\"\n\n"
	} else {
		part0 := eb.parts[0]
		if part0.contentTransferEncoding == "" {
			part0.contentTransferEncoding = "8bit"
		}
		if part0.contentType == "" {
			part0.contentType = "text/plain"
		}
		output += "Content-Type: " + part0.contentType + ";charset=\"utf-8\"\n" +
			"Content-Transfer-Encoding: " + part0.contentTransferEncoding + "\n\n"
	}

	for i := 0; i < partsLen; i++ {
		part := eb.parts[i]

		if part.body == "" {
			return "", errors.ErrNoPartBody
		}
		if part.contentTransferEncoding == "" {
			part.contentTransferEncoding = "8bit"
		}
		if part.contentType == "" {
			part.contentType = "text/plain"
		}
		if partsLen > 1 {
			output += boundary + "\n" +
				"Content-Type: " + part.contentType + "; charset=\"utf-8\"\n" +
				"Content-Transfer-Encoding: " + part.contentTransferEncoding + "\n" +
				"\n"
		}
		output += part.body
		if partsLen-i > 1 {
			output += "\n"
		}
	}

	return output, nil
}

func correctLineBreak(message string) string {
	return strings.ReplaceAll(message, "\\n", clrf)
}

func (eb *EmailBuilder) SendEmail(imapServer string) error {
	eb.boundary = RandomString(30)

	if eb.from == "" {
		eb.from = config.EmailAddress
	}

	if eb.to == "" {
		return errors.ErrNoToParameter
	}

	if eb.subject == "" {
		return errors.ErrNoSubject
	}

	bodyMessage, err := eb.parseMultiPart()
	if err != nil {
		return err
	}

	r := strings.NewReader(correctLineBreak("From: " + eb.from + "\n" +
		"To: " + eb.to + "\n" +
		"Subject: " + eb.subject + "\n" +
		"MIME-Version: 1.0\n" +
		bodyMessage))

	return smtp.SendMail(imapServer, auth, eb.from, []string{eb.to}, r)
}
