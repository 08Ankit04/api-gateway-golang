# API Gateway Golang

## Overview

This project implements an API Gateway using Go, providing features such as routing, authentication, rate limiting, and logging. It uses gRPC for communication with microservices, Envoy as a reverse proxy, Redis for rate limiting, and JWT for authentication. The project is containerized with Docker and orchestrated using Kubernetes.

## Features

- **Routing:** Routes requests to appropriate microservices.
- **Authentication:** Uses JWT tokens for securing endpoints.
- **Rate Limiting:** Prevents abuse by limiting the number of requests.
- **Logging:** Logs all requests for monitoring and debugging.
- **Monitoring:** Integrates with Prometheus for metrics collection.

## Technologies

- **Go**: Programming language for the API Gateway.
- **gRPC**: Communication protocol for microservices.
- **Envoy**: Reverse proxy and load balancer.
- **Redis**: Data store for rate limiting.
- **JWT**: JSON Web Tokens for authentication.
- **Prometheus**: Monitoring and metrics collection.
- **Docker**: Containerization.
- **Kubernetes**: Container orchestration.

## Getting Started

### Prerequisites

- Go 1.16+
- Docker
- Kubernetes
- Redis
- Envoy
- Prometheus

### Installation

1. **Clone the Repository:**
    ```sh
    git clone https://github.com/api-gateway-golang/go-api-gateway.git
    cd github.com/api-gateway-golang
    ```

2. **Set Up Environment Variables:**
    Create a `.env` file in the root directory and add the following variables:
    ```env
    JWT_SECRET=your_jwt_secret
    REDIS_ADDR=localhost:6379
    ```

3. **Install Dependencies:**
    ```sh
    go mod tidy
    ```

4. **Build the Project:**
    ```sh
    go build -o api-gateway .
    ```

### Running the Application

1. **Start Redis:**
    ```sh
    docker run --name redis -p 6379:6379 -d redis
    ```

2. **Start Envoy:**
    Refer to the Envoy configuration in the `envoy.yaml` file and start Envoy:
    ```sh
    envoy -c envoy.yaml
    ```

3. **Run the API Gateway:**
    ```sh
    ./api-gateway
    ```

4. **Deploy with Kubernetes:**
    Apply the Kubernetes manifests:
    ```sh
    kubectl apply -f k8s/
    ```

### Configuration

The main configuration files include:

- **config.yaml:** Configuration for the API Gateway.
- **envoy.yaml:** Configuration for Envoy proxy.
- **k8s/**: Kubernetes manifests for deployment.

### Usage

- **Authentication:** Include a JWT token in the `Authorization` header of your requests.
- **Rate Limiting:** The gateway limits the number of requests based on the configured policy.
- **Logging:** Logs are output to the console and can be viewed for debugging.

### Project Structure

```

api-gateway-golang/
├── cmd/
│ └── main.go # Entry point of the application
├── internal/
│ ├── auth/ # Authentication logic
│ ├── rate_limit/ # Rate limiting logic
│ ├── router/ # Request routing logic
│ └── logger/ # Logging logic
├── k8s/ # Kubernetes manifests
├── config.yaml # Configuration file
├── envoy.yaml # Envoy configuration
└── README.md # Project documentation

```

### Contributing

Contributions are welcome! Please open an issue or submit a pull request.

### License

This project is licensed under the MIT License.

### Acknowledgments

- Inspired by various open-source projects and Go community contributions.