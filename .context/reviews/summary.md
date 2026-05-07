# Review Summary

**Last updated:** 2026-05-06T22:41:37Z

## Findings
- Any change to `schema.sql` must be accompanied by a matching update to the schema smoke test (`db_test.go` or equivalent) so that missing tables/indexes are caught in CI. | count: 1 | status: tracked | sources: review-20260506-195853.md | stories: story-002
- Every new foreign key should have a negative-path test verifying constraint enforcement. | count: 1 | status: tracked | sources: review-20260506-195853.md | stories: story-002
- Every parser helper with ≥2 branching formats observed in real data should have direct unit tests. | count: 1 | status: tracked | sources: review-20260506-213858.md | stories: story-003
- Hierarchy linkers must have explicit negative-path tests for wrong-type parent references. | count: 1 | status: tracked | sources: review-20260506-213858.md | stories: story-003
- Keep schema DDL grouped by type (tables first, then indexes, then views, etc.). | count: 1 | status: tracked | sources: review-20260506-195853.md | stories: story-002
- Schema changes that add columns to existing tables require explicit migration handling beyond `CREATE TABLE IF NOT EXISTS`. | count: 1 | status: tracked | sources: review-20260506-195853.md | stories: story-002
- Struct fields that are written in one pass and never read in subsequent passes should be eliminated or inlined. | count: 1 | status: tracked | sources: review-20260506-213858.md | stories: story-003
- Synthetic/ fallback ID generation should be tested for every entity type that can trigger it, not just the most common one. | count: 1 | status: tracked | sources: review-20260506-213858.md | stories: story-003
