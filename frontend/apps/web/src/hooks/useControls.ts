import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import { buildQueryString } from "@/lib/utils";
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
  const path = buildQueryString({ framework_id: filters.frameworkId, search: filters.search });
  return `/api/v1/controls${path}`;
}

export function useControls(filters: ControlsFilters = {}) {
  return useQuery({
    queryKey: [CONTROLS_KEY, filters.frameworkId ?? "", filters.search ?? ""],
    queryFn: () => api.getAndUnwrap<FrameworkControl[]>(buildControlsPath(filters)),
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
    queryFn: () => api.getAndUnwrap<ControlLinkedRisk[]>(`/api/v1/controls/${controlId}/risks`),
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
