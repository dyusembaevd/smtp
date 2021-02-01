package smtp

import (
	"strconv"
	"strings"

	"github.com/mohito22/tcp"
)

// AUTH (true when AUTH method is not supported on currnet port) SSL (true when SSL or TLS is required on current port)

type SMTPConfig struct {
	TCPConfig tcp.TCPConfig
	Status    int
	Connected bool
	AUTH      bool
	SSL       bool
}

func CheckPort25(hostname string) (SMTPConfig, error) {
	smtp, err := CheckPort(hostname, "25")
	return smtp, err
}

func CheckPort465(hostname string) (SMTPConfig, error) {
	smtp, err := CheckPort(hostname, "465")
	return smtp, err
}

func CheckPort(hostname, port string) (SMTPConfig, error) {
	tcp := tcp.NewConfig()
	smtp := NewConfig()
	err := tcp.Connect(hostname, port)

	if err != nil {
		return smtp, err
	}

	smtp.SetConfigInfo(tcp)
	smtp.SMTPSummaryCheck()
	smtp.TCPConfig.CloseConnection()
	return smtp, nil
}

func (s *SMTPConfig) SMTPSummaryCheck() {
	s.Connected = true
	resp := s.SendRequest("AUTH LOGIN")
	status, err := strconv.Atoi(resp[:3])
	if err == nil {
		s.Status = status
		if status == 500 {
			s.AUTH = true
		}
	}
	if strings.Contains(strings.ToLower(resp), "available only with tls or ssl") {
		s.SSL = true
	}
}

func (s *SMTPConfig) SendRequest(req string) string {
	var resp string
	tcp := s.TCPConfig
	_ = tcp.ReadTCPMessage()
	tcp.WriteTCPMessage([]byte(req))
	resp = string(tcp.ReadTCPMessage())
	return resp
}

func (s *SMTPConfig) SetConfigInfo(tcp *tcp.TCPConfig) {
	s.TCPConfig = *tcp
}

func NewConfig() SMTPConfig {
	return SMTPConfig{}
}
