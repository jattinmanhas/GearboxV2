"use client"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { AlertMessage } from "@/components/ui/alert-message"
import Link from "next/link"
import { useState } from "react"
import { useRouter } from "next/navigation"
import { authApi, ApiError } from "@/lib/api"
import { LoginFormData, FormErrors, formDataToLoginRequest, loginFormSchema } from "@/lib/types"
import { useUserStore } from "@/lib/stores/user-store"

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
      if (response.data && response.data.user) {
        const user = response.data.user
        const userData = {
          id: user.id || 1,
          username: user.username || formData.username,
          email: user.email || '',
          firstName: user.firstName || '',
          middleName: user.middleName || '',
          lastName: user.lastName || '',
          avatar: user.avatar || '',
          role: user.role || 'user',
          createdAt: user.createdAt || new Date().toISOString(),
          updatedAt: user.updatedAt || new Date().toISOString(),
        }
        
        // Debug: Log the user data being stored
        console.log("User data to store:", userData)
        
        // Store user data in Zustand store
        login(userData)
        
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
            <a
              href="#"
              className="ml-auto text-sm underline-offset-4 hover:underline"
            >
              Forgot your password?
            </a>
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
        <Button variant="outline" className="w-full" type="button" disabled={isLoading}>
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <path
              d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"
              fill="currentColor"
            />
          </svg>
          Login with GitHub
        </Button>
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
