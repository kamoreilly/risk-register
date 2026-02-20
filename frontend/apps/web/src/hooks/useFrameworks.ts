import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type {
  Framework,
  RiskFrameworkControl,
  LinkControlInput,
  CreateFrameworkInput,
} from '@/types/framework';

const FRAMEWORKS_KEY = 'frameworks';
const CONTROLS_KEY = 'controls';

// List all frameworks
export function useFrameworks() {
  return useQuery({
    queryKey: [FRAMEWORKS_KEY],
    queryFn: async () => {
      const response = await api.get<{ data: Framework[] }>('/api/v1/frameworks');
      return response.data;
    },
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
}

// Get controls for a specific risk
export function useRiskControls(riskId: string) {
  return useQuery({
    queryKey: [CONTROLS_KEY, riskId],
    queryFn: async () => {
      const response = await api.get<{ data: RiskFrameworkControl[] }>(
        `/api/v1/risks/${riskId}/controls`
      );
      return response.data;
    },
    enabled: !!riskId,
  });
}

// Link a control to a risk
export function useLinkControl(riskId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: LinkControlInput) =>
      api.post<RiskFrameworkControl>(`/api/v1/risks/${riskId}/controls`, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [CONTROLS_KEY, riskId] });
    },
  });
}

// Unlink a control from a risk
export function useUnlinkControl(riskId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (controlId: string) =>
      api.delete(`/api/v1/risks/${riskId}/controls/${controlId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [CONTROLS_KEY, riskId] });
    },
  });
}

// Create a new framework (admin only)
export function useCreateFramework() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreateFrameworkInput) =>
      api.post<Framework>('/api/v1/frameworks', input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [FRAMEWORKS_KEY] });
    },
  });
}
