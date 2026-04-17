import type { Metadata } from "next"
import { Geologica } from "next/font/google"
import Script from "next/script"
import "./globals.css"

const geologica = Geologica({
    variable: "--font-geologica",
    subsets: ["latin"],
})

export const metadata: Metadata = {
    title: "Go do List",
    description: "To-do list usando React/NextJS, Ws e Golang!",
}

export default function RootLayout({ children, }: Readonly<{ children: React.ReactNode; }>) {
    return (
        <html lang="pt-br">
            <head>
                <Script 
                    id=""
                />
            </head>
            <body className={`${geologica.variable} antialiased`}>
                {children}
            </body>
        </html>
    )
}
