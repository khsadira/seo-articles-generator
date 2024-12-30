
FROM golang:1.23

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o .

EXPOSE 8080

# Run
CMD ["./seo-articles-generator"]