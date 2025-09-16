import Fastify from 'fastify';
import cors from '@fastify/cors';
import helmet from '@fastify/helmet';
import rateLimit from '@fastify/rate-limit';
import swagger from '@fastify/swagger';
import swaggerUi from '@fastify/swagger-ui';
import cookie from '@fastify/cookie';
import { blogRoutes } from './routes/blog-routes';
import { categoryRoutes } from './routes/category-routes';
import config from './config';

const fastify = Fastify({
  logger: {
    level: process.env.LOG_LEVEL || 'info',
    transport: process.env.NODE_ENV === 'development' ? {
      target: 'pino-pretty',
      options: {
        colorize: true
      }
    } : undefined
  }
});

// Register plugins
async function registerPlugins() {
  // CORS
  await fastify.register(cors, {
    origin: config.cors.origin,
    credentials: true
  });

  // Cookie support
  await fastify.register(cookie, {
    secret: process.env.COOKIE_SECRET || 'your-cookie-secret'
  });

  // Security headers
  await fastify.register(helmet, {
    contentSecurityPolicy: false
  });

  // Rate limiting
  await fastify.register(rateLimit, {
    max: 100,
    timeWindow: '1 minute'
  });

  // Swagger documentation
  await fastify.register(swagger, {
    swagger: {
      info: {
        title: 'Blog Service API',
        description: 'API documentation for the blog service',
        version: '1.0.0'
      },
      host: `${config.host}:${config.port}`,
      schemes: ['http', 'https'],
      consumes: ['application/json'],
      produces: ['application/json'],
      tags: [
        { name: 'Blog Posts', description: 'Blog post related endpoints' },
        { name: 'Categories', description: 'Category related endpoints' }
      ]
    }
  });

  await fastify.register(swaggerUi, {
    routePrefix: '/docs',
    uiConfig: {
      docExpansion: 'list',
      deepLinking: false
    }
  });
}

// Register routes
async function registerRoutes() {
  // Health check
  fastify.get('/health', {
    schema: {
      description: 'Health check endpoint',
      tags: ['Health'],
      response: {
        200: {
          type: 'object',
          properties: {
            status: { type: 'string' },
            timestamp: { type: 'string' },
            uptime: { type: 'number' }
          }
        }
      }
    }
  }, async (request, reply) => {
    return {
      status: 'healthy',
      timestamp: new Date().toISOString(),
      uptime: process.uptime()
    };
  });

  // API routes
  await fastify.register(blogRoutes, { prefix: '/api/v1' });
  await fastify.register(categoryRoutes, { prefix: '/api/v1' });
}

// Error handler
fastify.setErrorHandler((error, request, reply) => {
  fastify.log.error(error);
  
  if (error.validation) {
    return reply.status(400).send({
      success: false,
      message: 'Validation error',
      errors: error.validation
    });
  }

  return reply.status(500).send({
    success: false,
    message: 'Internal server error'
  });
});

// Not found handler
fastify.setNotFoundHandler((request, reply) => {
  return reply.status(404).send({
    success: false,
    message: 'Route not found'
  });
});

// Setup function to initialize the server
async function setup() {
  await registerPlugins();
  await registerRoutes();
}

// Start server function
async function start() {
  try {
    await setup();
    
    await fastify.listen({ 
      port: config.port, 
      host: config.host 
    });
    
    fastify.log.info(`Blog service is running on http://${config.host}:${config.port}`);
    fastify.log.info(`API documentation available at http://${config.host}:${config.port}/docs`);
  } catch (error) {
    fastify.log.error({ error }, 'Failed to start server');
    throw error;
  }
}

// Graceful shutdown function
async function stop() {
  try {
    fastify.log.info('Shutting down gracefully...');
    await fastify.close();
  } catch (error) {
    fastify.log.error({ error }, 'Error during shutdown');
    throw error;
  }
}

export { fastify, setup, start, stop };
export default fastify;
