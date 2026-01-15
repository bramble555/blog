/**
 * 广告 API 模块
 * 
 * @description 提供广告相关操作的 API 方法（CRUD）。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires ./axios
 */
import service from './axios';

/**
 * 获取广告列表
 * 
 * @param {Object} params - 查询参数
 * @returns {Promise} Axios 响应 Promise
 */
export const getAdverts = (params) => {
    return service.get('/adverts', { params });
};

// 上传广告图片
export const uploadAdvertImage = (data) => {
    return service.post('/adverts', data, {
        headers: { 'Content-Type': 'multipart/form-data' }
    });
};

/**
 * 创建广告
 * 
 * @param {Object} data - 广告数据
 * @returns {Promise} Axios 响应 Promise
 */
export const createAdvert = (data) => {
    return service.post('/adverts', data);
};

/**
 * 更新广告
 * 
 * @param {string} id - 广告 ID
 * @param {Object} data - 广告数据
 * @returns {Promise} Axios 响应 Promise
 */
export const updateAdvert = (id, data) => {
    return service.put('/adverts', { sn: id, ...data });
};

/**
 * 删除广告
 * 
 * @param {Array<string>} ids - 广告 ID 列表
 * @returns {Promise} Axios 响应 Promise
 */
export const deleteAdverts = (ids) => {
    return service.delete('/adverts', { data: { sn_list: (ids || []).map(String) } });
};

export const updateAdvertShow = (sn, is_show) => {
    return service.put('/adverts', { sn, is_show });
};
