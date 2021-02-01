package smtp

import (
	"log"
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

func CheckPort25(hostname string) (*SMTPConfig, error) {
	smtp, err := CheckPort(hostname, "25")
	return smtp, err
}

func CheckPort465(hostname string) (*SMTPConfig, error) {
	smtp, err := CheckPort(hostname, "465")
	return smtp, err
}

func CheckPort(hostname, port string) (*SMTPConfig, error) {
	tcp := tcp.NewConfig()
	smtp := NewConfig()
	err := tcp.Connect(hostname, port)

	if err != nil {
		return smtp, err
	}

	smtp.SMTPSummaryCheck(tcp)
	return smtp, nil
}

func (s *SMTPConfig) SMTPSummaryCheck(tcp *tcp.TCPConfig) {
	resp := s.SendRequest(tcp, "AUTH LOGIN")
	if len(resp) == 0 {
		return
	}
	status, err := strconv.Atoi(resp[:3])
	s.Connected = true
	if err == nil {
		s.Status = status
		if status == 500 {
			s.AUTH = true
		}
	}
	if strings.Contains(strings.ToLower(resp), "available only with ssl or tls") {
		s.SSL = true
	}
}

func (s *SMTPConfig) SendRequest(tcp *tcp.TCPConfig, req string) string {
	var resp string
	asd := []byte(req + "\n")
	tcp.ReadTCPMessage()
	if err := tcp.WriteTCPMessage(asd); err != nil {
		log.Printf("Could not write msg: %s", err)
		return ""
	}
	resp = string(tcp.ReadTCPMessage())
	return resp
}

func NewConfig() *SMTPConfig {
	return &SMTPConfig{}
}
