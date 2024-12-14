package smtp

import (
	"errors"
	"io"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	r "github.com/end1essrage/listener-smtp-rmq/rmq"
	"github.com/sirupsen/logrus"
)

type Server struct {
	RmqClient *r.Client
}

func NewServer(RmqClient *r.Client) *Server {
	return &Server{RmqClient: RmqClient}
}

// Создается на каждый HELO EHLO запрос, то есть при каждом новом запросе на расслылку
func (s *Server) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{RmqClient: s.RmqClient}, nil
}

// A Session is returned after successful login.
type Session struct {
	Recepents []string
	Body      string
	RmqClient *r.Client
}

// AuthMechanisms returns a slice of available auth mechanisms; only PLAIN is
// supported in this example.
func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

// Auth is the handler for supported authenticators.
func (s *Session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		if username != "username" || password != "password" {
			return errors.New("Invalid username or password")
		}
		return nil
	}), nil
}

// Вызывается на каждый запрос MAIL FROM:
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	logrus.Println("Mail from:", from)
	return nil
}

// Вызывается на каждый запрос RCPT TO:
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	//Удостовериться в формате и что на каждого получателя новый запрос
	s.Recepents = append(s.Recepents, to)
	logrus.Println("Rcpt to:", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if b, err := io.ReadAll(r); err != nil {
		return err
	} else {
		s.Body = string(b)
		logrus.Println("Data:", string(b))
	}
	return nil
}

func (s *Session) Reset() {}

// Вызывается при . ? или может вызываться reset?
func (s *Session) Logout() error {
	//format message
	msg := "message is: " + s.Body + " FOR : " + s.Recepents[0]
	s.RmqClient.SendSting(msg)
	return nil
}

// It can be tested manually with e.g. netcat:
//
//	> netcat -C localhost 25
//	EHLO localhost
//	AUTH PLAIN
//	AHVzZXJuYW1lAHBhc3N3b3Jk
//	MAIL FROM:<root@nsa.gov>
//	RCPT TO:<root@gchq.gov.uk>
//	DATA
//	Hey <3
//	.
