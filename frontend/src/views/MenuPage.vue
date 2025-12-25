<template>
  <div class="menu-page">
    <!-- 页面头部 -->
    <div class="page-header animate__animated animate__fadeInDown">
      <div class="header-content">
        <h1 class="page-title">
          <el-icon class="title-icon"><Coffee /></el-icon>
          咖啡菜单
        </h1>
        <p class="page-subtitle">精选优质咖啡，为您带来美好的味蕾体验</p>
      </div>
    </div>

    <!-- 分类筛选区域 - 现代标签式设计 -->
    <div class="filter-section animate__animated animate__fadeInUp">
      <div class="filter-wrapper">
        <!-- 分类标签 -->
        <div class="category-tabs">
          <button
            v-for="category in categories"
            :key="category.value"
            class="category-tab"
            :class="{ active: selectedCategory === category.value }"
            @click="selectedCategory = category.value; handleCategoryChange()"
          >
            <span class="tab-icon">{{ getCategoryIcon(category.value) }}</span>
            <span class="tab-label">{{ category.label }}</span>
          </button>
        </div>

        <!-- 排序下拉 -->
        <div class="sort-dropdown">
          <el-dropdown trigger="click" @command="(cmd) => { sortBy = cmd; handleSort() }">
            <button class="sort-btn">
              <el-icon><Sort /></el-icon>
              <span>{{ getSortLabel(sortBy) }}</span>
              <el-icon class="arrow"><ArrowDown /></el-icon>
            </button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="default" :class="{ active: sortBy === 'default' }">默认排序</el-dropdown-item>
                <el-dropdown-item command="price_asc" :class="{ active: sortBy === 'price_asc' }">价格从低到高</el-dropdown-item>
                <el-dropdown-item command="price_desc" :class="{ active: sortBy === 'price_desc' }">价格从高到低</el-dropdown-item>
                <el-dropdown-item command="popular" :class="{ active: sortBy === 'popular' }">最受欢迎</el-dropdown-item>
                <el-dropdown-item command="newest" :class="{ active: sortBy === 'newest' }">最新上架</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </div>

    <!-- 菜单内容区域 -->
    <div class="menu-content">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="6" animated />
      </div>

      <!-- 空状态 -->
      <div v-else-if="menuItems.length === 0" class="empty-state">
        <el-empty description="暂无商品">
          <template #image>
            <el-icon size="60"><Coffee /></el-icon>
          </template>
          <el-button type="primary" @click="resetFilters">重置筛选</el-button>
        </el-empty>
      </div>

      <!-- 菜单网格 - 使用增强版卡片 -->
      <div v-else class="menu-grid">
        <CoffeeCardEnhanced
          v-for="item in menuItems"
          :key="item.id"
          :item="item"
          :category-label="getCategoryLabel(item.category)"
          @add-to-cart="handleAddToCart"
          @view-detail="handleViewDetail"
          @notify-when-available="handleNotifyWhenAvailable"
          class="stagger-item"
        />
      </div>

      <!-- 加载更多 -->
      <div v-if="hasMore && !loading" class="load-more">
        <el-button
          type="primary"
          :loading="loadingMore"
          @click="loadMore"
          class="load-more-btn"
        >
          {{ loadingMore ? '加载中...' : '加载更多' }}
        </el-button>
      </div>
    </div>

    <!-- 商品详情弹窗 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="selectedItem?.name"
      width="600px"
      class="detail-dialog"
      destroy-on-close
    >
      <div v-if="selectedItem" class="item-detail">
        <div class="detail-image">
          <img :src="selectedItem.image_url || '/default-coffee.jpg'" :alt="selectedItem.name" />
        </div>
        <div class="detail-info">
          <h3>{{ selectedItem.name }}</h3>
          <p class="description">{{ selectedItem.description }}</p>
          <div class="price-info">
            <span class="price">¥{{ selectedItem.price }}</span>
            <span v-if="selectedItem.original_price" class="original-price">
              ¥{{ selectedItem.original_price }}
            </span>
          </div>
          <div class="category-tag">
            <el-tag type="info">{{ getCategoryLabel(selectedItem.category) }}</el-tag>
          </div>
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="detailDialogVisible = false">取消</el-button>
          <el-button
            type="primary"
            @click="handleAddToCartFromDetail"
            :loading="addingToCart"
          >
            <el-icon><ShoppingCart /></el-icon>
            加入购物车
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Coffee, ShoppingCart, Search, Sort, ArrowDown } from '@element-plus/icons-vue'
import { useCartStore } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import menuAPI from '@/api/menu'
import CoffeeCardEnhanced from '@/components/CoffeeCardEnhanced.vue'
import { useScrollReveal, useStaggerAnimation } from '@/composables/useScrollReveal'

// 状态管理
const cartStore = useCartStore()
const authStore = useAuthStore()

// 使用滚动动画
useScrollReveal({ threshold: 0.1, once: true })
useStaggerAnimation('.stagger-item', 50)

// 响应式数据
const loading = ref(false)
const loadingMore = ref(false)
const menuItems = ref([])
const categories = ref([])
const searchQuery = ref('')
const selectedCategory = ref('all')
const sortBy = ref('default')
const currentPage = ref(1)
const pageSize = ref(12)
const hasMore = ref(true)
const detailDialogVisible = ref(false)
const selectedItem = ref(null)
const addingToCart = ref(false)

// 分页数据
const pagination = reactive({
  current: 1,
  size: 12,
  total: 0
})

// 计算属性
const filteredItems = computed(() => {
  let items = menuItems.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    items = items.filter(item =>
      item.name.toLowerCase().includes(query) ||
      item.description.toLowerCase().includes(query)
    )
  }

  // 分类过滤
  if (selectedCategory.value !== 'all') {
    items = items.filter(item => item.category === selectedCategory.value)
  }

  // 排序
  switch (sortBy.value) {
    case 'price_asc':
      items.sort((a, b) => a.price - b.price)
      break
    case 'price_desc':
      items.sort((a, b) => b.price - a.price)
      break
    case 'popular':
      items.sort((a, b) => (b.order_count || 0) - (a.order_count || 0))
      break
    case 'newest':
      items.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
      break
  }

  return items
})

// 方法
const getMenuList = async (reset = false) => {
  if (reset) {
    currentPage.value = 1
    menuItems.value = []
    hasMore.value = true
  }

  loading.value = true

  try {
    const params = {
      page: currentPage.value,
      per_page: pageSize.value,
      category: selectedCategory.value === 'all' ? undefined : selectedCategory.value,
      keyword: searchQuery.value || undefined
    }

    const response = await menuAPI.getMenuList(params)

    if (reset) {
      menuItems.value = response.data.items || []
    } else {
      menuItems.value.push(...(response.data.items || []))
    }

    pagination.total = response.data.total || 0
    hasMore.value = menuItems.value.length < pagination.total

  } catch (error) {
    console.error('获取菜单列表失败:', error)
    ElMessage.error('获取菜单列表失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

const getCategories = async () => {
  try {
    const response = await menuAPI.getCategories()
    // 后端已经返回包含"全部"的分类列表，直接使用
    categories.value = response.data || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
    // 设置默认分类
    categories.value = [
      { value: 'all', label: '全部' },
      { value: 'coffee', label: '咖啡' },
      { value: 'tea', label: '茶饮' },
      { value: 'dessert', label: '甜点' },
      { value: 'other', label: '其他' }
    ]
  }
}

const handleSearch = (query) => {
  // 防抖处理
  clearTimeout(handleSearch.debounceTimer)
  handleSearch.debounceTimer = setTimeout(() => {
    getMenuList(true)
  }, 500)
}

const handleCategoryChange = () => {
  getMenuList(true)
}

const handleSort = () => {
  getMenuList(true)
}

const loadMore = async () => {
  if (loadingMore.value || !hasMore.value) return

  loadingMore.value = true
  currentPage.value++

  try {
    const params = {
      page: currentPage.value,
      per_page: pageSize.value,
      category: selectedCategory.value === 'all' ? undefined : selectedCategory.value,
      keyword: searchQuery.value || undefined
    }

    const response = await menuAPI.getMenuList(params)
    menuItems.value.push(...(response.data.items || []))
    hasMore.value = menuItems.value.length < pagination.total

  } catch (error) {
    console.error('加载更多失败:', error)
    ElMessage.error('加载更多失败，请稍后重试')
    currentPage.value-- // 恢复页码
  } finally {
    loadingMore.value = false
  }
}

const resetFilters = () => {
  searchQuery.value = ''
  selectedCategory.value = 'all'
  sortBy.value = 'default'
  getMenuList(true)
}

const handleAddToCart = async (item, quantity = 1) => {
  // 无需登录验证，直接添加到购物车
  try {
    await cartStore.addItem(item, quantity)
    ElMessage.success(`${item.name} 已添加到购物车`)
  } catch (error) {
    console.error('添加到购物车失败:', error)
    ElMessage.error('添加到购物车失败，请稍后重试')
  }
}

const handleViewDetail = (item) => {
  selectedItem.value = item
  detailDialogVisible.value = true
}

const handleAddToCartFromDetail = async () => {
  if (!selectedItem.value) return

  addingToCart.value = true

  try {
    await handleAddToCart(selectedItem.value)
    detailDialogVisible.value = false
  } finally {
    addingToCart.value = false
  }
}

const getCategoryLabel = (categoryValue) => {
  const category = categories.value.find(cat => cat.value === categoryValue)
  return category ? category.label : categoryValue
}

const handleNotifyWhenAvailable = (item) => {
  ElMessage.success('我们会在商品到货时通知您')
}

// 获取分类图标
const getCategoryIcon = (category) => {
  const icons = { all: '🍽️', coffee: '☕', tea: '🍵', dessert: '🍰', other: '🥤' }
  return icons[category] || '📦'
}

// 获取排序标签
const getSortLabel = (sort) => {
  const labels = { 
    default: '默认排序', 
    price_asc: '价格从低到高', 
    price_desc: '价格从高到低', 
    popular: '最受欢迎', 
    newest: '最新上架' 
  }
  return labels[sort] || '排序'
}

// 生命周期
onMounted(async () => {
  await Promise.all([
    getCategories(),
    getMenuList(true)
  ])
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';
@import '@/styles/mixins.scss';

.menu-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 30px 20px 20px;
  
  @include dark-mode {
    background: linear-gradient(135deg, $dark-bg 0%, $dark-bg-secondary 100%);
  }
  
  @media (max-width: 768px) {
    padding: 20px 15px 15px;
  }
}

.page-header {
  text-align: center;
  margin-bottom: 30px;

  .header-content {
    max-width: 800px;
    margin: 0 auto;
  }

  .page-title {
    font-size: 2.5rem;
    font-weight: bold;
    color: #2c3e50;
    margin-bottom: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 15px;
    
    @include dark-mode {
      color: $dark-text;
    }

    .title-icon {
      font-size: 2.2rem;
      color: #8b4513;
      
      @include dark-mode {
        color: $accent-color;
      }
    }
  }

  .page-subtitle {
    font-size: 1.1rem;
    color: #7f8c8d;
    margin: 0;
    
    @include dark-mode {
      color: rgba($dark-text, 0.7);
    }
  }
}

.filter-section {
  margin-bottom: 30px;
  max-width: 1200px;
  margin-left: auto;
  margin-right: auto;

  .filter-wrapper {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 20px;
    flex-wrap: wrap;
  }

  .category-tabs {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
    flex: 1;

    .category-tab {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 12px 20px;
      border: none;
      border-radius: 50px;
      background: rgba(255, 255, 255, 0.9);
      color: #666;
      font-size: 14px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
      
      @include dark-mode {
        background: rgba($dark-bg-secondary, 0.9);
        color: rgba($dark-text, 0.8);
      }

      .tab-icon {
        font-size: 18px;
      }

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 16px rgba(139, 69, 19, 0.2);
        background: rgba(139, 69, 19, 0.1);
        color: #8b4513;
      }

      &.active {
        background: linear-gradient(135deg, #8b4513 0%, #a0522d 100%);
        color: white;
        box-shadow: 0 4px 16px rgba(139, 69, 19, 0.4);
        transform: translateY(-2px);

        .tab-icon {
          transform: scale(1.1);
        }
      }
    }
  }

  .sort-dropdown {
    .sort-btn {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 12px 20px;
      border: none;
      border-radius: 50px;
      background: rgba(255, 255, 255, 0.9);
      color: #666;
      font-size: 14px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.3s ease;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
      
      @include dark-mode {
        background: rgba($dark-bg-secondary, 0.9);
        color: rgba($dark-text, 0.8);
      }

      &:hover {
        background: rgba(139, 69, 19, 0.1);
        color: #8b4513;
      }

      .arrow {
        font-size: 12px;
        transition: transform 0.3s ease;
      }
    }
  }
}

.menu-content {
  max-width: 1200px;
  margin: 0 auto;
}

.loading-container {
  padding: 40px;
}

.empty-state {
  padding: 60px 20px;
  text-align: center;
  
  :deep(.el-empty__description) {
    @include dark-mode {
      color: rgba($dark-text, 0.7);
    }
  }
  
  :deep(.el-icon) {
    @include dark-mode {
      color: rgba($dark-text, 0.5);
    }
  }
}

.menu-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 32px;
  margin-bottom: 40px;
  
  @media (max-width: 768px) {
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 24px;
  }
}

.load-more {
  text-align: center;
  margin-top: 40px;

  .load-more-btn {
    padding: 12px 30px;
    font-size: 1rem;
    border-radius: 25px;
  }
}

.detail-dialog {
  :deep(.el-dialog) {
    border-radius: 12px;
    
    @include dark-mode {
      background-color: $dark-bg-secondary;
      border: 1px solid $dark-border;
    }
  }
  
  :deep(.el-dialog__header) {
    @include dark-mode {
      border-bottom: 1px solid $dark-border;
    }
  }
  
  :deep(.el-dialog__title) {
    @include dark-mode {
      color: $dark-text;
    }
  }
  
  :deep(.el-dialog__body) {
    @include dark-mode {
      color: $dark-text;
    }
  }

  .item-detail {
    display: flex;
    gap: 20px;
    margin-bottom: 20px;

    .detail-image {
      flex-shrink: 0;
      width: 200px;
      height: 200px;
      border-radius: 8px;
      overflow: hidden;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .detail-info {
      flex: 1;

      h3 {
        margin: 0 0 15px 0;
        color: #2c3e50;
        font-size: 1.5rem;
        
        @include dark-mode {
          color: $dark-text;
        }
      }

      .description {
        color: #7f8c8d;
        line-height: 1.6;
        margin-bottom: 20px;
        
        @include dark-mode {
          color: rgba($dark-text, 0.8);
        }
      }

      .price-info {
        margin-bottom: 15px;

        .price {
          font-size: 1.8rem;
          font-weight: bold;
          color: #e74c3c;
          
          @include dark-mode {
            color: $accent-color;
          }
        }

        .original-price {
          font-size: 1.2rem;
          color: #95a5a6;
          text-decoration: line-through;
          margin-left: 10px;
          
          @include dark-mode {
            color: rgba($dark-text, 0.5);
          }
        }
      }

      .category-tag {
        margin-bottom: 20px;
      }
    }
  }

  .dialog-footer {
    .el-button {
      border-radius: 20px;
      padding: 10px 20px;
    }
  }
}

// 响应式设计 - 移动端适配
@media (max-width: 768px) {
  .menu-page {
    padding: 15px 10px;
  }

  .page-header {
    margin-bottom: 20px;
    
    .page-title {
      font-size: 1.8rem !important;
      gap: 10px;

      .title-icon {
        font-size: 1.5rem !important;
      }
    }
    
    .page-subtitle {
      font-size: 0.9rem;
    }
  }

  .filter-section {
    margin-bottom: 20px;
    
    .filter-wrapper {
      flex-direction: column;
      gap: 12px;
    }
    
    .category-tabs {
      width: 100%;
      overflow-x: auto;
      flex-wrap: nowrap;
      padding-bottom: 8px;
      -webkit-overflow-scrolling: touch;
      scrollbar-width: none;
      
      &::-webkit-scrollbar {
        display: none;
      }
      
      .category-tab {
        flex-shrink: 0;
        padding: 10px 16px;
        font-size: 13px;
        
        .tab-icon {
          font-size: 16px;
        }
      }
    }
    
    .sort-dropdown {
      width: 100%;
      
      .sort-btn {
        width: 100%;
        justify-content: center;
      }
    }
  }

  .menu-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }
  
  .load-more {
    margin-top: 24px;
    
    .load-more-btn {
      width: 100%;
      padding: 12px 20px;
    }
  }
}

// 超小屏幕适配 (< 480px)
@media (max-width: 480px) {
  .menu-page {
    padding: 10px 8px;
  }
  
  .page-header {
    .page-title {
      font-size: 1.5rem !important;
    }
    
    .page-subtitle {
      font-size: 0.85rem;
    }
  }
  
  .menu-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .filter-section .category-tabs .category-tab {
    padding: 8px 14px;
    font-size: 12px;
  }
}

// 商品详情弹窗移动端适配
@media (max-width: 768px) {
  .detail-dialog {
    :deep(.el-dialog) {
      width: 95% !important;
      margin: 10px auto;
      max-height: 90vh;
      overflow-y: auto;
    }
    
    .item-detail {
      flex-direction: column;
      gap: 16px;

      .detail-image {
        width: 100%;
        height: 200px;
      }
      
      .detail-info {
        h3 {
          font-size: 1.3rem;
        }
        
        .description {
          font-size: 0.9rem;
        }
        
        .price-info .price {
          font-size: 1.5rem;
        }
      }
    }
    
    .dialog-footer {
      display: flex;
      gap: 10px;
      
      .el-button {
        flex: 1;
      }
    }
  }
}
</style>