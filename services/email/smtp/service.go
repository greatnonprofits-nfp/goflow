package smtp

import (
	"strings"

	"github.com/greatnonprofits-nfp/goflow/flows"
	"github.com/greatnonprofits-nfp/goflow/utils/smtpx"
)

type service struct {
	smtpClient *smtpx.Client
}

// NewService creates a new SMTP email service
func NewService(smtpURL string) (flows.EmailService, error) {
	c, err := smtpx.NewClientFromURL(smtpURL)
	if err != nil {
		return nil, err
	}

	return &service{smtpClient: c}, nil
}

func (s *service) Send(session flows.Session, addresses []string, subject, body string, attachments []string) error {
	// sending blank emails is a good way to get flagged as a spammer so use placeholder if body is empty
	if strings.TrimSpace(body) == "" {
		body = "(empty body)"
	}

	m := smtpx.NewMessage(addresses, subject, body, "", attachments)
	return smtpx.Send(s.smtpClient, m)
}
