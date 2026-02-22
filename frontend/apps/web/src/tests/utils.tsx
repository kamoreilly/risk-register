import { render } from '@testing-library/react';
import { RouterProvider, createMemoryHistory } from '@tanstack/react-router';
import { createRouter as createTanStackRouter } from '@tanstack/react-router';
import { routeTree } from '../routeTree.gen';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

// Create a custom render function that includes providers
export function renderWithProviders(
  ui: ReactNode,
  { route = '/', preloadedState = {} } = {}
) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  const history = createMemoryHistory({
    initialEntries: [route],
  });

  const router = createTanStackRouter({
    routeTree,
    history,
    context: {
      queryClient,
    },
  });

  return {
    ...render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    ),
    router,
    queryClient,
  };
}

// Mock API responses
export const mockApiResponse = (data: any, status = 200) => {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => data,
    text: async () => JSON.stringify(data),
  };
};

// Mock API error
export const mockApiError = (message: string, status = 500) => {
  return {
    ok: false,
    status,
    json: async () => ({ error: message }),
    text: async () => JSON.stringify({ error: message }),
  };
};

// Test data factories
export const createTestRisk = (overrides = {}) => ({
  id: '1',
  title: 'Test Risk',
  description: 'This is a test risk',
  status: 'open',
  severity: 'medium',
  likelihood: 'medium',
  impact: 'medium',
  category_id: '1',
  owner_id: '1',
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
  ...overrides,
});

export const createTestUser = (overrides = {}) => ({
  id: '1',
  email: 'test@example.com',
  name: 'Test User',
  role: 'user',
  created_at: '2024-01-01T00:00:00Z',
  ...overrides,
});

export const createTestCategory = (overrides = {}) => ({
  id: '1',
  name: 'Test Category',
  description: 'Test category description',
  created_at: '2024-01-01T00:00:00Z',
  ...overrides,
});