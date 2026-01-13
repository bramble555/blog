import service from './axios';

export const getBanners = (params) => {
    return service.get('/images', { params });
};

export const getAllBanners = () => {
    return service.get('/images_detail');
};

export const uploadBanners = (formData) => {
    return service.post('/images', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
    });
};

export const deleteBanners = (snList) => {
    return service.delete('/images', { data: { sn_list: snList } });
};
