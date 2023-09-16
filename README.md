# Chat Broadcaster

The Chat Broadcaster is a specialized service designed to receive messages from various chat services via RabbitMQ, then broadcast these messages to users who have subscribed to specific chat rooms. This repository contains the code and configuration files for this service.

## Endpoints
| Endpoint                     | Method | Protected | Description         |
|------------------------------| ------ |-----------|---------------------|
| `/api/v1/subscribe/{roomId}` | GET    | YES       | Subscribe to a room |
| `/api/v1/healthz`            | GET    | NO        | Health check        |

## Related repositories

- [Chat App](https://github.com/sesaquecruz/react-chat-app)
- [Chat API](https://github.com/sesaquecruz/go-chat-api)
- [Chat Infra](https://github.com/sesaquecruz/k8s-chat-infra)
- [Chat Broadcaster Docker Hub](https://hub.docker.com/r/sesaquecruz/go-chat-broadcaster/tags)

## Tech Stack

- [Go](https://go.dev)
- [Gin](https://gin-gonic.com)
- [RabbitMQ](https://www.rabbitmq.com)
- [Redis](https://redis.io)

## Contributing

Contributions are welcome! If you find a bug or would like to suggest an enhancement, please make a fork, create a new branch with the bugfix or feature, and submit a pull request.

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE) file for more information.
