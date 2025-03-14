# Go-WebRTC

<!-- ![Go-WebRTC-Logo](./images/go-webrtc.png) -->
<img src="./images/go-webrtc.png" width="300" alt="Go-WebRTC-Logo" />


Golang WebRTC server implementation.

## How does it works?

- A user, will enter the platform, at root route, and can opt for stream or watch a stream. 

- When theres a user streaming its screen, other users, can then opt for watch the stream.

## How to run

Installing dependencies:
```bash
$ go mod tidy
```

### Building:
- Directly using Go:
```bash
$ go build -ldflags='-s -w' -trimpath -o ./bin/go-webrtc ./cmd
```

- Using Makefile:
```bash
$ make build
```

### Running:
```bash
./bin/go-webrtc
```

Specify the address and port to bind to through cli argument _-addr_. Like so:
```bash
./bin/go-webrtc -addr 192.168.2.44:3000
```

## _Tech Stack:_

### Backend:
- Golang HTTP Server (stdlib);
- gorilla/websocket;

### Frontend:
- HTML, Css and Javascript;
- MediaDevices Web-Api;
- WebSocket Web-Api;

### _License:_

This project is under MIT software license.