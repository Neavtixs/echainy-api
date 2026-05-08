FROM golang:1.25.1-alpine AS builder

WORKDIR /app/
COPY . .
RUN go build -o /app/main .

# --- RUNNER ---
FROM alpine:3.19
WORKDIR /app/

COPY --from=builder /app/main /app/

CMD /app/main