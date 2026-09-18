FROM golang:1.27.1-alpine AS build

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /calculator-api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /calculator-api /calculator-api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/calculator-api"]
