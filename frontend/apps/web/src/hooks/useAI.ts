import { useMutation } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type {
  SummarizeRequest,
  SummarizeResponse,
  DraftMitigationRequest,
  DraftMitigationResponse,
} from '@/types/ai';

export function useSummarize() {
  return useMutation({
    mutationFn: (input: SummarizeRequest) =>
      api.post<SummarizeResponse>('/api/v1/ai/summarize', input),
  });
}

export function useDraftMitigation() {
  return useMutation({
    mutationFn: (input: DraftMitigationRequest) =>
      api.post<DraftMitigationResponse>('/api/v1/ai/draft-mitigation', input),
  });
}
