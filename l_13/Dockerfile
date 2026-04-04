FROM golang:1.25-alpine AS build
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod tidy
ENV CGO_ENABLED=0
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app/main.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=build /app/main .
EXPOSE 8080
CMD ["./main"]