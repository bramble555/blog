/**
 * 文章 API 模块
 * 
 * @description 提供文章相关操作的 API 方法（CRUD、收藏）。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires ./axios
 */
import service from './axios';

/**
 * 获取文章列表
 * 
 * @param {Object} params - 查询参数 { page, size, keyword }
 * @returns {Promise} Axios 响应 Promise
 */
export const getArticles = (params) => {
    const query = { ...params };
    if (query.keyword) {
        query.title = query.keyword;
        delete query.keyword;
    }
    return service.get('/articles', { params: query });
};

/**
 * 获取文章详情
 * 
 * @param {string} id - 文章 ID
 * @returns {Promise} Axios 响应 Promise
 */
export const getArticle = (id) => {
    return service.get(`/articles/${id}`);
};

/**
 * 创建文章
 * 
 * @param {Object} data - 文章数据
 * @returns {Promise} Axios 响应 Promise
 */
export const createArticle = (data) => {
    return service.post('/articles', data);
};

/**
 * 更新文章
 * 
 * @param {string} sn - 文章序列号
 * @param {Object} data - 文章数据
 * @returns {Promise} Axios 响应 Promise
 */
export const updateArticle = (sn, data) => {
    return service.put(`/articles/${sn}`, data);
};

/**
 * 删除文章（批量）
 * 
 * @param {Array<string>} idList - 文章 ID 列表
 * @returns {Promise} Axios 响应 Promise
 */
export const deleteArticles = (idList) => {
    return service.delete('/articles', { data: { sn_list: idList } });
};

/**
 * 收藏文章
 * 
 * @param {string} sn - 文章序列号
 * @returns {Promise} Axios 响应 Promise
 */
export const collectArticle = (sn) => {
    return service.post('/articles/collects', { sn: String(sn) });
};

/**
 * 取消收藏文章（批量）
 * 
 * @param {Array<string>} snList - 文章序列号列表
 * @returns {Promise} Axios 响应 Promise
 */
export const deleteCollectArticle = (snList) => {
    return service.delete('/articles/collects', { data: { sn_list: snList } });
};

/**
 * 获取用户收藏的文章列表
 * 
 * @returns {Promise} Axios 响应 Promise
 */
export const getCollects = () => {
    return service.get('/articles/collects');
};
