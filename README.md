# Kaicode Backend

This is the backend service for Kaicode, a platform for collaborative coding and code execution.

## Overview

The backend consists of the following main components:

*   **Job Worker:** This component picks up code execution jobs from a queue (SQS) and executes them in isolated Docker containers. It ensures secure and sandboxed execution of user-submitted code.
*   **Yjs WebSocket Server:** This server provides real-time synchronization for collaborative coding using the Yjs framework. It allows multiple users to edit the same code simultaneously.
*   **Sync WebSocket Server:** This server handles synchronization of other data and application state between clients.

## Key Features

*   **Sandboxed Code Execution:** Executes code in isolated Docker containers for security.
*   **Real-time Collaboration:** Enables multiple users to code together in real-time.
*   **Scalable Architecture:** Designed to handle a large number of concurrent users and code execution jobs.

## Architecture

The backend uses a queue-based architecture with SQS for managing code execution jobs. The job worker processes these jobs by running them in Docker containers. The WebSocket servers provide real-time communication between clients.

## Getting Started

To get started with the Kaicode backend, please refer to the documentation for setting up the environment and running the services.

## License

[MIT License](LICENSE)

## Deployment steps
- Build it on the server using `go build . -o kaicode`
- Make sure it runs
- After having several issues with Nginx not forwading the headers correctly, and a lot of issues, I just setup caddy with three lines. 
`sudo nano /etc/caddy/Caddyfile`
```
kc-be.yashk.dev {
        reverse_proxy localhost:1234
}
```
 and it handled everything well, no seperate commands for ssl, no seperate stuff for websocket, just this.
- as this is a high level program that runs and orchestrate docker containers itself, its recommended you run this as a daemon, using systemd, i'll give the following example config that I did on my AWS Lightsail instance
`sudo nano /etc/systemd/system/kaicode.service`

```
[Unit]
Description=Kaicode backend
After=network.target

[Service]
WorkingDirectory=/home/admin/kaicode-backend
ExecStart=/home/admin/kaicode-backend/kaicode
Restart=on-failure
RestartSec=5
EnvironmentFile=/home/admin/kaicode-backend/.env
User=admin
Group=admin

[Install]
WantedBy=multi-user.target
```