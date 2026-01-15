/**
 * Banner API 模块
 * 
 * @description 提供 Banner（轮播图）相关操作的 API 方法（CRUD）。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires ./axios
 */
import service from './axios';

/**
 * 获取 Banner 列表
 * 
 * @param {Object} params - 查询参数
 * @returns {Promise} Axios 响应 Promise
 */
export const getBanners = (params) => {
    return service.get('/images', { params }); // Assuming banners use images endpoint or similar
};

/**
 * 上传 Banner
 * 
 * @param {FormData} data - 包含文件的 FormData
 * @returns {Promise} Axios 响应 Promise
 */
export const uploadBanners = (data) => {
    return service.post('/images', data, {
        headers: { 'Content-Type': 'multipart/form-data' }
    });
};

/**
 * 删除 Banner
 * 
 * @param {Array<string>} ids - Banner ID 列表
 * @returns {Promise} Axios 响应 Promise
 */
export const deleteBanners = (ids) => {
    return service.delete('/images', { data: { sn_list: (ids || []).map(String) } });
};

/**
 * 更新 Banner 信息
 * 
 * @param {Object} data - Banner 数据
 * @returns {Promise} Axios 响应 Promise
 */
export const updateBanner = (data) => {
    return service.put('/images', data);
};
