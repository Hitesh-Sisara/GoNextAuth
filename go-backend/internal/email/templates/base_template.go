// File: internal/email/templates/base_template.go

package templates

import (
	"fmt"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/config"
)

type BaseTemplate struct {
	config *config.Config
}

type EmailTemplate interface {
	GetTitle() string
	GetSubtitle() string
	RenderContent(data interface{}) string
	RenderText(data interface{}) string
}

func NewBaseTemplate(cfg *config.Config) *BaseTemplate {
	return &BaseTemplate{config: cfg}
}

// RenderHTML renders the complete HTML email using base template
func (b *BaseTemplate) RenderHTML(template EmailTemplate, data interface{}) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>%s</title>
			%s
		</head>
		<body>
			<div class="email-container">
				%s
				<div class="content">
					%s
				</div>
				%s
			</div>
		</body>
		</html>
	`, template.GetTitle(), b.getBaseStyles(), b.getEmailHeader(template.GetTitle(), template.GetSubtitle()),
		template.RenderContent(data), b.getEmailFooter())
}

// getBaseStyles returns common CSS styles for all emails
func (b *BaseTemplate) getBaseStyles() string {
	return fmt.Sprintf(`
		<style>
			* {
				margin: 0;
				padding: 0;
				box-sizing: border-box;
			}
			body {
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
				line-height: 1.6;
				color: #1f2937;
				background: linear-gradient(135deg, #f8fafc 0%%, #e2e8f0 100%%);
				padding: 20px;
			}
			.email-container {
				max-width: 600px;
				margin: 0 auto;
				background: #ffffff;
				border-radius: 16px;
				overflow: hidden;
				box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
			}
			.header {
				background: linear-gradient(135deg, %s 0%%, %s 100%%);
				padding: 40px 30px;
				text-align: center;
				color: white;
				position: relative;
			}
			.header::before {
				content: '';
				position: absolute;
				top: 0;
				left: 0;
				right: 0;
				bottom: 0;
				background: url('data:image/svg+xml;charset=UTF-8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="%%23ffffff" opacity="0.1"/><circle cx="75" cy="75" r="1" fill="%%23ffffff" opacity="0.1"/><circle cx="50" cy="10" r="0.5" fill="%%23ffffff" opacity="0.15"/><circle cx="20" cy="80" r="0.5" fill="%%23ffffff" opacity="0.15"/></pattern></defs><rect width="100%%" height="100%%" fill="url(%%23grain)"/></svg>');
				pointer-events: none;
			}
			.header h1 {
				font-size: 28px;
				font-weight: 700;
				margin-bottom: 8px;
				position: relative;
				z-index: 1;
			}
			.header p {
				font-size: 16px;
				opacity: 0.9;
				position: relative;
				z-index: 1;
			}
			.logo {
				width: 60px;
				height: 60px;
				margin: 0 auto 20px;
				background: rgba(255, 255, 255, 0.2);
				border-radius: 12px;
				display: flex;
				align-items: center;
				justify-content: center;
				font-size: 24px;
				font-weight: bold;
				position: relative;
				z-index: 1;
			}
			.content {
				padding: 40px 30px;
			}
			.otp-container {
				background: linear-gradient(135deg, #f8fafc 0%%, #f1f5f9 100%%);
				border: 2px dashed %s;
				border-radius: 12px;
				padding: 30px;
				text-align: center;
				margin: 30px 0;
				position: relative;
			}
			.otp-label {
				font-size: 14px;
				color: #64748b;
				margin-bottom: 10px;
				text-transform: uppercase;
				letter-spacing: 1px;
				font-weight: 600;
			}
			.otp-code {
				font-size: 36px;
				font-weight: 800;
				color: %s;
				letter-spacing: 8px;
				font-family: 'Courier New', monospace;
				text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
			}
			.info-box {
				background: linear-gradient(135deg, #fef3cd 0%%, #fde68a 100%%);
				border-left: 4px solid #f59e0b;
				border-radius: 8px;
				padding: 20px;
				margin: 25px 0;
			}
			.info-box.security {
				background: linear-gradient(135deg, #fef2f2 0%%, #fecaca 100%%);
				border-left-color: #ef4444;
			}
			.info-box h3 {
				font-size: 16px;
				margin-bottom: 10px;
				color: #92400e;
			}
			.info-box.security h3 {
				color: #b91c1c;
			}
			.info-box ul {
				list-style: none;
				margin: 0;
				padding: 0;
			}
			.info-box li {
				margin: 8px 0;
				padding-left: 20px;
				position: relative;
			}
			.info-box li:before {
				content: '•';
				color: #f59e0b;
				font-size: 18px;
				position: absolute;
				left: 0;
			}
			.info-box.security li:before {
				color: #ef4444;
			}
			.feature-grid {
				display: grid;
				gap: 20px;
				margin: 30px 0;
			}
			.feature-card {
				background: linear-gradient(135deg, #f8fafc 0%%, #f1f5f9 100%%);
				border: 1px solid #e2e8f0;
				border-radius: 12px;
				padding: 24px;
				transition: all 0.3s ease;
			}
			.feature-card:hover {
				transform: translateY(-2px);
				box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
			}
			.feature-icon {
				font-size: 32px;
				margin-bottom: 12px;
				display: block;
			}
			.feature-card h3 {
				font-size: 18px;
				margin-bottom: 8px;
				color: #1f2937;
			}
			.feature-card p {
				color: #64748b;
				font-size: 14px;
			}
			.cta-button {
				display: inline-block;
				background: linear-gradient(135deg, %s 0%%, %s 100%%);
				color: white !important;
				text-decoration: none;
				padding: 14px 32px;
				border-radius: 8px;
				font-weight: 600;
				font-size: 16px;
				text-align: center;
				margin: 20px 0;
				box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
				transition: all 0.3s ease;
			}
			.cta-button:hover {
				transform: translateY(-1px);
				box-shadow: 0 8px 15px -3px rgba(0, 0, 0, 0.1);
			}
			.footer {
				background: #f8fafc;
				padding: 30px;
				text-align: center;
				border-top: 1px solid #e2e8f0;
			}
			.footer p {
				color: #64748b;
				font-size: 13px;
				margin: 5px 0;
			}
			.social-links {
				margin: 20px 0;
			}
			.social-links a {
				display: inline-block;
				width: 36px;
				height: 36px;
				background: %s;
				color: white;
				text-decoration: none;
				border-radius: 50%%;
				line-height: 36px;
				margin: 0 8px;
				font-size: 16px;
			}
			@media only screen and (max-width: 600px) {
				.email-container {
					margin: 10px;
					border-radius: 12px;
				}
				.header, .content, .footer {
					padding: 24px 20px;
				}
				.otp-code {
					font-size: 28px;
					letter-spacing: 4px;
				}
				.header h1 {
					font-size: 24px;
				}
			}
		</style>
	`, b.config.Branding.PrimaryColor, b.config.Branding.SecondaryColor,
		b.config.Branding.PrimaryColor, b.config.Branding.PrimaryColor,
		b.config.Branding.PrimaryColor, b.config.Branding.SecondaryColor,
		b.config.Branding.PrimaryColor)
}

// getEmailHeader returns the common header HTML
func (b *BaseTemplate) getEmailHeader(title, subtitle string) string {
	logoSection := ""
	if b.config.Branding.LogoURL != "" {
		logoSection = fmt.Sprintf(`<img src="%s" alt="%s" style="width: 60px; height: 60px; margin-bottom: 20px;">`,
			b.config.Branding.LogoURL, b.config.Branding.Name)
	} else {
		logoSection = fmt.Sprintf(`<div class="logo">%s</div>`, string(b.config.Branding.Name[0]))
	}

	return fmt.Sprintf(`
		<div class="header">
			%s
			<h1>%s</h1>
			<p>%s</p>
		</div>
	`, logoSection, title, subtitle)
}

// getEmailFooter returns the common footer HTML
func (b *BaseTemplate) getEmailFooter() string {
	return fmt.Sprintf(`
		<div class="footer">
			<div class="social-links">
				<a href="%s">🌐</a>
				<a href="mailto:%s">✉️</a>
			</div>
			<p><strong>%s</strong></p>
			<p>This is an automated email. Please do not reply to this message.</p>
			<p>If you need assistance, contact us at <a href="mailto:%s" style="color: %s;">%s</a></p>
			<p style="margin-top: 20px; font-size: 12px;">
				© 2025 %s. All rights reserved.
			</p>
		</div>
	`, b.config.Branding.Website, b.config.Branding.SupportEmail, b.config.Branding.Name,
		b.config.Branding.SupportEmail, b.config.Branding.PrimaryColor, b.config.Branding.SupportEmail,
		b.config.Branding.Name)
}
