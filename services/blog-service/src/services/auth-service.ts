import jwt from 'jsonwebtoken';
import { Claims } from '../middleware/auth';

/**
 * Authentication service for JWT token validation and generation
 */
export class AuthService {
  private accessTokenSecret: string;
  private refreshTokenSecret: string;
  private accessTokenExpiry: string;
  private refreshTokenExpiry: string;

  constructor() {
    this.accessTokenSecret = process.env.JWT_ACCESS_SECRET || 'your-access-secret';
    this.refreshTokenSecret = process.env.JWT_REFRESH_SECRET || 'your-refresh-secret';
    this.accessTokenExpiry = process.env.JWT_ACCESS_EXPIRY || '15m';
    this.refreshTokenExpiry = process.env.JWT_REFRESH_EXPIRY || '7d';
  }

  /**
   * Validate access token and return claims
   */
  async validateAccessToken(token: string): Promise<Claims | null> {
    try {
      const decoded = jwt.verify(token, this.accessTokenSecret) as any;
      return {
        user_id: decoded.user_id,
        username: decoded.username,
        email: decoded.email,
        role: decoded.role,
        iat: decoded.iat,
        exp: decoded.exp,
      };
    } catch (error) {
      console.warn('Access token validation failed:', error);
      return null;
    }
  }

  /**
   * Validate refresh token and return claims
   */
  async validateRefreshToken(token: string): Promise<Claims | null> {
    try {
      const decoded = jwt.verify(token, this.refreshTokenSecret) as any;
      return {
        user_id: decoded.user_id,
        username: decoded.username,
        email: decoded.email,
        role: decoded.role,
        iat: decoded.iat,
        exp: decoded.exp,
      };
    } catch (error) {
      console.warn('Refresh token validation failed:', error);
      return null;
    }
  }

  /**
   * Generate access token from user information
   */
  async generateAccessTokenFromUser(user: Claims): Promise<string> {
    const payload = {
      user_id: user.user_id,
      username: user.username,
      email: user.email,
      role: user.role,
    };

    return jwt.sign(payload, this.accessTokenSecret, {
      expiresIn: this.accessTokenExpiry,
    } as jwt.SignOptions);
  }

  /**
   * Generate refresh token from user information
   */
  async generateRefreshTokenFromUser(user: Claims): Promise<string> {
    const payload = {
      user_id: user.user_id,
      username: user.username,
      email: user.email,
      role: user.role,
    };

    return jwt.sign(payload, this.refreshTokenSecret, {
      expiresIn: this.refreshTokenExpiry,
    } as jwt.SignOptions);
  }

  /**
   * Extract token from Authorization header
   */
  extractTokenFromHeader(authHeader?: string): string | null {
    if (authHeader && authHeader.startsWith('Bearer ')) {
      return authHeader.substring(7);
    }
    return null;
  }

  /**
   * Extract token from cookie
   */
  extractTokenFromCookie(cookies: Record<string, string>, cookieName: string): string | null {
    return cookies?.[cookieName] || null;
  }
}
