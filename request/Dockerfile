FROM golang:alpine
LABEL authors="MasterHahn04"

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main /app/

ENTRYPOINT ["/app/main"]