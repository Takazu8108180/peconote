import { useState } from 'react';
import { useCreateMemo } from '../hooks/useMemos';

interface Props {
  mode: 'create' | 'edit';
}

function MemoFormPage({ mode }: Props) {
  const [body, setBody] = useState('');
  const mutation = useCreateMemo();

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    mutation.mutate({ body, tags: [] });
  };

  return (
    <form onSubmit={onSubmit}>
      <textarea value={body} onChange={(e) => setBody(e.target.value)} />
      <button type="submit">{mode === 'create' ? 'Create' : 'Update'}</button>
    </form>
  );
}

export default MemoFormPage;
