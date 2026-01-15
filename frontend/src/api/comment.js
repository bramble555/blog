/**
 * 评论 API 模块
 * 
 * @description 提供评论相关操作的 API 方法（列表、创建、删除、点赞）。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires ./axios
 */
import service from './axios';

/**
 * 创建评论
 * 
 * @param {Object} data - 评论数据 { article_sn, parent_comment_sn, content }
 * @returns {Promise} Axios 响应 Promise
 */
export const createComment = (data) => {
    return service.post('/comments', data);
};

/**
 * 获取评论列表
 * 
 * @param {string} sn - 文章 SN
 * @param {Object} params - 分页参数 { page, size }
 * @returns {Promise} Axios 响应 Promise
 */
export const getComments = (sn, params) => {
    const query = { ...(params || {}) };
    if (sn !== undefined && sn !== null && sn !== '') {
        query.article_sn = sn;
    }
    return service.get('/comments', { params: query });
};

/**
 * 删除评论
 * 
 * @param {string} sn - 评论 SN
 * @returns {Promise} Axios 响应 Promise
 */
export const deleteComment = (sn) => {
    return service.delete('/comments', { data: { sn: String(sn) } });
};

/**
 * 评论点赞
 * 
 * @param {string} sn - 评论 SN
 * @returns {Promise} Axios 响应 Promise
 */
export const diggComment = (sn) => {
    return service.post('/comments/digg', { sn: String(sn) });
};

/**
 * 批量删除评论
 * 
 * @param {Array<string>} ids - 评论 ID 列表
 * @returns {Promise} Axios 响应 Promise
 */
export const removeCommentBatch = (ids) => {
    return service.delete('/comments', { data: { id_list: ids } });
};
