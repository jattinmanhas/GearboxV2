import { defineConfig } from 'drizzle-kit';

export default defineConfig({
  schema: './src/config/schema.ts',
  out: './migrations',
  dialect: 'postgresql',
  dbCredentials: {
    url: process.env.DATABASE_URL || 'postgresql://localhost:5432/blog_service',
  },
  verbose: true,
  strict: true,
});
