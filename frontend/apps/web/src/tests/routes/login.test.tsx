import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithProviders } from '../utils';

// Mock the API module
vi.mock('@/lib/api', () => ({
  api: {
    post: vi.fn(),
  },
}));

import { api } from '@/lib/api';

describe('Login Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders login form correctly', async () => {
    renderWithProviders(null, { route: '/login' });

    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument();
    expect(screen.getByText(/don't have an account\?/i)).toBeInTheDocument();
  });

  it('validates required fields', async () => {
    const user = userEvent.setup();
    renderWithProviders(null, { route: '/login' });

    const submitButton = screen.getByRole('button', { name: /sign in/i });
    await user.click(submitButton);

    expect(await screen.findByText(/email is required/i)).toBeInTheDocument();
    expect(await screen.findByText(/password is required/i)).toBeInTheDocument();
  });

  it('submits form with valid data', async () => {
    const user = userEvent.setup();
    const mockResponse = { token: 'test-token', user: { id: '1', email: 'test@example.com', name: 'Test User' } };
    
    vi.mocked(api.post).mockResolvedValueOnce({ data: mockResponse });
    
    renderWithProviders(null, { route: '/login' });

    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/password/i), 'password123');
    await user.click(screen.getByRole('button', { name: /sign in/i }));

    await waitFor(() => {
      expect(api.post).toHaveBeenCalledWith('/auth/login', {
        email: 'test@example.com',
        password: 'password123',
      });
    });
  });

  it('handles login failure', async () => {
    const user = userEvent.setup();
    const errorMessage = 'Invalid credentials';
    
    vi.mocked(api.post).mockRejectedValueOnce({
      response: { data: { error: errorMessage } },
    });
    
    renderWithProviders(null, { route: '/login' });

    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/password/i), 'wrongpassword');
    await user.click(screen.getByRole('button', { name: /sign in/i }));

    expect(await screen.findByText(errorMessage)).toBeInTheDocument();
  });

  it('navigates to register page', async () => {
    const user = userEvent.setup();
    renderWithProviders(null, { route: '/login' });

    const registerLink = screen.getByRole('link', { name: /sign up/i });
    await user.click(registerLink);

    // This would typically be verified by checking navigation state
    // For now, we just verify the link is present and has correct href
    expect(registerLink).toHaveAttribute('href', '/register');
  });
});