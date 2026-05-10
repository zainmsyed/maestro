import { fireEvent, render, screen, within } from '@testing-library/svelte';
import { tick } from 'svelte';
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
        featureId: null,
        featureTitle: '',
        isSynthetic: false,
        dateSource: '',
        storyPoints: null,
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
        featureId: null,
        featureTitle: '',
        isSynthetic: false,
        dateSource: 'imported',
        storyPoints: null,
      },
      {
        id: 'S-1',
        type: 'story',
        title: 'Story One',
        owner: 'Bob',
        sprint: 'Sprint 2',
        originalDate: '2026-05-02',
        committedDate: '',
        actualDate: '',
        slipEvents: 0,
        status: 'Active',
        health: 'blue',
        healthLabel: 'In progress / neutral',
        epicId: 'E-1',
        epicTitle: 'Alpha Epic',
        featureId: 'F-1',
        featureTitle: 'Feature One',
        isSynthetic: false,
        dateSource: 'imported',
        storyPoints: 5,
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
        featureId: null,
        featureTitle: '',
        isSynthetic: true,
        dateSource: 'imported',
        storyPoints: null,
      },
      {
        id: 'story-unassigned',
        type: 'story',
        title: 'Orphaned Story',
        owner: '',
        sprint: '',
        originalDate: '',
        committedDate: '',
        actualDate: '',
        slipEvents: 0,
        status: 'New',
        health: 'blue',
        healthLabel: 'In progress / neutral',
        epicId: 'epic-unassigned',
        epicTitle: 'Unassigned Epic',
        featureId: 'feature-unassigned',
        featureTitle: 'Unassigned Feature',
        isSynthetic: true,
        dateSource: 'imported',
        storyPoints: 3,
      },
    ],
  },
];

function renderTable(props: Record<string, any> = {}) {
  return render(DataTable, {
    props: {
      groups: baseGroups,
      sortKey: 'title',
      sortDirection: 'asc',
      loading: false,
      onSort: vi.fn(),
      onSaveDate: vi.fn(),
      onReassignFeature: vi.fn(),
      onReassignStory: vi.fn(),
      epicOptions: [],
      featureOptions: [],
      ...props,
    },
  });
}

describe('DataTable states', () => {
  it('renders loading state when loading is true', () => {
    renderTable({ loading: true });
    expect(screen.getByText(/loading list view/i)).toBeInTheDocument();
  });

  it('renders empty state when no groups match', () => {
    renderTable({ groups: [] });
    expect(screen.getByText(/no roadmap items match/i)).toBeInTheDocument();
  });
});

describe('DataTable interactions', () => {
  it('calls onSort when a header is clicked', async () => {
    const onSort = vi.fn();
    renderTable({ onSort });

    await fireEvent.click(screen.getByRole('button', { name: /owner/i }));
    expect(onSort).toHaveBeenCalledWith('owner');
  });

  it('supports inline feature date editing and save callback', async () => {
    const onSaveDate = vi.fn();
    renderTable({ onSaveDate });

    const featureTitle = screen.getByText('Feature One');
    const featureRow = featureTitle.closest('tr') as HTMLTableRowElement;
    const featureButton = within(featureRow).getByRole('button', { name: /set date/i });
    await fireEvent.click(featureButton);
    const input = screen.getByLabelText('Committed date for Feature One');
    await fireEvent.input(input, { target: { value: '2026-05-09' } });
    await fireEvent.click(screen.getByRole('button', { name: /^save$/i }));

    expect(onSaveDate).toHaveBeenCalledWith(expect.objectContaining({ id: 'F-1' }), '2026-05-09');
  });

  it('supports inline story date editing and save callback', async () => {
    const onSaveDate = vi.fn();
    renderTable({ onSaveDate });

    const storyTitle = screen.getByText('Story One');
    const storyRow = storyTitle.closest('tr') as HTMLTableRowElement;
    const storyButton = within(storyRow).getByRole('button', { name: /set date/i });
    await fireEvent.click(storyButton);
    const input = screen.getByLabelText('Committed date for Story One');
    await fireEvent.input(input, { target: { value: '2026-05-10' } });
    await fireEvent.click(screen.getByRole('button', { name: /^save$/i }));

    expect(onSaveDate).toHaveBeenCalledWith(expect.objectContaining({ id: 'S-1' }), '2026-05-10');
  });

  it('calls reassignment callback for synthetic feature rows', async () => {
    const onReassignFeature = vi.fn();
    const { container } = renderTable({
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

  it('calls reassignment callback for synthetic story rows', async () => {
    const onReassignStory = vi.fn();
    const { container } = renderTable({
      onReassignStory,
      featureOptions: [{ id: 'F-1', title: 'Feature One' }],
    });

    const reassignSelects = container.querySelectorAll('.reassign-label select');
    // Second reassign select is for the story row
    const storyReassign = reassignSelects[1] as HTMLSelectElement;
    await fireEvent.change(storyReassign, { target: { value: 'F-1' } });

    expect(onReassignStory).toHaveBeenCalledWith(
      expect.objectContaining({ id: 'story-unassigned' }),
      'F-1',
    );
  });

  it('renders story rows with story points and date source badge', () => {
    renderTable();

    expect(screen.getByText('Story One')).toBeInTheDocument();

    const storyTitle = screen.getByText('Story One');
    const storyRow = storyTitle.closest('tr') as HTMLTableRowElement;
    expect(within(storyRow).getByText('5')).toBeInTheDocument();
    expect(within(storyRow).getByText('imported')).toBeInTheDocument();
  });
});
