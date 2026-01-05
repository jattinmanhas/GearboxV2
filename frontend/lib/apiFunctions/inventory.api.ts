import { httpClient } from "./http-client";
import type { ApiResponse } from "../types";

export const inventoryApi = {
    async getInventory(params: string = ""): Promise<ApiResponse> {
        try {
            return await httpClient.get<ApiResponse>(
                `/inventory${params ? `?${params}` : ""}`
            );
        } catch (error) {
            console.error("Error fetching inventory:", error);
            return {
                success: false,
                message:
                    error instanceof Error ? error.message : "Failed to fetch inventory",
                data: null,
            };
        }
    },

    async getInventorySummary(): Promise<ApiResponse> {
        try {
            return await httpClient.get<ApiResponse>("/inventory/summary");
        } catch (error) {
            console.error("Error fetching inventory summary:", error);
            return {
                success: false,
                message:
                    error instanceof Error
                        ? error.message
                        : "Failed to fetch inventory summary",
                data: null,
            };
        }
    },

    async getInventoryAlerts(params: string = ""): Promise<ApiResponse> {
        try {
            return await httpClient.get<ApiResponse>(
                `/inventory/alerts${params ? `?${params}` : ""}`
            );
        } catch (error) {
            console.error("Error fetching inventory alerts:", error);
            return {
                success: false,
                message:
                    error instanceof Error
                        ? error.message
                        : "Failed to fetch inventory alerts",
                data: null,
            };
        }
    },

    async resolveInventoryAlert(id: number): Promise<ApiResponse> {
        try {
            return await httpClient.put<ApiResponse>(
                `/inventory/alerts/${id}/resolve`
            );
        } catch (error) {
            console.error("Error resolving inventory alert:", error);
            return {
                success: false,
                message:
                    error instanceof Error
                        ? error.message
                        : "Failed to resolve inventory alert",
                data: null,
            };
        }
    },

    async checkInventoryAlerts(): Promise<ApiResponse> {
        try {
            return await httpClient.post<ApiResponse>("/inventory/alerts/check");
        } catch (error) {
            console.error("Error checking inventory alerts:", error);
            return {
                success: false,
                message:
                    error instanceof Error
                        ? error.message
                        : "Failed to check inventory alerts",
                data: null,
            };
        }
    },
};
