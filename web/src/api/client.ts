import axios from 'axios';

const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE,
});

client.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      // handle unauthorized
    }
    if (err.response?.status && err.response.status >= 500) {
      console.error('Server error', err);
    }
    return Promise.reject(err);
  },
);

export default client;
