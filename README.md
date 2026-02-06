# AI-HOLOGRAM-NEW-BACKEND
[![en](https://img.shields.io/badge/lang-EN-red.svg)](README.md)
[![nl](https://img.shields.io/badge/lang-NL-blue.svg)](README.nl.md)

*Authors: Twan Wolthaus & Aiden van Wijnbergen*

Software to run background processes of the AI-Hologram. 
Uses include:
- Managing websocket communication
- Proxying requests to external AI services
- Handling 3D model generation and downloads

---
## Setup
To get started, you can clone the repository:
```sh
git clone https://github.com/technolab-admin/AI-hologram-new-backend.git
cd AI-hologram-new-backend
```

When you've cloned the repository, have a look at the `.env.example` file and create a copy of it named `.env`.
Place this file in the same directory and modify the environment variables where necessary.

### Docker
To run the application it's advisable to use Docker. To install Docker on your environment take a look
[here](https://www.docker.com/products/docker-desktop/).

After Docker is installed you can run the following command to start the application:
```sh
docker compose up --build
```
This will start the server and websocket-server on the ports you've set in your `.env` file.

### Run Tests
To run unit-tests on the WebSocket logic:
```sh
cd AI-hologram-new-backend
go test ./test/unit_test/
```

## Project Structure

- `cmd/server` – Application entry point
- `internal/http` – HTTP handlers and routing
- `internal/websockets` – WebSocket server and validation
- `internal/meshy` – Managing interaction with the Meshy API
- `internal/config` – Environment configuration



