import axios from 'axios';
import { ErrorCode } from '../utils/error-code';
import { ElMessage } from 'element-plus';

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
        const res = response.data;
        // If the response has a code property, we process it
        if (res.code !== undefined) {
            // Handle auth errors (10007: NeedLogin, 10008: InvalidAuth)
            if (res.code === 10007 || res.code === 10008) {
                localStorage.removeItem('token');
                localStorage.removeItem('username');
                localStorage.removeItem('role');
                localStorage.removeItem('sn');
                localStorage.removeItem('avatar');
                // Redirect to login if not already there
                if (!window.location.pathname.includes('/login')) {
                    ElMessage.error(ErrorCode[res.code] || 'Session expired, please login again');
                    setTimeout(() => {
                        window.location.href = '/login';
                    }, 1000);
                }
                return response;
            }

            // Translate error message if mapping exists
            if (ErrorCode[res.code]) {
                response.data.msg = ErrorCode[res.code];
            }
        }
        return response;
    },
    (error) => {
        return Promise.reject(error);
    }
);

export default service;
