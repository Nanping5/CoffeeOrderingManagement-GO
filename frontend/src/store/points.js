import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import pointsApi from '@/api/points'

export const usePointsStore = defineStore('points', () => {
  // 状态
  const pointsInfo = ref(null)
  const transactions = ref([])
  const loading = ref(false)
  const pagination = ref({
    page: 1,
    per_page: 20,
    total: 0,
    pages: 0
  })

  // 计算属性
  const totalPoints = computed(() => pointsInfo.value?.total_points || 0) // 当前可用积分
  const lifetimePoints = computed(() => pointsInfo.value?.lifetime_points || 0)
  const memberLevel = computed(() => pointsInfo.value?.member_level || 'bronze')
  
  // 会员等级显示名称
  const memberLevelName = computed(() => {
    const levelMap = {
      bronze: '铜牌会员',
      silver: '银牌会员',
      gold: '金牌会员',
      platinum: '白金会员'
    }
    return levelMap[memberLevel.value] || '铜牌会员'
  })

  // 下一等级信息
  const nextLevel = computed(() => pointsInfo.value?.next_level || null)
  
  // 等级权益
  const levelBenefits = computed(() => pointsInfo.value?.level_benefits || {})

  // 获取积分信息
  const fetchPointsInfo = async () => {
    try {
      loading.value = true
      const response = await pointsApi.getPointsInfo()
      if (response.success) {
        pointsInfo.value = response.data
      }
      return response
    } catch (error) {
      console.error('获取积分信息失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取积分历史
  const fetchTransactions = async (params = {}) => {
    try {
      loading.value = true
      const response = await pointsApi.getTransactions({
        page: pagination.value.page,
        per_page: pagination.value.per_page,
        ...params
      })

      if (response.success) {
        transactions.value = response.data.transactions || []
        pagination.value = {
          page: response.data.page || 1,
          per_page: response.data.per_page || 20,
          total: response.data.total || 0,
          pages: response.data.pages || 0
        }
      }
      return response
    } catch (error) {
      console.error('获取积分记录失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 计算积分抵扣
  const calculatePointsDiscount = async (orderData) => {
    try {
      const response = await pointsApi.calculateDiscount(orderData)
      return response
    } catch (error) {
      console.error('计算积分抵扣失败:', error)
      throw error
    }
  }

  // 刷新积分数据
  const refreshPointsData = async () => {
    await Promise.all([
      fetchPointsInfo(),
      fetchTransactions()
    ])
  }

  // 清除积分数据
  const clearPointsData = () => {
    pointsInfo.value = null
    transactions.value = []
    pagination.value = { page: 1, per_page: 20, total: 0, pages: 0 }
  }

  // 格式化积分变动
  const formatPointsChange = (change) => {
    if (change > 0) return `+${change}`
    return change.toString()
  }

  // 获取交易类型文本
  const getTransactionTypeText = (type) => {
    const typeMap = {
      earned: '获得',
      used: '使用',
      expired: '过期',
      refunded: '退还',
      signup_bonus: '注册奖励',
      birthday_bonus: '生日奖励'
    }
    return typeMap[type] || type
  }

  // 获取交易类型颜色
  const getTransactionTypeColor = (type) => {
    const colorMap = {
      earned: 'success',
      used: 'warning',
      expired: 'info',
      refunded: 'primary',
      signup_bonus: 'success',
      birthday_bonus: 'success'
    }
    return colorMap[type] || 'info'
  }

  return {
    pointsInfo,
    transactions,
    loading,
    pagination,
    totalPoints,
    lifetimePoints,
    memberLevel,
    memberLevelName,
    nextLevel,
    levelBenefits,
    fetchPointsInfo,
    fetchTransactions,
    calculatePointsDiscount,
    refreshPointsData,
    clearPointsData,
    formatPointsChange,
    getTransactionTypeText,
    getTransactionTypeColor
  }
})
