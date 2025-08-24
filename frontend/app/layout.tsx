import "./globals.css";
import React from "react";

export const metadata = {
  title: "Release Calendar",
  description: "Manage company releases",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <header className="border-b">
          <div className="container py-4 flex items-center justify-between">
            <div className="text-xl font-semibold">Release Calendar</div>
            <a className="btn" href="https://thatqa.com" target="_blank" rel="noreferrer">Docs</a>
          </div>
        </header>
        <main className="container py-6">{children}</main>
      </body>
    </html>
  );
}
