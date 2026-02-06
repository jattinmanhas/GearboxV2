import { httpClient } from "./http-client";

export const paymentApi = {
    async createPayment(paymentData: any): Promise<any> {
        return httpClient.post<any>("/payments", paymentData);
    },

    async processPayment(paymentId: number, paymentData: any): Promise<any> {
        return httpClient.post<any>("/payments/process", {
            payment_id: paymentId,
            ...paymentData,
        });
    },

    async listPayments(params?: any): Promise<any> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.keys(params).forEach(key => {
                if (params[key] !== undefined && params[key] !== null) {
                    searchParams.append(key, params[key]);
                }
            });
        }
        return httpClient.get<any>(`/payments?${searchParams.toString()}`);
    },

    async getPayment(id: number): Promise<any> {
        return httpClient.get<any>(`/payments/${id}`);
    },

    async updatePaymentStatus(id: number, statusData: any): Promise<any> {
        return httpClient.put<any>(`/payments/${id}/status`, statusData);
    },

    async refundPayment(refundData: any): Promise<any> {
        return httpClient.post<any>("/payments/refund", refundData);
    },

    async getPaymentSummary(params?: any): Promise<any> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.keys(params).forEach(key => {
                if (params[key] !== undefined && params[key] !== null) {
                    searchParams.append(key, params[key]);
                }
            });
        }
        return httpClient.get<any>(`/payments/summary?${searchParams.toString()}`);
    },
};
