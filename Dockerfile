FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .
RUN templ generate ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o htmx-server ./cmd/htmx-server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/htmx-server .

# Expose port to the outside world
EXPOSE 5000

#Command to run the executable
CMD [ "./htmx-server" ]
