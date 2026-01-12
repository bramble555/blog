import axios from 'axios';

const api = axios.create({
    baseURL: 'http://localhost:8080/api', // Requesting directly to backend port as per standard go-gin setups
    timeout: 5000,
});

// Interceptor to handle optional token if we had auth, but for now just basic
api.interceptors.request.use(
    (config) => {
        // const token = localStorage.getItem('token');
        // if (token) {
        //   config.headers.Authorization = `Bearer ${token}`;
        // }
        return config;
    },
    (error) => Promise.reject(error)
);

export default api;
