import request from './index'

// 导出对象和具名导出
const userAPI = {
  // 用户注册
  register(userData) {
    return request.post('/users/register', userData)
  },

  // 创建用户（管理员）
  createUser(userData) {
    return request.post('/users', userData)
  },

  // 获取用户列表（管理员）
  getUsers(params = {}) {
    return request.get('/users', params)
  },

  // 获取用户详情
  getUserDetail(id) {
    return request.get(`/users/${id}`)
  },

  // 更新用户信息
  updateUser(id, userData) {
    return request.put(`/users/${id}`, userData)
  },

  // 删除用户（管理员）
  deleteUser(id) {
    return request.delete(`/users/${id}`)
  },

  // 获取用户统计信息（管理员）
  getUserStats(params = {}) {
    return request.get('/users/statistics', params)
  },

  // 更新用户密码
  updatePassword(id, passwordData) {
    return request.put(`/users/${id}/password`, passwordData)
  },

  // 重置用户密码（管理员）
  resetUserPassword(id) {
    return request.post(`/users/${id}/reset-password`)
  },

  // 更新用户角色（管理员）
  updateUserRole(id, role) {
    return request.put(`/users/${id}/role`, { role })
  },

  // 搜索用户（管理员）
  searchUsers(params) {
    return request.get('/users/search', params)
  },

  // 禁用/启用用户（管理员）
  toggleUserStatus(id) {
    return request.put(`/users/${id}/toggle-status`)
  },

  // 获取用户订单历史
  getUserOrders(id, params = {}) {
    return request.get(`/users/${id}/orders`, params)
  },

  // 获取今日注册用户数量
  getTodayRegisteredCount() {
    return request.get('/users/today-registered')
  }
}

export default userAPI
export { userAPI }