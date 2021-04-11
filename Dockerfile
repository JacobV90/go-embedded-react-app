# Build the react app for production
FROM node:12.13.1 AS ui-builder
WORKDIR /app
COPY webapp/ .
RUN npm install && npm run build

# Build the Go server application
FROM golang:1.16.3 AS server-builder
WORKDIR /src
COPY go.mod .
COPY server.go .
# Copy over production build from ui-bulder stage
COPY --from=ui-builder /app/build /src/webapp/build
# Install Go module dependencies
RUN go get
# Build Go executable for linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o run_webapp 

# Run Go application in linux container
FROM alpine:latest
WORKDIR /root
# Copy over executable built by Go from server-builder stage
COPY --from=server-builder /src/run_webapp .
# Run executable
CMD [ "./run_webapp" ]
