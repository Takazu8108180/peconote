import { Routes, Route, Navigate } from 'react-router-dom';
import MemoListPage from './pages/MemoListPage';
import MemoFormPage from './pages/MemoFormPage';
import MemoDetailPage from './pages/MemoDetailPage';

function App() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/memos" replace />} />
      <Route path="/memos" element={<MemoListPage />} />
      <Route path="/memos/new" element={<MemoFormPage mode="create" />} />
      <Route path="/memos/:id" element={<MemoDetailPage />} />
    </Routes>
  );
}

export default App;
