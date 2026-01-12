import service from './axios';

export const getUsers = (params) => {
    return service.get('/users', { params });
};

export const updateUserRole = (data) => {
    return service.put('/user_role', data);
};

export const updateUserPassword = (data) => {
    return service.put('/user_password', data);
};

export const deleteUsers = (idList) => {
    return service.delete('/users', { data: { id_list: idList } });
};

export const login = (data) => {
    return service.post('/email_login', data); // Based on router analysis username login is handling email_login path?
};

export const logout = () => {
    return service.post('/logout');
};
