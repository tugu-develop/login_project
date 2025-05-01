FROM golang:latest

WORKDIR /app
COPY . .

RUN go mod download
EXPOSE 8080
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]