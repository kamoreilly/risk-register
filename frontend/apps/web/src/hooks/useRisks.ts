import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type {
  Category,
  CreateRiskInput,
  Risk,
  RiskListParams,
  RiskListResponse,
  UpdateRiskInput,
} from '@/types/risk';

const RISKS_KEY = ['risks'];
const CATEGORIES_KEY = ['categories'];

// Categories
export function useCategories() {
  return useQuery({
    queryKey: CATEGORIES_KEY,
    queryFn: () => api.get<Category[]>('/api/v1/categories'),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

// Risks list
export function useRisks(params?: RiskListParams) {
  const queryString = buildQueryString(params);
  return useQuery({
    queryKey: [...RISKS_KEY, params],
    queryFn: () => api.get<RiskListResponse>(`/api/v1/risks${queryString}`),
  });
}

// Single risk
export function useRisk(id: string) {
  return useQuery({
    queryKey: [...RISKS_KEY, id],
    queryFn: () => api.get<Risk>(`/api/v1/risks/${id}`),
    enabled: !!id,
  });
}

// Create risk
export function useCreateRisk() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreateRiskInput) =>
      api.post<Risk>('/api/v1/risks', input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: RISKS_KEY });
    },
  });
}

// Update risk
export function useUpdateRisk(id: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdateRiskInput) =>
      api.put<Risk>(`/api/v1/risks/${id}`, input),
    onSuccess: (updatedRisk) => {
      queryClient.setQueryData([...RISKS_KEY, id], updatedRisk);
      queryClient.invalidateQueries({ queryKey: RISKS_KEY });
    },
  });
}

// Delete risk
export function useDeleteRisk() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/api/v1/risks/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: RISKS_KEY });
    },
  });
}

// Helper to build query string from params
function buildQueryString(params?: RiskListParams): string {
  if (!params) return '';

  const searchParams = new URLSearchParams();

  if (params.status) searchParams.set('status', params.status);
  if (params.severity) searchParams.set('severity', params.severity);
  if (params.category_id) searchParams.set('category_id', params.category_id);
  if (params.owner_id) searchParams.set('owner_id', params.owner_id);
  if (params.search) searchParams.set('search', params.search);
  if (params.sort) searchParams.set('sort', params.sort);
  if (params.order) searchParams.set('order', params.order);
  if (params.page) searchParams.set('page', String(params.page));
  if (params.limit) searchParams.set('limit', String(params.limit));

  const queryString = searchParams.toString();
  return queryString ? `?${queryString}` : '';
}
