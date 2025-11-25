FROM golang:1.25.4-alpine AS build
RUN apk add --no-cache curl libstdc++ libgcc alpine-sdk npm curl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN npm install && \
    curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64-musl -o tailwindcss && \
    chmod +x tailwindcss && \
    curl -sL https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz migrate && \
    chmod +x migrate && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    go install github.com/a-h/templ/cmd/templ@latest 

RUN sqlc generate && \
    templ generate && \
    ./tailwindcss -i tailwind.css -o internal/server/assets/css/output.css -m

RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /app/migrate /app/migrate
# Required for the following command
# migrate -source file://migrations -database sqlite3://checkpoint.db up
COPY --from=build /app/migrations /app/migrations
EXPOSE ${PORT}

CMD ["migrate -source file:///app/migrations -database sqlite3://$(DB_ADDRESS) goto $(DB_MIGRATION_VERSION) && ./main"]



