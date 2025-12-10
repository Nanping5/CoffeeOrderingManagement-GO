import request from './index'

// 导出对象和具名导出
const menuAPI = {
  // 获取菜单列表
  getMenuItems(params = {}) {
    return request.get('/menu', params)
  },

  // 获取菜单列表（别名，兼容旧代码）
  getMenuList(params = {}) {
    return request.get('/menu', params)
  },

  // 根据分类获取菜单
  getMenuByCategory(category, params = {}) {
    return request.get(`/menu/category/${category}`, params)
  },

  // 获取菜单项详情
  getMenuItem(id) {
    return request.get(`/menu/${id}`)
  },

  // 搜索菜单项
  searchMenuItems(params) {
    return request.get('/menu/search', params)
  },

  // 获取菜单分类
  getCategories() {
    return request.get('/menu/categories')
  },

  // 获取热门商品
  getPopularItems(params = {}) {
    return request.get('/menu/popular', params)
  },

  // 创建菜单项（管理员）
  createMenuItem(data) {
    return request.post('/admin/menu', data)
  },

  // 更新菜单项（管理员）
  updateMenuItem(id, data) {
    return request.put(`/admin/menu/${id}`, data)
  },

  // 删除菜单项（管理员）
  deleteMenuItem(id) {
    return request.delete(`/admin/menu/${id}`)
  },

  // 切换菜单项可用状态（管理员）
  toggleMenuItemStatus(id) {
    return request.patch(`/admin/menu/${id}/toggle`)
  },

  // 上传菜单图片（管理员）
  uploadMenuItemImage(id, file) {
    const formData = new FormData()
    formData.append('image', file)
    return request.upload(`/admin/menu/${id}`, formData)
  },

  // 获取菜单统计信息（管理员）
  getMenuStatistics() {
    return request.get('/admin/menu/stats')
  }
}

export default menuAPI
export { menuAPI }