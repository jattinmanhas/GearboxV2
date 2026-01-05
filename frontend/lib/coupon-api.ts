import { httpClient } from "./apiFunctions/http-client";

export const couponApi = {
  // Coupon Management
  async getCoupons(params?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
    type?: string;
  }): Promise<any> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append("page", params.page.toString());
    if (params?.limit) queryParams.append("limit", params.limit.toString());
    if (params?.search) queryParams.append("search", params.search);
    if (params?.status) queryParams.append("status", params.status);
    if (params?.type) queryParams.append("type", params.type);

    const query = queryParams.toString();
    return httpClient.get<any>(`/coupons${query ? `?${query}` : ""}`);
  },

  async getCoupon(id: string): Promise<any> {
    return httpClient.get<any>(`/coupons/${id}`);
  },

  async createCoupon(couponData: any): Promise<any> {
    return httpClient.post<any>("/coupons", couponData);
  },

  async updateCoupon(id: string, couponData: any): Promise<any> {
    return httpClient.put<any>(`/coupons/${id}`, couponData);
  },

  async deleteCoupon(id: string): Promise<any> {
    return httpClient.delete<any>(`/coupons/${id}`);
  },

  async validateCoupon(couponCode: string): Promise<any> {
    return httpClient.post<any>("/coupons/validate", {
      coupon_code: couponCode,
    });
  },

  async getCouponUsage(params?: {
    page?: number;
    limit?: number;
    coupon_id?: string;
    start_date?: string;
    end_date?: string;
  }): Promise<any> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append("page", params.page.toString());
    if (params?.limit) queryParams.append("limit", params.limit.toString());
    if (params?.coupon_id) queryParams.append("coupon_id", params.coupon_id);
    if (params?.start_date) queryParams.append("start_date", params.start_date);
    if (params?.end_date) queryParams.append("end_date", params.end_date);

    const query = queryParams.toString();
    return httpClient.get<any>(`/coupons/usage${query ? `?${query}` : ""}`);
  },
};

