import "@/styles/globals.css"

export default function RootLayout({ children }) {
  return (
    <div  className={
        "min-h-screen bg-custom  antialiased text-white"
      }>
   {children}
    </div>

  )
}
