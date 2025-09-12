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
  UpdateProductVariantRequest
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
