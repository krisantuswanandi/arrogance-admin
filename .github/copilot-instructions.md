# GitHub Copilot Instructions for Arrogance Admin

## Project Overview

This is a terminal-based admin interface for managing Firebase users. It's built with TypeScript, Bun runtime, React with Ink for terminal UI, and Firebase Admin SDK.

## Tech Stack

- **Runtime**: Bun
- **Language**: TypeScript
- **UI Framework**: React with Ink (terminal-based UI)
- **Backend Services**: Firebase Admin SDK (Auth, Firestore)
- **Package Manager**: Bun

## Code Style & Conventions

### TypeScript

- Use strict TypeScript with proper typing
- Import types with `type` keyword: `import type { UserRecord } from "firebase-admin/auth"`
- Use modern ES modules syntax
- Prefer named exports over default exports where appropriate

### React/Ink Components

- Use functional components with hooks
- Follow React best practices for state management
- Use Ink's built-in components (Box, Text, etc.) for terminal UI
- Handle keyboard input with `useInput` hook
- Use `useApp` hook for app lifecycle management (e.g., exit)

### File Organization

- Components go in `src/components/`
- Pages go in `src/pages/`
- Utilities go in `src/utils/`
- Types are defined in `src/utils/types.ts`
- Main entry point is `index.ts`

### Firebase Integration

- Use Firebase Admin SDK for server-side operations
- Initialize Firebase in `src/utils/firebase.ts`
- Export helper functions for common Firebase operations
- Handle authentication and Firestore operations through utility functions

## Development Guidelines

### When adding new features:

1. Create reusable components in the `components/` directory
2. Add new pages in the `pages/` directory
3. Define new types in `src/utils/types.ts`
4. Add Firebase utility functions in `src/utils/firebase.ts`
5. Use proper error handling for Firebase operations
6. Ensure terminal UI is responsive and user-friendly

### When working with Firebase:

- Always handle async operations properly
- Use try-catch blocks for Firebase API calls
- Provide meaningful error messages for users
- Respect Firebase quotas and limits
- Use pagination for large datasets (users list, etc.)

### Terminal UI Best Practices:

- Use consistent spacing and layout with Ink's Box component
- Provide clear navigation instructions
- Handle keyboard shortcuts consistently
- Use colors appropriately for status indicators
- Ensure text is readable in terminal environments

### Performance Considerations:

- Implement pagination for large user lists
- Use React's useMemo and useCallback where appropriate
- Avoid unnecessary re-renders in terminal UI
- Cache Firebase data when appropriate

## Key Components Structure

### App.tsx

- Main application component with routing logic
- Handles global keyboard shortcuts (q to quit)
- Manages page state and navigation

### Pages

- `Users.tsx`: Lists Firebase users with pagination
- `User.tsx`: Shows individual user details
- Each page should handle its own state and props

### Utils

- `firebase.ts`: Firebase initialization and helper functions
- `types.ts`: TypeScript type definitions

## Commands

- `bun run dev`: Start the development server
- `bun install`: Install dependencies

## Environment Setup

- Requires `service-account.json` for Firebase Admin SDK
- Uses Bun as the JavaScript runtime
- TypeScript configuration in `tsconfig.json`

When suggesting code changes, always consider the terminal-based nature of the UI and ensure compatibility with the Ink framework and Bun runtime.
