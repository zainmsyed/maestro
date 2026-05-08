import { describe, expect, it } from 'vitest';
import { buildCsv, escapeCsv, type CsvRow } from './csvExport';

describe('escapeCsv', () => {
  it('returns plain text unchanged when no special characters', () => {
    expect(escapeCsv('hello')).toBe('hello');
    expect(escapeCsv(123)).toBe('123');
  });

  it('wraps values containing commas in double quotes', () => {
    expect(escapeCsv('a, b, c')).toBe('"a, b, c"');
  });

  it('wraps values containing double quotes and escapes them', () => {
    expect(escapeCsv('say "hello"')).toBe('"say ""hello"""');
  });

  it('wraps values containing newlines', () => {
    expect(escapeCsv('line1\nline2')).toBe('"line1\nline2"');
  });

  it('returns empty string for null and undefined', () => {
    expect(escapeCsv(null)).toBe('');
    expect(escapeCsv(undefined)).toBe('');
  });
});

describe('buildCsv', () => {
  it('returns empty string for empty row array', () => {
    expect(buildCsv([])).toBe('');
  });

  it('builds a simple CSV with headers', () => {
    const rows: CsvRow[] = [
      { name: 'Alice', age: 30 },
      { name: 'Bob', age: 25 },
    ];
    expect(buildCsv(rows)).toBe('name,age\nAlice,30\nBob,25');
  });

  it('handles rows with commas and quotes', () => {
    const rows: CsvRow[] = [
      { title: 'Hello, World', note: 'say "hi"' },
    ];
    expect(buildCsv(rows)).toBe('title,note\n"Hello, World","say ""hi"""');
  });

  it('handles rows with null/undefined values', () => {
    const rows: CsvRow[] = [
      { a: 'x', b: null, c: undefined },
    ];
    expect(buildCsv(rows)).toBe('a,b,c\nx,,');
  });
});
