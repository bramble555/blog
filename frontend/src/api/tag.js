/**
 * 标签 API 模块
 * 
 * @description 提供标签相关操作的 API 方法（CRUD）。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires ./axios
 */
import service from './axios';

/**
 * 获取标签列表
 * 
 * @param {Object} params - 查询参数
 * @returns {Promise} Axios 响应 Promise
 */
export const getTags = (params) => {
    return service.get('/tags', { params });
};

/**
 * 创建标签
 * 
 * @param {Object} data - 标签数据
 * @returns {Promise} Axios 响应 Promise
 */
export const createTag = (data) => {
    return service.post('/tags', data);
};

/**
 * 批量删除标签
 * 
 * @param {Array<string>} snList - 标签序列号列表
 * @returns {Promise} Axios 响应 Promise
 */
export const deleteTags = (snList) => {
    return service.delete('/tags', { data: { sn_list: snList } });
};
