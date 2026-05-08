# File Index

.pi/extensions/vazir-context/helpers.ts — Context injection, init, plan, and consolidation extension
.pi/extensions/vazir-context/index.ts — Context injection, init, plan, and consolidation extension
.pi/extensions/vazir-live-reload.ts — Pi extension hot-reload watcher for development
.pi/extensions/vazir-tracker/chrome.ts — Change tracker, diff, fix, and reset extension
.pi/extensions/vazir-tracker/index.ts — Change tracker, diff, fix, and reset extension
.pi/extensions/vazir-tracker/vcs.ts — Change tracker, diff, fix, and reset extension
.pi/lib/vazir-helpers.ts — Shared story frontmatter, VCS detection, and date helpers
.pi/skills/vazir-base/SKILL.md — Vazir baseline skill instructions
AGENTS.md — Cross-framework project guidance and working notes
package-lock.json — package-lock.json configuration file
package.json — package.json configuration file

cmd/maestro/main.go — Application entry point: parses CLI flags, opens SQLite DB, wires repositories, API server, and serves embedded SPA
internal/config/config.go — CLI flag parsing for port and database path with validation
internal/db/db.go — SQLite database initialization, schema migration, and PRAGMA configuration
internal/db/schema.sql — DDL for epics, features, stories, sprints, audit logs, import reports, and indexes
internal/models/audit.go — DateAuditLog model for tracking committed date changes
internal/models/epic.go — Epic domain model with date fields and synthetic flag
internal/models/feature.go — Feature domain model with epic FK, dates, story points, and date source
internal/models/sprint.go — Sprint domain model with start/end dates, team, and capacity
internal/models/story.go — Story domain model with feature FK, dates, story points, and date source
internal/models/import_report.go — ImportReport and candidate models for CSV import summaries
internal/db/db_test.go — Schema verification tests: table and index existence
internal/repository/repository.go — Repository factory that wires all SQLite repositories
internal/repository/helpers.go — Shared date/time formatting and scanning helpers for SQLite
internal/repository/epic.go — SQLite epic repository with CRUD and list operations
internal/repository/feature.go — SQLite feature repository with CRUD, list, date patch, and epic reassignment
internal/repository/story.go — SQLite story repository with CRUD, list, date patch, and feature reassignment
internal/repository/sprint.go — SQLite sprint repository with CRUD and list operations
internal/repository/audit.go — SQLite audit log repository for date change tracking
internal/repository/import_report.go — SQLite import report persistence with JSON serialization
internal/repository/metrics.go — SQLite metrics queries for deadline hit rate, scope creep, slip, and orphans
internal/repository/repository_test.go — Repository CRUD round-trip tests for all entity types
internal/importer/csv.go — CSV import orchestrator: two-pass parse, hierarchy linking, synthetic orphan handling
internal/importer/dates.go — Date parsing with multiple format support and ambiguous format detection
internal/importer/hierarchy.go — Iteration path parsing, assigned-to extraction, and story points parsing
internal/importer/validation.go — CSV header normalization, required column validation, and title column detection
internal/importer/types.go — Importer state machine, raw entity types, and report finalization
internal/importer/types_helper.go — Work item type normalization (epic/feature/story aliases)
internal/importer/importer_test.go — Full CSV import integration tests with fixture data
internal/api/server.go — HTTP server setup, JSON response helpers, and date parsing utilities
internal/api/routes.go — Route registration and sub-resource dispatch for stories, features, epics, metrics, audit, import
internal/api/epics.go — Epic list/get/create handlers with nested feature/story tree responses
internal/api/features.go — Feature list/get/create handlers with date patch and epic reassignment
internal/api/stories.go — Story list/get/create handlers with date patch and feature reassignment
internal/api/audit.go — Audit log list handler with batched date_source enrichment
internal/api/metrics.go — Metrics aggregation handler for deadline/scope creep and slip/orphan endpoints
internal/api/import.go — CSV upload handler and last import report retrieval
internal/api/api_test.go — Comprehensive HTTP API integration tests for all endpoints

frontend/index.html — Vite SPA HTML entry point
frontend/package.json — Frontend dependencies: Svelte, Vite, Vitest, Testing Library
frontend/package-lock.json — Locked dependency versions
frontend/tsconfig.json — TypeScript configuration for Svelte project
frontend/vite.config.ts — Vite build config with dev proxy to Go backend and Vitest setup
frontend/embed.go — Go embed directive for bundling frontend/dist into the binary
frontend/src/main.ts — Svelte application bootstrap mounting App to #app
frontend/src/app.css — Design system CSS variables and global layout styles
frontend/src/App.svelte — Root shell: sidebar nav, topbar, content router, status bar
frontend/src/lib/api.ts — HTTP client wrapper, API types, and endpoint helpers
frontend/src/stores/project.ts — Svelte writable store for project state and import metadata
frontend/src/stores/view.ts — Svelte writable store for primary view and roadmap mode
frontend/src/components/DatePicker.svelte — Reusable date input with optional sprint boundary snap
frontend/src/components/DateAssignmentRow.svelte — Row component for missing-date assignment UI
frontend/src/components/DropZone.svelte — Drag-and-drop CSV file upload component
frontend/src/components/SprintPreview.svelte — Sprint cadence review table with manual/config toggle
frontend/src/components/StepBar.svelte — Onboarding progress indicator (import → sprints → confirm)
frontend/src/components/DataTable.svelte — Sortable grouped data table with inline date editing
frontend/src/components/FilterBar.svelte — Filter dropdowns and group-by selector for List view
frontend/src/components/GroupBySelect.svelte — Group-by dropdown (Epic/Sprint/Owner)
frontend/src/lib/csvExport.ts — CSV string building and browser download helper
frontend/src/lib/listView.ts — View-model helpers: flatten, filter, sort, group, export for List view
frontend/src/screens/Onboarding.svelte — Multi-step import flow: upload, sprint review, date assignment, confirm
frontend/src/screens/DateAssignment.svelte — Post-import screen for batch-assigning missing target dates
frontend/src/screens/ListView.svelte — List view screen with filters, grouping, sorting, and inline editing
frontend/src/test/setup.ts — Vitest test setup: jest-dom matchers and Svelte cleanup
frontend/src/components/DataTable.test.ts — DataTable component tests for sort, edit, reassign, states
frontend/src/components/FilterBar.test.ts — FilterBar component tests for filter selection and chips
frontend/src/lib/listView.test.ts — Unit tests for flatten, filter, group, sort, and export helpers
frontend/src/lib/csvExport.test.ts — Unit tests for CSV escaping and row building edge cases
