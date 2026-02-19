import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type {
  CreateMitigationInput,
  Mitigation,
  UpdateMitigationInput,
} from '@/types/mitigation';

const MITIGATIONS_KEY = 'mitigations';

// List mitigations for a risk
export function useMitigations(riskId: string) {
  return useQuery({
    queryKey: [MITIGATIONS_KEY, riskId],
    queryFn: () => api.get<Mitigation[]>(`/api/v1/risks/${riskId}/mitigations`),
    enabled: !!riskId,
  });
}

// Get a single mitigation
export function useMitigation(riskId: string, id: string) {
  return useQuery({
    queryKey: [MITIGATIONS_KEY, riskId, id],
    queryFn: () => api.get<Mitigation>(`/api/v1/risks/${riskId}/mitigations/${id}`),
    enabled: !!riskId && !!id,
  });
}

// Create a mitigation
export function useCreateMitigation(riskId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: Omit<CreateMitigationInput, 'risk_id'>) =>
      api.post<Mitigation>(`/api/v1/risks/${riskId}/mitigations`, {
        ...input,
        risk_id: riskId,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [MITIGATIONS_KEY, riskId] });
    },
  });
}

// Update a mitigation
export function useUpdateMitigation(riskId: string, id: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdateMitigationInput) =>
      api.put<Mitigation>(`/api/v1/risks/${riskId}/mitigations/${id}`, input),
    onSuccess: (updatedMitigation) => {
      queryClient.setQueryData([MITIGATIONS_KEY, riskId, id], updatedMitigation);
      queryClient.invalidateQueries({ queryKey: [MITIGATIONS_KEY, riskId] });
    },
  });
}

// Delete a mitigation
export function useDeleteMitigation(riskId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      api.delete(`/api/v1/risks/${riskId}/mitigations/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [MITIGATIONS_KEY, riskId] });
    },
  });
}
