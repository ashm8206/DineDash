### DineDash

```
---
Author: Aishwarya Mali
Date:   09/16/24
---
```

A Microservice-Based POS System for Neighborhood Diners!

![Architecture](https://github.com/ashm8206/DineDash/blob/main/ui/github-assets/DineDash.png)

**Clean Layered Architecture**: The system is built with a modular, layered architecture where multiple microservices communicate via gRPC and asynchronously with RabbitMQ. Each service is separated into distinct Transport, Service, and Storage layers to ensure scalability and maintainability.

**Service Discovery and Integration**: Service discovery is seamlessly handled with Consul, enabling dynamic service registration and configuration. Integrated best-in-class third-party APIs, such as Stripe, following industry best practices for a composable and future-proof architecture.

**Enhanced Observability and Telemetry**: Leveraging distributed tracing tools like Jaeger, and provided comprehensive observability and telemetry, empowering proactive monitoring and quick identification of issues in real-time.

**Reliability and Fault Tolerance**: Built-in reliability features, such as dead-letter queues and robust logging through middleware, ensure system stability and smooth error handling for unexpected failures.
