# Stage 1: Build stage
FROM golang:1.21.3 AS builder

WORKDIR /backend/

COPY ./cmd /backend/cmd
COPY ./internal /backend/internal
COPY ./dev /backend/dev

RUN go mod init medods
RUN go mod tidy

RUN go build -o /backend/build ./cmd/medods/

# Stage 2: Final stage
FROM ubuntu:22.04

WORKDIR /backend

COPY --from=builder /backend/build /backend/build
COPY --from=builder /backend/dev/.env /backend/dev/.env
COPY --from=builder /backend/internal/app/migrations /backend/internal/app/migrations

CMD [ "/backend/build" ]
