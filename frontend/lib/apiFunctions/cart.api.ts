import { httpClient } from "./http-client";

export const cartApi = {
    // Cart Management
    async getCartBySession(currency: string = "INR"): Promise<any> {
        return httpClient.get<any>(
            `/carts?endpoint=session&currency=${currency}`
        );
    },

    async getOrCreateCart(currency: string = "INR"): Promise<any> {
        return httpClient.get<any>(
            `/carts?endpoint=get-or-create&currency=${currency}`
        );
    },

    async getCart(cartId: string): Promise<any> {
        return httpClient.get<any>(`/carts/${cartId}`);
    },

    async updateCart(cartId: string, data: any): Promise<any> {
        return httpClient.put<any>(`/carts/${cartId}`, data);
    },

    async deleteCart(cartId: string): Promise<any> {
        return httpClient.delete<any>(`/carts/${cartId}`);
    },

    // Cart Items
    async addItemToCart(cartId: string, itemData: any): Promise<any> {
        return httpClient.post<any>(`/carts/${cartId}/items`, itemData);
    },

    async getCartItems(cartId: string): Promise<any> {
        return httpClient.get<any>(`/carts/${cartId}/items`);
    },

    async updateCartItem(itemId: string, data: any): Promise<any> {
        return httpClient.put<any>(`/carts/items/${itemId}`, data);
    },

    async deleteCartItem(itemId: string): Promise<any> {
        return httpClient.delete<any>(`/carts/items/${itemId}`);
    },

    async clearCartItems(cartId: string): Promise<any> {
        return httpClient.delete<any>(`/carts/${cartId}/items`);
    },

    // Cart Summary
    async getCartSummary(cartId: string): Promise<any> {
        return httpClient.get<any>(`/carts/${cartId}/summary`);
    },

    // Cart Merging
    async mergeCarts(targetCartId: string, sourceCartId: string): Promise<any> {
        return httpClient.post<any>(`/carts/${targetCartId}/merge`, {
            source_cart_id: sourceCartId,
        });
    },

    // Cart Coupons
    async applyCouponToCart(cartId: string, couponCode: string): Promise<any> {
        return httpClient.post<any>(`/carts/${cartId}/coupons`, {
            coupon_code: couponCode,
        });
    },

    async removeCouponFromCart(cartId: string, couponCode: string): Promise<any> {
        return httpClient.delete<any>(`/carts/${cartId}/coupons`, {
            coupon_code: couponCode,
        });
    },

    async getCartCoupons(cartId: string): Promise<any> {
        return httpClient.get<any>(`/carts/${cartId}/coupons`);
    },
};
