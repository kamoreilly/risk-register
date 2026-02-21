import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { AuditLog } from '@/types/audit';

const AUDIT_KEY = 'audit';

export function useAuditLogs(riskId: string) {
  return useQuery({
    queryKey: [AUDIT_KEY, riskId],
    queryFn: async () => {
      const response = await api.get<{ data: AuditLog[] }>(`/api/v1/risks/${riskId}/audit`);
      return response.data;
    },
    enabled: !!riskId,
  });
}
