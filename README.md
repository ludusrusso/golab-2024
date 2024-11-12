# NATS JetStream Example

This example demonstrates how to use the NATS JetStream API to publish and consume messages.

## Prerequisites

#### Run the Example cluster

```bash
nats-server -js
```

## Generate Random Orders

```bash
go run main.go generate 20
```

## Process Orders

```bash
go run main.go process
```

## Compute Stats

```bash
go run main.go stats EU
```

## My Contacts

- [LinkedIn](https://www.linkedin.com/in/ludusrusso/)
- [BlueSky](https://bsky.app/profile/ludusrusso.bsky.social)
- [ludusrusso.dev](https://ludusrusso.dev)
- [GitHub](https://github.com/ludusrusso)
