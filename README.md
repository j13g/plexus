# plexus
Plexus is a WIP experimental microservice service architecture using NATS message queues.

## Overview
This is very much a work in progress, some things which can be found in this repo include:

### pkg/postbox
The core of the message passing framework, handles sending and receving envelopes via Inbox and Outbox, serialization,
request/reply patterns, etc.

### pkg/config
A config singleton to automatically load application configuration from the environment.
Enables a quasi-dependency-injection style as it contains references to the main components required in the application:
- Inbox & Outbox, the interfaces for message handling
- Cron, the interface for scheduling background tasks
- Etc

### pkg/mainutil
An easy way to start a new plexus binary, handles all the application setup and CLI entrypoints.

## Example Agent
WIP Example of an agent that listens for command messages, and heartbeats information about itself.
