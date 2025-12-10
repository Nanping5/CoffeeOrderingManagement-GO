import request from './index'

// 积分相关API
const pointsAPI = {
  // 获取积分信息
  getPointsInfo() {
    return request.get('/user/points')
  },

  // 获取积分历史记录
  getTransactions(params = {}) {
    return request.get('/user/points/transactions', params)
  },

  // 计算积分抵扣
  calculateDiscount(orderData) {
    return request.post('/orders/points-calculation', orderData)
  }
}

export default pointsAPI
