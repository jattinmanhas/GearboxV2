import { httpClient } from "./http-client";
import type {
    OAuthProvider,
    OAuthInitiateResponse,
    LinkedProvidersResponse,
    ApiResponse,
} from "../types";

export const oauthApi = {
    // Initiate OAuth flow
    async initiateOAuth(provider: OAuthProvider): Promise<OAuthInitiateResponse> {
        const response = await httpClient.get<any>(`/auth/oauth/${provider}`);
        // Backend wraps responses in { data: ... }; fall back to raw response for non-wrapped cases
        return response?.data ?? response;
    },

    // Get linked OAuth providers
    async getLinkedProviders(): Promise<LinkedProvidersResponse> {
        const data = await httpClient.get<{ data: LinkedProvidersResponse }>(
            "/auth/oauth/providers"
        );
        return data.data;
    },

    // Unlink OAuth provider
    async unlinkProvider(provider: OAuthProvider): Promise<ApiResponse> {
        return httpClient.delete<ApiResponse>(`/auth/oauth/unlink/${provider}`);
    },

    // Link OAuth provider (used when user is already logged in)
    async linkProvider(
        provider: OAuthProvider,
        code: string
    ): Promise<ApiResponse> {
        return httpClient.post<ApiResponse>(`/auth/oauth/link/${provider}`, {
            provider,
            code,
        });
    },
};
