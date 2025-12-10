import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import authApi from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  // 状态
  const profile = ref(null)
  const loading = ref(false)

  // 计算属性
  const fullName = computed(() => {
    if (!profile.value) return ''
    return `${profile.value.first_name || ''} ${profile.value.last_name || ''}`.trim()
  })

  const hasAvatar = computed(() => !!profile.value?.avatar_url)

  // 获取用户资料
  const fetchProfile = async () => {
    try {
      loading.value = true
      const response = await authApi.getProfile()
      if (response.success) {
        profile.value = response.data
      }
      return response
    } catch (error) {
      console.error('获取用户信息失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新用户资料
  const updateProfile = async (profileData) => {
    try {
      loading.value = true
      const response = await authApi.updateProfile(profileData)
      if (response.success) {
        profile.value = { ...profile.value, ...response.data }
      }
      return response
    } catch (error) {
      console.error('更新用户信息失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 修改密码
  const changePassword = async (passwordData) => {
    try {
      loading.value = true
      const response = await authApi.changePassword(passwordData)
      return response
    } catch (error) {
      console.error('修改密码失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 清除用户资料
  const clearProfile = () => {
    profile.value = null
  }

  return {
    profile,
    loading,
    fullName,
    hasAvatar,
    fetchProfile,
    updateProfile,
    changePassword,
    clearProfile
  }
})
