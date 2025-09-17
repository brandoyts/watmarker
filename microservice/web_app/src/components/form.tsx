import api from "@/lib/api"
import { useState, type ChangeEvent, type FormEvent } from "react"
import {
    Dropzone,
    DropzoneContent,
    DropzoneEmptyState,
} from "./ui/shadcn-io/dropzone"
import { Input } from "./ui/input"
import { Button } from "./ui/button"
import { Loader2Icon } from "lucide-react"
import { ImageZoom } from "./ui/shadcn-io/image-zoom"
import { MAXIMUM_IMAGE_SIZE, MINIMUM_IMAGE_SIZE } from "@/constants"
import { cn } from "@/lib/utils"
import { isAxiosError } from "axios"

interface FormErrors {
    file?: string
    watermark?: string
    general?: string
}

function Form() {
    const [files, setFiles] = useState<File[] | undefined>()
    const [watermark, setWatermark] = useState("")
    const [isLoading, setIsLoading] = useState(false)
    const [formErrors, setFormErrors] = useState<FormErrors>({})
    const [generatedImageUrl, setGeneratedImageUrl] = useState("")

    const handleDrop = (file: File[]) => {
        if (formErrors?.file) {
            setFormErrors((prev) => ({ ...prev, ["file"]: "" }))
        }
        setFiles(file)
    }

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        setFormErrors({})


        if (!files || files.length === 0) {
            setFormErrors({ file: "file is required." })
            return
        }

        if (watermark.length === 0) {
            setFormErrors({ watermark: "Watermark is required." })
            return
        }

        setIsLoading(true)

        const formData = new FormData()
        formData.append("file_data", files[0])
        formData.append("watermark_text", watermark)

        try {
            const { data } = await api.post("/watermark", formData)

            setGeneratedImageUrl(data.image_url)
        } catch (error: unknown) {
            if (isAxiosError(error)) {
                if (error.response?.status === 429) {
                    setFormErrors((prev) => ({
                        ...prev,
                        general: "Too many requests. Please try again later.",
                    }));
                } else {
                    setFormErrors((prev) => ({
                        ...prev,
                        general: "An unexpected error occurred.",
                    }));
                }
            } else {
                setFormErrors((prev) => ({
                    ...prev,
                    general: "An unexpected error occurred.",
                }));
            }
            return
        } finally {
            setIsLoading(false)
        }

        setFormErrors({})
    }

    const handleWatermarkChange = (e: ChangeEvent<HTMLInputElement>) => {
        if (formErrors?.watermark) {
            setFormErrors((prev) => ({ ...prev, ["watermark"]: "" }))
        }
        setWatermark(e.target.value)
    }
    return (
        <>
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
                    className={cn(
                        "w-full max-w-md sm:max-w-lg bg-white hover:bg-white border-gray-300 hover:border-blue-600 rounded-xl",
                        formErrors?.file && "border-red-500"
                    )}
                >
                    <DropzoneEmptyState />
                    <DropzoneContent />
                </Dropzone>
                <div className="flex flex-col sm:flex-row gap-3 w-full max-w-lg">
                    <Input
                        value={watermark}
                        onChange={handleWatermarkChange}
                        disabled={isLoading}
                        placeholder="watermark"
                        className={cn(
                            "bg-white flex-1 rounded-xl border border-gray-300 px-3 py-2 sm:py-3",
                            formErrors?.watermark && "border-red-500"
                        )}
                    />

                    <Button
                        disabled={isLoading}
                        className="cursor-pointer text-white px-4 sm:px-6 py-2 sm:py-3 rounded-xl bg-blue-500 hover:bg-blue-600 transition"
                    >
                        {isLoading && <Loader2Icon className="animate-spin" />}
                        Apply watermark
                    </Button>
                </div>
            </form>
            {formErrors?.general && (
                <p className="text-red-500">{formErrors.general}</p>
            )}

            {/* render generated image */}
            {generatedImageUrl && (
                <ImageZoom className="text-center">
                    <img
                        alt="Placeholder image"
                        className="h-auto w-55"
                        src={generatedImageUrl}
                    />
                    <Button variant="link">
                        <a href={generatedImageUrl}>Download</a>
                    </Button>
                </ImageZoom>
            )}
        </>
    )
}

export default Form
