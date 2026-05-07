# System Rules

## Rules
- Follow existing project conventions.
- Write directly to real project files.
- Ask before changing ambiguous areas.

## Learned Rules
- Keep the application modular and easy to maintain. Split files by responsibility; no file should grow so large that it becomes hard to scan. Use language-specific best practices (Go: small packages, explicit interfaces, table-driven tests; TypeScript/Svelte: component-per-file, colocated types, minimal shared mutable state).

## Project Configuration
- **Default backend port:** `9000` (pinned across Go backend, Vite dev proxy, and documentation).
- **Vite dev server port:** `9001` (frontend dev).
- **Go API base URL:** `http://localhost:9000/api`.
