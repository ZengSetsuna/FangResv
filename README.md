# FangResv

**FangResv** is a backend event management platform built with Go, Gin, and PostgreSQL. It enables users to register with email verification, create and join events, and view event details. The project emphasizes clean architecture, security, and ease of deployment using Docker.

## ✨ Features

- User registration with email verification (random code via SMTP)
- Login and authentication
- Create, browse, and join events
- View event details including organizer, participants, time, and venue
- RESTful API design for frontend integration
- Dockerized deployment and environment configuration

## 🛠 Tech Stack

- **Backend**: Go, Gin framework
- **Database**: PostgreSQL + `sqlc` (type-safe query generator)
- **Authentication**: Email verification, bcrypt password hashing
- **Mailer**: Go’s `net/smtp` for sending verification codes
- **Deployment**: Docker, Docker Compose
- **Environment Management**: `.env` files using `github.com/joho/godotenv`

