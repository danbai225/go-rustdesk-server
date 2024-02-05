import "@/styles/globals.css"
import { Inter as FontSans } from "next/font/google"

import { cn } from "@/lib/utils"

export const fontSans = FontSans({
  subsets: ["latin"],
  variable: "--font-sans",
})

export default function RootLayout({ children }) {
  return (
    <html lang="en" suppressHydrationWarning >
      <head />
      <body
        className={cn(
          "min-h-screen bg-custom  antialiased text-white",
          fontSans.variable
        )}
      >
       {children}
      </body>
    </html>
  )
}
