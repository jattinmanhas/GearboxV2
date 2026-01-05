import { httpClient } from "./http-client";
import type { RegisterRequest, LoginRequest, ApiResponse } from "../types";

export const authApi = {
    async register(userData: RegisterRequest): Promise<ApiResponse> {
        return httpClient.post<ApiResponse>("/auth/register", userData);
    },

    async login(credentials: LoginRequest): Promise<ApiResponse> {
        return httpClient.post<ApiResponse>("/auth/login", credentials);
    },

    async logout(): Promise<ApiResponse> {
        return httpClient.post<ApiResponse>("/auth/logout");
    },

    async refreshToken(): Promise<ApiResponse> {
        return httpClient.post<ApiResponse>("/auth/refresh");
    },

    async forgotPassword(
        email?: string,
        username?: string
    ): Promise<ApiResponse> {
        const body: { email?: string; username?: string } = {};
        if (email) body.email = email;
        if (username) body.username = username;

        return httpClient.post<ApiResponse>("/auth/forgot-password", body);
    },

    async resetPassword(
        token: string,
        newPassword: string
    ): Promise<ApiResponse> {
        return httpClient.post<ApiResponse>("/auth/reset-password", {
            token,
            new_password: newPassword,
        });
    },
};
