"use client";

import { useState } from "react";
import QRCodeDisplay from "./QRCodeDisplay";

// A simple utility to format dates
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString(undefined, {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

export default function URLForm() {
  const [longUrl, setLongUrl] = useState("");
  const [alias, setAlias] = useState("");
  const [context, setContext] = useState("");
  const [shortUrl, setShortUrl] = useState("");
  const [expiresAt, setExpiresAt] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setIsLoading(true);
    setError("");
    setShortUrl("");
    setExpiresAt("");

    try {
      const res = await fetch("http://localhost:3001/api/shorten", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          long_url: longUrl,
          alias: alias,
          context: context,
        }),
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(
          data.error || "Something went wrong. Please try again."
        );
      }

      setShortUrl(data.short_url);
      setExpiresAt(data.expires_at);
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unexpected error occurred.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(shortUrl);
    // You could add a small notification here to confirm the copy.
  };

  return (
    <>
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label
            htmlFor="longUrl"
            className="block text-sm font-medium text-gray-700 dark:text-gray-300"
          >
            Your Long URL
          </label>
          <div className="mt-1">
            <input
              id="longUrl"
              name="longUrl"
              type="url"
              required
              value={longUrl}
              onChange={(e) => setLongUrl(e.target.value)}
              className="block w-full px-4 py-3 text-gray-900 bg-gray-100 border-gray-300 rounded-md shadow-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              placeholder="https://example.com/a-very-long-url-to-shorten"
              disabled={isLoading}
            />
          </div>
        </div>

        <div>
          <label
            htmlFor="alias"
            className="block text-sm font-medium text-gray-700 dark:text-gray-300"
          >
            Custom Alias (Optional)
          </label>
          <div className="flex mt-1 rounded-md shadow-sm">
            <span className="inline-flex items-center px-3 text-sm text-gray-500 border border-r-0 border-gray-300 rounded-l-md bg-gray-50 dark:bg-gray-600 dark:border-gray-500 dark:text-gray-300">
              gochop.io/
            </span>
            <input
              id="alias"
              name="alias"
              type="text"
              value={alias}
              onChange={(e) => setAlias(e.target.value)}
              className="flex-1 block w-full min-w-0 px-3 py-2 text-gray-900 bg-gray-100 border-gray-300 rounded-none rounded-r-md dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              placeholder="my-custom-link"
              disabled={isLoading}
            />
          </div>
        </div>

        <div>
          <label
            htmlFor="context"
            className="block text-sm font-medium text-gray-700 dark:text-gray-300"
          >
            Tags or Context (Optional)
          </label>
          <div className="mt-1">
            <input
              id="context"
              name="context"
              type="text"
              value={context}
              onChange={(e) => setContext(e.target.value)}
              className="block w-full px-4 py-3 text-gray-900 bg-gray-100 border-gray-300 rounded-md shadow-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              placeholder="campaign-2024, social-media, newsletter"
              disabled={isLoading}
            />
          </div>
          <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
            Add tags or context to help organize and track your links
          </p>
        </div>

        <div>
          <button
            type="submit"
            className="flex justify-center w-full px-4 py-3 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
            disabled={isLoading}
          >
            {isLoading ? "Shortening..." : "Shorten"}
          </button>
        </div>
      </form>

      {error && (
        <div className="p-4 mt-4 text-sm text-red-700 bg-red-100 rounded-md dark:bg-red-900 dark:text-red-200">
          {error}
        </div>
      )}

      {shortUrl && (
        <div className="p-4 mt-4 bg-green-100 rounded-md dark:bg-green-900">
          <label className="block text-sm font-medium text-green-700 dark:text-green-200">
            Your Short URL:
          </label>
          <div className="flex items-center mt-2 space-x-2">
            <input
              type="text"
              readOnly
              value={shortUrl}
              className="flex-grow px-3 py-2 text-green-900 bg-white border border-green-300 rounded-md dark:bg-gray-800 dark:text-green-100 dark:border-green-600"
            />
            <button
              onClick={handleCopy}
              className="px-4 py-2 text-sm font-medium text-white bg-green-600 border border-transparent rounded-md hover:bg-green-700"
            >
              Copy
            </button>
          </div>
          <p className="mt-2 text-xs text-green-600 dark:text-green-300">
            Expires on: {formatDate(expiresAt)}
          </p>
        </div>
      )}

      {shortUrl && <QRCodeDisplay shortUrl={shortUrl} />}
    </>
  );
}
