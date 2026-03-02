import { httpClient } from "./http-client";
import type { ApiResponse } from "../types";

export const orderApi = {
    async getOrders(params: Record<string, string>): Promise<ApiResponse<any>> {
        const queryString = new URLSearchParams(params).toString();
        return httpClient.get<ApiResponse<any>>(`/orders?${queryString}`);
    },

    async createOrderFromCart(cartId: number, orderData: any): Promise<any> {
        return httpClient.post<any>(
            `/orders/from-cart?cart_id=${cartId}`,
            orderData
        );
    },

    async getOrder(orderId: number): Promise<any> {
        const response = await httpClient.get<ApiResponse<any>>(`/orders/${orderId}`);
        return response.data;
    },

    async createOrderPayment(orderId: number, paymentData: any): Promise<any> {
        return httpClient.post<any>(`/orders/${orderId}/payments`, paymentData);
    },

    async processOrderPayment(orderId: number, processData: any): Promise<any> {
        return httpClient.post<any>(`/orders/${orderId}/payments/process`, processData);
    },

    async updateOrderStatus(orderId: number, status: string, notes?: string): Promise<any> {
        return httpClient.put<any>(`/orders/${orderId}/status`, {
            status,
            notes,
        });
    },

    async updatePaymentStatus(orderId: number, paymentStatus: string): Promise<any> {
        // This might be handled differently depending on the backend
        // For now, mirroring the payment status update if applicable
        return httpClient.put<any>(`/orders/${orderId}`, {
            payment_status: paymentStatus,
        });
    },
};
