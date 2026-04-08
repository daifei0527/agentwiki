package email

import (
	"fmt"
	"log"

	"github.com/agentwiki/agentwiki/pkg/config"
)

// Service 邮件服务接口
type Service interface {
	// SendVerificationEmail 发送邮箱验证邮件
	SendVerificationEmail(email, token, publicKey string) error

	// SendWelcomeEmail 发送欢迎邮件
	SendWelcomeEmail(email, agentName string) error
}

// NewSMTPService 创建SMTP邮件服务
func NewSMTPService(cfg *config.Config) Service {
	return &smtpService{
		cfg: cfg,
	}
}

// NewMockService 创建模拟邮件服务（用于测试）
func NewMockService() Service {
	return &mockService{}
}

// smtpService SMTP邮件服务实现
type smtpService struct {
	cfg *config.Config
}

// SendVerificationEmail 发送邮箱验证邮件
func (s *smtpService) SendVerificationEmail(email, token, publicKey string) error {
	// 这里实现SMTP邮件发送逻辑
	// 由于是模拟实现，这里只打印日志
	log.Printf("[SMTP] 发送验证邮件到 %s, 令牌: %s, 公钥: %s", email, token, publicKey)
	return nil
}

// SendWelcomeEmail 发送欢迎邮件
func (s *smtpService) SendWelcomeEmail(email, agentName string) error {
	// 这里实现SMTP邮件发送逻辑
	// 由于是模拟实现，这里只打印日志
	log.Printf("[SMTP] 发送欢迎邮件到 %s, 智能体名称: %s", email, agentName)
	return nil
}

// mockService 模拟邮件服务实现（用于测试）
type mockService struct{}

// SendVerificationEmail 发送邮箱验证邮件（模拟实现）
func (m *mockService) SendVerificationEmail(email, token, publicKey string) error {
	fmt.Printf("[MOCK] 发送验证邮件到 %s, 令牌: %s, 公钥: %s\n", email, token, publicKey)
	return nil
}

// SendWelcomeEmail 发送欢迎邮件（模拟实现）
func (m *mockService) SendWelcomeEmail(email, agentName string) error {
	fmt.Printf("[MOCK] 发送欢迎邮件到 %s, 智能体名称: %s\n", email, agentName)
	return nil
}
