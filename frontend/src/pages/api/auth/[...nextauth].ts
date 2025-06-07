import NextAuth, { NextAuthOptions } from "next-auth";
import GoogleProvider from "next-auth/providers/google";
import CredentialsProvider from "next-auth/providers/credentials";
import PgAdapter from "@auth/pg-adapter";
import { Pool } from "pg";

// Create PostgreSQL connection pool
const pool = new Pool({
  connectionString:
    process.env.DATABASE_URL ||
    `postgres://gochop_user:gochop_password@localhost:5432/gochop`,
  ssl:
    process.env.NODE_ENV === "production"
      ? { rejectUnauthorized: false }
      : false,
});

export const authOptions: NextAuthOptions = {
  adapter: PgAdapter(pool),

  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID!,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
    }),
    CredentialsProvider({
      name: "credentials",
      credentials: {
        email: { label: "Email", type: "email" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        if (!credentials?.email || !credentials?.password) {
          return null;
        }

        // TODO: Add password verification logic here
        // This will be implemented when we create the backend auth endpoints
        try {
          const response = await fetch(
            `${process.env.NEXTAUTH_URL}/api/auth/verify-credentials`,
            {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({
                email: credentials.email,
                password: credentials.password,
              }),
            }
          );

          if (response.ok) {
            const user = await response.json();
            return {
              id: user.id,
              email: user.email,
              name: user.name,
              isAdmin: user.isAdmin || false,
            };
          }
        } catch (error) {
          console.error("Authentication error:", error);
        }

        return null;
      },
    }),
  ],

  callbacks: {
    async jwt({ token, user, account }) {
      // Persist the OAuth access_token and user id to the token right after signin
      if (account) {
        token.accessToken = account.access_token;
      }
      if (user) {
        token.userId = user.id;
        // Check if user is admin (this will be set during user creation)
        token.isAdmin = user.isAdmin || false;
      }
      return token;
    },

    async session({ session, token }) {
      // Send properties to the client
      if (token) {
        session.accessToken = token.accessToken as string;
        session.userId = token.userId as string;
        session.isAdmin = token.isAdmin as boolean;
      }
      return session;
    },

    async signIn() {
      // Allow sign in
      return true;
    },

    async redirect({ url, baseUrl }) {
      // Allows relative callback URLs
      if (url.startsWith("/")) return `${baseUrl}${url}`;
      // Allows callback URLs on the same origin
      else if (new URL(url).origin === baseUrl) return url;
      return baseUrl;
    },
  },

  pages: {
    signIn: "/login",
    error: "/auth/error",
    newUser: "/dashboard", // Redirect new users to dashboard
  },

  session: {
    strategy: "jwt",
    maxAge: 30 * 24 * 60 * 60, // 30 days
  },

  jwt: {
    maxAge: 30 * 24 * 60 * 60, // 30 days
  },

  secret: process.env.NEXTAUTH_SECRET,

  debug: process.env.NODE_ENV === "development",

  // Security configuration
  useSecureCookies: process.env.NODE_ENV === "production",
  cookies: {
    sessionToken: {
      name:
        process.env.NODE_ENV === "production"
          ? "__Secure-next-auth.session-token"
          : "next-auth.session-token",
      options: {
        httpOnly: true,
        sameSite: "lax",
        path: "/",
        secure: process.env.NODE_ENV === "production",
      },
    },
    callbackUrl: {
      name:
        process.env.NODE_ENV === "production"
          ? "__Secure-next-auth.callback-url"
          : "next-auth.callback-url",
      options: {
        sameSite: "lax",
        path: "/",
        secure: process.env.NODE_ENV === "production",
      },
    },
    csrfToken: {
      name:
        process.env.NODE_ENV === "production"
          ? "__Host-next-auth.csrf-token"
          : "next-auth.csrf-token",
      options: {
        httpOnly: true,
        sameSite: "lax",
        path: "/",
        secure: process.env.NODE_ENV === "production",
      },
    },
  },
};

export default NextAuth(authOptions);
