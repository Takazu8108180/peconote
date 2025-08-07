import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import client from '../api/client';

export interface Memo {
  id: string;
  body: string;
  tags: string[];
  createdAt: string;
}

export interface ListParams {
  page?: number;
  tag?: string;
}

export function useListMemos(params: ListParams) {
  return useQuery({
    queryKey: ['memos', params],
    queryFn: async () => {
      const { data } = await client.get('/memos', { params });
      return data;
    },
  });
}

export function useCreateMemo() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (memo: Omit<Memo, 'id' | 'createdAt'>) => {
      const { data } = await client.post('/memos', memo);
      return data as Memo;
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['memos'] });
    },
  });
}
