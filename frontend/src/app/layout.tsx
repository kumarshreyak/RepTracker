import type { Metadata } from "next";
import "./globals.css";
import { Providers } from "@/components/providers";
import { auth } from "@/auth";

export const metadata: Metadata = {
  title: "GymLog",
  description: "Track your workouts and progress",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const session = await auth();

  return (
    <html lang="en">
      <body
        className="font-sans antialiased"
        style={{
          fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
        }}
      >
        <Providers session={session}>
          {children}
        </Providers>
      </body>
    </html>
  );
}
