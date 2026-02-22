import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithProviders } from '../utils';

// Mock the API module
vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
    delete: vi.fn(),
  },
}));

// Mock the useRisks hook
vi.mock('@/hooks/useRisks', () => ({
  useRisks: () => ({
    risks: [
      {
        id: '1',
        title: 'Test Risk 1',
        description: 'Test risk description',
        status: 'open',
        severity: 'high',
        likelihood: 'medium',
        impact: 'high',
        category: { name: 'Security' },
        owner: { name: 'John Doe' },
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      },
      {
        id: '2',
        title: 'Test Risk 2',
        description: 'Another test risk',
        status: 'mitigated',
        severity: 'medium',
        likelihood: 'low',
        impact: 'medium',
        category: { name: 'Operational' },
        owner: { name: 'Jane Smith' },
        created_at: '2024-01-02T00:00:00Z',
        updated_at: '2024-01-02T00:00:00Z',
      },
    ],
    isLoading: false,
    error: null,
    refetch: vi.fn(),
  }),
}));

import { api } from '@/lib/api';

describe('Risks List Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders risks list with correct data', async () => {
    renderWithProviders(null, { route: '/app/risks' });

    expect(screen.getByText(/risks/i)).toBeInTheDocument();
    expect(screen.getByText(/test risk 1/i)).toBeInTheDocument();
    expect(screen.getByText(/test risk 2/i)).toBeInTheDocument();
    expect(screen.getByText(/security/i)).toBeInTheDocument();
    expect(screen.getByText(/operational/i)).toBeInTheDocument();
    expect(screen.getByText(/high/i)).toBeInTheDocument();
    expect(screen.getByText(/medium/i)).toBeInTheDocument();
  });

  it('renders create risk button', async () => {
    renderWithProviders(null, { route: '/app/risks' });

    const createButton = screen.getByRole('link', { name: /create risk/i });
    expect(createButton).toHaveAttribute('href', '/app/risks/new');
  });

  it('filters risks by status', async () => {
    const user = userEvent.setup();
    renderWithProviders(null, { route: '/app/risks' });

    const filterButton = screen.getByRole('button', { name: /filter/i });
    await user.click(filterButton);

    const openCheckbox = screen.getByLabelText(/open/i);
    await user.click(openCheckbox);

    // Should only show open risks
    expect(screen.getByText(/test risk 1/i)).toBeInTheDocument();
    expect(screen.queryByText(/test risk 2/i)).not.toBeInTheDocument();
  });

  it('searches risks by title', async () => {
    const user = userEvent.setup();
    renderWithProviders(null, { route: '/app/risks' });

    const searchInput = screen.getByPlaceholderText(/search risks/i);
    await user.type(searchInput, 'Test Risk 1');

    // Should only show matching risks
    expect(screen.getByText(/test risk 1/i)).toBeInTheDocument();
    expect(screen.queryByText(/test risk 2/i)).not.toBeInTheDocument();
  });

  it('deletes a risk', async () => {
    const user = userEvent.setup();
    vi.mocked(api.delete).mockResolvedValueOnce({ data: { success: true } });
    
    renderWithProviders(null, { route: '/app/risks' });

    const deleteButton = screen.getAllByRole('button', { name: /delete/i })[0];
    await user.click(deleteButton);

    const confirmButton = screen.getByRole('button', { name: /confirm/i });
    await user.click(confirmButton);

    await waitFor(() => {
      expect(api.delete).toHaveBeenCalledWith('/risks/1');
    });
  });

  it('navigates to risk detail', async () => {
    renderWithProviders(null, { route: '/app/risks' });

    const riskLink = screen.getByRole('link', { name: /test risk 1/i });
    expect(riskLink).toHaveAttribute('href', '/app/risks/1');
  });

  it('shows empty state when no risks', async () => {
    // Mock empty state
    vi.doMock('@/hooks/useRisks', () => ({
      useRisks: () => ({
        risks: [],
        isLoading: false,
        error: null,
        refetch: vi.fn(),
      }),
    }));

    renderWithProviders(null, { route: '/app/risks' });

    expect(screen.getByText(/no risks found/i)).toBeInTheDocument();
    expect(screen.getByText(/create your first risk/i)).toBeInTheDocument();
  });
});