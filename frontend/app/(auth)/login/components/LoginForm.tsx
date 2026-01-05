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
import { LoginFormData, FormErrors, formDataToLoginRequest, loginFormSchema } from "@/lib/types"
import { useUserStore } from "@/lib/stores/user-store"
import { OAuthButtonsGroup } from "@/components/auth/oauth-button"

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"form">) {
  const router = useRouter()
  const { login, setLoading, setError } = useUserStore()
  const [formData, setFormData] = useState<LoginFormData>({
    username: "",
    password: "",
  })
  const [errors, setErrors] = useState<FormErrors>({})
  const [isLoading, setIsLoading] = useState(false)
  const [submitError, setSubmitError] = useState<string>("")

  const validateForm = (): boolean => {
    try {
      loginFormSchema.parse(formData)
      setErrors({})
      return true
    } catch (error) {
      if (error instanceof Error) {
        const zodError = error as any
        const newErrors: FormErrors = {}

        if (zodError.errors) {
          zodError.errors.forEach((err: any) => {
            newErrors[err.path[0]] = err.message
          })
        }

        setErrors(newErrors)
      }
      return false
    }
  }

  const handleInputChange = (field: keyof LoginFormData, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }))
    // Clear error when user starts typing
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: "" }))
    }
    setSubmitError("")
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      return
    }

    setIsLoading(true)
    setSubmitError("")

    try {
      const response = await authApi.login(formDataToLoginRequest(formData))

      // Debug: Log the response to see the structure
      console.log("Login response:", response)

      // The API returns user data in response.data.user
      const userData = response?.data?.user

      if (userData) {
        const user = {
          id: userData.id || 1,
          username: userData.username || formData.username,
          email: userData.email || '',
          firstName: userData.firstName || '',
          middleName: userData.middleName || '',
          lastName: userData.lastName || '',
          avatar: userData.avatar || '',
          role: userData.role || 'user',
          createdAt: userData.createdAt || new Date().toISOString(),
          updatedAt: userData.updatedAt || new Date().toISOString(),
        }

        // Debug: Log the user data being stored
        console.log("User data to store:", user)

        // Store user data in Zustand store and handle cart merging
        await login(user)

        // Redirect to dashboard or home
        router.push("/")
      } else {
        console.error("No user data in response:", response)
        throw new Error("No user data received from server")
      }
    } catch (error) {
      console.error("Login error:", error)
      if (error instanceof ApiError) {
        setSubmitError(error.message)
        setError(error.message)
      } else {
        const errorMessage = "Login failed. Please try again."
        setSubmitError(errorMessage)
        setError(errorMessage)
      }
    } finally {
      setIsLoading(false)
      setLoading(false)
    }
  }

  return (
    <form className={cn("flex flex-col gap-6", className)} onSubmit={handleSubmit} {...props}>
      <div className="flex flex-col items-center gap-2 text-center">
        <h1 className="text-2xl font-bold">Login to your account</h1>
        <p className="text-muted-foreground text-sm text-balance">
          Enter your details below to login to your account
        </p>
      </div>

      <AlertMessage type="error" message={submitError} />

      <div className="grid gap-6">
        <div className="grid gap-3">
          <Label htmlFor="username">Username</Label>
          <Input
            className={cn("inherit dark:bg-neutral-800/50", errors.username && "border-red-500")}
            id="username"
            type="text"
            placeholder="test_user"
            value={formData.username}
            onChange={(e) => handleInputChange("username", e.target.value)}
            disabled={isLoading}
          />
          {errors.username && (
            <p className="text-sm text-red-600 dark:text-red-400">{errors.username}</p>
          )}
        </div>
        <div className="grid gap-3">
          <div className="flex items-center">
            <Label htmlFor="password">Password</Label>
            <Link
              href="/forgot-password"
              className="ml-auto text-sm underline-offset-4 hover:underline"
            >
              Forgot your password?
            </Link>
          </div>
          <Input
            className={cn("inherit dark:bg-neutral-800/50", errors.password && "border-red-500")}
            id="password"
            type="password"
            value={formData.password}
            onChange={(e) => handleInputChange("password", e.target.value)}
            disabled={isLoading}
          />
          {errors.password && (
            <p className="text-sm text-red-600 dark:text-red-400">{errors.password}</p>
          )}
        </div>
        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ? "Logging in..." : "Login"}
        </Button>
        <div className="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t">
          <span className="bg-background text-muted-foreground relative z-10 px-2">
            Or continue with
          </span>
        </div>
        <OAuthButtonsGroup />
      </div>
      <div className="text-center text-sm">
        Don&apos;t have an account?{" "}
        <Link href="/register" className="underline underline-offset-4">
          Sign up
        </Link>
      </div>
    </form>
  )
}
