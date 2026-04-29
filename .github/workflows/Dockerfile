# Builder stage
FROM ubuntu:22.04 as builder

# Set environment variable to suppress interactive prompts
ENV DEBIAN_FRONTEND=noninteractive

# Install tzdata and other dependencies
RUN apt-get update --fix-missing && \
    apt-get install -y tzdata wget curl build-essential ca-certificates && \
    ln -fs /usr/share/zoneinfo/Asia/Makassar /etc/localtime && \
    echo "Asia/Makassar" > /etc/timezone && \
    dpkg-reconfigure -f noninteractive tzdata

# Download and install specific version of Go
RUN wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz && \
    rm go1.22.2.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"

# Set the working directory
WORKDIR /app

# Copy application files
COPY . .

# Build the application
RUN /usr/local/go/bin/go build -o application cmd/main.go

# Final stage
FROM ubuntu:22.04

# Set environment variable to suppress interactive prompts
ENV DEBIAN_FRONTEND=noninteractive

# Install tzdata
RUN apt-get update --fix-missing && \
    apt-get install -y tzdata ca-certificates && \
    ln -fs /usr/share/zoneinfo/Asia/Makassar /etc/localtime && \
    echo "Asia/Makassar" > /etc/timezone && \
    dpkg-reconfigure -f noninteractive tzdata

# Set the timezone environment variable
ENV TZ=Asia/Makassar

# Set the working directory
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/application /app/application

# Run the application
CMD ["./application"]
