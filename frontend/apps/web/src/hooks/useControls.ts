import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  ControlLinkedRisk,
  CreateFrameworkControlInput,
  FrameworkControl,
  LinkControlInput,
  RiskFrameworkControl,
  UpdateFrameworkControlInput,
} from "@/types/framework";

const CONTROLS_KEY = "controls";
const CONTROL_LINKED_RISKS_KEY = "control-linked-risks";
const RISK_CONTROLS_KEY = "risk-controls";

interface ControlsFilters {
  frameworkId?: string;
  search?: string;
}

function buildControlsPath(filters: ControlsFilters = {}) {
  const params = new URLSearchParams();

  if (filters.frameworkId) {
    params.set("framework_id", filters.frameworkId);
  }
  if (filters.search) {
    params.set("search", filters.search);
  }

  const query = params.toString();
  return query ? `/api/v1/controls?${query}` : "/api/v1/controls";
}

export function useControls(filters: ControlsFilters = {}) {
  return useQuery({
    queryKey: [CONTROLS_KEY, filters.frameworkId ?? "", filters.search ?? ""],
    queryFn: async () => {
      const response = await api.get<{ data: FrameworkControl[] }>(
        buildControlsPath(filters),
      );
      return response.data;
    },
  });
}

export function useCreateControl() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreateFrameworkControlInput) =>
      api.post<FrameworkControl>("/api/v1/controls", input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [CONTROLS_KEY] });
    },
  });
}

export function useControlLinkedRisks(controlId: string, enabled = true) {
  return useQuery({
    queryKey: [CONTROL_LINKED_RISKS_KEY, controlId],
    queryFn: async () => {
      const response = await api.get<{ data: ControlLinkedRisk[] }>(
        `/api/v1/controls/${controlId}/risks`,
      );
      return response.data;
    },
    enabled: enabled && !!controlId,
  });
}

export function useUpdateControl() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      ...input
    }: { id: string } & UpdateFrameworkControlInput) =>
      api.put<FrameworkControl>(`/api/v1/controls/${id}`, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [CONTROLS_KEY] });
    },
  });
}

export function useDeleteControl() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/api/v1/controls/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [CONTROLS_KEY] });
    },
  });
}

export function useRiskControls(riskId: string) {
  return useQuery({
    queryKey: [RISK_CONTROLS_KEY, riskId],
    queryFn: async () => {
      const response = await api.get<{ data: RiskFrameworkControl[] }>(
        `/api/v1/risks/${riskId}/controls`,
      );
      return response.data;
    },
    enabled: !!riskId,
  });
}

export function useLinkControl(riskId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: LinkControlInput) =>
      api.post<RiskFrameworkControl>(`/api/v1/risks/${riskId}/controls`, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [RISK_CONTROLS_KEY, riskId] });
      queryClient.invalidateQueries({ queryKey: [CONTROL_LINKED_RISKS_KEY] });
    },
  });
}

export function useUnlinkControl(riskId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (controlId: string) =>
      api.delete(`/api/v1/risks/${riskId}/controls/${controlId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [RISK_CONTROLS_KEY, riskId] });
      queryClient.invalidateQueries({ queryKey: [CONTROL_LINKED_RISKS_KEY] });
    },
  });
}
