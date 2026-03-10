import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { Framework, CreateFrameworkInput } from "@/types/framework";

const FRAMEWORKS_KEY = "frameworks";

// List all frameworks
export function useFrameworks() {
  return useQuery({
    queryKey: [FRAMEWORKS_KEY],
    queryFn: async () => {
      const response = await api.get<{ data: Framework[] }>(
        "/api/v1/frameworks",
      );
      return response.data;
    },
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
}

// Create a new framework (admin only)
export function useCreateFramework() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreateFrameworkInput) =>
      api.post<Framework>("/api/v1/frameworks", input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [FRAMEWORKS_KEY] });
    },
  });
}

export function useUpdateFramework() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      ...input
    }: {
      id: string;
      name: string;
      description?: string;
    }) => api.put<Framework>(`/api/v1/frameworks/${id}`, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [FRAMEWORKS_KEY] });
    },
  });
}

export function useDeleteFramework() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/api/v1/frameworks/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [FRAMEWORKS_KEY] });
    },
  });
}
