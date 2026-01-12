import service from './axios';

export const getComments = (params) => {
    return service.get('/comments', { params });
};

export const createComment = (data) => {
    return service.post('/comments', data);
};

export const deleteComment = (id) => {
    return service.delete('/comments', { data: { id: id } });
};

export const digComment = (id) => {
    return service.post('/comments/digg', { id });
};
