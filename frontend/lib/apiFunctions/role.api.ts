import { httpClient } from "./http-client";
import type {
    Role,
    AssignRoleRequest,
    RemoveRoleRequest,
    CheckPermissionRequest,
    CheckPermissionResponse,
    ApiResponse,
} from "../types";

export const roleApi = {
    // Get all roles
    async getRoles(): Promise<Role[]> {
        const data = await httpClient.get<{ data: Role[] }>("/auth/roles");
        return data.data || [];
    },

    // Get current user's role
    async getMyRole(): Promise<Role> {
        return httpClient.get<Role>("/auth/roles/my-role");
    },

    // Get user's role
    async getUserRole(userId: number): Promise<Role> {
        return httpClient.get<Role>(`/auth/roles/user?user_id=${userId}`);
    },

    // Assign role to user (Editor+ only)
    async assignRole(assignmentData: AssignRoleRequest): Promise<ApiResponse> {
        return httpClient.post<ApiResponse>("/auth/roles/assign", assignmentData);
    },

    // Remove role from user (Editor+ only)
    async removeRole(removalData: RemoveRoleRequest): Promise<ApiResponse> {
        return httpClient.delete<ApiResponse>("/auth/roles/remove", removalData);
    },

    // Check permission
    async checkPermission(
        permissionData: CheckPermissionRequest
    ): Promise<CheckPermissionResponse> {
        const params = new URLSearchParams();
        params.append("permission", permissionData.permission);
        if (permissionData.resource) {
            params.append("resource", permissionData.resource);
        }

        return httpClient.get<CheckPermissionResponse>(
            `/auth/roles/check-permission?${params}`
        );
    },
};
