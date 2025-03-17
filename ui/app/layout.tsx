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
      <body>{children}</body>
    </html>
  );
}
