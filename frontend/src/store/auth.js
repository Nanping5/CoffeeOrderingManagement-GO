import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import authApi from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  // 从localStorage恢复用户信息
  const getStoredUser = () => {
    try {
      const storedUser = localStorage.getItem('user')
      return storedUser ? JSON.parse(storedUser) : null
    } catch {
      return null
    }
  }

  // 状态 - 初始化时从localStorage恢复
  const user = ref(getStoredUser())
  const token = ref(localStorage.getItem('token'))
  const refreshToken = ref(localStorage.getItem('refreshToken'))
  const loading = ref(false)
  const authCheckInProgress = ref(false)

  // 计算属性
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const userInfo = computed(() => user.value)

  // 方法
  const login = async (credentials) => {
    try {
      loading.value = true
      const response = await authApi.login(credentials)

      if (response.success && response.data) {
        // 后端返回格式: { success: true, data: { user, token } }
        const { user: userData, token: accessToken } = response.data
        
        // 保存token和用户信息
        token.value = accessToken
        user.value = userData

        // 持久化存储
        localStorage.setItem('token', accessToken)
        localStorage.setItem('user', JSON.stringify(userData))

        return { success: true, message: response.message }
      } else {
        return { success: false, errors: [response.message || '登录失败'] }
      }
    } catch (error) {
      console.error('登录失败:', error)
      return { success: false, errors: [error.message || '登录失败'] }
    } finally {
      loading.value = false
    }
  }

  const register = async (userData) => {
    try {
      loading.value = true
      const response = await authApi.register(userData)

      if (response.success && response.data) {
        // 后端返回格式: { success: true, data: { user, token } }
        const { user: userInfo, token: accessToken } = response.data
        
        // 保存token和用户信息
        token.value = accessToken
        user.value = userInfo

        // 持久化存储
        localStorage.setItem('token', accessToken)
        localStorage.setItem('user', JSON.stringify(userInfo))

        return { success: true, message: response.message }
      } else {
        return { success: false, errors: [response.message || '注册失败'] }
      }
    } catch (error) {
      console.error('注册失败:', error)
      return { success: false, errors: [error.message || '注册失败'] }
    } finally {
      loading.value = false
    }
  }

  const logout = () => {
    // 清除状态
    user.value = null
    token.value = null
    refreshToken.value = null

    // 清除本地存储
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')

    // 不调用logout API，避免重复请求
    // authApi.logout().catch(console.error)
  }

  const checkAuthStatus = async () => {
    try {
      // 防止重复认证检查
      if (authCheckInProgress.value) {
        return false
      }

      // 检查本地存储
      const storedToken = localStorage.getItem('token')
      const storedUser = localStorage.getItem('user')

      if (!storedToken || !storedUser) {
        // 没有存储的认证信息，直接返回false
        return false
      }

      // 如果有本地存储的信息，设置到响应式状态
      if (!token.value) {
        token.value = storedToken
        user.value = JSON.parse(storedUser)
      }

      // 如果是假token（管理员登录），直接返回true，不验证
      if (storedToken.startsWith('fake-admin-token-')) {
        return true
      }

      // 如果已经有用户信息，直接返回true，不需要验证
      // 这样可以避免在每次路由切换时都调用API
      if (user.value && token.value) {
        return true
      }

      // 防止重复验证
      if (authCheckInProgress.value) {
        return !!user.value
      }

      authCheckInProgress.value = true

      // 验证token有效性（可选，只在必要时调用）
      try {
        const response = await authApi.verifyToken()
        if (response.success) {
          user.value = response.user
          authCheckInProgress.value = false
          return true
        } else {
          // token无效，清除状态
          logout()
          authCheckInProgress.value = false
          return false
        }
      } catch (verifyError) {
        // 如果验证失败，但有本地存储的用户信息，仍然认为已登录
        // 让后续的API调用来处理Token过期的情况
        if (storedUser && storedToken) {
          authCheckInProgress.value = false
          return true
        }
        throw verifyError
      }
    } catch (error) {
      console.error('检查认证状态失败:', error)
      authCheckInProgress.value = false
      // 如果有本地存储的信息，不要清除，让用户继续使用
      // 如果Token真的无效，后续API调用会返回401
      if (localStorage.getItem('token') && localStorage.getItem('user')) {
        return true
      }
      // 只有在完全没有认证信息时才清除
      logout()
      return false
    }
  }

  const updateProfile = async (profileData) => {
    try {
      loading.value = true
      const response = await authApi.updateProfile(profileData)

      if (response.success) {
        // 后端返回格式: { success: true, data: {...} }
        const userData = response.data || response.user
        if (userData) {
          user.value = userData
          localStorage.setItem('user', JSON.stringify(userData))
        }
        return { success: true, message: response.message }
      } else {
        return { success: false, errors: [response.message || '更新失败'] }
      }
    } catch (error) {
      console.error('更新资料失败:', error)
      return { success: false, errors: [error.message || '更新资料失败'] }
    } finally {
      loading.value = false
    }
  }

  const changePassword = async (passwordData) => {
    try {
      loading.value = true
      const response = await authApi.changePassword(passwordData)

      if (response.success) {
        return { success: true, message: response.message }
      } else {
        return { success: false, errors: response.errors }
      }
    } catch (error) {
      console.error('修改密码失败:', error)
      return { success: false, errors: [error.message || '修改密码失败'] }
    } finally {
      loading.value = false
    }
  }

  const refreshAccessToken = async () => {
    try {
      if (!refreshToken.value) {
        return false
      }

      const response = await authApi.refreshToken()
      if (response.success) {
        token.value = response.access_token
        localStorage.setItem('token', response.access_token)
        return true
      } else {
        logout()
        return false
      }
    } catch (error) {
      console.error('刷新token失败:', error)
      logout()
      return false
    }
  }

  // 设置重定向路径
  const setRedirectPath = (path) => {
    localStorage.setItem('redirectPath', path)
  }

  // 获取重定向路径
  const getRedirectPath = () => {
    return localStorage.getItem('redirectPath') || '/'
  }

  // 清除重定向路径
  const clearRedirectPath = () => {
    localStorage.removeItem('redirectPath')
  }

  // 初始化时检查认证状态
  const initAuth = async () => {
    await checkAuthStatus()
  }

  // 管理员登录（JWT认证）
  const adminLogin = async (credentials) => {
    try {
      loading.value = true
      const response = await authApi.adminLogin(credentials)

      if (response.success && response.data) {
        const { user: userData, token: accessToken } = response.data
        
        token.value = accessToken
        user.value = userData

        localStorage.setItem('token', accessToken)
        localStorage.setItem('user', JSON.stringify(userData))

        return { success: true, message: response.message }
      } else {
        return { success: false, message: response.message || '登录失败' }
      }
    } catch (error) {
      console.error('管理员登录失败:', error)
      return { success: false, message: error.message || '登录失败' }
    } finally {
      loading.value = false
    }
  }

  // 管理员假登录（兼容旧版）
  const setAdminAuth = (adminUser) => {
    const fakeToken = 'fake-admin-token-' + Date.now()
    user.value = adminUser
    token.value = fakeToken
    localStorage.setItem('token', fakeToken)
    localStorage.setItem('user', JSON.stringify(adminUser))
  }

  return {
    // 状态
    user,
    token,
    refreshToken,
    loading,
    authCheckInProgress,

    // 计算属性
    isAuthenticated,
    isAdmin,
    userInfo,

    // 方法
    login,
    register,
    logout,
    checkAuthStatus,
    updateProfile,
    changePassword,
    refreshAccessToken,
    initAuth,
    setRedirectPath,
    getRedirectPath,
    clearRedirectPath,
    adminLogin,
    setAdminAuth
  }
})