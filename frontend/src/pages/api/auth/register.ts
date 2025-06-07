import { NextApiRequest, NextApiResponse } from "next";
import { Pool } from "pg";
import bcrypt from "bcryptjs";
import { rateLimit } from "./rate-limit";

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

// Rate limiting: 5 requests per 15 minutes
const rateLimiter = rateLimit({ maxRequests: 5, windowMs: 15 * 60 * 1000 });

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method !== "POST") {
    return res.status(405).json({ message: "Method not allowed" });
  }

  // Apply rate limiting
  await new Promise<void>((resolve) => {
    rateLimiter(req, res, () => resolve());
  });

  const { name, email, password } = req.body;

  // Validate input
  if (!name || !email || !password) {
    return res.status(400).json({
      message: "Name, email, and password are required",
    });
  }

  // Enhanced password validation
  if (password.length < 8) {
    return res.status(400).json({
      message: "Password must be at least 8 characters long",
    });
  }

  // Check password strength
  const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/;
  if (!passwordRegex.test(password)) {
    return res.status(400).json({
      message:
        "Password must contain at least one uppercase letter, one lowercase letter, and one number",
    });
  }

  // Validate email format
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(email)) {
    return res.status(400).json({
      message: "Please provide a valid email address",
    });
  }

  // Validate name length and content
  if (name.trim().length < 2 || name.trim().length > 100) {
    return res.status(400).json({
      message: "Name must be between 2 and 100 characters long",
    });
  }

  // Sanitize inputs
  const sanitizedName = name.trim();
  const sanitizedEmail = email.toLowerCase().trim();

  try {
    // Check if user already exists
    const existingUser = await pool.query(
      "SELECT id FROM users WHERE email = $1",
      [sanitizedEmail]
    );

    if (existingUser.rows.length > 0) {
      return res.status(400).json({
        message: "User with this email already exists",
      });
    }

    // Hash password
    const saltRounds = 12;
    const hashedPassword = await bcrypt.hash(password, saltRounds);

    // Create user
    const result = await pool.query(
      `INSERT INTO users (name, email, "emailVerified", image) 
       VALUES ($1, $2, NULL, NULL) 
       RETURNING id, name, email`,
      [sanitizedName, sanitizedEmail]
    );

    const user = result.rows[0];

    // Create account record for credentials authentication
    await pool.query(
      `INSERT INTO accounts ("userId", type, provider, "providerAccountId", password)
       VALUES ($1, 'credentials', 'credentials', $2, $3)`,
      [user.id, sanitizedEmail, hashedPassword]
    );

    // Return success (don't return password or sensitive data)
    res.status(201).json({
      message: "User created successfully",
      user: {
        id: user.id,
        name: user.name,
        email: user.email,
      },
    });
  } catch (error) {
    console.error("Registration error:", error);
    res.status(500).json({
      message: "Internal server error",
    });
  }
}
