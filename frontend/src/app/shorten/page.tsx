import URLForm from "@/components/URLForm";

export default function ShortenPage() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="w-full max-w-2xl p-8 space-y-8 bg-white rounded-lg shadow-md dark:bg-gray-800">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-900 dark:text-white">
            Shorten a Long URL
          </h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            Enter your URL below to create a short, shareable link.
          </p>
        </div>
        <URLForm />
      </div>
    </div>
  );
}
