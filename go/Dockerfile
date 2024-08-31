FROM golang:1.23-alpine

WORKDIR /app/

# COPY go.mod, go.sum and download the dependencies
COPY go.* ./
RUN go mod download

# COPY All things inside the project and build
COPY . .
RUN go build .

EXPOSE 7042
ENTRYPOINT [ "/app/test" ]