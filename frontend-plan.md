# 🧑‍💻 Frontend Development Plan: GoChop (Next.js + Cursor IDE)

This file outlines the development plan for the GoChop frontend, built with Next.js and managed in Cursor IDE.

---

## 🧰 Tools & Stack

- **Framework**: [Next.js 14+ (App Router)](https://nextjs.org)
- **IDE**: Cursor IDE (AI-native, VSCode-compatible)
- **Styling**: Tailwind CSS
- **UI Components**: shadcn/ui (headless UI built for Tailwind)
- **Charts (for Analytics)**: [Recharts](https://recharts.org/)
- **State Management**: React Context or Zustand (lightweight)
- **API Integration**: Axios or built-in fetch
- **Auth (Optional MVP scope)**: Clerk or simple token-based auth
- **QR Code Rendering**: `qrcode.react` or dynamic image URL from backend

---

## 📁 Folder Structure (App Router)

```bash
/app
  /shorten
    page.tsx         # Shorten URL form
  /analytics
    [id]/page.tsx    # Analytics view per shortened URL
  /dashboard
    page.tsx         # User dashboard (list of shortened URLs)
  layout.tsx
  page.tsx           # Landing page
/components
  URLForm.tsx
  URLCard.tsx
  QRCodeDisplay.tsx
  AnalyticsChart.tsx
/lib
  api.ts             # Central API utility
  types.ts
/styles
  globals.css
```

---

## 🚦 Milestone-Based Development Plan

### ✅ **Milestone 1: Project Setup**

**Goal**: Get Next.js project running inside Cursor IDE

- [x] Run: `npx create-next-app@latest gochop-frontend --app`
- [x] Add Tailwind CSS (`npx tailwindcss init -p`)
- [x] Install dependencies:

  ```bash
  pnpm add @shadcn/ui tailwind-variants axios recharts qrcode.react
  ```

- [x] Set up global layout and base styling
- [x] Commit to Git and initialize GitHub repo

---

### ✅ **Milestone 2: Landing Page (Marketing)**

**Goal**: Simple page with intro, features, and call-to-action

- [x] Hero section with name "GoChop"
- [x] Features: Context-aware URLs, QR codes, real-time analytics
- [x] CTA: "Shorten a URL" button ➝ /shorten

---

### ✅ **Milestone 3: Shorten URL Interface**

**Goal**: Form to create shortened URLs

- [x] Create `/shorten` page.
- [x] Implement a form with:
  - [x] Long URL input
  - [x] Optional alias input
  - [x] Optional tags or context input
- [x] Submit form data to the backend API (`/api/shorten`).
- [x] Display the result from the API, including:
  - [x] Short URL
  - [x] Copy button
  - [x] QR Code preview (using `QRCodeDisplay.tsx` component).

---

### ✅ **Milestone 4: Dashboard View**

**Goal**: Show list of previously shortened links for a user.

- [x] Create `/dashboard` page.
- [x] Fetch and display a list of URLs.
- [x] For each URL, show:
  - [x] Clicks
  - [x] Creation date
  - [x] Expiration date
- [x] Add a button to navigate to the analytics page for each link (`/analytics/[id]`).

---

### ✅ **Milestone 5: Analytics Page**

**Goal**: Visualize usage data for a specific shortened link.

- [x] Create dynamic route `/analytics/[id]`.
- [x] Use `Recharts` to display charts for:
  - [x] Clicks over time
  - [x] Referral sources
  - [x] Geographic data
- [x] Display the QR Code for the link.
- [ ] Add options to update link properties (e.g., expiration, access control).

---

### ✅ **Milestone 6: NextAuth.js Authentication System**

**Goal**: Implement secure, production-ready authentication using NextAuth.js with social login support.

- [ ] **NextAuth.js Setup**:
  - [ ] Install NextAuth.js and required dependencies
  - [ ] Configure NextAuth.js with PostgreSQL adapter
  - [ ] Set up database schema for NextAuth.js (users, accounts, sessions)
  - [ ] Configure environment variables for auth providers
- [ ] **Authentication Providers**:
  - [ ] Set up Google OAuth provider for social login
  - [ ] Set up GitHub OAuth provider for developer-friendly auth
  - [ ] Configure email/password provider for traditional auth
  - [ ] Add magic link authentication option
- [ ] **Frontend Implementation**:
  - [ ] Create authentication pages (/login, /register)
  - [ ] Implement NextAuth session provider and hooks
  - [ ] Create protected route wrapper component
  - [ ] Add authentication UI components (login/logout buttons)
  - [ ] Implement user profile management page
- [ ] **Session Management**:
  - [ ] Configure session strategy and JWT tokens
  - [ ] Implement session-aware API calls to Go backend
  - [ ] Add automatic token refresh handling
  - [ ] Create middleware for protected routes
- [ ] **User Experience**:
  - [ ] Add loading states for authentication
  - [ ] Implement proper error handling for auth failures
  - [ ] Create user onboarding flow for new registrations
  - [ ] Add social login buttons with proper branding

**Technical Implementation Details:**

**NextAuth.js Configuration:**

```typescript
// pages/api/auth/[...nextauth].ts
import NextAuth from "next-auth";
import GoogleProvider from "next-auth/providers/google";
import GitHubProvider from "next-auth/providers/github";
import EmailProvider from "next-auth/providers/email";
import { PostgresAdapter } from "@next-auth/postgres-adapter";

export default NextAuth({
  adapter: PostgresAdapter(process.env.DATABASE_URL),
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET,
    }),
    GitHubProvider({
      clientId: process.env.GITHUB_ID,
      clientSecret: process.env.GITHUB_SECRET,
    }),
    EmailProvider({
      server: process.env.EMAIL_SERVER,
      from: process.env.EMAIL_FROM,
    }),
  ],
  callbacks: {
    async jwt({ token, user, account }) {
      if (user) {
        token.userId = user.id;
        token.isAdmin = user.isAdmin || false;
      }
      return token;
    },
    async session({ session, token }) {
      session.userId = token.userId;
      session.isAdmin = token.isAdmin;
      return session;
    },
  },
  pages: {
    signIn: "/login",
    signUp: "/register",
    error: "/auth/error",
  },
});
```

**Protected API Calls:**

```typescript
// lib/api.ts
import { getSession } from "next-auth/react";

export const authenticatedFetch = async (
  url: string,
  options: RequestInit = {}
) => {
  const session = await getSession();

  return fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${session?.accessToken}`,
      "Content-Type": "application/json",
    },
  });
};
```

**Environment Variables Needed:**

```env
# NextAuth
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-secret-key

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# GitHub OAuth
GITHUB_ID=your-github-client-id
GITHUB_SECRET=your-github-client-secret

# Database
DATABASE_URL=postgresql://user:password@localhost:5432/gochop

# Email (for magic links)
EMAIL_SERVER=smtp://username:password@smtp.example.com:587
EMAIL_FROM=noreply@gochop.io
```

---

### ✅ **Milestone 7: Advanced Link Creation**

**Goal**: Enhance the form to support unique features.

- [ ] **Update Shorten URL Form (`/shorten`)**:
  - [ ] Add a UI to define rules for **Context-Aware Redirects** (e.g., "if user is on mobile, send to X").
  - [ ] Implement an interface for **A/B Testing** to add multiple destination URLs with traffic weights.
  - [ ] Add input fields for **Password Protection** and **Max Clicks** (self-destructing links).
- [ ] Create a password entry modal/page for users accessing protected links.

---

### ✅ **Milestone 8: Enhanced Dashboard & Analytics**

**Goal**: Display and manage advanced link features.

- [ ] **Update Dashboard (`/dashboard`)**:
  - [ ] Add icons/tags to indicate if a link is context-aware, an A/B test, or password-protected.
  - [ ] Implement an "Edit" flow to allow users to change a link's destination, which is key for **Dynamic QR Codes**.
- [ ] **Update Analytics Page (`/analytics/[id]`)**:
  - [ ] For A/B tests, display a comparison view showing the performance of each destination URL.
  - [ ] Add filters to the analytics to segment data by context (e.g., view clicks from "Desktop" only).

---

### ✅ **Milestone 9: Polish & Deploy**

**Goal**: Final touch-ups and deployment.

- [ ] Add favicon/logo.
- [ ] Ensure the application is responsive and works on mobile views.
- [ ] Set up environment variables for the backend API URL.
- [ ] Deploy the application to **Vercel**.

---

## 🔗 Frontend <–> Backend Integration

### `.env.local` Example

```env
NEXT_PUBLIC_API_BASE=https://api.gochop.io
```

### API Utility (`lib/api.ts`)

```ts
export const shortenUrl = async (data) => {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE}/api/shorten`, {
    method: "POST",
    body: JSON.stringify(data),
    headers: { "Content-Type": "application/json" },
  });
  return res.json();
};
```

### QR Code Integration

```tsx
<img src={`${process.env.NEXT_PUBLIC_API_BASE}/api/qrcode/${shortCode}`} />
```

---

## 🧪 Testing

- [ ] Write component tests using Jest.
- [ ] Write end-to-end tests using Playwright.

---

## 🪪 Authentication Notes

The authentication system is now implemented as a core feature (Milestone 6) rather than optional post-MVP. This ensures proper user isolation, security, and scalability from the start.
