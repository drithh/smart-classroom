# Use the official Golang base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local code into the container at the working directory
COPY . .

# Download and install any dependencies
RUN go mod download

# Build the application
RUN go build -o main .

# Expose a port (replace 8080 with your application's port)
EXPOSE 3000

# Command to run the application
CMD ["./main"]
