import { 
  RegisterRequest, 
  LoginRequest, 
  ApiResponse, 
  Category, 
  Product, 
  CreateCategoryRequest, 
  UpdateCategoryRequest, 
  CreateProductRequest, 
  UpdateProductRequest, 
  ListProductsResponse, 
  ListCategoriesResponse,
  UserProfile,
  Address,
  CreateAddressRequest,
  UpdateAddressRequest,
  UpdateProfileRequest
} from './types'

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

export const productApi = {
  // Categories
  async getCategories(page: number = 1, limit: number = 10, search?: string): Promise<ListCategoriesResponse> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
      ...(search && { search }),
    })
    
    const response = await fetch(`${API_BASE_URL}/products/categories?${params}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    
    return handleResponse<ListCategoriesResponse>(response)
  },

  async getCategory(id: number): Promise<Category> {
    const response = await fetch(`${API_BASE_URL}/products/categories/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    
    return handleResponse<Category>(response)
  },

  async createCategory(categoryData: CreateCategoryRequest): Promise<Category> {
    const response = await fetch(`${API_BASE_URL}/products/categories`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(categoryData),
    })
    
    return handleResponse<Category>(response)
  },

  async updateCategory(id: number, categoryData: UpdateCategoryRequest): Promise<Category> {
    const response = await fetch(`${API_BASE_URL}/products/categories/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(categoryData),
    })
    
    return handleResponse<Category>(response)
  },

  async deleteCategory(id: number): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/products/categories/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    
    return handleResponse<ApiResponse>(response)
  },

  // Products
  async getProducts(page: number = 1, limit: number = 10, search?: string, categoryId?: number): Promise<ListProductsResponse> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
      ...(search && { search }),
      ...(categoryId && { category_id: categoryId.toString() }),
    })
    
    const response = await fetch(`${API_BASE_URL}/products?${params}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    
    return handleResponse<ListProductsResponse>(response)
  },

  async getProduct(id: number): Promise<Product> {
    const response = await fetch(`${API_BASE_URL}/products/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    
    return handleResponse<Product>(response)
  },

  async createProduct(productData: CreateProductRequest): Promise<Product> {
    const response = await fetch(`${API_BASE_URL}/products`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(productData),
    })
    
    return handleResponse<Product>(response)
  },

  async updateProduct(id: number, productData: UpdateProductRequest): Promise<Product> {
    const response = await fetch(`${API_BASE_URL}/products/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(productData),
    })
    
    return handleResponse<Product>(response)
  },

  async deleteProduct(id: number): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/products/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    
    return handleResponse<ApiResponse>(response)
  },
}

export const profileApi = {
  // User Profile
  async getProfile(): Promise<UserProfile> {
    const response = await fetch(`${API_BASE_URL}/auth/profile`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<UserProfile>(response)
  },

  async updateProfile(profileData: UpdateProfileRequest): Promise<UserProfile> {
    const response = await fetch(`${API_BASE_URL}/auth/profile`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(profileData),
    })
    
    return handleResponse<UserProfile>(response)
  },

  // Addresses
  async getAddresses(): Promise<Address[]> {
    const response = await fetch(`${API_BASE_URL}/auth/addresses`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<Address[]>(response)
  },

  async getAddress(id: number): Promise<Address> {
    const response = await fetch(`${API_BASE_URL}/auth/addresses/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<Address>(response)
  },

  async createAddress(addressData: CreateAddressRequest): Promise<Address> {
    const response = await fetch(`${API_BASE_URL}/auth/addresses`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(addressData),
    })
    
    return handleResponse<Address>(response)
  },

  async updateAddress(id: number, addressData: UpdateAddressRequest): Promise<Address> {
    const response = await fetch(`${API_BASE_URL}/auth/addresses/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(addressData),
    })
    
    return handleResponse<Address>(response)
  },

  async deleteAddress(id: number): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/addresses/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<ApiResponse>(response)
  },
}
