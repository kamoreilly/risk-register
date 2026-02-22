import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from '@tanstack/react-router';

import { api } from '@/lib/api';
import type { AuthResponse, LoginInput, RegisterInput, User } from '@/types/auth';

const AUTH_KEY = ['auth', 'me'];

export function useAuth() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const { data: user, isLoading, error } = useQuery({
    queryKey: AUTH_KEY,
    queryFn: () => api.get<User>('/api/v1/auth/me'),
    retry: false,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });

  const loginMutation = useMutation({
    mutationFn: (input: LoginInput) =>
      api.post<AuthResponse>('/api/v1/auth/login', input),
    onSuccess: (data) => {
      localStorage.setItem('token', data.token);
      queryClient.setQueryData(AUTH_KEY, data.user);
      navigate({ to: '/app' });
    },
  });

  const registerMutation = useMutation({
    mutationFn: (input: RegisterInput) =>
      api.post<AuthResponse>('/api/v1/auth/register', input),
    onSuccess: (data) => {
      localStorage.setItem('token', data.token);
      queryClient.setQueryData(AUTH_KEY, data.user);
      navigate({ to: '/app' });
    },
  });

  const logout = () => {
    localStorage.removeItem('token');
    queryClient.clear();
    navigate({ to: '/login' });
  };

  return {
    user,
    isLoading,
    error,
    isAuthenticated: !!user,
    login: loginMutation.mutate,
    loginError: loginMutation.error,
    isLoginLoading: loginMutation.isPending,
    register: registerMutation.mutate,
    registerError: registerMutation.error,
    isRegisterLoading: registerMutation.isPending,
    logout,
  };
}
