import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { User } from '@/types/auth';

const USERS_KEY = ['users'];

export function useUsers() {
  return useQuery({
    queryKey: USERS_KEY,
    queryFn: () => api.get<User[]>('/api/v1/users'),
    staleTime: 5 * 60 * 1000,
  });
}
