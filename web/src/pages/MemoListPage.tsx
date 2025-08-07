import { useListMemos } from '../hooks/useMemos';

function MemoListPage() {
  const { data } = useListMemos({});

  return (
    <div>
      <h1>Memos</h1>
      <ul>
        {data?.items?.map((m: any) => (
          <li key={m.id}>{m.body}</li>
        ))}
      </ul>
    </div>
  );
}

export default MemoListPage;
