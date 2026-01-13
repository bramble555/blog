import service from './axios';

export const getComments = (params) => {
    return service.get('/comments', { params });
};

export const createComment = (data) => {
    return service.post('/comments', data);
};

export const deleteComment = (sn) => {
    return service.delete('/comments', { data: { sn: sn } });
};

export const digComment = (sn) => {
    return service.post('/comments/digg', { sn: sn });
};
