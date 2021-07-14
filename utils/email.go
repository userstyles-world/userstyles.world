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
	clrf = "\r\n"
)

type emailBuilder struct {
	to       string
	from     string
	subject  string
	boundary string
	parts    []mimePart
}

type mimePart struct {
	contentType             string
	contentTransferEncoding string
	body                    string
}

func NewEmail() *emailBuilder {
	return &emailBuilder{}
}

func (eb *emailBuilder) SetTo(to string) *emailBuilder {
	eb.to = to
	return eb
}

func (eb *emailBuilder) SetFrom(from string) *emailBuilder {
	eb.from = from
	return eb
}

func (eb *emailBuilder) SetSubject(subject string) *emailBuilder {
	eb.subject = subject
	return eb
}

func NewPart() *mimePart {
	return &mimePart{}
}

func (mp *mimePart) SetContentType(contentType string) *mimePart {
	mp.contentType = contentType
	return mp
}

func (mp *mimePart) SetContentTransferEncoding(contentTransferEncoding string) *mimePart {
	mp.contentTransferEncoding = contentTransferEncoding
	return mp
}

func (mp *mimePart) SetBody(body string) *mimePart {
	mp.body = body
	return mp
}

func (eb *emailBuilder) AddPart(part mimePart) *emailBuilder {
	eb.parts = append(eb.parts, part)
	return eb
}

func (eb *emailBuilder) parseMultiPart() (string, error) {
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
		output += part.body + "\n\n"
	}

	return output, nil
}

func correctLineBreak(message string) string {
	return string(strings.ReplaceAll(message, "\\n", clrf))
}

func (eb *emailBuilder) SendEmail() error {
	eb.boundary = UnsafeString(RandStringBytesMaskImprSrcUnsafe(30))

	if eb.from == "" {
		eb.from = config.EMAIL_ADDRESS
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

	r := strings.NewReader(correctLineBreak("To: " + eb.to + "\n" +
		"From:" + eb.from + "\n" +
		"Subject:" + eb.subject + "\n" +
		"MIME-Version: 1.0\n" +
		bodyMessage))

	return smtp.SendMail("mail.userstyles.world:587", auth, config.EMAIL_ADDRESS, []string{eb.to}, r)
}
