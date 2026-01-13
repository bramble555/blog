import service from './axios';

export const getArticles = (params) => {
    return service.get('/articles', { params });
};

export const getArticle = (sn) => {
    return service.get(`/articles/${sn}`);
};

export const createArticle = (data) => {
    return service.post('/articles', data);
};

export const updateArticle = (sn, data) => {
    return service.put(`/articles/${sn}`, data);
};

export const deleteArticles = (snList) => {
    return service.delete('/articles', { data: { sn_list: snList } });
};

export const collectArticle = (sn) => {
    return service.post('/articles/collects', { sn: String(sn) });
};

export const deleteCollectArticle = (snList) => {
    return service.delete('/articles/collects', { data: { sn_list: snList } });
};

export const getCollects = () => {
    return service.get('/articles/collects');
};
