# Conventions Documentation

This document outlines the conventions used in the Digital Bank project for Git branching, commit messages, package naming, diagram creation, and logging practices. Adhering to these conventions ensures consistency and clarity across the codebase and documentation. Any questions or suggestions for improvements to these conventions can be raised in the meeting.

## Git

### Branching Conventions

- Use the following branch naming conventions:
  - `main`: The main production branch.
  - `<service>-main`: The main branch for each microservice (e.g., `auth-service-main`).
  - `<service>-dev`: The development branch for each microservice (e.g., `auth-service-dev`).
  - `feature/<service>/<short-description>`: Feature branches for new features (e.g., `feature/auth-service/jwt-auth`).
  - `bugfix/<service>/<short-description>`: Bugfix branches for fixing bugs (e.g., `bugfix/auth-service/token-expiry`).
  - `hotfix/<service>/<short-description>`: Hotfix branches for urgent fixes in production (e.g., `hotfix/auth-service/login-issue`).
  - `release/<service>/<version>`: Release branches for preparing new versions (e.g., `release/auth-service/v1.2.0`).

### Commit Message Conventions

- Use the following format for commit messages:

  ```
  <type>(<scope>): <subject>

  <body>

  <footer>
  ```

- **Type**: The type of change being made. Common types include:
  - `feat`: A new feature
  - `fix`: A bug fix
  - `refactor`: Code refactoring without changing functionality
  - `chore`: Maintenance tasks that do not affect the application code. (build process, dependencies, etc.)
  - `style`: Code style changes (formatting, missing semicolons, etc.)
  - `docs`: Documentation changes
  - `test`: Adding or updating tests
- **Scope**: The area of the codebase affected (e.g., `auth-service`, `http`, `usecases`, etc.)
- **Subject**: A brief description of the change (max 50 characters).
- **Body**: A more detailed explanation of the change, if necessary (wrap at 72 characters).
- **Footer**: Any relevant issue references or breaking changes.
- Use the imperative mood in the subject line (e.g., "Add feature" instead of "Added feature" or "Adds feature").
- Limit the subject line to 50 characters and the body to 72 characters per line.
- Separate the subject from the body with a blank line.
- Reference issues and pull requests in the footer using keywords like "Closes #123" or "Fixes #456".
- For breaking changes, include a "BREAKING CHANGE:" section in the footer with a description of the change and its impact.
- Example commit message:

  ```
  feat(auth-service): implement JWT authentication middleware

  Add a new middleware to handle JWT authentication for incoming requests.
  This middleware verifies the token and extracts user information.

  Closes #42
  ```

## Package Naming Conventions

- Use lowercase letters for package names.
- Use descriptive names that reflect the package's purpose (e.g., `http`, `middleware`, `usecases`, `jwt`, `model`).
- Avoid using underscores or mixed case in package names.
- Ensure package names are singular (e.g., `handler` instead of `handlers`).
- Align package names with their directory structure for clarity.
- Example package names:
  - `http`: For HTTP handlers and related code.
  - `middleware`: For HTTP middleware components.
  - `usecase`: For business logic and use case implementations.
  - `jwt`: For JWT-related functionality.
  - `model`: For data models and entities.

## Diagram Conventions

- Use PlantUML for creating diagrams.
- Follow a consistent style for all diagrams, including colors, fonts, and shapes.
- Clearly label all components, containers, and relationships in the diagrams.
- Use descriptive names for components and containers that reflect their purpose.
- Include a brief description of each component or container in the diagram.

## Logging Conventions

- Use a structured logging library (e.g., Zerolog, Zap) for consistent log formatting.
- Include relevant context in log messages (e.g., request IDs, user IDs).
- Use appropriate log levels (e.g., Info, Debug, Error) based on the severity of the event.
- Ensure log messages are clear and informative, avoiding ambiguity.
- Example log message structure:

  ```json
  {
    "timestamp": "2025-11-03T08:15:30Z",
    "level": "INFO",
    "service": "ledger-service",
    "instance_id": "ledger-5f4a3b7c9d",
    "trace_id": "abc123xyz",
    "user_id": "12345",
    "action": "transfer_request",
    "status": "success",
    "details": {
      "amount": 10.0,
      "currency": "USDT",
      "to_user_id": "67890"
    }
  }
  ```

- Avoid logging sensitive information (e.g., passwords, tokens).
