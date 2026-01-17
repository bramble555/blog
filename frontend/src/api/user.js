/**
 * 用户 API 模块
 * 
 * @description 提供用户管理相关操作的 API 方法（列表、角色更新、密码更新、删除、登录、登出）。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires ./axios
 */
import service from './axios';

/**
 * 获取用户列表
 * 
 * @param {Object} params - 查询参数
 * @returns {Promise} Axios 响应 Promise
 */
export const getUsers = (params) => {
    return service.get('/users', { params });
};

/**
 * 更新用户角色
 * 
 * @param {Object} data - 数据 { user_sn, role }
 * @returns {Promise} Axios 响应 Promise
 */
export const updateUserRole = (data) => {
    return service.put('/user_role', data);
};

/**
 * 更新用户密码
 * 
 * @param {Object} data - 密码更新数据
 * @returns {Promise} Axios 响应 Promise
 */
export const updateUserPassword = (data) => {
    return service.put('/user_password', data);
};

/**
 * 删除单个用户
 * 
 * @param {string} sn - 用户 SN
 * @returns {Promise} Axios 响应 Promise
 */
export const deleteUser = (sn) => {
    return service.delete('/users', { data: { sn } });
};

/**
 * 用户登录
 * 
 * @param {Object} data - 登录凭证
 * @returns {Promise} Axios 响应 Promise
 */
export const login = (data) => {
    return service.post('/user_login', data); // Based on router analysis username login is handling user_login path?
};

export const register = (data) => {
    return service.post('/register', data);
};

export const sendRegisterCode = (data) => {
    return service.post('/register_code', data);
};

/**
 * 用户登出
 * 
 * @returns {Promise} Axios 响应 Promise
 */
export const logout = () => {
    return service.post('/logout');
};

/**
 * 绑定邮箱
 * 
 * @param {Object} data - { email, code }
 * @returns {Promise} Axios 响应 Promise
 */
export const bindEmail = (data) => {
    return service.post('/user_bind_email', data);
};

/**
 * 选择用户头像 Banner
 * 
 * @param {string|number} bannerSN - Banner SN
 * @returns {Promise} Axios 响应 Promise
 */
export const selectUserBanner = (bannerSN) => {
    return service.put('/user/banner/select', {
        banner_sn: String(bannerSN)
    });
};
