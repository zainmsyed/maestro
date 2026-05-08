import { fireEvent, render, screen, within } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import DataTable from './DataTable.svelte';
import type { ListGroup } from '../lib/listView';

const baseGroups: ListGroup[] = [
  {
    key: 'E-1',
    label: 'Alpha Epic',
    rows: [
      {
        id: 'E-1',
        type: 'epic',
        title: 'Alpha Epic',
        owner: 'Alice',
        sprint: 'Sprint 1',
        originalDate: '2026-05-01',
        committedDate: '2026-05-03',
        actualDate: '',
        slipEvents: 0,
        status: 'Active',
        health: 'blue',
        healthLabel: 'Neutral',
        epicId: 'E-1',
        epicTitle: 'Alpha Epic',
        isSynthetic: false,
      },
      {
        id: 'F-1',
        type: 'feature',
        title: 'Feature One',
        owner: 'Bob',
        sprint: 'Sprint 2',
        originalDate: '2026-05-02',
        committedDate: '',
        actualDate: '',
        slipEvents: 1,
        status: 'Active',
        health: 'amber',
        healthLabel: 'Has slip events',
        epicId: 'E-1',
        epicTitle: 'Alpha Epic',
        isSynthetic: false,
      },
    ],
  },
  {
    key: 'epic-unassigned',
    label: 'Unassigned Epic',
    rows: [
      {
        id: 'feature-unassigned',
        type: 'feature',
        title: 'Unassigned Feature',
        owner: '',
        sprint: '',
        originalDate: '',
        committedDate: '',
        actualDate: '',
        slipEvents: 0,
        status: 'Active',
        health: 'blue',
        healthLabel: 'Neutral',
        epicId: 'epic-unassigned',
        epicTitle: 'Unassigned Epic',
        isSynthetic: true,
      },
    ],
  },
];

describe('DataTable states', () => {
  it('renders loading state when loading is true', () => {
    render(DataTable, {
      groups: baseGroups,
      sortKey: 'title',
      sortDirection: 'asc',
      loading: true,
      onSort: vi.fn(),
      onSaveDate: vi.fn(),
      onReassignFeature: vi.fn(),
      epicOptions: [],
    });

    expect(screen.getByText(/loading list view/i)).toBeInTheDocument();
  });

  it('renders empty state when no groups match', () => {
    render(DataTable, {
      groups: [],
      sortKey: 'title',
      sortDirection: 'asc',
      loading: false,
      onSort: vi.fn(),
      onSaveDate: vi.fn(),
      onReassignFeature: vi.fn(),
      epicOptions: [],
    });

    expect(screen.getByText(/no roadmap items match/i)).toBeInTheDocument();
  });
});

describe('DataTable interactions', () => {
  it('calls onSort when a header is clicked', async () => {
    const onSort = vi.fn();
    render(DataTable, {
      groups: baseGroups,
      sortKey: 'title',
      sortDirection: 'asc',
      onSort,
      onSaveDate: vi.fn(),
      onReassignFeature: vi.fn(),
      epicOptions: [],
    });

    await fireEvent.click(screen.getByRole('button', { name: /owner/i }));
    expect(onSort).toHaveBeenCalledWith('owner');
  });

  it('supports inline feature date editing and save callback', async () => {
    const onSaveDate = vi.fn();
    const { container } = render(DataTable, {
      groups: baseGroups,
      sortKey: 'title',
      sortDirection: 'asc',
      onSort: vi.fn(),
      onSaveDate,
      onReassignFeature: vi.fn(),
      epicOptions: [],
    });

    const featureTitle = screen.getByText('Feature One');
    const featureRow = featureTitle.closest('tr') as HTMLTableRowElement;
    const featureButton = within(featureRow).getByRole('button', { name: /set date/i });
    await fireEvent.click(featureButton);
    const input = screen.getByLabelText('Committed date for Feature One');
    await fireEvent.input(input, { target: { value: '2026-05-09' } });
    await fireEvent.click(screen.getByRole('button', { name: /^save$/i }));

    expect(onSaveDate).toHaveBeenCalledWith(expect.objectContaining({ id: 'F-1' }), '2026-05-09');
  });

  it('calls reassignment callback for synthetic feature rows', async () => {
    const onReassignFeature = vi.fn();
    const { container } = render(DataTable, {
      groups: baseGroups,
      sortKey: 'title',
      sortDirection: 'asc',
      onSort: vi.fn(),
      onSaveDate: vi.fn(),
      onReassignFeature,
      epicOptions: [{ id: 'E-1', title: 'Alpha Epic' }],
    });

    const reassignSelect = container.querySelector('.reassign-label select') as HTMLSelectElement;
    await fireEvent.change(reassignSelect, { target: { value: 'E-1' } });

    expect(onReassignFeature).toHaveBeenCalledWith(
      expect.objectContaining({ id: 'feature-unassigned' }),
      'E-1',
    );
  });
});
