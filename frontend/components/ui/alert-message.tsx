import { Alert, AlertDescription } from "@/components/ui/alert"
import { CheckCircle, AlertCircle, Info } from "lucide-react"

interface AlertMessageProps {
  type: "success" | "error" | "info"
  message: string
  className?: string
}

export function AlertMessage({ type, message, className }: AlertMessageProps) {
  if (!message) return null

  const variants = {
    success: {
      variant: "default" as const,
      icon: CheckCircle,
      containerClass:
        "border-green-300 bg-green-50 text-green-800 dark:border-green-800 dark:bg-green-900/30 dark:text-green-200",
      iconClass: "text-green-600 dark:text-green-400",
    },
    error: {
      variant: "destructive" as const,
      icon: AlertCircle,
      containerClass:
        "border-red-300 bg-red-50 text-red-800 dark:border-red-800 dark:bg-red-900/30 dark:text-red-200",
      iconClass: "text-red-600 dark:text-red-400",
    },
    info: {
      variant: "default" as const,
      icon: Info,
      containerClass:
        "border-blue-300 bg-blue-50 text-blue-800 dark:border-blue-800 dark:bg-blue-900/30 dark:text-blue-200",
      iconClass: "text-blue-600 dark:text-blue-400",
    },
  }

  const { variant, icon: Icon, containerClass, iconClass } = variants[type]

  return (
    <Alert
      variant={variant}
      className={`flex items-center gap-2 rounded-xl border p-3 shadow-sm ${containerClass} ${className}`}
    >
      <Icon className={`h-5 w-5 shrink-0 ${iconClass}`} />
      <AlertDescription className="text-sm mx-2 mt-1.5">{message}</AlertDescription>
    </Alert>
  )
}
