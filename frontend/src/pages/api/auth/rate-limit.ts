import { NextApiRequest, NextApiResponse } from "next";

// Simple in-memory rate limiting (for production, use Redis)
const requests = new Map();

export interface RateLimitConfig {
  maxRequests: number;
  windowMs: number;
}

export function rateLimit(config: RateLimitConfig) {
  return (req: NextApiRequest, res: NextApiResponse, next: () => void) => {
    const ip =
      req.headers["x-forwarded-for"] || req.socket.remoteAddress || "unknown";
    const key = `${ip}-${req.url}`;
    const now = Date.now();

    if (!requests.has(key)) {
      requests.set(key, { count: 1, resetTime: now + config.windowMs });
      return next();
    }

    const requestData = requests.get(key);

    if (now > requestData.resetTime) {
      // Reset the window
      requests.set(key, { count: 1, resetTime: now + config.windowMs });
      return next();
    }

    if (requestData.count >= config.maxRequests) {
      return res.status(429).json({
        message: "Too many requests. Please try again later.",
        retryAfter: Math.ceil((requestData.resetTime - now) / 1000),
      });
    }

    requestData.count++;
    return next();
  };
}

// Clean up old entries periodically
setInterval(() => {
  const now = Date.now();
  for (const [key, data] of requests.entries()) {
    if (now > data.resetTime) {
      requests.delete(key);
    }
  }
}, 60000); // Clean up every minute
