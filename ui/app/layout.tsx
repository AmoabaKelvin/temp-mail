import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Is-Temp",
  description:
    "Receive emails anonymously with our free temporary email address generator.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        {children}
        {process.env.NODE_ENV === 'production' && (
          <script
            defer
            src="https://static.cloudflareinsights.com/beacon.min.js"
            data-cf-beacon='{"token": "48d1fad733f04b8f9418856a6870d907"}'
          />
        )}
      </body>
    </html>
  );
}
