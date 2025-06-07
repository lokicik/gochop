"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
} from "recharts";
import QRCodeDisplay from "../../../components/QRCodeDisplay";

// Type definitions matching the backend API response
interface AnalyticsData {
  short_code: string;
  total_clicks: number;
  clicks_by_date: DailyClickData[];
  top_referrers: ReferrerData[];
  top_user_agents: UserAgentData[];
  geographic_data: GeographicData[];
}

interface DailyClickData {
  date: string;
  clicks: number;
}

interface ReferrerData {
  referrer: string;
  clicks: number;
}

interface UserAgentData {
  user_agent: string;
  clicks: number;
}

interface GeographicData {
  country: string;
  region: string;
  city: string;
  clicks: number;
}

// Colors for charts
const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042", "#8884D8"];

export default function AnalyticsPage() {
  const params = useParams();
  const shortCode = params?.id as string;

  const [analytics, setAnalytics] = useState<AnalyticsData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (shortCode) {
      fetchAnalytics();
    }
  }, [shortCode]);

  const fetchAnalytics = async () => {
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

      // Fetch analytics with authentication
      const response = await fetch(
        `http://localhost:3001/api/admin/analytics/${shortCode}`,
        {
          headers: {
            Authorization: `Bearer ${tokenData.token}`,
            "Content-Type": "application/json",
          },
        }
      );

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error("Link not found");
        }
        throw new Error("Failed to fetch analytics");
      }

      const data = await response.json();
      setAnalytics(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="text-xl text-gray-600 dark:text-gray-400">
          Loading analytics...
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col justify-center items-center min-h-screen">
        <div className="text-xl text-red-600 dark:text-red-400 mb-4">
          Error: {error}
        </div>
        <Link
          href="/dashboard"
          className="text-indigo-600 hover:text-indigo-500 dark:text-indigo-400"
        >
          ← Back to Dashboard
        </Link>
      </div>
    );
  }

  if (!analytics) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="text-xl text-gray-600 dark:text-gray-400">
          No analytics data available
        </div>
      </div>
    );
  }

  const shortUrl = `http://localhost:3001/${shortCode}`;

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
                Analytics for /{shortCode}
              </h1>
              <p className="mt-2 text-gray-600 dark:text-gray-400">
                Detailed insights for your shortened link
              </p>
            </div>
            <Link
              href="/dashboard"
              className="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
            >
              ← Back to Dashboard
            </Link>
          </div>

          {/* Short URL Display */}
          <div className="mt-4 p-4 bg-white dark:bg-gray-800 rounded-lg shadow">
            <div className="flex items-center justify-between">
              <div>
                <code className="text-lg font-mono text-indigo-600 dark:text-indigo-400">
                  {shortUrl}
                </code>
              </div>
              <button
                onClick={() => navigator.clipboard.writeText(shortUrl)}
                className="ml-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                title="Copy to clipboard"
              >
                <svg
                  className="w-5 h-5"
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
          </div>
        </div>

        {/* Stats Overview */}
        <div className="grid grid-cols-1 md:grid-cols-5 gap-6 mb-8">
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <div className="text-3xl font-bold text-indigo-600 dark:text-indigo-400">
              {analytics.total_clicks}
            </div>
            <div className="text-sm text-gray-600 dark:text-gray-400">
              Total Clicks
            </div>
          </div>
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <div className="text-3xl font-bold text-green-600 dark:text-green-400">
              {analytics.clicks_by_date?.length || 0}
            </div>
            <div className="text-sm text-gray-600 dark:text-gray-400">
              Active Days
            </div>
          </div>
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <div className="text-3xl font-bold text-yellow-600 dark:text-yellow-400">
              {analytics.top_referrers?.length || 0}
            </div>
            <div className="text-sm text-gray-600 dark:text-gray-400">
              Referrer Sources
            </div>
          </div>
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <div className="text-3xl font-bold text-purple-600 dark:text-purple-400">
              {analytics.top_user_agents?.length || 0}
            </div>
            <div className="text-sm text-gray-600 dark:text-gray-400">
              User Agents
            </div>
          </div>
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <div className="text-3xl font-bold text-teal-600 dark:text-teal-400">
              {analytics.geographic_data?.length || 0}
            </div>
            <div className="text-sm text-gray-600 dark:text-gray-400">
              Locations
            </div>
          </div>
        </div>

        {/* Charts Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
          {/* Clicks Over Time */}
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              Clicks Over Time (Last 30 Days)
            </h3>
            {analytics.clicks_by_date && analytics.clicks_by_date.length > 0 ? (
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={[...analytics.clicks_by_date].reverse()}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis
                    dataKey="date"
                    tickFormatter={(value) =>
                      new Date(value).toLocaleDateString()
                    }
                  />
                  <YAxis />
                  <Tooltip
                    labelFormatter={(value) =>
                      new Date(value).toLocaleDateString()
                    }
                    labelStyle={{ color: "#374151" }}
                  />
                  <Line
                    type="monotone"
                    dataKey="clicks"
                    stroke="#8884d8"
                    strokeWidth={2}
                  />
                </LineChart>
              </ResponsiveContainer>
            ) : (
              <div className="text-center text-gray-500 dark:text-gray-400 py-8">
                No click data available
              </div>
            )}
          </div>

          {/* Top Referrers */}
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              Top Referrers
            </h3>
            {analytics.top_referrers && analytics.top_referrers.length > 0 ? (
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={analytics.top_referrers.slice(0, 5)}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis
                    dataKey="referrer"
                    angle={-45}
                    textAnchor="end"
                    height={80}
                  />
                  <YAxis />
                  <Tooltip labelStyle={{ color: "#374151" }} />
                  <Bar dataKey="clicks" fill="#8884d8" />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="text-center text-gray-500 dark:text-gray-400 py-8">
                No referrer data available
              </div>
            )}
          </div>

          {/* User Agents Pie Chart */}
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              Top User Agents
            </h3>
            {analytics.top_user_agents &&
            analytics.top_user_agents.length > 0 ? (
              <ResponsiveContainer width="100%" height={300}>
                <PieChart>
                  <Pie
                    data={analytics.top_user_agents.slice(0, 5)}
                    cx="50%"
                    cy="50%"
                    labelLine={false}
                    label={({ user_agent, percent }) =>
                      `${
                        percent > 10 ? user_agent.substring(0, 20) + "..." : ""
                      }`
                    }
                    outerRadius={80}
                    fill="#8884d8"
                    dataKey="clicks"
                  >
                    {analytics.top_user_agents
                      .slice(0, 5)
                      .map((entry, index) => (
                        <Cell
                          key={`cell-${index}`}
                          fill={COLORS[index % COLORS.length]}
                        />
                      ))}
                  </Pie>
                  <Tooltip formatter={(value, name) => [value, name]} />
                </PieChart>
              </ResponsiveContainer>
            ) : (
              <div className="text-center text-gray-500 dark:text-gray-400 py-8">
                No user agent data available
              </div>
            )}
          </div>

          {/* Geographic Data */}
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              Geographic Distribution
            </h3>
            {analytics.geographic_data &&
            analytics.geographic_data.length > 0 ? (
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={analytics.geographic_data.slice(0, 10)}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis
                    dataKey="country"
                    angle={-45}
                    textAnchor="end"
                    height={80}
                  />
                  <YAxis />
                  <Tooltip
                    labelStyle={{ color: "#374151" }}
                    formatter={(value, name, props) => [
                      value,
                      `${props.payload.city}, ${props.payload.region}, ${props.payload.country}`,
                    ]}
                  />
                  <Bar dataKey="clicks" fill="#10B981" />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="text-center text-gray-500 dark:text-gray-400 py-8">
                No geographic data available
              </div>
            )}
          </div>

          {/* QR Code */}
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              QR Code
            </h3>
            <div className="flex justify-center">
              <QRCodeDisplay shortUrl={shortUrl} />
            </div>
          </div>
        </div>

        {/* Data Tables */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Referrers Table */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
            <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                All Referrers
              </h3>
            </div>
            <div className="overflow-y-auto max-h-64">
              {analytics.top_referrers && analytics.top_referrers.length > 0 ? (
                <table className="w-full text-sm">
                  <thead className="bg-gray-50 dark:bg-gray-700">
                    <tr>
                      <th className="px-6 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                        Referrer
                      </th>
                      <th className="px-6 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                        Clicks
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                    {analytics.top_referrers.map((referrer, index) => (
                      <tr key={index}>
                        <td className="px-6 py-4 text-gray-900 dark:text-gray-100">
                          {referrer.referrer}
                        </td>
                        <td className="px-6 py-4 text-gray-900 dark:text-gray-100">
                          {referrer.clicks}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              ) : (
                <div className="p-6 text-center text-gray-500 dark:text-gray-400">
                  No referrer data available
                </div>
              )}
            </div>
          </div>

          {/* User Agents Table */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
            <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                All User Agents
              </h3>
            </div>
            <div className="overflow-y-auto max-h-64">
              {analytics.top_user_agents &&
              analytics.top_user_agents.length > 0 ? (
                <table className="w-full text-sm">
                  <thead className="bg-gray-50 dark:bg-gray-700">
                    <tr>
                      <th className="px-6 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                        User Agent
                      </th>
                      <th className="px-6 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                        Clicks
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                    {analytics.top_user_agents.map((userAgent, index) => (
                      <tr key={index}>
                        <td className="px-6 py-4 text-gray-900 dark:text-gray-100">
                          {userAgent.user_agent}
                        </td>
                        <td className="px-6 py-4 text-gray-900 dark:text-gray-100">
                          {userAgent.clicks}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              ) : (
                <div className="p-6 text-center text-gray-500 dark:text-gray-400">
                  No user agent data available
                </div>
              )}
            </div>
          </div>

          {/* Geographic Data Table */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
            <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                Geographic Locations
              </h3>
            </div>
            <div className="overflow-y-auto max-h-64">
              {analytics.geographic_data &&
              analytics.geographic_data.length > 0 ? (
                <table className="w-full text-sm">
                  <thead className="bg-gray-50 dark:bg-gray-700">
                    <tr>
                      <th className="px-6 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                        Location
                      </th>
                      <th className="px-6 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                        Clicks
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                    {analytics.geographic_data.map((location, index) => (
                      <tr key={index}>
                        <td className="px-6 py-4 text-gray-900 dark:text-gray-100">
                          {location.city}, {location.region}, {location.country}
                        </td>
                        <td className="px-6 py-4 text-gray-900 dark:text-gray-100">
                          {location.clicks}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              ) : (
                <div className="p-6 text-center text-gray-500 dark:text-gray-400">
                  No geographic data available
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
