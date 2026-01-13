import service from './axios';

export const getAdverts = (params) => {
    return service.get('/adverts', { params });
};

export const createAdvert = (data) => {
    return service.post('/adverts', data);
};

export const deleteAdverts = (snList) => {
    return service.delete('/adverts', { data: { sn_list: snList } });
};
