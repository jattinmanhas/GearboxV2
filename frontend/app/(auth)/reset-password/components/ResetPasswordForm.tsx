"use client"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { AlertMessage } from "@/components/ui/alert-message"
import Link from "next/link"
import { useState, useEffect } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { authApi, ApiError } from "@/lib/apiFunctions"
import { PasswordStrength } from "@/app/(auth)/register/components/password-strength"

export function ResetPasswordForm({
  className,
  ...props
}: React.ComponentProps<"form">) {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [token, setToken] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [errors, setErrors] = useState<{ password?: string; confirmPassword?: string }>({})
  const [isLoading, setIsLoading] = useState(false)
  const [submitError, setSubmitError] = useState<string>("")
  const [success, setSuccess] = useState(false)

  useEffect(() => {
    // Get token from URL query parameter
    const tokenParam = searchParams.get("token")
    if (tokenParam) {
      setToken(tokenParam)
    } else {
      setSubmitError("Invalid or missing reset token. Please request a new password reset link.")
    }
  }, [searchParams])

  const validatePassword = (password: string): boolean => {
    if (!password) {
      setErrors(prev => ({ ...prev, password: "Password is required" }))
      return false
    }
    if (password.length < 6) {
      setErrors(prev => ({ ...prev, password: "Password must be at least 6 characters" }))
      return false
    }
    // Check for uppercase, lowercase, and number
    const hasUpper = /[A-Z]/.test(password)
    const hasLower = /[a-z]/.test(password)
    const hasNumber = /[0-9]/.test(password)

    if (!hasUpper || !hasLower || !hasNumber) {
      setErrors(prev => ({
        ...prev,
        password: "Password must contain at least one uppercase letter, one lowercase letter, and one number"
      }))
      return false
    }
    setErrors(prev => ({ ...prev, password: undefined }))
    return true
  }

  const validateConfirmPassword = (confirmPassword: string, password: string): boolean => {
    if (!confirmPassword) {
      setErrors(prev => ({ ...prev, confirmPassword: "Please confirm your password" }))
      return false
    }
    if (confirmPassword !== password) {
      setErrors(prev => ({ ...prev, confirmPassword: "Passwords do not match" }))
      return false
    }
    setErrors(prev => ({ ...prev, confirmPassword: undefined }))
    return true
  }

  const handleInputChange = (field: "password" | "confirmPassword", value: string) => {
    if (field === "password") {
      setPassword(value)
      if (errors.password) {
        setErrors(prev => ({ ...prev, password: undefined }))
      }
      // Re-validate confirm password if it's already filled
      if (confirmPassword) {
        validateConfirmPassword(confirmPassword, value)
      }
    } else {
      setConfirmPassword(value)
      if (errors.confirmPassword) {
        setErrors(prev => ({ ...prev, confirmPassword: undefined }))
      }
    }
    setSubmitError("")
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitError("")

    if (!token) {
      setSubmitError("Invalid or missing reset token. Please request a new password reset link.")
      return
    }

    const isPasswordValid = validatePassword(password)
    const isConfirmPasswordValid = validateConfirmPassword(confirmPassword, password)

    if (!isPasswordValid || !isConfirmPasswordValid) {
      return
    }

    setIsLoading(true)

    try {
      const response = await authApi.resetPassword(token, password)

      if (response.message) {
        setSuccess(true)
        // Redirect to login after 3 seconds
        setTimeout(() => {
          router.push("/login")
        }, 3000)
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
      <form className={cn("flex flex-col gap-6", className)} onSubmit={(e) => e.preventDefault()} {...props}>
        <div className="flex flex-col items-center gap-2 text-center">
          <h1 className="text-2xl font-bold">Password reset successful!</h1>
          <p className="text-muted-foreground text-sm text-balance">
            Your password has been reset successfully. Redirecting to login...
          </p>
        </div>

        <AlertMessage type="success" message="Password reset successfully!" />

        <div className="text-center text-sm">
          <Link href="/login" className="underline underline-offset-4">
            Go to login
          </Link>
        </div>
      </form>
    )
  }

  return (
    <form className={cn("flex flex-col gap-6", className)} onSubmit={handleSubmit} {...props}>
      <div className="flex flex-col items-center gap-2 text-center">
        <h1 className="text-2xl font-bold">Reset your password</h1>
        <p className="text-muted-foreground text-sm text-balance">
          Enter your new password below.
        </p>
      </div>

      <AlertMessage type="error" message={submitError} />

      <div className="grid gap-6">
        <div className="grid gap-3">
          <Label htmlFor="password">New Password</Label>
          <Input
            className={cn("inherit dark:bg-neutral-800/50", errors.password && "border-red-500")}
            id="password"
            type="password"
            value={password}
            onChange={(e) => handleInputChange("password", e.target.value)}
            disabled={isLoading}
            required
          />
          <PasswordStrength password={password} />
          {errors.password && (
            <p className="text-sm text-red-600 dark:text-red-400">{errors.password}</p>
          )}
        </div>
        <div className="grid gap-3">
          <Label htmlFor="confirmPassword">Confirm Password</Label>
          <Input
            className={cn("inherit dark:bg-neutral-800/50", errors.confirmPassword && "border-red-500")}
            id="confirmPassword"
            type="password"
            value={confirmPassword}
            onChange={(e) => handleInputChange("confirmPassword", e.target.value)}
            disabled={isLoading}
            required
          />
          {errors.confirmPassword && (
            <p className="text-sm text-red-600 dark:text-red-400">{errors.confirmPassword}</p>
          )}
        </div>
        <Button type="submit" className="w-full" disabled={isLoading || !token}>
          {isLoading ? "Resetting..." : "Reset password"}
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

