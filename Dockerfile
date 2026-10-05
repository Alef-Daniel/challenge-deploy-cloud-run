FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

FROM gcr.io/distroless/static-debian12

COPY --from=build /server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
