import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { DashboardSummary, ReviewListResponse } from '@/types/dashboard';

const DASHBOARD_KEY = ['dashboard'];

export function useDashboardSummary() {
  return useQuery({
    queryKey: [...DASHBOARD_KEY, 'summary'],
    queryFn: () => api.get<DashboardSummary>('/api/v1/dashboard/summary'),
  });
}

export function useUpcomingReviews(days = 30) {
  return useQuery({
    queryKey: [...DASHBOARD_KEY, 'reviews', 'upcoming', days],
    queryFn: async () => {
      const response = await api.get<ReviewListResponse>(`/api/v1/dashboard/reviews/upcoming?days=${days}`);
      return response;
    },
  });
}

export function useOverdueReviews() {
  return useQuery({
    queryKey: [...DASHBOARD_KEY, 'reviews', 'overdue'],
    queryFn: async () => {
      const response = await api.get<ReviewListResponse>('/api/v1/dashboard/reviews/overdue');
      return response;
    },
  });
}
