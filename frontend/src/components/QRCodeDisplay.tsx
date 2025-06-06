"use client";

interface QRCodeDisplayProps {
  shortUrl: string;
}

const QRCodeDisplay = ({ shortUrl }: QRCodeDisplayProps) => {
  if (!shortUrl) {
    return null;
  }

  // Extract the short code from the full URL (e.g., http://localhost:3001/xYz123)
  const shortCode = shortUrl.split("/").pop();

  if (!shortCode) {
    return <p>Could not generate QR code.</p>;
  }

  const qrCodeUrl = `http://localhost:3001/api/qrcode/${shortCode}`;

  return (
    <div className="p-4 mt-6 text-center bg-gray-100 border-t border-gray-200 rounded-b-lg dark:bg-gray-700 dark:border-gray-600">
      <h3 className="text-lg font-medium text-gray-900 dark:text-white">
        Your QR Code
      </h3>
      <div className="flex justify-center mt-2">
        <img src={qrCodeUrl} alt="QR Code" className="w-48 h-48" />
      </div>
      <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
        Scan this with your phone to open the link.
      </p>
    </div>
  );
};

export default QRCodeDisplay;
