# Build stage
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /build

# Copy package files
COPY frontend/package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY frontend/ ./

# Build Next.js application
RUN npm run build

# Production stage
FROM node:20-alpine

# Set working directory
WORKDIR /app

# Copy package files
COPY frontend/package*.json ./

# Install production dependencies only
RUN npm ci --only=production

# Copy built application
COPY --from=builder /build/.next ./.next
COPY --from=builder /build/public ./public
COPY --from=builder /build/next.config.ts ./
COPY --from=builder /build/tsconfig.json ./

# Expose port
EXPOSE 3000

# Set environment to production
ENV NODE_ENV=production

# Run the application
CMD ["npm", "start"]

