import { z } from 'zod'

// Zod validation schemas
export const registerFormSchema = z.object({
  firstName: z.string().min(1, "First name is required").max(50, "First name must be less than 50 characters"),
  middleName: z.string().max(50, "Middle name must be less than 50 characters").optional(),
  lastName: z.string().min(1, "Last name is required").max(50, "Last name must be less than 50 characters"),
  username: z.string().min(3, "Username must be at least 3 characters").max(50, "Username must be less than 50 characters"),
  email: z.string().email("Please enter a valid email address"),
  password: z.string()
    .min(6, "Password must be at least 6 characters")
    .regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/, "Password must contain uppercase, lowercase, and number"),
})

export const loginFormSchema = z.object({
  username: z.string().min(1, "Username is required"),
  password: z.string().min(1, "Password is required"),
})

// Form data structure (camelCase for frontend)
export type RegisterFormData = z.infer<typeof registerFormSchema>
export type LoginFormData = z.infer<typeof loginFormSchema>

// API request structure (snake_case for backend)
export interface RegisterRequest {
  username: string
  password: string
  email: string
  first_name: string
  middle_name: string
  last_name: string
}

export interface LoginRequest {
  username: string
  password: string
}

// API response types
export interface ApiResponse<T = any> {
  success: boolean
  message: string
  data?: T
  errors?: string[]
}

// Form validation types
export interface FormErrors {
  [key: string]: string
}

// Utility function to convert form data to API request
export function formDataToRegisterRequest(formData: RegisterFormData): RegisterRequest {
  return {
    username: formData.username.trim(),
    password: formData.password,
    email: formData.email.trim(),
    first_name: formData.firstName.trim(),
    middle_name: formData.middleName?.trim() || "",
    last_name: formData.lastName.trim(),
  }
}

export function formDataToLoginRequest(formData: LoginFormData): LoginRequest {
  return {
    username: formData.username.trim(),
    password: formData.password,
  }
}
