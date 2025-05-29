# net-cat TCP Chat Server

This project implements a multi-room TCP chat server in Go. The core logic is in `server.go`, which manages chat rooms, client connections, and command handling.

## Features

- **Multiple Chat Rooms:** Users can join or create rooms with `/join ROOM_NAME`.
- **Nicknames:** Users set/change their nickname with `/nick NAME`. Name changes are broadcast to the room.
- **Room Listing:** `/rooms` lists all available rooms.
- **Message History:** Messages in the "general" room are saved and replayed to new joiners.
- **Help Command:** `/help` lists available commands.
- **Graceful Client Exit:** `/quit` disconnects the user and notifies the room.
- **Logging:** Server logs client connections and disconnections.

## Usage

Start the server:
```
go run main.go
```
or
```
go build && ./net-cat $PORT
```

## Commands

- `/nick NAME` — Set your nickname (max 20 chars)
- `/join ROOM_NAME` — Join or create a room (max 30 chars)
- `/rooms` — List all rooms
- `/msg MESSAGE` — Send a message to the current room
- `/quit` — Leave the chat
- `/help` — Show available commands

## File Structure

- `server.go` — Main server logic (rooms, clients, commands)
- `client.go` — Client connection and input handling
- `room.go` — Room struct and broadcast logic
- `command.go` — Command types and struct
- `utils.go` — Utility functions (port, shutdown, etc.)
- `main.go` — Entry point

## Notes

- When joining the "general" room, previous messages are shown.
- Nickname and room name limits are enforced.
- Server logs are printed to the console.

---