import { httpClient } from "./http-client";

export const paymentApi = {
    async getPaymentMethods(): Promise<any> {
        return httpClient.get<any>("/payment-methods");
    },

    async getPaymentGateways(): Promise<any> {
        return httpClient.get<any>("/payment-gateways");
    },

    async createPayment(paymentData: any): Promise<any> {
        return httpClient.post<any>("/payments", paymentData);
    },

    async processPayment(paymentId: number, paymentData: any): Promise<any> {
        return httpClient.post<any>("/payments/process", {
            payment_id: paymentId,
            ...paymentData,
        });
    },
};
