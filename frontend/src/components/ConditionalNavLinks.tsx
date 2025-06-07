"use client";

import Link from "next/link";
import { useSession } from "next-auth/react";
import AuthButton from "./AuthButton";

export default function ConditionalNavLinks() {
  const { data: session, status } = useSession();

  return (
    <div className="flex items-center space-x-4">
      {/* Only show protected routes if user is authenticated */}
      {status === "authenticated" && session && (
        <>
          <Link
            href="/shorten"
            className="text-gray-700 dark:text-gray-300 hover:text-indigo-600 dark:hover:text-indigo-400 px-3 py-2 rounded-md text-sm font-medium"
          >
            Shorten URL
          </Link>
          <Link
            href="/dashboard"
            className="text-gray-700 dark:text-gray-300 hover:text-indigo-600 dark:hover:text-indigo-400 px-3 py-2 rounded-md text-sm font-medium"
          >
            Dashboard
          </Link>
        </>
      )}

      {/* Always show auth button */}
      <AuthButton />
    </div>
  );
}
