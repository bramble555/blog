import service from './axios';

export const getMessagesAll = (params) => {
    return service.get('/messages_all', { params });
};

export const getMyMessages = () => {
    return service.get('/messages');
};

export const getMessageRecord = (userSN) => {
    return service.get('/messages_record', { data: { user_sn: userSN } }); 
};

export const sendMessage = (data) => {
    return service.post('/messages', data);
};
