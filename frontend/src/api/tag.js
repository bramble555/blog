import service from './axios';

export const getTags = (params) => {
    return service.get('/tags', { params });
};

export const createTag = (data) => {
    return service.post('/tags', data);
};

export const deleteTags = (idList) => {
    return service.delete('/tags', { data: { id_list: idList } });
};
