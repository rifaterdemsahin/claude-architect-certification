# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY cmd/ cmd/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# ── Runtime stage (scratch) ───────────────────────────────────────────────────
FROM scratch
WORKDIR /app

# Copy binary
COPY --from=build /app/server .

# Copy assets the server reads at runtime from its working directory
COPY navigation_config.json .
COPY templates/ templates/
COPY shared/ shared/
COPY 5_Symbols/ 5_Symbols/
COPY markdown_renderer.html .
COPY robots.txt .
COPY sitemap.xml .

EXPOSE 8080
CMD ["/app/server"]
