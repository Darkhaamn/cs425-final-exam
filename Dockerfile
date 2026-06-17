# ---- build stage ----
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/server .

# ---- run stage ----
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/server .
EXPOSE 8080
ENTRYPOINT ["/app/server"]
