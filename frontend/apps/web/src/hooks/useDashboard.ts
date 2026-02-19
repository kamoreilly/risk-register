import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { DashboardSummary } from '@/types/dashboard';

const DASHBOARD_KEY = ['dashboard'];

export function useDashboardSummary() {
  return useQuery({
    queryKey: [...DASHBOARD_KEY, 'summary'],
    queryFn: () => api.get<DashboardSummary>('/api/v1/dashboard/summary'),
  });
}
