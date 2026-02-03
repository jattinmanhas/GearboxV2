import { ApiError } from "./api-error";
import { API_BASE_URL } from "./constants";

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

let getAccessToken: (() => string | null) | null = null;
let setAccessToken: ((token: string | null) => void) | null = null;

export function setAuthTokenHandlers(
    getToken: () => string | null,
    setToken: (token: string | null) => void
) {
    // Initialize local cache from external source
    localAccessToken = getToken();

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
            const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
                method: "POST",
                credentials: "include", // Send refresh token cookie
            });

            if (response.ok) {
                const data = await response.json();
                const newAccessToken = data.data?.access_token || data.access_token;

                if (newAccessToken && setAccessToken) {
                    setAccessToken(newAccessToken);
                    return newAccessToken;
                }
            }

            // Refresh failed, clear token
            if (setAccessToken) {
                setAccessToken(null);
            }
            return null;
        } catch (error) {
            console.error("Token refresh failed:", error);
            if (setAccessToken) {
                setAccessToken(null);
            }
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
            const newToken = await refreshAccessToken();

            if (newToken) {
                // Token refreshed successfully, throw a special error to retry
                throw new ApiError("TOKEN_REFRESHED", 401, null, true);
            } else {
                // Refresh failed, clear user state and redirect
                const { useUserStore } = await import("../stores/user-store");
                const userStore = useUserStore.getState();

                if (userStore.isAuthenticated) {
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

        throw new ApiError(errorMessage, response.status, data?.errors);
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
