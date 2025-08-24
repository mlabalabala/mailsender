# mailsender

---

```
Email Sender - Send emails via SMTP with SSL/TLS

Usage:
  mailsender -s "subject" -c "content"

Required Flags:
  -s string    Email subject
  -c string    Email content/body
  -h           Show this help message

Required Environment Variables:
  EMAIL_FROM          Sender email address (e.g., your_email@gmail.com)
  EMAIL_PASSWORD      Email password or app-specific password
  EMAIL_TO            Recipient email address
  EMAIL_SMTP_HOST     SMTP server hostname (e.g., smtp.gmail.com)
  EMAIL_SMTP_PORT     SMTP server port (e.g., 465 for SSL, 587 for TLS)

Examples:
  # Set environment variables first:
  export EMAIL_FROM="your_email@gmail.com"
  export EMAIL_PASSWORD="your_app_password"
  export EMAIL_TO="recipient@example.com"
  export EMAIL_SMTP_HOST="smtp.gmail.com"
  export EMAIL_SMTP_PORT="465"

  # Then send email:
  mailsender -s "Test Subject" -c "This is a test email"

Common SMTP Servers:
  Gmail:      smtp.gmail.com:465 (SSL) or :587 (TLS)
  Outlook:    smtp.office365.com:587
  QQ Mail:    smtp.qq.com:465
  163 Mail:   smtp.163.com:465
```