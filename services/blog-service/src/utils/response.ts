import { FastifyReply } from 'fastify';

/**
 * Standardized API response structure matching Go services
 */
export interface APIResponse {
  timestamp: string;
  status: number;
  success: boolean;
  message?: string;
  data?: any;
  error?: any;
}

/**
 * Response utility class for consistent API responses
 */
export class ResponseHelper {
  /**
   * Send a successful response
   */
  static ok(reply: FastifyReply, message: string, data?: any): void {
    const response: APIResponse = {
      timestamp: new Date().toISOString(),
      status: 200,
      success: true,
      message,
      data,
    };
    reply.status(200).send(response);
  }

  /**
   * Send a created response
   */
  static created(reply: FastifyReply, message: string, data?: any): void {
    const response: APIResponse = {
      timestamp: new Date().toISOString(),
      status: 201,
      success: true,
      message,
      data,
    };
    reply.status(201).send(response);
  }

  /**
   * Send an error response
   */
  static error(reply: FastifyReply, status: number, message: string, error?: any): void {
    const errorPayload: any = { message };
    if (error) {
      errorPayload.detail = error instanceof Error ? error.message : error;
    }

    const response: APIResponse = {
      timestamp: new Date().toISOString(),
      status,
      success: false,
      message,
      error: errorPayload,
    };
    reply.status(status).send(response);
  }

  /**
   * Send a bad request response (400)
   */
  static badRequest(reply: FastifyReply, message: string, data?: any): void {
    const response: APIResponse = {
      timestamp: new Date().toISOString(),
      status: 400,
      success: false,
      message,
      error: data,
    };
    reply.status(400).send(response);
  }

  /**
   * Send a not found response (404)
   */
  static notFound(reply: FastifyReply, message: string, data?: any): void {
    const response: APIResponse = {
      timestamp: new Date().toISOString(),
      status: 404,
      success: false,
      message,
      error: data,
    };
    reply.status(404).send(response);
  }

  /**
   * Send an unauthorized response (401)
   */
  static unauthorized(reply: FastifyReply, message: string, data?: any): void {
    const response: APIResponse = {
      timestamp: new Date().toISOString(),
      status: 401,
      success: false,
      message,
      error: data,
    };
    reply.status(401).send(response);
  }

  /**
   * Send a forbidden response (403)
   */
  static forbidden(reply: FastifyReply, message: string, data?: any): void {
    const response: APIResponse = {
      timestamp: new Date().toISOString(),
      status: 403,
      success: false,
      message,
      error: data,
    };
    reply.status(403).send(response);
  }

  /**
   * Send an internal server error response (500)
   */
  static internalServerError(reply: FastifyReply, message: string, data?: any): void {
    const response: APIResponse = {
      timestamp: new Date().toISOString(),
      status: 500,
      success: false,
      message,
      error: data,
    };
    reply.status(500).send(response);
  }
}
