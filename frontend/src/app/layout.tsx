import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "GymLog",
  description: "Your personal fitness tracking application",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className="font-sans antialiased"
        style={{
          fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
        }}
      >
        {children}
      </body>
    </html>
  );
}
