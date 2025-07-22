// File: internal/email/templates/welcome.go

package templates

import "fmt"

type WelcomeTemplate struct{}

func NewWelcomeTemplate() *WelcomeTemplate {
	return &WelcomeTemplate{}
}

func (w *WelcomeTemplate) GetTitle() string {
	return "Welcome"
}

func (w *WelcomeTemplate) GetSubtitle() string {
	return "You're now part of our family"
}

func (w *WelcomeTemplate) RenderContent(data interface{}) string {
	d := data.(WelcomeEmailData)

	return fmt.Sprintf(`
		<h2 style="color: #1f2937; margin-bottom: 20px;">🎉 Welcome Aboard, %s!</h2>
		<p style="font-size: 16px; color: #4b5563; margin-bottom: 30px;">
			Your email has been successfully verified and your %s account is now active! 
			We're thrilled to have you join our community.
		</p>

		<div class="feature-grid">
			<div class="feature-card">
				<span class="feature-icon">🚀</span>
				<h3>Get Started</h3>
				<p>Explore our platform and discover all the amazing features we have built for you.</p>
			</div>
			<div class="feature-card">
				<span class="feature-icon">📚</span>
				<h3>Learn & Grow</h3>
				<p>Access our comprehensive documentation, tutorials, and resources to maximize your experience.</p>
			</div>
			<div class="feature-card">
				<span class="feature-icon">💬</span>
				<h3>Get Support</h3>
				<p>Our dedicated support team is here to help you succeed every step of the way.</p>
			</div>
		</div>

		<div style="text-align: center; margin: 35px 0;">
			<a href="%s" class="cta-button">Start Exploring →</a>
		</div>

		<div class="info-box">
			<h3>🌟 What's Next?</h3>
			<ul>
				<li>Complete your profile to personalize your experience</li>
				<li>Explore our features and discover what works best for you</li>
				<li>Join our community and connect with other users</li>
				<li>Check out our getting started guide for tips and tricks</li>
			</ul>
		</div>

		<p style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
			Questions? We're here to help! Reach out to us at 
			<a href="mailto:%s" style="color: %s; text-decoration: none;">%s</a>
		</p>
	`, d.FirstName, d.BrandName, d.Website, d.SupportEmail, d.PrimaryColor, d.SupportEmail)
}

func (w *WelcomeTemplate) RenderText(data interface{}) string {
	d := data.(WelcomeEmailData)

	return fmt.Sprintf(`
		Welcome to %s!

		Hi %s, we're excited to have you on board!

		Your email has been successfully verified and your %s account is now active! 
		We're thrilled to have you join our community.

		What's Next?
		🚀 Get Started - Explore our platform and discover all the amazing features we have built for you
		📚 Learn & Grow - Access our comprehensive documentation, tutorials, and resources  
		💬 Get Support - Our dedicated support team is here to help you succeed

		Next Steps:
		- Complete your profile to personalize your experience
		- Explore our features and discover what works best for you
		- Join our community and connect with other users
		- Check out our getting started guide for tips and tricks

		Visit us at: %s

		Questions? We're here to help! Reach out to us at %s

		Thanks for choosing %s!
		This is an automated email from %s. Please do not reply to this message.
	`, d.BrandName, d.FirstName, d.BrandName, d.Website, d.SupportEmail, d.BrandName, d.BrandName)
}
