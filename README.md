# GoNextAuth: Production-Ready Go + Next.js Authentication Starter

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.18%2B-blue.svg)](https://golang.org)
[![Next.js Version](https://img.shields.io/badge/Next.js-14%2B-black.svg)](https://nextjs.org)
[![GitHub stars](https://img.shields.io/github/stars/Hitesh-Sisara/GoNextAuth?style=social)](https://github.com/Hitesh-Sisara/GoNextAuth/stargazers)

**GoNextAuth** is a feature-rich, secure, and production-ready starter kit for building modern web applications. It combines a powerful **Go (Golang)** backend with a sleek **Next.js 14** frontend, providing a complete authentication system out of the box.

This project is designed to save you weeks of development time by providing a solid foundation for your next project, with a strong focus on security, performance, and developer experience.

## ✨ Features

This starter kit is packed with features to get you up and running immediately.

### 🔐 Backend (Go)

- **Complete Authentication Flows:**
  - Multi-step email & password registration with OTP verification.
  - Secure password-based login.
  - Passwordless login with email OTP.
  - Robust "Forgot Password" flow with OTP verification.
- **Social Login:** Seamless Google OAuth 2.0 integration with account merging.
- **JWT Authentication:** Secure, stateless authentication using HS256 signed Access and Refresh tokens.
- **Robust Security:**
  - **Rate Limiting:** Protects against brute-force attacks with strict limits on auth endpoints.
  - **CORS & Security Headers:** Pre-configured middleware for CORS, CSP, XSS protection, and more.
  - **Graceful Shutdown:** Ensures no requests are dropped during deployments.
  - **Request Deduplication:** Prevents duplicate requests for sensitive operations like logout and OAuth callbacks.
  - **Secure Password Hashing:** Uses `bcrypt` for hashing user passwords.
- **Database:**
  - PostgreSQL integration using `pgxpool` for high performance.
  - Automatic database migrations on startup.
- **Email Service:**
  - Beautiful, responsive, and customizable email templates.
  - Sends Welcome, Verification, Password Reset, and Login OTP emails via AWS SES.
- **User Management & Auditing:**
  - User profile management endpoints.
  - Detailed user activity logging (logins, signups, password resets, etc.).
- **Developer Experience:**
  - **API Documentation:** Auto-generated Swagger (OpenAPI) documentation.
  - **Configuration:** Centralized, environment-variable-driven configuration.
  - **Background Jobs:** Goroutines for cleaning up expired OTPs and old activity logs.

### 🚀 Frontend (Next.js)

- **Modern Tech Stack:** Next.js 14 (App Router), React 18, TypeScript, and Tailwind CSS.
- **Beautiful, Responsive UI:**
  - Professionally designed components from **Shadcn/UI**.
  - Responsive and accessible design for all devices.
  - Light and Dark mode support.
- **State Management:** Global, persistent authentication state managed with **Zustand**.
- **Advanced Authentication Handling:**
  - **Protected Routes:** Middleware-based protection for dashboard and other private pages.
  - **Guest-Only Routes:** Automatically redirects authenticated users from login/signup pages.
  - **Automatic Token Refresh:** Axios interceptors handle token expiration and refresh seamlessly.
  - **Secure Token Storage:** Manages JWTs in cookies and localStorage.
- **Seamless User Experience:**
  - **Multi-Step Forms:** Intuitive, step-by-step flows for Signup and Forgot Password.
  - **Google OAuth:** One-click Google Sign-In with server-side validation.
  - **User-Friendly Feedback:** Toast notifications for all actions using `react-sonner`.
  - **Comprehensive Form Validation:** Real-time validation on all input fields.

## 🛠️ Tech Stack

| Category     | Technology                                                                                                                                                                                                                                                                |
| ------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Backend**  | [Go](https://golang.org/), [Gin](https://gin-gonic.com/), [PostgreSQL](https://www.postgresql.org/), [pgx](https://github.com/jackc/pgx), [JWT](https://github.com/golang-jwt/jwt), [AWS SES](https://aws.amazon.com/ses/), [Swagger](https://swagger.io/)                |
| **Frontend** | [Next.js](https://nextjs.org/), [React](https://react.dev/), [TypeScript](https://www.typescriptlang.org/), [Tailwind CSS](https://tailwindcss.com/), [Zustand](https://github.com/pmndrs/zustand), [Shadcn/UI](https://ui.shadcn.com/), [Axios](https://axios-http.com/) |

## 🚀 Getting Started

Follow these instructions to get the project up and running on your local machine.

### Prerequisites

- **Go**: Version 1.18 or higher.
- **Node.js**: Version 18.x or higher (with `npm` or `bun`).
- **PostgreSQL**: A running instance. Using Docker is recommended for ease of setup.
- **AWS SES Credentials**: For sending transactional emails.
- **Google OAuth Credentials**: For Google Sign-In.

### 1. Clone the Repository

```bash
git clone https://github.com/Hitesh-Sisara/GoNextAuth.git
cd GoNextAuth
```

### 2. Backend Setup (Go)

1.  **Navigate to the backend directory:**

    ```bash
    cd go-backend
    ```

2.  **Configure your environment.** Rename `.env.example` to `.env` and fill it out with your credentials. The file is commented to guide you on what's required (database connection, JWT secret, email provider, and Google credentials).

3.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

4.  **Run the server:**
    ```bash
    go run cmd/server/main.go
    ```

The backend server will start on `http://localhost:8080`. API documentation is available at `http://localhost:8080/docs/index.html`.

### 3. Frontend Setup (Next.js)

1.  **Navigate to the frontend directory:**

    ```bash
    cd ../nextjs-frontend
    ```

2.  **Configure your environment.** Rename `.env.local.example` to `.env.local` and add your Google Client ID and API URL.

3.  **Install dependencies:**

    ```bash
    bun install
    # or npm install
    ```

4.  **Run the development server:**
    ```bash
    bun dev
    # or npm run dev
    ```

The frontend application will be available at `http://localhost:3000`.

## 🔒 Security

Security is a top priority for this project. Here are some of the key security features implemented:

- **HTTPS in Production:** Assumed and configured via cookies (`secure: true`).
- **Password Policy:** Enforces strong passwords on the frontend and backend.
- **HTTP Security Headers:** Middleware adds `CSP`, `X-Frame-Options`, `X-XSS-Protection`, etc., to mitigate common web vulnerabilities.
- **Rate Limiting:** Protects against brute-force and denial-of-service attacks.
- **Secure Token Storage:** Securely stores JWTs using cookies and localStorage.
- **Google OAuth `state` Parameter:** Used to prevent Cross-Site Request Forgery (CSRF) attacks during the OAuth flow.
- **Input Sanitization & Validation:** All user input is validated and sanitized.

## 🤝 Contributing

We welcome contributions of all kinds! Whether you're fixing a bug, improving documentation, or adding a new feature, your help is appreciated.

### How to Contribute

1.  **Fork the repository** to your own GitHub account.
2.  **Create a new branch** for your changes (`git checkout -b feature/your-feature-name`).
3.  **Make your changes** and commit them with clear, descriptive messages.
4.  **Push your changes** to your fork (`git push origin feature/your-feature-name`).
5.  **Open a Pull Request** against the `main` branch of this repository.

### Reporting Security Issues

If you discover a security vulnerability, please do **not** open a public issue. Instead, send a private email to our security contact. We will address it as a top priority. **(hitesh.sisara.dev@gmail.com)**.

## 📄 License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for more details.
