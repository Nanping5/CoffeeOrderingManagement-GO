import request from './index'

const authAPI = {
  // 用户登录
  login(credentials) {
    return request.post('/auth/login', credentials)
  },

  // 用户注册
  register(userData) {
    return request.post('/auth/register', userData)
  },

  // 刷新Token
  refreshToken() {
    return request.post('/auth/refresh')
  },

  // 登出
  logout() {
    return request.post('/auth/logout')
  },

  // 获取用户资料
  getProfile() {
    return request.get('/user/profile')
  },

  // 更新用户资料
  updateProfile(profileData) {
    return request.put('/user/profile', profileData)
  },

  // 修改密码
  changePassword(passwordData) {
    return request.put('/user/password', passwordData)
  },

  // 验证Token
  verifyToken() {
    return request.get('/auth/verify-token')
  },

  // 管理员登录
  adminLogin(credentials) {
    return request.post('/auth/admin/login', credentials)
  }
}

export default authAPI
