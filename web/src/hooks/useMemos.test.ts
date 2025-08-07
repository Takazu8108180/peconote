import { describe, it, expect } from 'vitest';
import { useListMemos } from './useMemos';

describe('useListMemos', () => {
  it('should expose query key', () => {
    const params = { page: 1 };
    const result = useListMemos(params);
    expect(result.queryKey).toEqual(['memos', params]);
  });
});
