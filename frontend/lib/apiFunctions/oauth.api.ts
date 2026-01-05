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
        return httpClient.get<OAuthInitiateResponse>(`/auth/oauth/${provider}`);
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
