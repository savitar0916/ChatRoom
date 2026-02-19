FROM golang:1.23 AS builder
WORKDIR /app

COPY chat/go.mod chat/go.sum ./
RUN go mod download

COPY chat/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o chat ./main.go

# ===== Runner =====
FROM gcr.io/distroless/base-debian12 AS runner
WORKDIR /app

COPY --from=builder /app/chat .

# 把 config.yaml 打包進去
COPY chat/config.yaml ./config.yaml

# 把前端頁面打包進去
COPY chat/public ./public

EXPOSE 8000
CMD ["./chat", "serve"]