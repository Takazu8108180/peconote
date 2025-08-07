import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import client from '../api/client';

function MemoDetailPage() {
  const { id } = useParams();
  const { data } = useQuery({
    queryKey: ['memo', id],
    queryFn: async () => {
      const { data } = await client.get(`/memos/${id}`);
      return data;
    },
    enabled: !!id,
  });

  if (!data) return <div>Loading...</div>;
  return (
    <div>
      <h1>Memo Detail</h1>
      <p>{data.body}</p>
    </div>
  );
}

export default MemoDetailPage;
