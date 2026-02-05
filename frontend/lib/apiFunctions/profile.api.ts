import { httpClient } from "./http-client";
import type {
    UserProfile,
    Address,
    CreateAddressRequest,
    UpdateAddressRequest,
    UpdateProfileRequest,
    ApiResponse,
} from "../types";

export const profileApi = {
    // User Profile
    async getProfile(): Promise<UserProfile> {
        const response = await httpClient.get<ApiResponse<UserProfile>>("/auth/profile");
        return response.data as UserProfile;
    },

    async updateProfile(profileData: UpdateProfileRequest): Promise<UserProfile> {
        const response = await httpClient.put<ApiResponse<UserProfile>>("/auth/profile", profileData);
        return response.data as UserProfile;
    },

    // Addresses
    async getAddresses(): Promise<Address[]> {
        const responseData = await httpClient.get<any>("/auth/addresses");

        // Extract addresses array from response
        const addresses = Array.isArray(responseData)
            ? responseData
            : responseData?.data
                ? Array.isArray(responseData.data)
                    ? responseData.data
                    : []
                : [];

        // Map backend address structure to frontend Address type
        return addresses.map((addr: any) => ({
            id: addr.id,
            user_id: addr.user_id,
            type: (addr.address_type || addr.type) as Address["type"],
            first_name: addr.first_name,
            last_name: addr.last_name,
            company: addr.company,
            address_line_1: addr.address_line_1,
            address_line_2: addr.address_line_2,
            city: addr.city,
            state: addr.state,
            postal_code: addr.postal_code,
            country: addr.country,
            phone_number: addr.phone || addr.phone_number,
            is_default: addr.is_default,
            created_at: addr.created_at,
            updated_at: addr.updated_at,
        }));
    },

    async getAddress(id: number): Promise<Address> {
        const response = await httpClient.get<any>(`/auth/addresses/${id}`);
        const addr = response.data || response;

        // Map backend address structure to frontend Address type
        return {
            id: addr.id,
            user_id: addr.user_id,
            type: (addr.address_type || addr.type) as Address["type"],
            first_name: addr.first_name,
            last_name: addr.last_name,
            company: addr.company,
            address_line_1: addr.address_line_1,
            address_line_2: addr.address_line_2,
            city: addr.city,
            state: addr.state,
            postal_code: addr.postal_code,
            country: addr.country,
            phone_number: addr.phone || addr.phone_number,
            is_default: addr.is_default,
            created_at: addr.created_at,
            updated_at: addr.updated_at,
        };
    },

    async createAddress(addressData: CreateAddressRequest): Promise<Address> {
        // Map frontend address structure to backend format
        const backendData = {
            address_type: addressData.type,
            first_name: addressData.first_name,
            last_name: addressData.last_name,
            company: addressData.company,
            address_line_1: addressData.address_line_1,
            address_line_2: addressData.address_line_2,
            city: addressData.city,
            state: addressData.state,
            country: addressData.country,
            postal_code: addressData.postal_code,
            phone: addressData.phone_number,
            is_default: addressData.is_default,
        };

        const response = await httpClient.post<any>("/auth/addresses", backendData);
        const addr = response.data || response;

        // Map backend address structure to frontend Address type
        return {
            id: addr.id,
            user_id: addr.user_id,
            type: (addr.address_type || addr.type) as Address["type"],
            first_name: addr.first_name,
            last_name: addr.last_name,
            company: addr.company,
            address_line_1: addr.address_line_1,
            address_line_2: addr.address_line_2,
            city: addr.city,
            state: addr.state,
            postal_code: addr.postal_code,
            country: addr.country,
            phone_number: addr.phone || addr.phone_number,
            is_default: addr.is_default,
            created_at: addr.created_at,
            updated_at: addr.updated_at,
        };
    },

    async updateAddress(
        id: number,
        addressData: UpdateAddressRequest
    ): Promise<Address> {
        // Map frontend address structure to backend format
        const backendData: any = {};
        if (addressData.type !== undefined)
            backendData.address_type = addressData.type;
        if (addressData.first_name !== undefined)
            backendData.first_name = addressData.first_name;
        if (addressData.last_name !== undefined)
            backendData.last_name = addressData.last_name;
        if (addressData.company !== undefined)
            backendData.company = addressData.company;
        if (addressData.address_line_1 !== undefined)
            backendData.address_line_1 = addressData.address_line_1;
        if (addressData.address_line_2 !== undefined)
            backendData.address_line_2 = addressData.address_line_2;
        if (addressData.city !== undefined) backendData.city = addressData.city;
        if (addressData.state !== undefined) backendData.state = addressData.state;
        if (addressData.country !== undefined)
            backendData.country = addressData.country;
        if (addressData.postal_code !== undefined)
            backendData.postal_code = addressData.postal_code;
        if (addressData.phone_number !== undefined)
            backendData.phone = addressData.phone_number;
        if (addressData.is_default !== undefined)
            backendData.is_default = addressData.is_default;

        const response = await httpClient.put<any>(
            `/auth/addresses/${id}`,
            backendData
        );
        const addr = response.data || response;

        // Map backend address structure to frontend Address type
        return {
            id: addr.id,
            user_id: addr.user_id,
            type: (addr.address_type || addr.type) as Address["type"],
            first_name: addr.first_name,
            last_name: addr.last_name,
            company: addr.company,
            address_line_1: addr.address_line_1,
            address_line_2: addr.address_line_2,
            city: addr.city,
            state: addr.state,
            postal_code: addr.postal_code,
            country: addr.country,
            phone_number: addr.phone || addr.phone_number,
            is_default: addr.is_default,
            created_at: addr.created_at,
            updated_at: addr.updated_at,
        };
    },

    async deleteAddress(id: number): Promise<ApiResponse> {
        return httpClient.delete<ApiResponse>(`/auth/addresses/${id}`);
    },
};
