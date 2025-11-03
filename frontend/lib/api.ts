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
  ProductFilters,
  CategoryFilters,
  UserProfile,
  Address,
  CreateAddressRequest,
  UpdateAddressRequest,
  UpdateProfileRequest,
  ProductVariant,
  CreateProductVariantRequest,
  UpdateProductVariantRequest,
  User,
  UpdateUserRequest,
  ChangePasswordRequest,
  ListUsersResponse,
  UserFilters,
  Role,
  AssignRoleRequest,
  RemoveRoleRequest,
  CheckPermissionRequest,
  CheckPermissionResponse,
  OAuthProvider,
  OAuthInitiateResponse,
  LinkedProvidersResponse
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

export async function handleResponse<T>(response: Response): Promise<T> {
  // Try to parse JSON, but handle cases where response might not be JSON
  let data: any = {}
  const contentType = response.headers.get('content-type')
  const isJson = contentType?.includes('application/json')
  
  try {
    // Read the response text once (can only be read once)
    const text = await response.text()
    
    if (text.trim()) {
      if (isJson) {
        try {
          data = JSON.parse(text)
        } catch (jsonError) {
          // If JSON parsing fails even though content-type says JSON, use the text as message
          console.error('Failed to parse JSON response:', jsonError)
          data = { message: text || `HTTP ${response.status}: ${response.statusText}` }
        }
      } else {
        // Not JSON content type, use text as message
        data = { message: text || `HTTP ${response.status}: ${response.statusText}` }
      }
    } else {
      // Empty response
      data = { message: `HTTP ${response.status}: ${response.statusText}` }
    }
  } catch (readError) {
    console.error('Failed to read response:', readError)
    // If reading fails, create a basic error object
    data = { 
      message: `HTTP ${response.status}: ${response.statusText}`,
      error: 'Failed to read response'
    }
  }
  
  if (!response.ok) {
    // Handle 401 Unauthorized responses more intelligently
    if (response.status === 401) {
      // Import user store dynamically to avoid circular dependencies
      const { useUserStore } = await import('./stores/user-store')
      const userStore = useUserStore.getState()
      
      // Only clear user state if we have a user and the error indicates actual auth failure
      if (userStore.isAuthenticated && userStore.user) {
        // Check if this is a genuine authentication failure vs other 401 scenarios
        const isAuthFailure = await isGenuineAuthFailure(response, data)
        
        if (isAuthFailure) {
          console.log('Genuine authentication failure detected, clearing user state')
          userStore.clearUser()
          
          // Only redirect to login if not already there and not on a public page
          if (typeof window !== 'undefined' && 
              !window.location.pathname.startsWith('/login') && 
              !window.location.pathname.startsWith('/register') &&
              !window.location.pathname.startsWith('/about') &&
              !window.location.pathname.startsWith('/contact') &&
              !window.location.pathname.startsWith('/blog') &&
              !window.location.pathname.startsWith('/shop')) {
            console.log('Redirecting to login page due to authentication failure')
            window.location.href = '/login'
          }
        } else {
          console.log('401 error not related to authentication, keeping user state')
        }
      }
    }
    
    // Extract detailed error message from backend response
    let errorMessage = data?.message || 'Request failed'
    
    // Check for detailed error information in the response
    if (data?.error?.detail) {
      errorMessage = data.error.detail
    } else if (data?.error?.message) {
      errorMessage = data.error.message
    } else if (data?.errors && Array.isArray(data.errors)) {
      errorMessage = data.errors.join(', ')
    } else if (data?.details) {
      // Handle nested error details (from Next.js API routes)
      try {
        const details = typeof data.details === 'string' ? JSON.parse(data.details) : data.details
        if (details?.error?.detail) {
          errorMessage = details.error.detail
        } else if (details?.error?.message) {
          errorMessage = details.error.message
        } else if (details?.message) {
          errorMessage = details.message
        }
      } catch {
        // If parsing fails, use the error field if it's a string
        if (typeof data?.error === 'string') {
          errorMessage = data.error
        }
      }
    } else if (typeof data?.error === 'string') {
      errorMessage = data.error
    }
    
    // Ensure errorMessage is always a valid string
    if (!errorMessage || typeof errorMessage !== 'string') {
      errorMessage = `HTTP ${response.status}: ${response.statusText}`
    }
    
    throw new ApiError(
      errorMessage,
      response.status,
      data?.errors
    )
  }
  
  // Extract data from wrapped response structure
  // Backend returns: { timestamp, status, success, message, data, error }
  // However, product/cart services use this format too, and frontend expects response.data.products
  // So we need to be smart about extraction:
  // - If data.data exists and is a plain object (not an array), check if it looks like a nested response
  // - Only extract if it's clearly a single object (like profile data), not a structured response
  // - For structured responses (like { products: [...] }), keep the wrapper so frontend can access .data.products
  
  if (
    data && 
    typeof data === 'object' && 
    'data' in data && 
    data.data !== undefined &&
    ('timestamp' in data || 'status' in data || 'success' in data)
  ) {
    // Only extract if data.data is NOT a structured response with nested properties
    // that the frontend expects (like products, categories, etc.)
    const extractedData = data.data
    
    // If extractedData is an array, keep the wrapper (frontend expects response.data)
    if (Array.isArray(extractedData)) {
      return data as T
    }
    
    // If extractedData has properties like 'products', 'categories', 'orders', etc.,
    // it's a structured response - keep the wrapper
    // Also check for pagination properties that indicate a list response
    // Also check for cart-specific properties (session_id, currency) that indicate cart responses
    // Also check for auth-related properties (user, message together) that indicate login/register responses
    if (
      typeof extractedData === 'object' && 
      extractedData !== null &&
      ('products' in extractedData || 'categories' in extractedData || 'orders' in extractedData ||
       'items' in extractedData || 'cart' in extractedData || 'carts' in extractedData ||
       'addresses' in extractedData || 'coupons' in extractedData ||
       'total' in extractedData || 'total_pages' in extractedData || 'page' in extractedData ||
       'subtotal' in extractedData || 'total_amount' in extractedData || 'discount_amount' in extractedData ||
       'session_id' in extractedData || 'currency' in extractedData || 'cart_id' in extractedData ||
       'user' in extractedData)
    ) {
      return data as T
    }
    
    // Otherwise, extract the data (for simple objects like profile)
    return extractedData as T
  }
  
  // If no wrapper, return the entire response (for backward compatibility)
  return data as T
}

/**
 * Determines if a 401 error represents a genuine authentication failure
 * vs other scenarios like insufficient permissions or missing tokens
 */
async function isGenuineAuthFailure(response: Response, data: any): Promise<boolean> {
  try {
    // Check if the error message indicates token expiration or invalid token
    const errorMessage = data.message || data.error?.message || ''
    const lowerErrorMessage = errorMessage.toLowerCase()
    
    // Common patterns that indicate genuine auth failure
    // Based on actual backend error messages from auth service and JWT library
    const authFailurePatterns = [
      // From auth service login failures
      'invalid credentials',
      
      // From JWT validation errors (golang-jwt/jwt v5)
      'invalid token',
      'failed to parse token',
      'token is expired',
      'token is malformed',
      'signature is invalid',
      'unexpected signing method',
      'token is not valid yet',
      'token used before issued',
      'token is unverifiable',
      'key is of invalid type',
      'claims are invalid',
      'token is blacklisted',
      'token is not active',
      'token is not valid',
      'token validation failed',
      'token parse error',
      'token verification failed',
      
      // From middleware auth failures
      'refresh token required',
      'invalid refresh token',
      
      // Generic auth failure patterns
      'unauthorized',
      'authentication failed',
      'session expired',
      'jwt expired',
      'jwt invalid',
      'token not found'
    ]
    
    // If error message contains auth failure patterns, it's likely genuine
    if (authFailurePatterns.some(pattern => lowerErrorMessage.includes(pattern))) {
      return true
    }
    
    // If no specific error message, try to validate the current token
    const { useUserStore } = await import('./stores/user-store')
    const userStore = useUserStore.getState()
    
    if (userStore.isAuthenticated && userStore.user) {
      // Try to validate the token by making a simple auth request
      try {
        const validationResponse = await fetch('/api/v1/auth/profile', {
          method: 'GET',
          credentials: 'include'
        })
        
        // If profile request also fails with 401, it's likely a genuine auth failure
        if (validationResponse.status === 401) {
          return true
        }
      } catch (error) {
        // If validation request fails entirely, assume it's not an auth issue
        console.log('Token validation request failed, assuming not auth failure:', error)
        return false
      }
    }
    
    // Default to not clearing user state for ambiguous 401 errors
    return false
  } catch (error) {
    console.error('Error determining auth failure type:', error)
    // On error, be conservative and don't clear user state
    return false
  }
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
      credentials: 'include', // Important: include cookies for authentication
    })
    
    return handleResponse<ApiResponse>(response)
  },
}

export const productApi = {
  // Categories
  async getCategories(filters: CategoryFilters = {}): Promise<ListCategoriesResponse> {
    const params = new URLSearchParams()
    
    // Pagination
    if (filters.page) params.append('page', filters.page.toString())
    if (filters.limit) params.append('limit', filters.limit.toString())
    
    // Search and filters
    if (filters.search) params.append('search', filters.search)
    if (filters.parent_id) params.append('parent_id', filters.parent_id.toString())
    if (filters.is_active !== undefined) params.append('is_active', filters.is_active.toString())
    
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
  async getProducts(filters: ProductFilters = {}): Promise<ListProductsResponse> {
    const params = new URLSearchParams()
    
    // Pagination
    if (filters.page) params.append('page', filters.page.toString())
    if (filters.limit) params.append('limit', filters.limit.toString())
    
    // Search and filters
    if (filters.search) params.append('search', filters.search)
    if (filters.category_id) params.append('category_id', filters.category_id.toString())
    if (filters.is_active !== undefined) params.append('is_active', filters.is_active.toString())
    if (filters.is_digital !== undefined) params.append('is_digital', filters.is_digital.toString())
    if (filters.min_price !== undefined) params.append('min_price', filters.min_price.toString())
    if (filters.max_price !== undefined) params.append('max_price', filters.max_price.toString())
    if (filters.in_stock !== undefined) params.append('in_stock', filters.in_stock.toString())
    if (filters.on_sale !== undefined) params.append('on_sale', filters.on_sale.toString())
    if (filters.tags && filters.tags.length > 0) params.append('tags', filters.tags.join(','))
    if (filters.sort_by) params.append('sort_by', filters.sort_by)
    if (filters.sort_order) params.append('sort_order', filters.sort_order)
    
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
    
    const data = await handleResponse<{data: Product}>(response)
    return data.data
  },

  async createProduct(productData: CreateProductRequest): Promise<Product> {
    console.log('[API] Creating product with data:', productData)
    console.log('[API] Images being sent:', productData.images)
    
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

  // Product Variants
  async getProductVariants(productId: number): Promise<ProductVariant[]> {
    const response = await fetch(`${API_BASE_URL}/products/${productId}/variants`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    
    const result = await handleResponse<{data: ProductVariant[]}>(response)
    return result.data
  },

  async getProductVariantsWithInventory(productId: number): Promise<ProductVariant[]> {
    const response = await fetch(`${API_BASE_URL}/products/${productId}/variants-with-inventory`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    
    const result = await handleResponse<{data: ProductVariant[]}>(response)
    console.log('API response for variants:', result)
    console.log('Extracted data:', result.data)
    return result.data
  },

  async getProductVariant(variantId: number): Promise<ProductVariant> {
    const response = await fetch(`${API_BASE_URL}/products/variants/${variantId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    
    const result = await handleResponse<{data: ProductVariant}>(response)
    return result.data
  },

  async createProductVariant(productId: number, variantData: CreateProductVariantRequest): Promise<ProductVariant> {
    const response = await fetch(`${API_BASE_URL}/products/${productId}/variants`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(variantData),
    })
    
    const result = await handleResponse<{data: ProductVariant}>(response)
    return result.data
  },

  async updateProductVariant(variantId: number, variantData: UpdateProductVariantRequest): Promise<ProductVariant> {
    const response = await fetch(`${API_BASE_URL}/products/variants/${variantId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(variantData),
    })
    
    const result = await handleResponse<{data: ProductVariant}>(response)
    return result.data
  },

  async deleteProductVariant(variantId: number): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/products/variants/${variantId}`, {
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
    
    const responseData = await handleResponse<any>(response)
    
    // Extract addresses array from response (could be wrapped in data field or direct array)
    const addresses = Array.isArray(responseData) 
      ? responseData 
      : (responseData?.data ? (Array.isArray(responseData.data) ? responseData.data : []) : [])
    
    // Map backend address structure to frontend Address type
    return addresses.map((addr: any) => ({
      id: addr.id,
      user_id: addr.user_id,
      type: (addr.address_type || addr.type) as 'billing' | 'shipping',
      first_name: addr.first_name,
      last_name: addr.last_name,
      company: addr.company,
      address_line_1: addr.address_line_1,
      address_line_2: addr.address_line_2,
      city: addr.city,
      state: addr.state,
      postal_code: addr.postal_code,
      country: addr.country,
      phone_number: addr.phone || addr.phone_number,
      is_default: addr.is_default,
      created_at: addr.created_at,
      updated_at: addr.updated_at,
    }))
  },

  async getAddress(id: number): Promise<Address> {
    const response = await fetch(`${API_BASE_URL}/auth/addresses/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    const addr = await handleResponse<any>(response)
    
    // Map backend address structure to frontend Address type
    return {
      id: addr.id,
      user_id: addr.user_id,
      type: (addr.address_type || addr.type) as 'billing' | 'shipping',
      first_name: addr.first_name,
      last_name: addr.last_name,
      company: addr.company,
      address_line_1: addr.address_line_1,
      address_line_2: addr.address_line_2,
      city: addr.city,
      state: addr.state,
      postal_code: addr.postal_code,
      country: addr.country,
      phone_number: addr.phone || addr.phone_number,
      is_default: addr.is_default,
      created_at: addr.created_at,
      updated_at: addr.updated_at,
    }
  },

  async createAddress(addressData: CreateAddressRequest): Promise<Address> {
    // Map frontend address structure to backend format
    const backendData = {
      address_type: addressData.type,
      first_name: addressData.first_name,
      last_name: addressData.last_name,
      company: addressData.company,
      address_line_1: addressData.address_line_1,
      address_line_2: addressData.address_line_2,
      city: addressData.city,
      state: addressData.state,
      country: addressData.country,
      postal_code: addressData.postal_code,
      phone: addressData.phone_number,
      is_default: addressData.is_default,
    }
    
    const response = await fetch(`${API_BASE_URL}/auth/addresses`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(backendData),
    })
    
    const addr = await handleResponse<any>(response)
    
    // Map backend address structure to frontend Address type
    return {
      id: addr.id,
      user_id: addr.user_id,
      type: (addr.address_type || addr.type) as 'billing' | 'shipping',
      first_name: addr.first_name,
      last_name: addr.last_name,
      company: addr.company,
      address_line_1: addr.address_line_1,
      address_line_2: addr.address_line_2,
      city: addr.city,
      state: addr.state,
      postal_code: addr.postal_code,
      country: addr.country,
      phone_number: addr.phone || addr.phone_number,
      is_default: addr.is_default,
      created_at: addr.created_at,
      updated_at: addr.updated_at,
    }
  },

  async updateAddress(id: number, addressData: UpdateAddressRequest): Promise<Address> {
    // Map frontend address structure to backend format
    const backendData: any = {}
    if (addressData.type !== undefined) backendData.address_type = addressData.type
    if (addressData.first_name !== undefined) backendData.first_name = addressData.first_name
    if (addressData.last_name !== undefined) backendData.last_name = addressData.last_name
    if (addressData.company !== undefined) backendData.company = addressData.company
    if (addressData.address_line_1 !== undefined) backendData.address_line_1 = addressData.address_line_1
    if (addressData.address_line_2 !== undefined) backendData.address_line_2 = addressData.address_line_2
    if (addressData.city !== undefined) backendData.city = addressData.city
    if (addressData.state !== undefined) backendData.state = addressData.state
    if (addressData.country !== undefined) backendData.country = addressData.country
    if (addressData.postal_code !== undefined) backendData.postal_code = addressData.postal_code
    if (addressData.phone_number !== undefined) backendData.phone = addressData.phone_number
    if (addressData.is_default !== undefined) backendData.is_default = addressData.is_default
    
    const response = await fetch(`${API_BASE_URL}/auth/addresses/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(backendData),
    })
    
    const addr = await handleResponse<any>(response)
    
    // Map backend address structure to frontend Address type
    return {
      id: addr.id,
      user_id: addr.user_id,
      type: (addr.address_type || addr.type) as 'billing' | 'shipping',
      first_name: addr.first_name,
      last_name: addr.last_name,
      company: addr.company,
      address_line_1: addr.address_line_1,
      address_line_2: addr.address_line_2,
      city: addr.city,
      state: addr.state,
      postal_code: addr.postal_code,
      country: addr.country,
      phone_number: addr.phone || addr.phone_number,
      is_default: addr.is_default,
      created_at: addr.created_at,
      updated_at: addr.updated_at,
    }
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

// Cart API
export const cartApi = {
  // Cart Management
  async getCartBySession(currency: string = 'INR'): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts?endpoint=session&currency=${currency}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async getOrCreateCart(currency: string = 'INR'): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts?endpoint=get-or-create&currency=${currency}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async getCart(cartId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/${cartId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async updateCart(cartId: string, data: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/${cartId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(data),
    })
    
    return handleResponse<any>(response)
  },

  async deleteCart(cartId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/${cartId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  // Cart Items
  async addItemToCart(cartId: string, itemData: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/${cartId}/items`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(itemData),
    })
    
    return handleResponse<any>(response)
  },

  async getCartItems(cartId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/${cartId}/items`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async updateCartItem(itemId: string, data: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/items/${itemId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(data),
    })
    
    return handleResponse<any>(response)
  },

  async deleteCartItem(itemId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/items/${itemId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async clearCartItems(cartId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/${cartId}/items`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  // Cart Summary
  async getCartSummary(cartId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/${cartId}/summary`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  // Cart Merging
  async mergeCarts(targetCartId: string, sourceCartId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/${targetCartId}/merge`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ source_cart_id: sourceCartId }),
    })
    
    return handleResponse<any>(response)
  },

  // Cart Coupons
  async applyCouponToCart(cartId: string, couponCode: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/${cartId}/coupons`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ coupon_code: couponCode }),
    })
    
    return handleResponse<any>(response)
  },

  async removeCouponFromCart(cartId: string, couponCode: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/${cartId}/coupons`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ coupon_code: couponCode }),
    })
    
    return handleResponse<any>(response)
  },

  async getCartCoupons(cartId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/carts/${cartId}/coupons`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },
}

// Wishlist API
export const wishlistApi = {
  // Wishlist Management
  async getWishlists(): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/wishlists`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async createWishlist(data: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/wishlists`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(data),
    })
    
    return handleResponse<any>(response)
  },

  async getWishlist(wishlistId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/wishlists/${wishlistId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async updateWishlist(wishlistId: string, data: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/wishlists/${wishlistId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(data),
    })
    
    return handleResponse<any>(response)
  },

  async deleteWishlist(wishlistId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/wishlists/${wishlistId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  // Wishlist Items
  async getWishlistItems(wishlistId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/wishlists/${wishlistId}/items`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },


  async addItemToWishlist(wishlistId: string, itemData: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/wishlists/${wishlistId}/items`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(itemData),
    })
    
    return handleResponse<any>(response)
  },

  async removeItemFromWishlist(itemId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/wishlists/items/${itemId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async deleteWishlistItem(itemId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/wishlists/items/${itemId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async moveItemToCart(itemId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/wishlists/items/${itemId}/move-to-cart`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },
}

// Export coupon API
export { couponApi } from './coupon-api'

// User Management API
export const userApi = {
  // Get all users (Admin only)
  async getUsers(filters: UserFilters = {}): Promise<ListUsersResponse> {
    const params = new URLSearchParams()
    
    // Pagination
    if (filters.page) params.append('page', filters.page.toString())
    if (filters.limit) params.append('limit', filters.limit.toString())
    
    // Search and filters
    if (filters.search) params.append('search', filters.search)
    if (filters.is_active !== undefined) params.append('is_active', filters.is_active.toString())
    if (filters.role_id) params.append('role_id', filters.role_id.toString())
    
    console.log('API getUsers called with filters:', filters)
    console.log('URL params:', params.toString())
    
    const response = await fetch(`${API_BASE_URL}/auth/users?${params}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    const data = await handleResponse<{data: any}>(response)
    
    // Handle both response formats:
    // 1. New format: { data: { users: [], total: 0, page: 1, limit: 10, total_pages: 0 } }
    // 2. Old format: { data: [user1, user2, ...] }
    let users = []
    let total = 0
    let page = 1
    let limit = 10
    let total_pages = 0
    
    if (Array.isArray(data.data)) {
      // Old format - data is an array of users
      users = data.data
      total = data.data.length
      page = filters.page || 1
      limit = filters.limit || 10
      total_pages = Math.ceil(total / limit)
    } else if (data.data && typeof data.data === 'object') {
      // New format - data is an object with pagination
      users = data.data.users || []
      total = data.data.total || 0
      page = data.data.page || 1
      limit = data.data.limit || 10
      total_pages = data.data.total_pages || 0
    }
    
    // Transform the response to match the expected structure
    return {
      data: {
        users,
        total,
        page,
        limit,
        total_pages,
      }
    }
  },

  // Get user by ID
  async getUser(id: number): Promise<User> {
    const response = await fetch(`${API_BASE_URL}/auth/user/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<User>(response)
  },

  // Update user
  async updateUser(id: number, userData: UpdateUserRequest): Promise<User> {
    const response = await fetch(`${API_BASE_URL}/auth/user/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(userData),
    })
    
    return handleResponse<User>(response)
  },

  // Delete user
  async deleteUser(id: number): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/user/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<ApiResponse>(response)
  },

  // Change password
  async changePassword(id: number, passwordData: ChangePasswordRequest): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/user/${id}/change-password`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(passwordData),
    })
    
    return handleResponse<ApiResponse>(response)
  },

  // Logout all devices
  async logoutAll(id: number): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/logout-all`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ user_id: id }),
    })
    
    return handleResponse<ApiResponse>(response)
  },

  // Cleanup expired tokens (Admin only)
  async cleanupExpiredTokens(): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/cleanup-expired-tokens`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<ApiResponse>(response)
  },
}

// Role Management API
export const roleApi = {
  // Get all roles
  async getRoles(): Promise<Role[]> {
    const response = await fetch(`${API_BASE_URL}/auth/roles`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    const data = await handleResponse<{data: Role[]}>(response)
    return data.data || []
  },

  // Get current user's role
  async getMyRole(): Promise<Role> {
    const response = await fetch(`${API_BASE_URL}/auth/roles/my-role`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<Role>(response)
  },

  // Get user's role
  async getUserRole(userId: number): Promise<Role> {
    const response = await fetch(`${API_BASE_URL}/auth/roles/user?user_id=${userId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<Role>(response)
  },

  // Assign role to user (Editor+ only)
  async assignRole(assignmentData: AssignRoleRequest): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/roles/assign`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(assignmentData),
    })
    
    return handleResponse<ApiResponse>(response)
  },

  // Remove role from user (Editor+ only)
  async removeRole(removalData: RemoveRoleRequest): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/roles/remove`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(removalData),
    })
    
    return handleResponse<ApiResponse>(response)
  },

  // Check permission
  async checkPermission(permissionData: CheckPermissionRequest): Promise<CheckPermissionResponse> {
    const params = new URLSearchParams()
    params.append('permission', permissionData.permission)
    if (permissionData.resource) {
      params.append('resource', permissionData.resource)
    }
    
    const response = await fetch(`${API_BASE_URL}/auth/roles/check-permission?${params}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<CheckPermissionResponse>(response)
  },
}

// OAuth API
export const oauthApi = {
  // Initiate OAuth flow
  async initiateOAuth(provider: OAuthProvider): Promise<OAuthInitiateResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/oauth/${provider}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    const data = await handleResponse<{data: OAuthInitiateResponse}>(response)
    return data.data
  },

  // Get linked OAuth providers
  async getLinkedProviders(): Promise<LinkedProvidersResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/oauth/providers`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    const data = await handleResponse<{data: LinkedProvidersResponse}>(response)
    return data.data
  },

  // Unlink OAuth provider
  async unlinkProvider(provider: OAuthProvider): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/oauth/unlink/${provider}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<ApiResponse>(response)
  },

  // Link OAuth provider (used when user is already logged in)
  async linkProvider(provider: OAuthProvider, code: string): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/oauth/link/${provider}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ provider, code }),
    })
    
    return handleResponse<ApiResponse>(response)
  },
}

// Order API
export const orderApi = {
  async createOrderFromCart(cartId: number, orderData: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/orders/from-cart?cart_id=${cartId}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(orderData),
    })
    
    return handleResponse<any>(response)
  },

  async getOrder(orderId: number): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/orders/${orderId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },
}

// Payment API
export const paymentApi = {
  async getPaymentMethods(): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/payment-methods`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async getPaymentGateways(): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/payment-gateways`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async createPayment(paymentData: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/payments`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(paymentData),
    })
    
    return handleResponse<any>(response)
  },

  async processPayment(paymentId: number, paymentData: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/payments/process`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({
        payment_id: paymentId,
        ...paymentData,
      }),
    })
    
    return handleResponse<any>(response)
  },
}
