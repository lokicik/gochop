# 🧩 Backend Plan (Golang)

This file outlines the development plan for the GoChop backend, built with Go and the Fiber framework.

---

## 🔧 Stack

- **Language**: Go 1.22+
- **Web Framework**: [Fiber](https://gofiber.io/)
- **Database**: PostgreSQL (for shortlinks, users, analytics)
- **Cache**: Redis (for shortlink → URL mapping, expiring keys)
- **QR Generator**: Go QR library (`github.com/skip2/go-qrcode`)
- **API Spec**: REST

---

### 📁 Backend Folder Structure

```
/cmd
  /gochop-server
    main.go
/internal
  /handlers        # API endpoints
  /services        # Business logic
  /models          # DB models
  /db              # Connection + queries
  /qr              # QR code logic
  /analytics       # Tracking
  /middleware
```

---

## 🚦 Backend Milestones

### ✅ Phase 1: Core Shortening API

- [x] **Endpoint: `POST /api/shorten`**
  - [x] Accept `long_url`, optional `alias` and `context`.
  - [x] Generate a unique short code.
  - [x] Store the mapping in Redis and PostgreSQL.
  - [x] Return the short code, full short URL, and expiration date.
- [x] **Endpoint: `GET /{shortCode}`**
  - [x] Fetch the original URL from Redis (or PostgreSQL as a fallback).
  - [x] Log request metadata (IP, referrer, user agent) for analytics.
  - [x] Redirect to the original URL.

### ✅ Phase 2: QR Code Endpoint

- [x] **Endpoint: `GET /api/qrcode/{shortCode}`**
  - [x] Generate a QR code image (PNG/SVG) for the given short code.
  - [x] Implement caching for the generated QR code to improve performance.

### ✅ Phase 3: Analytics Logging

- [x] Create middleware to intercept requests.
  - [x] Log IP address, timestamp, referrer, and device information.
  - [x] Write analytics data asynchronously to the database to avoid blocking.
- [ ] (Optional) Consider using a message queue like Kafka for high-throughput logging.

### ✅ Phase 4: Admin APIs

- [x] **Endpoint: `GET /api/links`**
  - [x] Fetch all links for a specific user.
  - [x] Include metadata and click count for each link.
- [x] **Endpoint: `GET /api/analytics/{shortCode}`**
  - [x] Provide aggregated analytics data for the frontend charts.

### ✅ Phase 5: Access Control + Expiry

- [x] Implement Time-to-Live (TTL) on Redis keys for automatic link expiration.
- [x] (Optional) Add JWT-based middleware for securing endpoints.
- [x] (Future) Implement IP whitelisting/blacklisting capabilities.
- [x] Add input validation for URLs, aliases, and context
- [x] Fix cryptographic random number generation
- [x] Make base URL configurable via environment variable

### ✅ Phase 6: NextAuth.js Integration & User Management

**Goal**: Integrate with NextAuth.js frontend auth system for secure, scalable authentication.

- [x] **Database Schema (NextAuth Compatible)**:
  - [x] Create NextAuth.js compatible users table schema
  - [x] Add accounts and sessions tables for NextAuth.js
  - [x] Add user_id foreign key to links table
  - [x] Add password field to accounts table for credentials auth
  - [x] Fix database field casing for NextAuth adapter compatibility
  - [x] Add database migration support
- [x] **NextAuth Session Validation**:
  - [x] Create middleware to validate NextAuth.js JWT tokens
  - [x] Extract user ID and admin status from validated tokens
  - [x] Add session-based authentication for all protected endpoints
  - [x] Remove custom JWT middleware and dev token endpoints
- [x] **User Context & Authorization**:
  - [x] Update link creation to require authentication and set user_id
  - [x] Update GetAllLinks to filter by authenticated user
  - [x] Update analytics endpoints to check link ownership
  - [x] Add admin-only endpoints for user management
- [x] **API Integration**:
  - [x] Create user profile endpoint compatible with NextAuth session
  - [x] Add CORS configuration for Next.js API routes
  - [x] Implement proper error handling for unauthenticated requests
  - [x] Create credentials verification for email/password auth
  - [x] Add user registration API endpoint with password hashing

**Technical Implementation Details:**

```go
// Example middleware for NextAuth JWT validation
func NextAuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := extractTokenFromHeader(c)
        userID, isAdmin, err := validateNextAuthToken(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "Invalid session"})
        }
        c.Locals("userID", userID)
        c.Locals("isAdmin", isAdmin)
        return c.Next()
    }
}
```

**Database Schema Updates:**

```sql
-- NextAuth.js compatible users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified TIMESTAMPTZ,
    image TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    is_admin BOOLEAN DEFAULT FALSE
);

-- Add user_id to existing links table
ALTER TABLE links ADD COLUMN user_id UUID REFERENCES users(id);
```

---

## ✨ **Advanced & Unique Features**

### ✅ **Phase 7: Context-Aware Redirects**

- [ ] **Modify `POST /api/shorten`**:
  - [ ] Allow an array of targets instead of a single `long_url`.
  - [ ] Each target should specify conditions (e.g., `device: "mobile"`, `location: "DE"`, `language: "en-US"`).
- [ ] **Enhance `GET /{shortCode}` Logic**:
  - [ ] Evaluate the user's context (User-Agent, IP address, Accept-Language header).
  - [ ] Match the context against the defined rules to determine the correct redirect destination.

### ✅ **Phase 8: A/B Testing**

- [ ] **Modify `POST /api/shorten`**:
  - [ ] Allow multiple target URLs, each with a `weight` for traffic distribution.
- [ ] **Enhance `GET /{shortCode}` Logic**:
  - [ ] Use a weighted-random algorithm to select a destination URL.
- [ ] **Enhance Analytics**:
  - [ ] Track clicks and conversions for each variant separately.
  - [ ] Provide comparison data via the `GET /api/analytics/{shortCode}` endpoint.

### ✅ **Phase 9: Enhanced Security & Link Control**

- [ ] **Password Protection**:
  - [ ] Add a `password` field to the shorten request.
  - [ ] Create a new endpoint `POST /api/verify-password/{shortCode}` to check the password before redirecting.
- [ ] **Self-Destructing Links**:
  - [ ] Add a `max_clicks` field to the shorten request.
  - [ ] Decrement a counter in Redis on each click and disable the link when the count reaches zero.
- [ ] **Editable Destinations (for Dynamic QR Codes)**:
  - [ ] Create a new endpoint `PUT /api/links/{shortCode}` to update the target URL(s) for an existing short link.

---

## 🧪 Testing

- [ ] Write unit tests for handlers.
- [ ] Write tests for business logic in services.
- [ ] Write tests for database queries.

---

## 🚀 Deployment Plan

- **Backend**: Deploy to Fly.io, Railway, or DigitalOcean App Platform.
- **Database**: Use Supabase Postgres or a managed RDS instance.
- **Cache**: Use Upstash Redis or a self-hosted instance.

---

## 🪪 Licensing

- [ ] Add an MIT License.
- [ ] Create a `CONTRIBUTING.md` with guidelines for contributors.
- [ ] Set up an open roadmap using GitHub Projects.
