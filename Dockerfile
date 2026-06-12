# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS build
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY go.mod ./
RUN go mod download
COPY cmd/ cmd/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# ── Runtime stage (scratch) ───────────────────────────────────────────────────
FROM scratch
WORKDIR /app

# CA certs — required for TLS calls to Axiom, Supabase etc. from scratch image
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=build /app/server .

# Copy assets the server reads at runtime from its working directory
COPY navigation_config.json .
COPY favicon.png .
COPY shared/ shared/
COPY 1_Real_Unknown/ 1_Real_Unknown/
COPY 2_Environment/ 2_Environment/
COPY 3_Simulation/ 3_Simulation/
COPY 4_Formula/ 4_Formula/
COPY 5_Symbols/ 5_Symbols/
COPY 6_Semblance/ 6_Semblance/
COPY 7_Testing_Known/ 7_Testing_Known/
COPY markdown_renderer.html .
COPY robots.txt .
COPY sitemap.xml .

EXPOSE 8080
CMD ["/app/server"]
