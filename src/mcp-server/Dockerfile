# ==============================================================================
# Build Stage
# ==============================================================================
FROM node:20-alpine AS builder

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
FROM node:20-alpine AS runner

WORKDIR /app

# Install production dependencies only
COPY package*.json ./
RUN npm ci --only=production

# Copy compiled files from builder stage
COPY --from=builder /app/build ./build

# Set environment defaults
ENV NODE_ENV=production
ENV MCP_LOG_LEVEL=info

# Use stdio transport as default
ENTRYPOINT ["node", "build/index.js"]
