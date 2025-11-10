# Build stage
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /build

# Copy package files
COPY services/blog-service/package*.json ./

# Install dependencies (including dev dependencies for build)
RUN npm ci

# Copy source code
COPY services/blog-service/ ./

# Build TypeScript to JavaScript
RUN npm run build

# Production stage
FROM node:20-alpine

# Set working directory
WORKDIR /app

# Copy package files
COPY services/blog-service/package*.json ./

# Install dependencies (including tsx for running TypeScript)
RUN npm ci --only=production && \
    npm install -g tsx

# Copy built application
COPY --from=builder /build/dist ./dist

# Copy source files (for tsx runtime)
COPY --from=builder /build/src ./src

# Copy migrations
COPY --from=builder /build/migrations ./migrations

# Copy tsconfig
COPY --from=builder /build/tsconfig.json ./

# Expose port
EXPOSE 8084

# Set environment
ENV NODE_ENV=production

# Run the application using tsx
CMD ["tsx", "src/index.ts"]

