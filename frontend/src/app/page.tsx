"use client";

import Link from "next/link";
import { useSession } from "next-auth/react";

export default function HomePage() {
  const { data: session, status } = useSession();

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="text-center">
        <h1 className="text-6xl font-bold text-gray-900 dark:text-white">
          GoChop
        </h1>
        <p className="mt-4 text-xl text-gray-600 dark:text-gray-300">
          The context-aware, intelligent URL shortener of the future.
        </p>
        <div className="mt-8 space-x-4">
          {status === "authenticated" && session ? (
            <Link
              href="/shorten"
              className="px-8 py-3 text-lg font-semibold text-white bg-indigo-600 rounded-md hover:bg-indigo-700"
            >
              Shorten a URL
            </Link>
          ) : (
            <>
              <Link
                href="/register"
                className="px-8 py-3 text-lg font-semibold text-white bg-indigo-600 rounded-md hover:bg-indigo-700"
              >
                Get Started
              </Link>
              <Link
                href="/login"
                className="px-8 py-3 text-lg font-semibold text-indigo-600 bg-white border border-indigo-600 rounded-md hover:bg-indigo-50 dark:text-indigo-400 dark:bg-gray-800 dark:border-indigo-400 dark:hover:bg-gray-700"
              >
                Sign In
              </Link>
            </>
          )}
        </div>
      </div>
      <div className="mt-16 text-center">
        <h2 className="text-3xl font-bold text-gray-900 dark:text-white">
          Features
        </h2>
        <div className="flex-wrap items-center justify-center mt-6 space-y-4 md:flex md:space-y-0 md:space-x-8">
          <div className="p-4 bg-white rounded-lg shadow-md dark:bg-gray-800">
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
              Context-Aware Links
            </h3>
            <p className="mt-2 text-gray-600 dark:text-gray-400">
              Redirect based on location, device, or time.
            </p>
          </div>
          <div className="p-4 bg-white rounded-lg shadow-md dark:bg-gray-800">
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
              Dynamic QR Codes
            </h3>
            <p className="mt-2 text-gray-600 dark:text-gray-400">
              Update the destination, not the code.
            </p>
          </div>
          <div className="p-4 bg-white rounded-lg shadow-md dark:bg-gray-800">
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
              Real-Time Analytics
            </h3>
            <p className="mt-2 text-gray-600 dark:text-gray-400">
              Track every click and gain insights.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
