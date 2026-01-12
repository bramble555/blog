import service from './axios';

export const getMessagesAll = (params) => {
    return service.get('/messages_all', { params });
};

export const getMyMessages = () => {
    return service.get('/messages');
};

export const getMessageRecord = (userId) => {
    return service.get('/messages_record', { data: { user_id: userId } }); // Note: Controller uses ShouldBindJSON for GET in MessageRecordHandler but standard GET doesn't have body. Axios might need workaround or backend is unconventional.
    // Actually MessageRecordHandler uses ShouldBindJSON. GET with body is non-standard but possible.
};

export const sendMessage = (data) => {
    return service.post('/messages', data);
};
