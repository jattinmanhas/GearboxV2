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

// Product and Category types for dashboard
export interface Category {
  id: number
  name: string
  description: string
  slug: string
  parent_id?: number
  is_active: boolean
  sort_order: number
  image_url: string
  meta_title: string
  meta_description: string
  created_at: string
  updated_at: string
}

export interface Product {
  id: number
  name: string
  description: string
  short_description: string
  sku: string
  price: number
  compare_price: number
  cost_price: number
  weight: number
  dimensions: string
  is_active: boolean
  is_digital: boolean
  requires_shipping: boolean
  taxable: boolean
  track_quantity: boolean
  min_quantity: number
  max_quantity: number
  meta_title: string
  meta_description: string
  tags: string
  category_ids?: number[]
  category_names?: string[]
  created_at: string
  updated_at: string
}

export interface ProductVariant {
  id: number
  product_id: number
  name: string
  sku: string
  price: number
  compare_price: number
  cost_price: number
  weight: number
  is_active: boolean
  position: number
  created_at: string
  updated_at: string
  // Inventory information
  quantity: number
  available_quantity: number
  is_in_stock: boolean
}

export interface CreateProductVariantRequest {
  product_id: number
  name: string
  sku: string
  price: number
  compare_price?: number
  cost_price?: number
  weight?: number
  is_active?: boolean
  position?: number
}

export interface UpdateProductVariantRequest {
  name?: string
  sku?: string
  price?: number
  compare_price?: number
  cost_price?: number
  weight?: number
  is_active?: boolean
  position?: number
}

export interface CreateCategoryRequest {
  name: string
  description: string
  slug: string
  parent_id?: number
  is_active: boolean
  sort_order: number
  image_url: string
  meta_title: string
  meta_description: string
}

export interface UpdateCategoryRequest {
  name?: string
  description?: string
  slug?: string
  parent_id?: number
  is_active?: boolean
  sort_order?: number
  image_url?: string
  meta_title?: string
  meta_description?: string
}

export interface CreateProductRequest {
  name: string
  description: string
  short_description?: string
  sku: string
  price: number
  compare_price?: number
  cost_price?: number
  weight?: number
  dimensions?: string
  is_active: boolean
  is_digital: boolean
  requires_shipping: boolean
  taxable: boolean
  track_quantity: boolean
  min_quantity?: number
  max_quantity?: number
  meta_title?: string
  meta_description?: string
  tags?: string
  category_ids?: number[]
}

export interface UpdateProductRequest {
  name?: string
  description?: string
  short_description?: string
  sku?: string
  price?: number
  compare_price?: number
  cost_price?: number
  weight?: number
  dimensions?: string
  is_active?: boolean
  is_digital?: boolean
  requires_shipping?: boolean
  taxable?: boolean
  track_quantity?: boolean
  min_quantity?: number
  max_quantity?: number
  meta_title?: string
  meta_description?: string
  tags?: string
  category_ids?: number[]
}
export interface ListProductsResponse {
  data: {
    products: Product[]
    total: number
    page: number
    limit: number
    total_pages: number
  }
}

export interface ProductFilters {
  search?: string
  category_id?: number
  is_active?: boolean
  is_digital?: boolean
  min_price?: number
  max_price?: number
  in_stock?: boolean
  tags?: string[]
  sort_by?: string
  sort_order?: 'asc' | 'desc'
  page?: number
  limit?: number
}

export interface CategoryFilters {
  search?: string
  parent_id?: number
  is_active?: boolean
  page?: number
  limit?: number
}

export interface ListCategoriesResponse {
  data: {
    categories: Category[]
    total: number
    page: number
    limit: number
    total_pages: number
  }
}

// User Profile and Address types
export interface UserProfile {
  id: number
  username: string
  email: string
  first_name: string
  middle_name: string
  last_name: string
  phone_number: string
  date_of_birth: string
  avatar: string
  created_at: string
  updated_at: string
}

export interface Address {
  id: number
  user_id: number
  type: 'billing' | 'shipping'
  first_name: string
  last_name: string
  company?: string
  address_line_1: string
  address_line_2?: string
  city: string
  state: string
  postal_code: string
  country: string
  phone_number?: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface CreateAddressRequest {
  type: 'billing' | 'shipping'
  first_name: string
  last_name: string
  company?: string
  address_line_1: string
  address_line_2?: string
  city: string
  state: string
  postal_code: string
  country: string
  phone_number?: string
  is_default?: boolean
}

export interface UpdateAddressRequest {
  type?: 'billing' | 'shipping'
  first_name?: string
  last_name?: string
  company?: string
  address_line_1?: string
  address_line_2?: string
  city?: string
  state?: string
  postal_code?: string
  country?: string
  phone_number?: string
  is_default?: boolean
}

export interface UpdateProfileRequest {
  first_name?: string
  middle_name?: string
  last_name?: string
  phone_number?: string
  date_of_birth?: string
  avatar?: string
}

// Inventory Management Types
export interface Inventory {
  id: number
  product_id: number
  product_variant_id?: number
  quantity: number
  reserved_quantity: number
  available_quantity: number
  min_stock_level: number
  max_stock_level?: number
  reorder_point: number
  last_restocked: string
  created_at: string
  updated_at: string
  // Additional fields for display
  product_name: string
  product_sku: string
  variant_name?: string
  variant_sku?: string
}

export interface CreateInventoryRequest {
  product_id: number
  product_variant_id?: number
  quantity: number
  min_stock_level: number
  max_stock_level?: number
  reorder_point: number
}

export interface UpdateInventoryRequest {
  quantity?: number
  min_stock_level?: number
  max_stock_level?: number
  reorder_point?: number
}

export interface StockMovement {
  id: number
  product_id: number
  product_variant_id?: number
  movement_type: 'in' | 'out' | 'adjustment' | 'transfer'
  quantity: number
  previous_quantity: number
  new_quantity: number
  reference?: string
  reference_type?: string
  reason?: string
  notes?: string
  created_by?: number
  created_at: string
  // Additional fields for display
  product_name?: string
  product_sku?: string
  variant_name?: string
  variant_sku?: string
}

export interface StockMovementRequest {
  product_id: number
  product_variant_id?: number
  movement_type: 'in' | 'out' | 'adjustment' | 'transfer'
  quantity: number
  reference?: string
  reference_type?: string
  reason?: string
  notes?: string
  created_by?: number
}

export interface StockReservation {
  id: number
  product_id: number
  product_variant_id?: number
  order_id: number
  quantity: number
  expires_at: string
  created_at: string
  // Additional fields for display
  product_name?: string
  product_sku?: string
  variant_name?: string
  variant_sku?: string
}

export interface ReserveStockRequest {
  product_id: number
  product_variant_id?: number
  order_id: number
  quantity: number
  expires_at: string
}

export interface ReleaseStockRequest {
  reservation_id?: number
  order_id?: number
}

export interface InventoryAlert {
  id: number
  product_id: number
  product_variant_id?: number
  alert_type: 'low_stock' | 'out_of_stock' | 'overstock'
  current_quantity: number
  threshold_quantity: number
  is_resolved: boolean
  resolved_at?: string
  created_at: string
  // Additional fields for display
  product_name?: string
  product_sku?: string
  variant_name?: string
  variant_sku?: string
}

export interface InventorySummary {
  total_products: number
  total_variants: number
  total_quantity: number
  total_reserved: number
  total_available: number
  low_stock_items: number
  out_of_stock_items: number
  total_value: number
  average_stock_level: number
}

export interface ListInventoryRequest {
  product_id?: number
  product_variant_id?: number
  low_stock?: boolean
  out_of_stock?: boolean
  page?: number
  limit?: number
}

export interface ListInventoryResponse {
  inventory: Inventory[]
  total: number
  page: number
  limit: number
  total_pages: number
}

export interface ListStockMovementsRequest {
  product_id?: number
  product_variant_id?: number
  movement_type?: 'in' | 'out' | 'adjustment' | 'transfer'
  start_date?: string
  end_date?: string
  page?: number
  limit?: number
}

export interface ListStockMovementsResponse {
  movements: StockMovement[]
  total: number
  page: number
  limit: number
  total_pages: number
}

export interface BulkStockUpdateItem {
  product_id: number
  product_variant_id?: number
  quantity: number
  movement_type: 'in' | 'out' | 'adjustment'
  reason: string
  notes?: string
}

export interface BulkStockUpdateRequest {
  updates: BulkStockUpdateItem[]
}

export interface BulkStockUpdateResponse {
  updated_items: number
  failed_items: Array<{
    product_id: number
    product_variant_id?: number
    error: string
  }>
  success: boolean
}

// User Management Types
export interface User {
  id: number
  username: string
  email: string
  first_name: string
  middle_name?: string
  last_name: string
  phone_number?: string
  date_of_birth?: string
  avatar?: string
  gender?: string
  is_active: boolean
  is_deleted: boolean
  created_at: string
  updated_at: string
  last_login?: string
  role_id?: number
  role?: string
}

export interface UpdateUserRequest {
  username?: string
  email?: string
  first_name?: string
  middle_name?: string
  last_name?: string
  phone_number?: string
  date_of_birth?: string
  avatar?: string
  is_active?: boolean
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
  confirm_password: string
}

export interface ListUsersResponse {
  data: {
    users: User[]
    total: number
    page: number
    limit: number
    total_pages: number
  }
}

export interface UserFilters {
  search?: string
  is_active?: boolean
  role_id?: number
  page?: number
  limit?: number
}

// Role Management Types (simplified - roles are fixed)
export interface Role {
  id: number
  name: string
  description: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface AssignRoleRequest {
  user_id: number
  role_id: number
}

export interface CreateUserRequest {
  username: string
  email: string
  password: string
  first_name: string
  middle_name?: string
  last_name: string
  phone_number?: string
  date_of_birth?: string
  avatar?: string
  is_active?: boolean
}

export interface CreateRoleRequest {
  name: string
  description: string
  is_active?: boolean
}

export interface UpdateRoleRequest {
  name?: string
  description?: string
  is_active?: boolean
}

export interface ListRolesResponse {
  data: {
    roles: Role[]
    total: number
    page: number
    limit: number
    total_pages: number
  }
}

export interface RoleFilters {
  search?: string
  is_active?: boolean
  page?: number
  limit?: number
}

export interface RemoveRoleRequest {
  user_id: number
  role_id: number
}

export interface CheckPermissionRequest {
  permission: string
  resource?: string
}

export interface CheckPermissionResponse {
  has_permission: boolean
  required_role: string
  user_role: string
}

// Fixed roles in the system
export const FIXED_ROLES = {
  'admin': 'Administrator',
  'editor': 'Editor',
  'moderator': 'Moderator',
  'viewer': 'Viewer',
  'user': 'User',
} as const

export type RoleName = keyof typeof FIXED_ROLES
