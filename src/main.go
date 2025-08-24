package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

// Config email configuration struct
type Config struct {
	FromName string
	From     string
	Password string
	To       string
	Subject  string
	Body     string
	SMTPHost string
	SMTPPort int
}

// showHelp displays usage information
func showHelp() {
	fmt.Println("Email Sender - Send emails via SMTP with SSL/TLS")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  mailsender -s \"subject\" -c \"content\"")
	fmt.Println("")
	fmt.Println("Required Flags:")
	fmt.Println("  -s string    Email subject")
	fmt.Println("  -c string    Email content/body")
	fmt.Println("  -h           Show this help message")
	fmt.Println("")
	fmt.Println("Required Environment Variables:")
	fmt.Println("  EMAIL_FROM_NAME     Sender email display name")
	fmt.Println("  EMAIL_FROM          Sender email address (e.g., your_email@gmail.com)")
	fmt.Println("  EMAIL_PASSWORD      Email password or app-specific password")
	fmt.Println("  EMAIL_TO            Recipient email address")
	fmt.Println("  EMAIL_SMTP_HOST     SMTP server hostname (e.g., smtp.gmail.com)")
	fmt.Println("  EMAIL_SMTP_PORT     SMTP server port (e.g., 465 for SSL, 587 for TLS)")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  # Set environment variables first:")
	fmt.Println("  export EMAIL_FROM_NAME=\"display name\"")
	fmt.Println("  export EMAIL_FROM=\"your_email@gmail.com\"")
	fmt.Println("  export EMAIL_PASSWORD=\"your_app_password\"")
	fmt.Println("  export EMAIL_TO=\"recipient@example.com\"")
	fmt.Println("  export EMAIL_SMTP_HOST=\"smtp.gmail.com\"")
	fmt.Println("  export EMAIL_SMTP_PORT=\"465\"")
	fmt.Println("")
	fmt.Println("  # Then send email:")
	fmt.Println("  mailsender -s \"Test Subject\" -c \"This is a test email\"")
	fmt.Println("")
	fmt.Println("Common SMTP Servers:")
	fmt.Println("  Gmail:      smtp.gmail.com:465 (SSL) or :587 (TLS)")
	fmt.Println("  Outlook:    smtp.office365.com:587")
	fmt.Println("  QQ Mail:    smtp.qq.com:465")
	fmt.Println("  163 Mail:   smtp.163.com:465")
}

func parseFlags() *Config {
	// Command line arguments
	subject := flag.String("s", "", "Email subject (required)")
	body := flag.String("c", "", "Email content (required)")
	help := flag.Bool("h", false, "Show help message")

	flag.Parse()

	// Show help if -h flag is provided
	if *help {
		showHelp()
		os.Exit(0)
	}

	// Read configuration from environment variables
	from := getEnvOrFail("EMAIL_FROM", "Sender email address")
	fromName := os.Getenv("EMAIL_FROM_NAME")
	if "" == fromName {
		fromName = from
	}
	password := getEnvOrFail("EMAIL_PASSWORD", "Email password")
	to := getEnvOrFail("EMAIL_TO", "Recipient email address")
	smtpHost := getEnvOrFail("EMAIL_SMTP_HOST", "Mail server host")
	smtpPortStr := getEnvOrFail("EMAIL_SMTP_PORT", "Mail server port (SSL 465, TLS 587)")

	// Convert port to integer
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatalf("SMTP port format error: %v", err)
	}

	// Validate required parameters
	if *subject == "" {
		fmt.Println("Error: Email subject parameter -s is required")
		fmt.Println("Use -h flag for help")
		os.Exit(1)
	}
	if *body == "" {
		fmt.Println("Error: Email content parameter -c is required")
		fmt.Println("Use -h flag for help")
		os.Exit(1)
	}

	return &Config{
		From:     from,
		FromName: fromName,
		Password: password,
		To:       to,
		Subject:  *subject,
		Body:     *body,
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
	}
}

// getEnvOrFail get value from environment variable, fail if not set
func getEnvOrFail(envVar, description string) string {
	value := os.Getenv(envVar)
	if value == "" {
		log.Fatalf("Error: Environment variable %s must be set (%s)", envVar, description)
	}
	return value
}

// validateConfig validate configuration
func validateConfig(config *Config) error {
	if config.From == "" {
		return fmt.Errorf("sender email cannot be empty")
	}
	if config.Password == "" {
		return fmt.Errorf("email password cannot be empty")
	}
	if config.To == "" {
		return fmt.Errorf("recipient email cannot be empty")
	}
	if config.Subject == "" {
		return fmt.Errorf("email subject cannot be empty")
	}
	if config.Body == "" {
		return fmt.Errorf("email content cannot be empty")
	}
	return nil
}

func sendEmail(config *Config) {
	// Authentication information
	auth := smtp.PlainAuth("", config.From, config.Password, config.SMTPHost)

	// TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         config.SMTPHost,
	}

	// Connect to server
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", config.SMTPHost, config.SMTPPort), tlsConfig)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, config.SMTPHost)
	if err != nil {
		log.Fatalf("Failed to create SMTP client: %v", err)
	}
	defer client.Close()

	// Authenticate
	if err = client.Auth(auth); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Set sender
	if err = client.Mail(config.From); err != nil {
		log.Fatalf("Failed to set sender: %v", err)
	}

	// Set recipient
	if err = client.Rcpt(config.To); err != nil {
		log.Fatalf("Failed to set recipient: %v", err)
	}

	// Send email content
	w, err := client.Data()
	if err != nil {
		log.Fatalf("Failed to create data writer: %v", err)
	}

	// 构建完整的邮件头，包含显示名称
	headers := []string{
		fmt.Sprintf("From: %s<%s>", config.FromName, config.From),
		fmt.Sprintf("To: %s", config.To),
		fmt.Sprintf("Subject: %s", config.Subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
	}

	messageStr := strings.Join(headers, "\r\n") + config.Body + "\r\n"
	message := []byte(messageStr)

	if _, err = w.Write(message); err != nil {
		log.Fatalf("Failed to write message: %v", err)
	}
	if err = w.Close(); err != nil {
		log.Fatalf("Failed to close writer: %v", err)
	}

	// Quit gracefully
	if err = client.Quit(); err != nil {
		log.Fatalf("Failed to quit gracefully: %v", err)
	}
}

func main() {
	// Parse command line arguments
	config := parseFlags()

	// Validate configuration
	if err := validateConfig(config); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Send email
	sendEmail(config)
}
