import { FastifyRequest, FastifyReply } from 'fastify';
import { AuthService } from '../services/auth-service';

/**
 * JWT Claims structure matching Go services
 */
export interface Claims {
  user_id: number;
  username: string;
  email: string;
  role: string;
  iat?: number;
  exp?: number;
}

/**
 * Authentication middleware for protecting routes
 */
export class AuthMiddleware {
  private authService: AuthService;

  constructor(authService: AuthService) {
    this.authService = authService;
  }

  /**
   * Middleware that requires authentication
   */
  requireAuth = async (request: FastifyRequest, reply: FastifyReply) => {
    try {
      // Extract token from Authorization header or cookie
      let token = this.extractTokenFromHeader(request);
      if (!token) {
        token = this.extractTokenFromCookie(request, 'access_token');
      }

      if (!token) {
        return reply.status(401).send({
          timestamp: new Date().toISOString(),
          status: 401,
          success: false,
          message: 'Authentication required',
          error: { message: 'No token provided' }
        });
      }

      // Validate access token
      let claims = await this.authService.validateAccessToken(token);
      
      if (!claims) {
        // Try refresh token if access token is invalid
        const refreshToken = this.extractTokenFromCookie(request, 'refresh_token');
        if (!refreshToken) {
          return reply.status(401).send({
            timestamp: new Date().toISOString(),
            status: 401,
            success: false,
            message: 'Authentication required',
            error: { message: 'Invalid access token and no refresh token provided' }
          });
        }

        // Validate refresh token and generate new access token
        const refreshClaims = await this.authService.validateRefreshToken(refreshToken);
        if (!refreshClaims) {
          return reply.status(401).send({
            timestamp: new Date().toISOString(),
            status: 401,
            success: false,
            message: 'Authentication required',
            error: { message: 'Invalid refresh token' }
          });
        }

        // Generate new access token
        const newAccessToken = await this.authService.generateAccessTokenFromUser(refreshClaims);
        
        // Set new access token in cookie
        reply.cookie('access_token', newAccessToken, {
          path: '/',
          httpOnly: true,
          secure: true,
          sameSite: 'strict',
          maxAge: 900 // 15 minutes
        });

        claims = refreshClaims;
      }

      // Add claims to request object for use in handlers
      (request as any).user = claims;
      (request as any).claims = claims;

    } catch (error) {
      return reply.status(401).send({
        timestamp: new Date().toISOString(),
        status: 401,
        success: false,
        message: 'Authentication failed',
        error: { message: error instanceof Error ? error.message : 'Unknown error' }
      });
    }
  };

  /**
   * Optional authentication middleware - doesn't require auth but adds user if present
   */
  optionalAuth = async (request: FastifyRequest, reply: FastifyReply) => {
    try {
      // Extract token from Authorization header or cookie
      let token = this.extractTokenFromHeader(request);
      if (!token) {
        token = this.extractTokenFromCookie(request, 'access_token');
      }

      if (token) {
        const claims = await this.authService.validateAccessToken(token);
        if (claims) {
          (request as any).user = claims;
          (request as any).claims = claims;
        }
      }
    } catch (error) {
      // Silently fail for optional auth
      console.warn('Optional auth failed:', error);
    }
  };

  /**
   * Role-based access control middleware
   */
  requireRole = (requiredRole: string) => {
    return async (request: FastifyRequest, reply: FastifyReply) => {
      const claims = (request as any).claims as Claims;
      
      if (!claims) {
        return reply.status(401).send({
          timestamp: new Date().toISOString(),
          status: 401,
          success: false,
          message: 'Authentication required',
          error: { message: 'No user claims found' }
        });
      }

      if (!this.canAccess(claims.role, requiredRole)) {
        return reply.status(403).send({
          timestamp: new Date().toISOString(),
          status: 403,
          success: false,
          message: 'Insufficient permissions',
          error: { message: `Role '${claims.role}' cannot access '${requiredRole}' resources` }
        });
      }
    };
  };

  /**
   * Admin-only middleware
   */
  requireAdmin = async (request: FastifyRequest, reply: FastifyReply) => {
    return this.requireRole('admin')(request, reply);
  };

  /**
   * Editor or higher middleware
   */
  requireEditor = async (request: FastifyRequest, reply: FastifyReply) => {
    return this.requireRole('editor')(request, reply);
  };

  /**
   * Extract token from Authorization header
   */
  private extractTokenFromHeader(request: FastifyRequest): string | null {
    const authHeader = request.headers.authorization;
    if (authHeader && authHeader.startsWith('Bearer ')) {
      return authHeader.substring(7);
    }
    return null;
  }

  /**
   * Extract token from cookie
   */
  private extractTokenFromCookie(request: FastifyRequest, cookieName: string): string | null {
    return (request.cookies as Record<string, string>)?.[cookieName] || null;
  }

  /**
   * Check if user role can access required role
   */
  private canAccess(userRole: string, requiredRole: string): boolean {
    const roleLevels: Record<string, number> = {
      'user': 1,
      'editor': 2,
      'admin': 3,
    };

    const userLevel = roleLevels[userRole] || 0;
    const requiredLevel = roleLevels[requiredRole] || 0;

    return userLevel >= requiredLevel;
  }
}

/**
 * Helper functions to extract user information from request
 */
export function getUserFromRequest(request: FastifyRequest): Claims | null {
  return (request as any).user || null;
}

export function getClaimsFromRequest(request: FastifyRequest): Claims | null {
  return (request as any).claims || null;
}

export function getUserIdFromRequest(request: FastifyRequest): number | null {
  const claims = getClaimsFromRequest(request);
  return claims?.user_id || null;
}
