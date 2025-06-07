"use client";

import { useState, useEffect } from "react";
import Link from "next/link";

// Type definitions matching the backend API response
interface LinkInfo {
  id: number;
  short_code: string;
  long_url: string;
  context: string;
  created_at: string;
  expires_at: string;
  click_count: number;
}

// Utility function to format dates
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

// Utility function to truncate long URLs
const truncateUrl = (url: string, maxLength: number = 50) => {
  return url.length > maxLength ? url.substring(0, maxLength) + "..." : url;
};

export default function DashboardPage() {
  const [links, setLinks] = useState<LinkInfo[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    fetchLinks();
  }, []);

  const fetchLinks = async () => {
    try {
      setIsLoading(true);

      // Get admin token for development
      const tokenResponse = await fetch(
        "http://localhost:3001/api/auth/dev-token"
      );
      if (!tokenResponse.ok) {
        throw new Error("Failed to get auth token");
      }
      const tokenData = await tokenResponse.json();

      // Fetch links with authentication
      const response = await fetch("http://localhost:3001/api/admin/links", {
        headers: {
          Authorization: `Bearer ${tokenData.token}`,
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error("Failed to fetch links");
      }

      const data = await response.json();
      setLinks(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = (shortCode: string) => {
    const shortUrl = `http://localhost:3001/${shortCode}`;
    navigator.clipboard.writeText(shortUrl);
    // You could add a toast notification here
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="text-xl text-gray-600 dark:text-gray-400">
          Loading your links...
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="text-xl text-red-600 dark:text-red-400">
          Error: {error}
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            Your Dashboard
          </h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            Manage and track your shortened links
          </p>
          <div className="mt-4">
            <Link
              href="/shorten"
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              Create New Link
            </Link>
          </div>
        </div>

        {/* Links Table */}
        {links.length === 0 ? (
          <div className="text-center py-12">
            <div className="text-gray-500 dark:text-gray-400 text-lg mb-4">
              No links found
            </div>
            <Link
              href="/shorten"
              className="text-indigo-600 hover:text-indigo-500 dark:text-indigo-400"
            >
              Create your first shortened link
            </Link>
          </div>
        ) : (
          <div className="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
            <ul className="divide-y divide-gray-200 dark:divide-gray-700">
              {links.map((link) => (
                <li key={link.id} className="px-6 py-4">
                  <div className="flex items-center justify-between">
                    <div className="flex-1 min-w-0">
                      {/* Short URL */}
                      <div className="flex items-center mb-2">
                        <code className="text-sm font-mono text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-900/30 px-2 py-1 rounded">
                          gochop.io/{link.short_code}
                        </code>
                        <button
                          onClick={() => copyToClipboard(link.short_code)}
                          className="ml-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                          title="Copy to clipboard"
                        >
                          <svg
                            className="w-4 h-4"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                            />
                          </svg>
                        </button>
                      </div>

                      {/* Long URL */}
                      <div className="text-sm text-gray-900 dark:text-gray-100 mb-1">
                        <span className="font-medium">→ </span>
                        <a
                          href={link.long_url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
                          title={link.long_url}
                        >
                          {truncateUrl(link.long_url)}
                        </a>
                      </div>

                      {/* Context/Tags */}
                      {link.context && (
                        <div className="text-xs text-gray-500 dark:text-gray-400 mb-2">
                          <span className="font-medium">Tags:</span>{" "}
                          {link.context}
                        </div>
                      )}

                      {/* Metadata */}
                      <div className="flex flex-wrap items-center text-xs text-gray-500 dark:text-gray-400 space-x-4">
                        <span>Created: {formatDate(link.created_at)}</span>
                        <span>Expires: {formatDate(link.expires_at)}</span>
                        <span className="flex items-center">
                          <svg
                            className="w-3 h-3 mr-1"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                            />
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                            />
                          </svg>
                          {link.click_count} clicks
                        </span>
                      </div>
                    </div>

                    {/* Actions */}
                    <div className="flex items-center space-x-2">
                      <Link
                        href={`/analytics/${link.short_code}`}
                        className="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-xs font-medium rounded text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                      >
                        Analytics
                      </Link>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}
