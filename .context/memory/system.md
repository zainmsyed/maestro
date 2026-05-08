# System Rules

## Rules
- Follow existing project conventions.
- Write directly to real project files.
- Ask before changing ambiguous areas.

## Learned Rules
- Keep the application modular and easy to maintain. Split files by responsibility; no file should grow so large that it becomes hard to scan. Use language-specific best practices (Go: small packages, explicit interfaces, table-driven tests; TypeScript/Svelte: component-per-file, colocated types, minimal shared mutable state).
- **Svelte reactivity — pass reactive dependencies explicitly.** When a `$:` reactive statement calls a function, Svelte's compiler only tracks variables that appear *directly in the `$:` line*, not those referenced inside the function body. If the function reads reactive state (stores, props, local variables), pass that state as explicit arguments. Failing to do so causes silent stale UI after async data loads.

## Project Configuration
- **Default backend port:** `9000` (pinned across Go backend, Vite dev proxy, and documentation).
- **Vite dev server port:** `9001` (frontend dev).
- **Go API base URL:** `http://localhost:9000/api`.
