"use client"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { AlertMessage } from "@/components/ui/alert-message"
import Link from "next/link"
import { useState } from "react"
import { useRouter } from "next/navigation"
import { authApi, ApiError } from "@/lib/apiFunctions"

export function ForgotPasswordForm({
  className,
  ...props
}: React.ComponentProps<"form">) {
  const router = useRouter()
  const [identifier, setIdentifier] = useState("")
  const [identifierType, setIdentifierType] = useState<"email" | "username">("email")
  const [errors, setErrors] = useState<{ identifier?: string }>({})
  const [isLoading, setIsLoading] = useState(false)
  const [submitError, setSubmitError] = useState<string>("")
  const [success, setSuccess] = useState(false)

  const validateIdentifier = (value: string, type: "email" | "username"): boolean => {
    if (!value) {
      setErrors({ identifier: `${type === "email" ? "Email" : "Username"} is required` })
      return false
    }
    if (type === "email") {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(value)) {
        setErrors({ identifier: "Please enter a valid email address" })
        return false
      }
    } else {
      // Username validation: 3-30 characters, alphanumeric + underscore + hyphen
      if (value.length < 3 || value.length > 30) {
        setErrors({ identifier: "Username must be between 3 and 30 characters" })
        return false
      }
      const usernameRegex = /^[a-zA-Z0-9_-]+$/
      if (!usernameRegex.test(value)) {
        setErrors({ identifier: "Username can only contain letters, numbers, underscores, and hyphens" })
        return false
      }
    }
    setErrors({})
    return true
  }

  const handleInputChange = (value: string) => {
    setIdentifier(value)
    if (errors.identifier) {
      setErrors({})
    }
    setSubmitError("")
    // Auto-detect if it's an email or username
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (emailRegex.test(value)) {
      setIdentifierType("email")
    } else {
      setIdentifierType("username")
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitError("")
    setSuccess(false)

    if (!validateIdentifier(identifier, identifierType)) {
      return
    }

    setIsLoading(true)

    try {
      const response = await authApi.forgotPassword(
        identifierType === "email" ? identifier : undefined,
        identifierType === "username" ? identifier : undefined
      )

      if (response.message) {
        // Always show success message (security best practice - don't reveal if email exists)
        setSuccess(true)
      } else {
        setSubmitError("An error occurred. Please try again.")
      }
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(error.message || "An error occurred. Please try again.")
      } else {
        setSubmitError("An unexpected error occurred. Please try again.")
      }
    } finally {
      setIsLoading(false)
    }
  }

  if (success) {
    return (
      <div className={cn("flex flex-col gap-6", className)} {...props}>
        <div className="flex flex-col items-center gap-2 text-center">
          <h1 className="text-2xl font-bold">Check your email</h1>
          <p className="text-muted-foreground text-sm text-balance">
            If an account with that email or username exists, we've sent you a password reset link.
          </p>
        </div>

        <AlertMessage type="success" message="Password reset email sent successfully!" />

        <div className="text-center text-sm">
          <Link href="/login" className="underline underline-offset-4">
            Back to login
          </Link>
        </div>
      </div>
    )
  }

  return (
    <form className={cn("flex flex-col gap-6", className)} onSubmit={handleSubmit} {...props}>
      <div className="flex flex-col items-center gap-2 text-center">
        <h1 className="text-2xl font-bold">Forgot your password?</h1>
        <p className="text-muted-foreground text-sm text-balance">
          Enter your email address or username and we'll send you a link to reset your password.
        </p>
      </div>

      <AlertMessage type="error" message={submitError} />

      <div className="grid gap-6">
        <div className="grid gap-3">
          <Label htmlFor="identifier">Email or Username</Label>
          <Input
            className={cn("inherit dark:bg-neutral-800/50", errors.identifier && "border-red-500")}
            id="identifier"
            type="text"
            placeholder="user@example.com or username"
            value={identifier}
            onChange={(e) => handleInputChange(e.target.value)}
            disabled={isLoading}
            required
          />
          {errors.identifier && (
            <p className="text-sm text-red-600 dark:text-red-400">{errors.identifier}</p>
          )}
          <p className="text-xs text-muted-foreground">
            We'll automatically detect if you entered an email or username
          </p>
        </div>
        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ? "Sending..." : "Send reset link"}
        </Button>
      </div>
      <div className="text-center text-sm">
        Remember your password?{" "}
        <Link href="/login" className="underline underline-offset-4">
          Back to login
        </Link>
      </div>
    </form>
  )
}

