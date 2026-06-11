# ==============================================================================
# Build Stage
# ==============================================================================
FROM node:22-alpine AS builder

WORKDIR /app

# Copy dependency mappings
COPY package*.json tsconfig.json ./

# Install development dependencies
RUN npm ci

# Copy source code and build
COPY src ./src
RUN npm run build

# ==============================================================================
# Production Stage
# ==============================================================================
FROM node:22-alpine AS runner

WORKDIR /app

# Install production dependencies only
COPY package*.json ./
RUN npm ci --only=production

# Copy compiled files from builder stage
COPY --from=builder /app/dist ./dist

# Set environment defaults
ENV NODE_ENV=production

# Expose standard port for SSE web server
EXPOSE 8080

# Execute server
ENTRYPOINT ["node", "dist/index.js"]
