FROM golang:1.23.0-alpine

LABEL maintainer="jononl3adama@gmail.com" \
    version="1.0" \
    description="A simple forum application " 

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o myapp .

EXPOSE 8080

CMD ["./myapp"]