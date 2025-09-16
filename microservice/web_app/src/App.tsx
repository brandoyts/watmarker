import {
  Dropzone,
  DropzoneContent,
  DropzoneEmptyState,
} from "@/components/ui/shadcn-io/dropzone"
import { MAXIMUM_IMAGE_SIZE, MINIMUM_IMAGE_SIZE } from "@/constants"

import { BackgroundBeams } from "@/components/ui/shadcn-io/background-beams"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Loader2Icon } from "lucide-react"
import { useState } from "react"

function App() {
  const [files, setFiles] = useState<File[] | undefined>()
  const [isLoading, setIsLoading] = useState(false)
  const handleDrop = (files: File[]) => {
    console.log(files)
    setFiles(files)
  }

  const handleSubmit = () => {
    // simulate data fetching
    setIsLoading(true)
    setTimeout(() => {
      setIsLoading(false)
    }, 3000)

  }

  return (
    <>
      <div className="h-screen w-screen bg-white relative flex flex-col items-center justify-start antialiased pt-12 sm:pt-16 lg:pt-20">
        <BackgroundBeams className="absolute inset-0" />

        <h1
          className="
            relative z-10 
            text-3xl sm:text-5xl md:text-7xl lg:text-8xl 
            bg-clip-text text-transparent 
            bg-gradient-to-r from-blue-400 via-blue-500 to-indigo-600
            text-center font-sans font-extrabold tracking-tight leading-tight
            mb-4 sm:mb-6
          "
        >
          Watmarker
        </h1>
        <form
          action="#"
          onSubmit={handleSubmit}
          className="max-w-2xl mx-auto p-4 flex flex-col items-center gap-4 sm:gap-5 relative w-full"
        >
          <Dropzone
            disabled={isLoading}
            accept={{ "image/*": [] }}
            maxFiles={1}
            maxSize={MAXIMUM_IMAGE_SIZE}
            minSize={MINIMUM_IMAGE_SIZE}
            onDrop={handleDrop}
            onError={console.error}
            src={files}
            className="w-full max-w-md sm:max-w-lg bg-white hover:bg-white border-gray-300 hover:border-blue-600 rounded-xl "
          >
            <DropzoneEmptyState />
            <DropzoneContent />
          </Dropzone>
          <div className="flex flex-col sm:flex-row gap-3 w-full max-w-lg">
            <Input disabled={isLoading} placeholder="watermark" className="bg-white flex-1 rounded-xl border border-gray-300 px-3 py-2 sm:py-3" />
            <Button disabled={isLoading} className="cursor-pointer text-white px-4 sm:px-6 py-2 sm:py-3 rounded-xl bg-blue-500 hover:bg-blue-600 transition">
              {isLoading && <Loader2Icon className="animate-spin" />}
              Apply watermark
            </Button>
          </div>
        </form>
      </div>
    </>
  )
}

export default App
