import { httpClient } from "./http-client";
import type { ApiResponse } from "../types";

export const cartApi = {
    // Cart Management
    async getCartBySession(currency: string = "INR"): Promise<any> {
        const response = await httpClient.get<ApiResponse<any>>(
            `/carts?endpoint=session&currency=${currency}`
        );
        return response.data;
    },

    async getOrCreateCart(currency: string = "INR"): Promise<any> {
        const response = await httpClient.get<ApiResponse<any>>(
            `/carts?endpoint=get-or-create&currency=${currency}`
        );
        return response.data;
    },

    async getCart(cartId: string): Promise<any> {
        const response = await httpClient.get<ApiResponse<any>>(`/carts/${cartId}`);
        return response.data;
    },

    async updateCart(cartId: string, data: any): Promise<any> {
        const response = await httpClient.put<ApiResponse<any>>(`/carts/${cartId}`, data);
        return response.data;
    },

    async deleteCart(cartId: string): Promise<any> {
        const response = await httpClient.delete<ApiResponse<any>>(`/carts/${cartId}`);
        return response.data;
    },

    // Cart Items
    async addItemToCart(cartId: string, itemData: any): Promise<any> {
        const response = await httpClient.post<ApiResponse<any>>(`/carts/${cartId}/items`, itemData);
        return response.data;
    },

    async getCartItems(cartId: string): Promise<any> {
        const response = await httpClient.get<ApiResponse<any>>(`/carts/${cartId}/items`);
        return response.data;
    },

    async updateCartItem(itemId: string, data: any): Promise<any> {
        const response = await httpClient.put<ApiResponse<any>>(`/carts/items/${itemId}`, data);
        return response.data;
    },

    async deleteCartItem(itemId: string): Promise<any> {
        const response = await httpClient.delete<ApiResponse<any>>(`/carts/items/${itemId}`);
        return response.data;
    },

    async clearCartItems(cartId: string): Promise<any> {
        const response = await httpClient.delete<ApiResponse<any>>(`/carts/${cartId}/items`);
        return response.data;
    },

    // Cart Summary
    async getCartSummary(cartId: string): Promise<any> {
        const response = await httpClient.get<ApiResponse<any>>(`/carts/${cartId}/summary`);
        return response.data;
    },

    // Cart Merging
    async mergeCarts(targetCartId: string, sourceCartId: string): Promise<any> {
        const response = await httpClient.post<ApiResponse<any>>(`/carts/${targetCartId}/merge`, {
            source_cart_id: sourceCartId,
        });
        return response.data;
    },

    // Cart Coupons
    async applyCouponToCart(cartId: string, couponCode: string): Promise<any> {
        const response = await httpClient.post<ApiResponse<any>>(`/carts/${cartId}/coupons`, {
            coupon_code: couponCode,
        });
        return response.data;
    },

    async removeCouponFromCart(cartId: string, couponCode: string): Promise<any> {
        const response = await httpClient.delete<ApiResponse<any>>(`/carts/${cartId}/coupons`, {
            coupon_code: couponCode,
        });
        return response.data;
    },

    async getCartCoupons(cartId: string): Promise<any> {
        const response = await httpClient.get<ApiResponse<any>>(`/carts/${cartId}/coupons`);
        return response.data;
    },
};
