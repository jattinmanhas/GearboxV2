import { handleResponse } from './api'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000/api/v1'

export const couponApi = {
  // Coupon Management
  async getCoupons(params?: {
    page?: number
    limit?: number
    search?: string
    status?: string
    type?: string
  }): Promise<any> {
    const queryParams = new URLSearchParams()
    if (params?.page) queryParams.append('page', params.page.toString())
    if (params?.limit) queryParams.append('limit', params.limit.toString())
    if (params?.search) queryParams.append('search', params.search)
    if (params?.status) queryParams.append('status', params.status)
    if (params?.type) queryParams.append('type', params.type)
    
    const response = await fetch(`${API_BASE_URL}/coupons?${queryParams.toString()}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async getCoupon(id: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/coupons/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async createCoupon(couponData: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/coupons`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(couponData),
    })
    
    return handleResponse<any>(response)
  },

  async updateCoupon(id: string, couponData: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/coupons/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(couponData),
    })
    
    return handleResponse<any>(response)
  },

  async deleteCoupon(id: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/coupons/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },

  async validateCoupon(couponCode: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/coupons/validate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ coupon_code: couponCode }),
    })
    
    return handleResponse<any>(response)
  },

  async getCouponUsage(params?: {
    page?: number
    limit?: number
    coupon_id?: string
    start_date?: string
    end_date?: string
  }): Promise<any> {
    const queryParams = new URLSearchParams()
    if (params?.page) queryParams.append('page', params.page.toString())
    if (params?.limit) queryParams.append('limit', params.limit.toString())
    if (params?.coupon_id) queryParams.append('coupon_id', params.coupon_id)
    if (params?.start_date) queryParams.append('start_date', params.start_date)
    if (params?.end_date) queryParams.append('end_date', params.end_date)
    
    const response = await fetch(`${API_BASE_URL}/coupons/usage?${queryParams.toString()}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    return handleResponse<any>(response)
  },
}
