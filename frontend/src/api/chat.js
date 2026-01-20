// Chat Group is WebSocket based.
// Base URL for websocket
export const WS_URL = 'ws://localhost:8080/api/chat_groups';

import service from './axios';

/**
 * 获取群聊历史记录
 * @param {Object} params { page, size }
 */
export const getChatHistory = (params) => {
    return service.get('/chat_groups_records', { params });
}

export const uploadChatImage = (data) => {
    return service.post('/chat_groups_images', data, {
        headers: { 'Content-Type': 'multipart/form-data' }
    });
}
