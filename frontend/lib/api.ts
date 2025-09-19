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
  CheckPermissionResponse
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
  const data = await response.json()
  
  if (!response.ok) {
    // Handle 401 Unauthorized responses by clearing user state
    if (response.status === 401) {
      // Import user store dynamically to avoid circular dependencies
      const { useUserStore } = await import('./stores/user-store')
      const userStore = useUserStore.getState()
      
      // Clear user state if user is currently authenticated
      if (userStore.isAuthenticated) {
        console.log('Token expired or invalid, clearing user state')
        userStore.clearUser()
        
        // Optional: Redirect to login page if not already there
        if (typeof window !== 'undefined' && !window.location.pathname.startsWith('/login')) {
          console.log('Redirecting to login page due to expired token')
          window.location.href = '/login'
        }
      }
    }
    
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
