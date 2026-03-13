import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type {
  IncidentCategory,
  CreateIncidentInput,
  Incident,
  IncidentListParams,
  IncidentListResponse,
  IncidentRisk,
  LinkIncidentRiskInput,
  UpdateIncidentInput,
} from '@/types/incident';

const INCIDENTS_KEY = ['incidents'];
const INCIDENT_CATEGORIES_KEY = ['incident-categories'];
const INCIDENT_RISKS_KEY = ['incident-risks'];

// Incident Categories
export function useIncidentCategories() {
  return useQuery({
    queryKey: INCIDENT_CATEGORIES_KEY,
    queryFn: () => api.get<IncidentCategory[]>('/api/v1/incident-categories'),
    staleTime: 5 * 60 * 1000,
  });
}

export function useCreateIncidentCategory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: { name: string; description?: string }) =>
      api.post<IncidentCategory>('/api/v1/incident-categories', input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: INCIDENT_CATEGORIES_KEY });
    },
  });
}

// Incidents List
export function useIncidents(params?: IncidentListParams) {
  const queryString = buildQueryString(params);
  return useQuery({
    queryKey: [...INCIDENTS_KEY, params],
    queryFn: () => api.get<IncidentListResponse>(`/api/v1/incidents${queryString}`),
  });
}

// Single Incident
export function useIncident(id: string) {
  return useQuery({
    queryKey: [...INCIDENTS_KEY, id],
    queryFn: () => api.get<Incident>(`/api/v1/incidents/${id}`),
    enabled: !!id,
  });
}

// Create Incident
export function useCreateIncident() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreateIncidentInput) =>
      api.post<Incident>('/api/v1/incidents', input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: INCIDENTS_KEY });
    },
  });
}

// Update Incident
export function useUpdateIncident(id: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdateIncidentInput) =>
      api.put<Incident>(`/api/v1/incidents/${id}`, input),
    onSuccess: (updatedIncident) => {
      queryClient.setQueryData([...INCIDENTS_KEY, id], updatedIncident);
      queryClient.invalidateQueries({ queryKey: INCIDENTS_KEY });
    },
  });
}

// Delete Incident
export function useDeleteIncident() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/api/v1/incidents/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: INCIDENTS_KEY });
    },
  });
}

// Incident-Risk Links
export function useIncidentRisks(incidentId: string) {
  return useQuery({
    queryKey: [...INCIDENT_RISKS_KEY, incidentId],
    queryFn: () => api.get<IncidentRisk[]>(`/api/v1/incidents/${incidentId}/risks`),
    enabled: !!incidentId,
  });
}

export function useLinkIncidentRisk(incidentId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: LinkIncidentRiskInput) =>
      api.post<IncidentRisk>(`/api/v1/incidents/${incidentId}/risks`, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [...INCIDENT_RISKS_KEY, incidentId] });
    },
  });
}

export function useUnlinkIncidentRisk(incidentId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (riskId: string) =>
      api.delete(`/api/v1/incidents/${incidentId}/risks/${riskId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [...INCIDENT_RISKS_KEY, incidentId] });
    },
  });
}

// Incident Audit Logs
export function useIncidentAuditLogs(incidentId: string) {
  return useQuery({
    queryKey: ['incident-audit', incidentId],
    queryFn: async () => {
      const response = await api.get<{ data: import('@/types/audit').AuditLog[] }>(`/api/v1/incidents/${incidentId}/audit`);
      return response.data;
    },
    enabled: !!incidentId,
  });
}

// Helper to build query string
function buildQueryString(params?: IncidentListParams): string {
  if (!params) return '';

  const searchParams = new URLSearchParams();

  if (params.status) searchParams.set('status', params.status);
  if (params.priority) searchParams.set('priority', params.priority);
  if (params.category_id) searchParams.set('category_id', params.category_id);
  if (params.assignee_id) searchParams.set('assignee_id', params.assignee_id);
  if (params.search) searchParams.set('search', params.search);
  if (params.sort) searchParams.set('sort', params.sort);
  if (params.order) searchParams.set('order', params.order);
  if (params.page) searchParams.set('page', String(params.page));
  if (params.limit) searchParams.set('limit', String(params.limit));

  const queryString = searchParams.toString();
  return queryString ? `?${queryString}` : '';
}
