# syntax=docker/dockerfile:1
FROM node:22-alpine AS frontend-build

WORKDIR /src/web
COPY web/package.json web/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile --allow-build=esbuild --allow-build=@clerk/shared
COPY web/ ./
RUN pnpm build

FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /service ./cmd

FROM alpine:3.22
RUN addgroup -S app && adduser -S app -G app
USER app
COPY --from=build /service /service
COPY --from=frontend-build --chown=app:app /src/web/dist /web
EXPOSE 8080
ENTRYPOINT ["/service"]
