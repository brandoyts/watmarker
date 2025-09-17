import { Route, Routes } from "react-router-dom"

import AppTitle from "./components/app-title"
import { BackgroundBeams } from "@/components/ui/shadcn-io/background-beams"
import Home from "./pages/home"

function App() {
  return (
    <main className="h-screen w-screen bg-white relative flex flex-col items-center justify-start antialiased pt-12 sm:pt-16 lg:pt-20">
      <BackgroundBeams className="absolute inset-0" />
      <AppTitle />
      <Routes>
        <Route path="/" element={<Home />} />
      </Routes>
    </main>
  )
}

export default App
