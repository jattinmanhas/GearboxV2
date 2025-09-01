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
import { authApi, ApiError } from "@/lib/api"
import { RegisterFormData, FormErrors, formDataToRegisterRequest, registerFormSchema } from "@/lib/types"

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
        <Button variant="outline" className="w-full" type="button" disabled={isLoading}>
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <path
              d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"
              fill="currentColor"
            />
          </svg>
          Sign up with GitHub
        </Button>
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
