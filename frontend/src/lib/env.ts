// Environment variable validation
export function validateEnvVars() {
  const requiredEnvVars = ["NEXTAUTH_SECRET", "NEXTAUTH_URL", "DATABASE_URL"];

  const missingVars = requiredEnvVars.filter(
    (varName) => !process.env[varName]
  );

  if (missingVars.length > 0) {
    throw new Error(
      `Missing required environment variables: ${missingVars.join(", ")}`
    );
  }

  // Validate NEXTAUTH_SECRET is strong enough
  const secret = process.env.NEXTAUTH_SECRET;
  if (secret && secret.length < 32) {
    throw new Error("NEXTAUTH_SECRET must be at least 32 characters long");
  }

  // Warn about using development defaults in production
  if (process.env.NODE_ENV === "production") {
    if (process.env.DATABASE_URL?.includes("localhost")) {
      console.warn("⚠️  WARNING: Using localhost database URL in production!");
    }

    if (secret === "your-super-secret-key-here") {
      throw new Error(
        "❌ SECURITY: Change default NEXTAUTH_SECRET in production!"
      );
    }
  }
}

// Validate environment on module load
if (typeof window === "undefined") {
  try {
    validateEnvVars();
    console.log("✅ Environment variables validated");
  } catch (error) {
    console.error("❌ Environment validation failed:", error);
    if (process.env.NODE_ENV === "production") {
      process.exit(1);
    }
  }
}
