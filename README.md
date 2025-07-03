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
