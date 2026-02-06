import { ApiError } from "./api-error";
import { API_BASE_URL } from "./constants";
import type { User } from "../stores/user-store";

interface RequestOptions {
    headers?: Record<string, string>;
    credentials?: RequestCredentials;
    body?: any;
}

// Global variable to track if we're currently refreshing
let isRefreshing = false;
let refreshPromise: Promise<string | null> | null = null;

// Local cache to avoid React state race conditions
let localAccessToken: string | null = null;

let getAccessToken: (() => string | null) | null = () => localAccessToken;
let setAccessToken: ((token: string | null) => void) | null = null;

export function setAuthTokenHandlers(
    getToken: () => string | null,
    setToken: (token: string | null) => void
) {
    // Initialize local cache from external source without wiping an existing token
    const initialToken = getToken();
    if (initialToken) {
        localAccessToken = initialToken;
    }

    // Wrap the getter to check local cache first
    getAccessToken = () => {
        return localAccessToken || getToken();
    };

    // Wrap the setter to update local cache immediately
    setAccessToken = (token: string | null) => {
        localAccessToken = token;
        setToken(token);
    };
}

async function refreshAccessToken(): Promise<string | null> {
    // If already refreshing, return the existing promise
    if (isRefreshing && refreshPromise) {
        return refreshPromise;
    }

    isRefreshing = true;
    refreshPromise = (async () => {
        try {
            console.log("[HTTPClient] Attempting to refresh access token using HTTP-only cookie...");
            const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
                method: "POST",
                credentials: "include", // Send refresh token cookie
            });

            if (response.ok) {
                const data = await response.json();
                const newAccessToken = data.data?.access_token || data.access_token;
                const refreshedUser = data.data?.user || data.user;

                if (newAccessToken) {
                    console.log("[HTTPClient] Access token refreshed successfully");

                    // Always update local cache, even if handlers aren't registered yet
                    localAccessToken = newAccessToken;
                    if (setAccessToken) {
                        setAccessToken(newAccessToken);
                    }

                    // Update in-memory user details if provided
                    if (typeof window !== "undefined" && refreshedUser) {
                        try {
                            const { useUserStore } = await import("../stores/user-store");
                            const existingUser = useUserStore.getState().user;
                            const mergedUser: User = {
                                id: refreshedUser.id ?? existingUser?.id ?? 0,
                                username: refreshedUser.username ?? existingUser?.username ?? "",
                                email: refreshedUser.email ?? existingUser?.email ?? "",
                                firstName: refreshedUser.firstName ?? existingUser?.firstName ?? "",
                                middleName: refreshedUser.middleName ?? existingUser?.middleName ?? "",
                                lastName: refreshedUser.lastName ?? existingUser?.lastName ?? "",
                                avatar: refreshedUser.avatar ?? existingUser?.avatar ?? "",
                                role: refreshedUser.role ?? existingUser?.role ?? "user",
                                createdAt: refreshedUser.createdAt ?? existingUser?.createdAt ?? new Date().toISOString(),
                                updatedAt: refreshedUser.updatedAt ?? existingUser?.updatedAt ?? new Date().toISOString(),
                            };
                            useUserStore.getState().setUser(mergedUser);
                        } catch (storeError) {
                            console.warn("[HTTPClient] Failed to update user store after refresh:", storeError);
                        }
                    }

                    return newAccessToken;
                }
                console.warn("[HTTPClient] Refresh succeeded but no access token was returned");
            } else {
                console.warn(`[HTTPClient] Token refresh failed with status: ${response.status}`);

                // Only clear token if we got a specific 401/403 (unauthorized/forbidden)
                // If it's a 500 or other error, we don't want to forcefully log the user out 
                // as it might be a transient server issue.
                if ((response.status === 401 || response.status === 403) && setAccessToken) {
                    console.warn("[HTTPClient] Session expired or invalid, clearing access token");
                    setAccessToken(null);
                }
            }
            return null;
        } catch (error) {
            console.error("[HTTPClient] Token refresh catch block error:", error);
            // Don't clear access token on network/other errors to avoid unintended logouts
            return null;
        } finally {
            isRefreshing = false;
            refreshPromise = null;
        }
    })();

    return refreshPromise;
}

async function handleResponse<T>(
    response: Response,
    isRetry: boolean = false
): Promise<T> {
    const contentType = response.headers.get("content-type");
    const isJson = contentType?.includes("application/json");

    let data: any = {};

    try {
        const text = await response.text();
        if (text.trim() && isJson) {
            data = JSON.parse(text);
        } else if (text.trim()) {
            data = { message: text };
        }
    } catch (error) {
        console.error("Failed to parse response:", error);
        data = { message: `HTTP ${response.status}: ${response.statusText}` };
    }

    if (!response.ok) {
        // Handle 401 Unauthorized - try to refresh token
        if (response.status === 401 && !isRetry) {
            console.log("[HTTPClient] 401 Unauthorized detected, attempting refresh...");
            const newToken = await refreshAccessToken();

            if (newToken) {
                console.log("[HTTPClient] Refresh successful, retrying original request");
                // Token refreshed successfully, throw a special error to retry
                throw new ApiError("TOKEN_REFRESHED", 401, null, true);
            } else {
                console.warn("[HTTPClient] Refresh failed or not possible");

                // Only clear user state and redirect if the store thinks we are authenticated
                // and we've confirmed the refresh token is truly invalid/expired.
                // We import dynamically to avoid circular dependencies.
                const { useUserStore } = await import("../stores/user-store");
                const userStore = useUserStore.getState();

                // IMPORTANT: We only clear the user if the refresh attempt confirmed 
                // the session is invalid. If refreshAccessToken returned null without 
                // clearing setAccessToken(null), we treat it as a transient error.
                if (userStore.isAuthenticated && localAccessToken === null) {
                    console.warn("[HTTPClient] Forcefully logging out user due to invalid session");
                    userStore.clearUser();

                    // Redirect to login if not already on public pages
                    if (
                        typeof window !== "undefined" &&
                        !window.location.pathname.startsWith("/login") &&
                        !window.location.pathname.startsWith("/register") &&
                        !window.location.pathname.startsWith("/about") &&
                        !window.location.pathname.startsWith("/contact") &&
                        !window.location.pathname.startsWith("/blog") &&
                        !window.location.pathname.startsWith("/shop")
                    ) {
                        window.location.href = "/login";
                    }
                }
            }
        }

        // Extract error message
        const errorMessage =
            data?.message ||
            data?.error?.detail ||
            data?.error?.message ||
            (Array.isArray(data?.errors) ? data.errors.join(", ") : "") ||
            `HTTP ${response.status}: ${response.statusText}`;

        console.warn(`[HTTPClient] API Error (${response.status}): ${errorMessage}`);
        return {
            error: true,
            status: response.status,
            message: errorMessage,
            details: data?.errors
        } as unknown as T;
    }

    // Backend returns structured responses, so just return the data as-is
    return data as T;
}

async function makeRequest<T>(
    url: string,
    options: RequestInit,
    isRetry: boolean = false
): Promise<T> {
    // Add Authorization header if access token exists
    const token = getAccessToken?.();
    if (token) {
        options.headers = {
            ...options.headers,
            Authorization: `Bearer ${token}`,
        };
    }
    console.log(`[HTTPClient] Request to ${url} | Token present: ${!!token}`, token ? `(Length: ${token.length})` : '');

    try {
        const response = await fetch(url, options);
        return await handleResponse<T>(response, isRetry);
    } catch (error) {
        // If token was refreshed, retry the request once
        if (error instanceof ApiError && error.shouldRetry && !isRetry) {
            return makeRequest<T>(url, options, true);
        }
        throw error;
    }
}

export const httpClient = {
    async get<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
        return makeRequest<T>(`${API_BASE_URL}${endpoint}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                ...options.headers,
            },
            credentials: options.credentials || "include",
        });
    },

    async post<T>(
        endpoint: string,
        body?: any,
        options: RequestOptions = {}
    ): Promise<T> {
        return makeRequest<T>(`${API_BASE_URL}${endpoint}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                ...options.headers,
            },
            credentials: options.credentials || "include",
            body: body ? JSON.stringify(body) : undefined,
        });
    },

    async put<T>(
        endpoint: string,
        body?: any,
        options: RequestOptions = {}
    ): Promise<T> {
        return makeRequest<T>(`${API_BASE_URL}${endpoint}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                ...options.headers,
            },
            credentials: options.credentials || "include",
            body: body ? JSON.stringify(body) : undefined,
        });
    },

    async patch<T>(
        endpoint: string,
        body?: any,
        options: RequestOptions = {}
    ): Promise<T> {
        return makeRequest<T>(`${API_BASE_URL}${endpoint}`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                ...options.headers,
            },
            credentials: options.credentials || "include",
            body: body ? JSON.stringify(body) : undefined,
        });
    },

    async delete<T>(
        endpoint: string,
        body?: any,
        options: RequestOptions = {}
    ): Promise<T> {
        return makeRequest<T>(`${API_BASE_URL}${endpoint}`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
                ...options.headers,
            },
            credentials: options.credentials || "include",
            body: body ? JSON.stringify(body) : undefined,
        });
    },
};
