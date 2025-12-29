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
    // 清理空字符串字段，避免后端解析time.Time失败
    const cleanedData = {}
    for (const [key, value] of Object.entries(profileData)) {
      // 跳过空字符串
      if (value === '') {
        continue
      }

      // 处理 birth_date 字段：将 YYYY-MM-DD 格式转换为 ISO 格式
      if (key === 'birth_date' && value) {
        // 确保日期格式正确（YYYY-MM-DD），然后转换为 ISO 格式
        if (value.match(/^\d{4}-\d{2}-\d{2}$/)) {
          cleanedData[key] = new Date(value + 'T00:00:00+08:00').toISOString()
        }
      } else {
        cleanedData[key] = value
      }
    }
    return request.put('/user/profile', cleanedData)
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
