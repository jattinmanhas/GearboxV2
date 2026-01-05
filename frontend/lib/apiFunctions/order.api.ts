import { httpClient } from "./http-client";

export const orderApi = {
    async createOrderFromCart(cartId: number, orderData: any): Promise<any> {
        return httpClient.post<any>(
            `/orders/from-cart?cart_id=${cartId}`,
            orderData
        );
    },

    async getOrder(orderId: number): Promise<any> {
        return httpClient.get<any>(`/orders/${orderId}`);
    },
};
