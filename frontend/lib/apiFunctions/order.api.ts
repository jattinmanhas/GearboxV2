import { httpClient } from "./http-client";
import type { ApiResponse } from "../types";

export const orderApi = {
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
