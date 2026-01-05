"use client"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { AlertMessage } from "@/components/ui/alert-message"
import { PasswordStrength } from "@/app/(auth)/register/components/password-strength"
import Link from "next/link"
import { useState } from "react"
import { useRouter } from "next/navigation"
import { authApi, ApiError } from "@/lib/apiFunctions"
import { RegisterFormData, FormErrors, formDataToRegisterRequest, registerFormSchema } from "@/lib/types"
import { OAuthButtonsGroup } from "@/components/auth/oauth-button"

export function RegisterForm({
  className,
  ...props
}: React.ComponentProps<"form">) {
  const router = useRouter()
  const [formData, setFormData] = useState<RegisterFormData>({
    firstName: "",
    middleName: "",
    lastName: "",
    username: "",
    email: "",
    password: "",
  })
  const [errors, setErrors] = useState<FormErrors>({})
  const [isLoading, setIsLoading] = useState(false)
  const [submitError, setSubmitError] = useState<string>("")

  const validateForm = (): boolean => {
    try {
      registerFormSchema.parse(formData)
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

  const handleInputChange = (field: keyof RegisterFormData, value: string) => {
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
      await authApi.register(formDataToRegisterRequest(formData))

      // Registration successful - redirect to login
      router.push("/login?message=Registration successful! Please log in.")
    } catch (error) {
      console.error("Registration error:", error)
      if (error instanceof ApiError) {
        setSubmitError(error.message)
      } else {
        setSubmitError("Registration failed. Please try again.")
      }
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <form className={cn("flex flex-col gap-6", className)} onSubmit={handleSubmit} {...props}>
      <div className="flex flex-col items-center gap-2 text-center">
        <h1 className="text-2xl font-bold">Create your account</h1>
        <p className="text-muted-foreground text-sm text-balance">
          Enter your details below to create your account
        </p>
      </div>

      <AlertMessage type="error" message={submitError} />

      <div className="grid gap-6">
        {/* Name fields at the top with more space */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="grid gap-2">
            <Label htmlFor="firstName">First Name</Label>
            <Input
              className={cn("inherit dark:bg-neutral-800/50", errors.firstName && "border-red-500")}
              id="firstName"
              type="text"
              placeholder="John"
              value={formData.firstName}
              onChange={(e) => handleInputChange("firstName", e.target.value)}
              disabled={isLoading}
            />
            {errors.firstName && (
              <p className="text-sm text-red-600 dark:text-red-400">{errors.firstName}</p>
            )}
          </div>
          <div className="grid gap-2">
            <Label htmlFor="middleName">Middle Name</Label>
            <Input
              className="inherit dark:bg-neutral-800/50"
              id="middleName"
              type="text"
              placeholder="M"
              value={formData.middleName}
              onChange={(e) => handleInputChange("middleName", e.target.value)}
              disabled={isLoading}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="lastName">Last Name</Label>
            <Input
              className={cn("inherit dark:bg-neutral-800/50", errors.lastName && "border-red-500")}
              id="lastName"
              type="text"
              placeholder="Doe"
              value={formData.lastName}
              onChange={(e) => handleInputChange("lastName", e.target.value)}
              disabled={isLoading}
            />
            {errors.lastName && (
              <p className="text-sm text-red-600 dark:text-red-400">{errors.lastName}</p>
            )}
          </div>
        </div>

        <div className="grid gap-4">
          <div className="grid gap-2">
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
          <div className="grid gap-2">
            <Label htmlFor="email">Email</Label>
            <Input
              className={cn("inherit dark:bg-neutral-800/50", errors.email && "border-red-500")}
              id="email"
              type="email"
              placeholder="user@example.com"
              value={formData.email}
              onChange={(e) => handleInputChange("email", e.target.value)}
              disabled={isLoading}
            />
            {errors.email && (
              <p className="text-sm text-red-600 dark:text-red-400">{errors.email}</p>
            )}
          </div>
          <div className="grid gap-2">
            <Label htmlFor="password">Password</Label>
            <Input
              className={cn("inherit dark:bg-neutral-800/50", errors.password && "border-red-500")}
              id="password"
              type="password"
              value={formData.password}
              onChange={(e) => handleInputChange("password", e.target.value)}
              disabled={isLoading}
            />
            <PasswordStrength password={formData.password} />
            {errors.password && (
              <p className="text-sm text-red-600 dark:text-red-400">{errors.password}</p>
            )}
          </div>
        </div>

        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ? "Creating Account..." : "Create Account"}
        </Button>
        <div className="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t">
          <span className="bg-background text-muted-foreground relative z-10 px-2">
            Or continue with
          </span>
        </div>
        <OAuthButtonsGroup />
      </div>
      <div className="text-center text-sm">
        Already have an account?{" "}
        <Link href="/login" className="underline underline-offset-4">
          Sign in
        </Link>
      </div>
    </form>
  )
}
