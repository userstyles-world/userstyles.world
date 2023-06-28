package email

import (
	"bufio"
	"bytes"
	"net"
	"net/textproto"
	"strings"
	"testing"
)

var sendMailServer = `220 hello world
502 EH?
250 mx.google.com at your service
250 Sender ok
250 Receiver ok
354 Go ahead
250 Data ok
221 Goodbye
`

var sendMailClient = `EHLO localhost
HELO localhost
MAIL FROM:<test@example.com>
RCPT TO:<other@example.com>
DATA
From: test@example.com
To: other@example.com
Subject: SendMail test
MIME-Version: 1.0
Content-Type: text/plain;charset="utf-8"
Content-Transfer-Encoding: 8bit

SendMail is working for me.
.
QUIT
`

func TestSendMail(t *testing.T) {
	server := strings.Join(strings.Split(sendMailServer, "\n"), "\r\n")
	client := strings.Join(strings.Split(sendMailClient, "\n"), "\r\n")
	var cmdbuf bytes.Buffer
	bcmdbuf := bufio.NewWriter(&cmdbuf)
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Unable to create listener: %v", err)
	}
	defer l.Close()

	// prevent data race on bcmdbuf
	done := make(chan struct{})
	go func(data []string) {
		defer close(done)

		conn, err := l.Accept()
		if err != nil {
			t.Errorf("Accept error: %v", err)
			return
		}
		defer conn.Close()

		tc := textproto.NewConn(conn)
		for i := 0; i < len(data) && data[i] != ""; i++ {
			_ = tc.PrintfLine(data[i])

			for len(data[i]) >= 4 && data[i][3] == '-' {
				i++
				_ = tc.PrintfLine(data[i])
			}
			if data[i] == "221 Goodbye" {
				return
			}
			read := false
			for !read || data[i] == "354 Go ahead" {
				msg, err := tc.ReadLine()
				_, _ = bcmdbuf.Write([]byte(msg + "\r\n"))
				read = true
				if err != nil {
					t.Errorf("Read error: %v", err)
					return
				}
				if data[i] == "354 Go ahead" && msg == "." {
					break
				}
			}
		}
	}(strings.Split(server, "\r\n"))

	plainParty := NewPart().
		SetBody("SendMail is working for me.")

	err = NewEmail().
		SetFrom("test@example.com").
		SetTo("other@example.com>\n\rDATA\r\nInjected message body\r\n.\r\nQUIT\r\n").
		SetSubject("SendMail test").
		AddPart(*plainParty).
		SendEmail(l.Addr().String())
	if err == nil {
		t.Errorf("Expected SendMail to be rejected due to a message injection attempt")
	}
	err = NewEmail().
		SetFrom("test@example.com").
		SetTo("other@example.com").
		SetSubject("SendMail test").
		AddPart(*plainParty).
		SendEmail(l.Addr().String())
	if err != nil {
		t.Errorf("%v", err)
	}

	<-done
	err = bcmdbuf.Flush()
	if err != nil {
		t.Errorf("Flush error: %v", err)
	}
	actualcmds := cmdbuf.String()
	if client != actualcmds {
		t.Errorf("Got:\n%s\nExpected:\n%s", actualcmds, client)
	}
}
