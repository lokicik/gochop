# 🚀 GoChop - Intelligent URL Shortener

<!-- ![GoChop Banner](https://your-image-host.com/gochop-banner.png) -->

<!-- TODO: Replace with a real banner/screenshot of the app -->

**GoChop** is a high-performance, open-source URL shortener built for power users. It goes beyond simple redirects with advanced features like context-aware links, A/B testing, and dynamic QR codes.

<!-- **[➡️ Live Demo](https://gochop.your-domain.com)** -->
<!-- TODO: Replace with your deployed app link -->

---

## ✨ Features

GoChop is more than just a URL shortener. It's a powerful link management tool.

- **⚡️ High-Performance Backend:** Built with **Go (Fiber)** for maximum speed and efficiency.
- **Modern Frontend:** A sleek, responsive UI built with **Next.js** and **Tailwind CSS**.
- **Context-Aware Redirects:** Redirect users to different URLs based on their device, location, or language.
- **A/B Testing:** Distribute traffic between multiple destination URLs from a single short link to test performance.
- **Dynamic QR Codes:** The destination of a QR code can be updated at any time, even after it's been printed.
- **Custom Aliases:** Create memorable, branded short links.
- **Link Expiration:** Links automatically expire after a configurable duration.
- **Enterprise-Ready Stack:** Uses **PostgreSQL** for data persistence and **Redis** for high-speed caching, all managed with **Docker**.

## 🛠 Tech Stack

This project is a monorepo containing both the backend and frontend.

| Component    | Technology                                        |
| ------------ | ------------------------------------------------- |
| **Backend**  | Go, Fiber, PostgreSQL, Redis                      |
| **Frontend** | Next.js, React, TypeScript, Tailwind CSS          |
| **DevOps**   | Docker, Docker Compose, GitHub Actions (optional) |
| **Database** | PostgreSQL (Persistence)                          |
| **Cache**    | Redis (High-speed lookups & caching)              |

---

## 📂 Project Structure

The project is structured as a monorepo:

```
/
├── backend/          # Go Fiber REST API
│   ├── cmd/
│   │   └── gochop-server/ # Main application entrypoint
│   └── internal/         # All internal business logic
├── frontend/         # Next.js & Tailwind CSS UI
│   ├── src/
│   │   └── app/          # App Router pages and components
│   └── public/
├── .github/          # (Optional) GitHub Actions workflows
├── .gitignore
├── docker-compose.yml  # Defines the Postgres & Redis services
├── gochop-backend.md   # Backend development plan
├── gochop-frontend.md  # Frontend development plan
└── README.md           # You are here!
```

---

## 🏁 Getting Started

Follow these instructions to get the project running on your local machine for development and testing.

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.22+ recommended)
- [Node.js](https://nodejs.org/en/download/) (version 18+ recommended)
- [Docker](https://www.docker.com/products/docker-desktop) and Docker Compose

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/gochop.git
cd gochop
```

### 2. Configure Environment Variables

The backend requires environment variables for database connections.

- Navigate to the `backend` directory.
- Create a `.env` file by copying the example:
  ```bash
  cd backend
  cp .env.example .env
  ```
- The default values in `.env` are pre-configured to work with the `docker-compose.yml` file. You shouldn't need to change them for local development.

### 3. Start the Services

The database (Postgres) and cache (Redis) are managed by Docker Compose.

- From the **project root directory**, run:
  ```bash
  docker-compose up -d
  ```
  This will start the required containers in the background.

### 4. Run the Backend Server

- In a new terminal, navigate to the `backend` directory and start the Go server:
  ```bash
  cd backend
  go run cmd/gochop-server/main.go
  ```
  The backend API will be running at `http://localhost:3001`.

### 5. Run the Frontend App

- In another new terminal, navigate to the `frontend` directory, install dependencies, and start the development server:
  ```bash
  cd frontend
  npm install
  npm run dev
  ```
  The frontend application will be available at `http://localhost:3000`.

You should now have the full application stack running locally!

---

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

<!-- You will need to create a LICENSE file for this link to work -->
