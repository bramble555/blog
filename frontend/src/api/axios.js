import axios from 'axios';

const service = axios.create({
	baseURL: 'http://localhost:8080/api',
	timeout: 5000,
	withCredentials: true,
});

service.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers['token'] = token;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

service.interceptors.response.use(
    (response) => {
        return response;
    },
    (error) => {
        return Promise.reject(error);
    }
);

export default service;
