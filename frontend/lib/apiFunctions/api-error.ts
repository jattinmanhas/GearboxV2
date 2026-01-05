export class ApiError extends Error {
    constructor(
        message: string,
        public status: number,
        public errors?: string[]
    ) {
        super(message);
        this.name = "ApiError";
    }
}
