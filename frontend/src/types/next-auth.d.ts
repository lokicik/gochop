import { DefaultSession, DefaultUser } from "next-auth";

declare module "next-auth" {
  interface Session extends DefaultSession {
    accessToken?: string;
    userId: string;
    isAdmin: boolean;
  }

  interface User extends DefaultUser {
    isAdmin?: boolean;
  }
}

declare module "next-auth/jwt" {
  interface JWT {
    accessToken?: string;
    userId: string;
    isAdmin: boolean;
  }
}
