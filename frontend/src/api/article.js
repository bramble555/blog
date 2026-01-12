import service from './axios';

export const getArticles = (params) => {
    return service.get('/articles', { params });
};

export const getArticle = (id) => {
    return service.get(`/articles/${id}`);
};

export const createArticle = (data) => {
    return service.post('/articles', data);
};

export const updateArticle = (id, data) => {
    return service.put(`/articles/${id}`, data);
};

export const deleteArticles = (idList) => {
    return service.delete('/articles', { data: { id_list: idList } });
};

export const collectArticle = (id) => {
    return service.post('/articles/collects', { id: String(id) });
};

export const deleteCollectArticle = (idList) => {
    return service.delete('/articles/collects', { data: { id_list: idList } });
};
