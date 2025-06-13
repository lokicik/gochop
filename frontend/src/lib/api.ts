import { getSession } from "next-auth/react";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE || "http://localhost:3001";

// Helper function to get NextAuth session token for backend API calls
export async function getAuthHeaders(): Promise<HeadersInit> {
  const session = await getSession();

  if (!session?.accessToken) {
    throw new Error("No valid session found");
  }

  return {
    Authorization: `Bearer ${session.accessToken}`,
    "Content-Type": "application/json",
  };
}

// Authenticated fetch wrapper
export async function authenticatedFetch(
  endpoint: string,
  options: RequestInit = {}
): Promise<Response> {
  const headers = await getAuthHeaders();

  return fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      ...headers,
      ...options.headers,
    },
  });
}

// API functions
export const api = {
  // User endpoints
  async getUserProfile() {
    return authenticatedFetch("/api/user/profile");
  },

  async getUserStats() {
    return authenticatedFetch("/api/user/stats");
  },

  async getUserLinks() {
    return authenticatedFetch("/api/user/links");
  },

  // Admin endpoints
  async getAllLinks() {
    return authenticatedFetch("/api/admin/links");
  },

  async getAnalytics(shortCode: string) {
    return authenticatedFetch(`/api/admin/analytics/${shortCode}`);
  },

  // Authenticated endpoints
  async shortenUrl(data: {
    long_url: string;
    alias?: string;
    context?: string;
  }) {
    return authenticatedFetch("/api/user/shorten", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  async getQRCode(shortCode: string) {
    return fetch(`${API_BASE_URL}/api/qrcode/${shortCode}`);
  },
};
