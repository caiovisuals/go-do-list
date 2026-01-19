import type { Metadata } from "next"
import { Geologica } from "next/font/google"
import "./globals.css"

const geologica = Geologica({
    variable: "--font-geologica",
    subsets: ["latin"],
})

export const metadata: Metadata = {
    title: "go-do-list",
    description: "To-do list usando React e Go!",
}

export default function RootLayout({ children, }: Readonly<{ children: React.ReactNode; }>) {
    return (
        <html lang="pt-br">
            <body className={`${geologica.variable} antialiased`}>
                {children}
            </body>
        </html>
    )
}
