import { ApiError } from "./api-error";
import { API_BASE_URL } from "./constants";

interface RequestOptions {
    headers?: Record<string, string>;
    credentials?: RequestCredentials;
    body?: any;
}

async function handleResponse<T>(response: Response): Promise<T> {
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
        // Handle 401 Unauthorized - clear user state and redirect
        if (response.status === 401) {
            const { useUserStore } = await import("../stores/user-store");
            const userStore = useUserStore.getState();

            if (userStore.isAuthenticated && userStore.user) {
                const errorMessage = data?.message || "";
                const isAuthFailure = errorMessage.toLowerCase().includes("token") ||
                    errorMessage.toLowerCase().includes("unauthorized") ||
                    errorMessage.toLowerCase().includes("authentication");

                if (isAuthFailure) {
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

export const httpClient = {
    async get<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                ...options.headers,
            },
            credentials: options.credentials || "include",
        });

        return handleResponse<T>(response);
    },

    async post<T>(
        endpoint: string,
        body?: any,
        options: RequestOptions = {}
    ): Promise<T> {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                ...options.headers,
            },
            credentials: options.credentials || "include",
            body: body ? JSON.stringify(body) : undefined,
        });

        return handleResponse<T>(response);
    },

    async put<T>(
        endpoint: string,
        body?: any,
        options: RequestOptions = {}
    ): Promise<T> {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                ...options.headers,
            },
            credentials: options.credentials || "include",
            body: body ? JSON.stringify(body) : undefined,
        });

        return handleResponse<T>(response);
    },

    async delete<T>(
        endpoint: string,
        body?: any,
        options: RequestOptions = {}
    ): Promise<T> {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
                ...options.headers,
            },
            credentials: options.credentials || "include",
            body: body ? JSON.stringify(body) : undefined,
        });

        return handleResponse<T>(response);
    },
};
