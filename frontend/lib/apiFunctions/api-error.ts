export class ApiError extends Error {
    constructor(
        message: string,
        public status: number,
        public errors?: string[] | null,
        public shouldRetry: boolean = false
    ) {
        super(message);
        this.name = "ApiError";
    }
}
