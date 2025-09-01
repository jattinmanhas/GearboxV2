import { RegisterRequest, LoginRequest, ApiResponse } from './types'

const API_BASE_URL = '/api/v1'

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public errors?: string[]
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  const data = await response.json()
  
  if (!response.ok) {
    // Extract detailed error message from backend response
    let errorMessage = data.message || 'Request failed'
    
    // Check for detailed error information in the response
    if (data.error?.detail) {
      errorMessage = data.error.detail
    } else if (data.error?.message) {
      errorMessage = data.error.message
    } else if (data.errors && Array.isArray(data.errors)) {
      errorMessage = data.errors.join(', ')
    }
    
    throw new ApiError(
      errorMessage,
      response.status,
      data.errors
    )
  }
  
  return data
}

export const authApi = {
  async register(userData: RegisterRequest): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(userData),
    })
    
    return handleResponse<ApiResponse>(response)
  },

  async login(credentials: LoginRequest): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    })
    
    return handleResponse<ApiResponse>(response)
  },

  async logout(): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/logout`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // Important: include cookies for authentication
    })
    
    return handleResponse<ApiResponse>(response)
  },

  async refreshToken(): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    
    return handleResponse<ApiResponse>(response)
  },
}
