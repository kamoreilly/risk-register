import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import { renderWithProviders } from '../utils';

// Mock the API module
vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
  },
}));

// Mock the useDashboard hook
vi.mock('@/hooks/useDashboard', () => ({
  useDashboard: () => ({
    stats: {
      totalRisks: 10,
      openRisks: 5,
      highSeverityRisks: 2,
      mitigatedRisks: 3,
    },
    recentRisks: [
      {
        id: '1',
        title: 'Test Risk',
        status: 'open',
        severity: 'high',
        created_at: '2024-01-01T00:00:00Z',
      },
    ],
    isLoading: false,
    error: null,
  }),
}));

import { api } from '@/lib/api';

describe('Dashboard Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders dashboard statistics', async () => {
    renderWithProviders(null, { route: '/app' });

    expect(screen.getByText(/dashboard/i)).toBeInTheDocument();
    expect(screen.getByText(/10/)).toBeInTheDocument(); // Total risks
    expect(screen.getByText(/5/)).toBeInTheDocument(); // Open risks
    expect(screen.getByText(/2/)).toBeInTheDocument(); // High severity risks
    expect(screen.getByText(/3/)).toBeInTheDocument(); // Mitigated risks
  });

  it('renders recent risks section', async () => {
    renderWithProviders(null, { route: '/app' });

    expect(screen.getByText(/recent risks/i)).toBeInTheDocument();
    expect(screen.getByText(/test risk/i)).toBeInTheDocument();
    expect(screen.getByText(/high/i)).toBeInTheDocument();
  });

  it('renders loading state', async () => {
    // Mock loading state
    vi.doMock('@/hooks/useDashboard', () => ({
      useDashboard: () => ({
        stats: null,
        recentRisks: [],
        isLoading: true,
        error: null,
      }),
    }));

    renderWithProviders(null, { route: '/app' });

    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
  });

  it('renders error state', async () => {
    // Mock error state
    vi.doMock('@/hooks/useDashboard', () => ({
      useDashboard: () => ({
        stats: null,
        recentRisks: [],
        isLoading: false,
        error: 'Failed to load dashboard data',
      }),
    }));

    renderWithProviders(null, { route: '/app' });

    expect(screen.getByText(/failed to load dashboard data/i)).toBeInTheDocument();
  });

  it('navigates to risks page when view all risks is clicked', async () => {
    renderWithProviders(null, { route: '/app' });

    const viewAllLink = screen.getByRole('link', { name: /view all risks/i });
    expect(viewAllLink).toHaveAttribute('href', '/app/risks');
  });
});