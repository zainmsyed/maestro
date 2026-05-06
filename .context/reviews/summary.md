# Review Summary

**Last updated:** 2026-05-06T20:05:38Z

## Findings
- Any change to `schema.sql` must be accompanied by a matching update to the schema smoke test (`db_test.go` or equivalent) so that missing tables/indexes are caught in CI. | count: 1 | status: tracked | sources: review-20260506-195853.md | stories: story-002
- Every new foreign key should have a negative-path test verifying constraint enforcement. | count: 1 | status: tracked | sources: review-20260506-195853.md | stories: story-002
- Keep schema DDL grouped by type (tables first, then indexes, then views, etc.). | count: 1 | status: tracked | sources: review-20260506-195853.md | stories: story-002
- Schema changes that add columns to existing tables require explicit migration handling beyond `CREATE TABLE IF NOT EXISTS`. | count: 1 | status: tracked | sources: review-20260506-195853.md | stories: story-002
