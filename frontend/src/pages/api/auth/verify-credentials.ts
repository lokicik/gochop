import { NextApiRequest, NextApiResponse } from "next";
import { Pool } from "pg";
import bcrypt from "bcryptjs";

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

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method !== "POST") {
    return res.status(405).json({ message: "Method not allowed" });
  }

  const { email, password } = req.body;

  if (!email || !password) {
    return res.status(400).json({ message: "Email and password are required" });
  }

  try {
    // Find user and their password hash
    const result = await pool.query(
      `SELECT u.id, u.name, u.email, a.password 
       FROM users u 
       JOIN accounts a ON u.id = a."userId" 
       WHERE u.email = $1 AND a.type = 'credentials'`,
      [email.toLowerCase()]
    );

    if (result.rows.length === 0) {
      return res.status(401).json({ message: "Invalid credentials" });
    }

    const user = result.rows[0];

    // Verify password
    const isValidPassword = await bcrypt.compare(password, user.password);

    if (!isValidPassword) {
      return res.status(401).json({ message: "Invalid credentials" });
    }

    // Check if user is admin
    const adminEmails = (process.env.ADMIN_EMAILS || "admin@gochop.io").split(
      ","
    );
    const isAdmin = adminEmails.includes(user.email.toLowerCase());

    // Return user data (excluding password)
    res.status(200).json({
      id: user.id,
      name: user.name,
      email: user.email,
      isAdmin,
    });
  } catch (error) {
    console.error("Credential verification error:", error);
    res.status(500).json({ message: "Internal server error" });
  }
}
