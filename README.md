# Arrogance Admin

A TUI (Terminal User Interface) admin application built with [Charm](https://github.com/charmbracelet/bubbletea) and Firebase.

## Features

- Terminal-based user interface
- Firebase authentication and Firestore database integration
- User management capabilities

## Prerequisites

- Go 1.23.3 or later
- Firebase project with Authentication and Firestore enabled

## Setup Firebase

1. Go to the [Firebase Console](https://console.firebase.google.com/)
2. Create a new project or use an existing one
3. Enable Authentication and Firestore services
4. Generate a new private key for service account:
   - Go to Project Settings > Service accounts
   - Click "Generate new private key"
   - Save the JSON file
5. Place the service account JSON file in one of these locations:
   - `~/.config/arrogance/service-account.json` (recommended)
   - In the project directory as `service-account.json`
   - Or set the `FIREBASE_SERVICE_ACCOUNT` environment variable to the path of your service account file

## Running the Application

```bash
# Run the application
go run .
```

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Project Structure

- `main.go`: Main application entry point and TUI logic
- `firebase/`: Firebase integration
  - `firebase.go`: Firebase initialization
  - `auth.go`: Authentication service
  - `firestore.go`: Firestore database service

## License

MIT
