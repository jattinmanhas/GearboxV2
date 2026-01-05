import { httpClient } from "./http-client";

export const wishlistApi = {
    // Wishlist Management
    async getWishlists(): Promise<any> {
        return httpClient.get<any>("/wishlists");
    },

    async createWishlist(data: any): Promise<any> {
        return httpClient.post<any>("/wishlists", data);
    },

    async getWishlist(wishlistId: string): Promise<any> {
        return httpClient.get<any>(`/wishlists/${wishlistId}`);
    },

    async updateWishlist(wishlistId: string, data: any): Promise<any> {
        return httpClient.put<any>(`/wishlists/${wishlistId}`, data);
    },

    async deleteWishlist(wishlistId: string): Promise<any> {
        return httpClient.delete<any>(`/wishlists/${wishlistId}`);
    },

    // Wishlist Items
    async getWishlistItems(wishlistId: string): Promise<any> {
        return httpClient.get<any>(`/wishlists/${wishlistId}/items`);
    },

    async addItemToWishlist(wishlistId: string, itemData: any): Promise<any> {
        return httpClient.post<any>(`/wishlists/${wishlistId}/items`, itemData);
    },

    async removeItemFromWishlist(itemId: string): Promise<any> {
        return httpClient.delete<any>(`/wishlists/items/${itemId}`);
    },

    async deleteWishlistItem(itemId: string): Promise<any> {
        return httpClient.delete<any>(`/wishlists/items/${itemId}`);
    },

    async moveItemToCart(itemId: string): Promise<any> {
        return httpClient.post<any>(`/wishlists/items/${itemId}/move-to-cart`);
    },
};
