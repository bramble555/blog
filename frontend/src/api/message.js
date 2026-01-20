/**
 * 消息 API 模块
 * 
 * @description 提供消息系统相关操作的 API 方法。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires ./axios
 */
import service from './axios';

/**
 * 获取所有消息（管理员）
 * 
 * @param {Object} params - 查询参数
 * @returns {Promise} Axios 响应 Promise
 */
export const getMessagesAll = (params) => {
    return service.get('/messages_all', { params });
};

/**
 * 获取当前用户的消息
 * 
 * @param {Object} params - 查询参数
 * @returns {Promise} Axios 响应 Promise
 */
export const getMyMessages = (params) => {
    return service.get('/messages', { params });
};

export const getSentMessages = (params) => {
    return service.get('/messages_sent', { params });
};

/**
 * 获取与特定用户的聊天记录
 * 
 * @param {string} userSN - 用户 SN
 * @returns {Promise} Axios 响应 Promise
 */
export const getMessageRecord = (userSN) => {
    return service.get('/messages_record', { data: { user_sn: userSN } });
};

/**
 * 发送消息
 * 
 * @param {Object} data - 消息数据
 * @returns {Promise} Axios 响应 Promise
 */
export const sendMessage = (data) => {
    return service.post('/message/send', data);
};

export const broadcastMessage = (data) => {
    return service.post('/message/broadcast', data);
};
export const readMessage = (sn) => {
    return service.put('/message/read', { sn: String(sn) });
};
