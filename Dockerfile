FROM golang:1.15.3-alpine3.12 as build_env

# Copy the source from the current directory to the Working Directory inside the container
WORKDIR /app

ENV GOPRIVATE=github.com/Checkmarx/*

#Copy go mod and sum files
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/kics cmd/console/main.go 

#runtime image
FROM scratch

COPY --from=build_env /app/bin/kics /app/bin/kics
COPY --from=build_env /app/assets/ /app/bin/assets/

WORKDIR /app/bin

# Command to run the executable
ENTRYPOINT ["/app/bin/kics"]
