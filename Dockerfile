FROM golang:1.24.4-alpine AS build
RUN apk add --no-cache alpine-sdk npm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN npm install && \
    curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64-musl -o tailwindcss && \
    chmod +x tailwindcss && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    go install github.com/a-h/templ/cmd/templ@latest 

RUN sqlc generate && \
    templ generate && \
    ./tailwindcss -i tailwind.css -o internal/server/assets/css/output.css -m

RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]


