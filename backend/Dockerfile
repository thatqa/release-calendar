# Build
FROM golang:1.23-alpine AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/app ./cmd

# Run
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /app/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
CMD ["serve"]
