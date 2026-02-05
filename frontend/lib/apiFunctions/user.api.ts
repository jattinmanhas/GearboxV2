import { httpClient } from "./http-client";
import type {
    User,
    UpdateUserRequest,
    ChangePasswordRequest,
    ListUsersResponse,
    UserFilters,
    ApiResponse,
} from "../types";

export const userApi = {
    // Get all users (Admin only)
    async getUsers(filters: UserFilters = {}): Promise<ListUsersResponse> {
        const params = new URLSearchParams();

        if (filters.page) params.append("page", filters.page.toString());
        if (filters.limit) params.append("limit", filters.limit.toString());
        if (filters.search) params.append("search", filters.search);
        if (filters.is_active !== undefined)
            params.append("is_active", filters.is_active.toString());
        if (filters.role_id) params.append("role_id", filters.role_id.toString());

        console.log("API getUsers called with filters:", filters);
        console.log("URL params:", params.toString());

        const query = params.toString();
        const data = await httpClient.get<{ data: any }>(
            `/auth/users${query ? `?${query}` : ""}`
        );

        // Handle both response formats:
        // 1. New format: { data: { users: [], total: 0, page: 1, limit: 10, total_pages: 0 } }
        // 2. Old format: { data: [user1, user2, ...] }
        let users = [];
        let total = 0;
        let page = 1;
        let limit = 10;
        let total_pages = 0;

        if (Array.isArray(data.data)) {
            // Old format - data is an array of users
            users = data.data;
            total = data.data.length;
            page = filters.page || 1;
            limit = filters.limit || 10;
            total_pages = Math.ceil(total / limit);
        } else if (data.data && typeof data.data === "object") {
            // New format - data is an object with pagination
            users = data.data.users || [];
            total = data.data.total || 0;
            page = data.data.page || 1;
            limit = data.data.limit || 10;
            total_pages = data.data.total_pages || 0;
        }

        return {
            data: {
                users,
                total,
                page,
                limit,
                total_pages,
            },
        };
    },

    // Get user by ID
    async getUser(id: number): Promise<User> {
        const response = await httpClient.get<ApiResponse<User>>(`/auth/user/${id}`);
        return response.data as User;
    },

    // Update user
    async updateUser(id: number, userData: UpdateUserRequest): Promise<User> {
        const response = await httpClient.put<ApiResponse<User>>(`/auth/user/${id}`, userData);
        return response.data as User;
    },

    // Delete user
    async deleteUser(id: number): Promise<ApiResponse> {
        return httpClient.delete<ApiResponse>(`/auth/user/${id}`);
    },

    // Change password
    async changePassword(
        id: number,
        passwordData: ChangePasswordRequest
    ): Promise<ApiResponse> {
        return httpClient.post<ApiResponse>(
            `/auth/user/${id}/change-password`,
            passwordData
        );
    },

    // Logout all devices
    async logoutAll(id: number): Promise<ApiResponse> {
        return httpClient.post<ApiResponse>("/auth/logout-all", { user_id: id });
    },

    // Cleanup expired tokens (Admin only)
    async cleanupExpiredTokens(): Promise<ApiResponse> {
        return httpClient.post<ApiResponse>("/auth/cleanup-expired-tokens");
    },
};
