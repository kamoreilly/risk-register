import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { AnalyticsResponse, Granularity } from '@/types/analytics';

export const ANALYTICS_KEY = ['analytics'];

export function useAnalytics(granularity: Granularity = 'monthly') {
  return useQuery({
    queryKey: [...ANALYTICS_KEY, granularity],
    queryFn: () =>
      api.get<AnalyticsResponse>(`/api/v1/analytics?granularity=${granularity}`),
  });
}
