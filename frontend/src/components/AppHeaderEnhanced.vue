<template>
  <header 
    class="app-header-enhanced"
    :class="{ 
      'is-scrolled': isScrolled,
      'is-hidden': isHidden
    }"
  >
    <!-- 进度条 -->
    <div class="scroll-progress-bar" :style="{ width: scrollProgress + '%' }"></div>

    <div class="header-container">
      <!-- Logo 和标题 -->
      <div class="header-brand" @click="goHome">
        <div class="logo-wrapper">
          <el-icon :size="32" class="logo-icon"><Coffee /></el-icon>
          <div class="logo-ripple"></div>
        </div>
        <div class="brand-info">
          <h1 class="brand-title gradient-text">咖啡点餐</h1>
          <p class="brand-subtitle">Coffee Ordering</p>
        </div>
      </div>

      <!-- 搜索框 -->
      <transition name="search-expand">
        <div v-if="showSearch || isSearchExpanded" class="header-search">
          <el-input
            v-model="searchQuery"
            placeholder="搜索您喜欢的咖啡..."
            :prefix-icon="Search"
            clearable
            @focus="handleSearchFocus"
            @blur="handleSearchBlur"
            @input="handleSearchInput"
            class="search-input glass-input"
          >
            <template #suffix>
              <el-button
                v-if="searchQuery"
                :icon="Close"
                circle
                size="small"
                @click="clearSearch"
                class="clear-search-btn"
              />
            </template>
          </el-input>
          
          <!-- 搜索建议 -->
          <transition name="dropdown-fade">
            <div v-if="showSuggestions && suggestions.length > 0" class="search-suggestions glass-panel">
              <div
                v-for="(item, index) in suggestions"
                :key="item.id"
                class="suggestion-item"
                :class="{ 'is-active': activeSuggestionIndex === index }"
                @click="selectSuggestion(item)"
                @mouseenter="activeSuggestionIndex = index"
              >
                <el-icon class="suggestion-icon"><Coffee /></el-icon>
                <span class="suggestion-name">{{ item.name }}</span>
                <span class="suggestion-price">¥{{ item.price }}</span>
              </div>
            </div>
          </transition>
        </div>
      </transition>

      <!-- 导航菜单 -->
      <nav class="header-nav">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="nav-item glass-button"
          :class="{ 'is-active': isActive(item.path) }"
        >
          <el-icon :size="18">
            <component :is="item.icon" />
          </el-icon>
          <span class="nav-text">{{ item.label }}</span>
        </router-link>
      </nav>

      <!-- 右侧操作区 -->
      <div class="header-actions">
        <!-- 主题切换 -->
        <ThemeToggle />

        <!-- 购物车按钮 -->
        <div class="cart-button-wrapper">
          <button
            class="cart-button glass-button"
            @click="toggleCart"
            :class="{ 'has-items': cartStore.totalItems > 0 }"
          >
            <el-icon :size="20"><ShoppingCart /></el-icon>
            <transition name="badge-bounce">
              <span v-if="cartStore.totalItems > 0" class="cart-badge">
                {{ cartStore.totalItems }}
              </span>
            </transition>
          </button>
          
          <!-- 购物车价格提示 -->
          <transition name="price-slide">
            <div v-if="cartStore.totalItems > 0" class="cart-price-tip glass-panel">
              ¥{{ cartStore.totalPrice.toFixed(2) }}
            </div>
          </transition>
        </div>

        <!-- 用户按钮 -->
        <div class="user-button-wrapper">
          <el-dropdown v-if="authStore.isAuthenticated && !authStore.isAdmin" trigger="click">
            <button class="user-button glass-button">
              <el-avatar :size="28" :src="authStore.userInfo?.avatar_url">
                <el-icon><User /></el-icon>
              </el-avatar>
            </button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="goToOrders">
                  <el-icon><List /></el-icon> 我的订单
                </el-dropdown-item>
                <el-dropdown-item @click="goToProfile">
                  <el-icon><User /></el-icon> 个人中心
                </el-dropdown-item>
                <el-dropdown-item @click="goToMember">
                  <el-icon><Trophy /></el-icon> 会员中心
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">
                  <el-icon><SwitchButton /></el-icon> 退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <router-link v-else-if="!authStore.isAdmin" to="/auth/login" class="login-btn glass-button">
            <el-icon><User /></el-icon>
            <span class="login-text">登录</span>
          </router-link>
        </div>

        <!-- 移动端菜单按钮 -->
        <button
          class="mobile-menu-btn glass-button"
          @click="toggleMobileMenu"
          :class="{ 'is-active': showMobileMenu }"
        >
          <span class="menu-icon">
            <span></span>
            <span></span>
            <span></span>
          </span>
        </button>
      </div>
    </div>

    <!-- 移动端菜单 -->
    <transition name="mobile-menu-slide">
      <div v-if="showMobileMenu" class="mobile-menu glass-panel">
        <nav class="mobile-nav">
          <router-link
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            class="mobile-nav-item"
            :class="{ 'is-active': isActive(item.path) }"
            @click="closeMobileMenu"
          >
            <el-icon :size="20">
              <component :is="item.icon" />
            </el-icon>
            <span>{{ item.label }}</span>
          </router-link>
          
          <!-- 用户相关菜单 -->
          <template v-if="authStore.isAuthenticated && !authStore.isAdmin">
            <router-link to="/user/orders" class="mobile-nav-item" @click="closeMobileMenu">
              <el-icon :size="20"><List /></el-icon>
              <span>我的订单</span>
            </router-link>
            <router-link to="/user/profile" class="mobile-nav-item" @click="closeMobileMenu">
              <el-icon :size="20"><User /></el-icon>
              <span>个人中心</span>
            </router-link>
            <router-link to="/user/member" class="mobile-nav-item" @click="closeMobileMenu">
              <el-icon :size="20"><Trophy /></el-icon>
              <span>会员中心</span>
            </router-link>
            <div class="mobile-nav-item logout" @click="handleLogout">
              <el-icon :size="20"><SwitchButton /></el-icon>
              <span>退出登录</span>
            </div>
          </template>
          <template v-else-if="!authStore.isAdmin">
            <router-link to="/auth/login" class="mobile-nav-item" @click="closeMobileMenu">
              <el-icon :size="20"><User /></el-icon>
              <span>登录/注册</span>
            </router-link>
          </template>
        </nav>
      </div>
    </transition>
  </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Coffee, Search, ShoppingCart, Close, Menu as MenuIcon, HomeFilled, User, Trophy, SwitchButton, List } from '@element-plus/icons-vue'
import { useCartStore } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import { usePointsStore } from '@/store/points'
import ThemeToggle from './ThemeToggle.vue'

const router = useRouter()
const route = useRoute()
const cartStore = useCartStore()
const authStore = useAuthStore()
const pointsStore = usePointsStore()

// 响应式数据
const isScrolled = ref(false)
const isHidden = ref(false)
const scrollProgress = ref(0)
const lastScrollY = ref(0)
const searchQuery = ref('')
const isSearchExpanded = ref(false)
const showSuggestions = ref(false)
const suggestions = ref([])
const activeSuggestionIndex = ref(-1)
const showMobileMenu = ref(false)

// 导航项
const navItems = [
  { path: '/menu', label: '菜单', icon: 'Coffee' },
  { path: '/cart', label: '购物车', icon: 'ShoppingCart' }
]

// 计算属性
const showSearch = computed(() => {
  return route.path === '/menu'
})

// 方法
const handleScroll = () => {
  const currentScrollY = window.scrollY
  
  // 滚动进度
  const windowHeight = document.documentElement.scrollHeight - window.innerHeight
  scrollProgress.value = (currentScrollY / windowHeight) * 100
  
  // 滚动状态
  isScrolled.value = currentScrollY > 50
  
  // 自动隐藏/显示
  if (currentScrollY > lastScrollY.value && currentScrollY > 200) {
    isHidden.value = true
  } else {
    isHidden.value = false
  }
  
  lastScrollY.value = currentScrollY
}

const goHome = () => {
  router.push('/menu')
}

const isActive = (path) => {
  return route.path === path
}

const handleSearchFocus = () => {
  isSearchExpanded.value = true
}

const handleSearchBlur = () => {
  setTimeout(() => {
    if (!searchQuery.value) {
      isSearchExpanded.value = false
    }
    showSuggestions.value = false
  }, 200)
}

const handleSearchInput = (value) => {
  if (value) {
    // 模拟搜索建议（实际应该调用API）
    suggestions.value = [
      { id: 1, name: '美式咖啡', price: 18 },
      { id: 2, name: '拿铁咖啡', price: 22 },
      { id: 3, name: '卡布奇诺', price: 24 }
    ].filter(item => item.name.includes(value))
    
    showSuggestions.value = suggestions.value.length > 0
  } else {
    showSuggestions.value = false
  }
}

const clearSearch = () => {
  searchQuery.value = ''
  showSuggestions.value = false
}

const selectSuggestion = (item) => {
  searchQuery.value = item.name
  showSuggestions.value = false
  // 触发搜索或跳转到商品详情
}

const toggleCart = () => {
  // 触发购物车侧边栏
  window.dispatchEvent(new CustomEvent('toggle-cart-sidebar'))
}

const toggleMobileMenu = () => {
  showMobileMenu.value = !showMobileMenu.value
}

const closeMobileMenu = () => {
  showMobileMenu.value = false
}

const goToOrders = () => {
  router.push('/user/orders')
}

const goToProfile = () => {
  router.push('/user/profile')
}

const goToMember = () => {
  router.push('/user/member')
}

const handleLogout = () => {
  authStore.logout()
  pointsStore.clearPointsData()
  ElMessage.success('已退出登录')
  router.push('/menu')
}

// 生命周期
onMounted(() => {
  window.addEventListener('scroll', handleScroll, { passive: true })
  handleScroll()
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

// 监听路由变化关闭移动菜单
watch(() => route.path, () => {
  closeMobileMenu()
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';
@import '@/styles/mixins.scss';

.app-header-enhanced {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: $z-sticky;
  @include glass-effect(0.8, 20px);
  box-shadow: $shadow-sm;
  transition: all $transition-base $ease-out;
  
  &.is-scrolled {
    box-shadow: $shadow-md;
    
    .header-brand {
      .logo-icon {
        transform: scale(0.9);
      }
      
      .brand-subtitle {
        opacity: 0;
        height: 0;
      }
    }
  }
  
  &.is-hidden {
    transform: translateY(-100%);
  }
  
  @include dark-mode {
    @include glass-dark(0.8, 20px);
  }
}

.scroll-progress-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  height: 3px;
  background: $primary-gradient;
  transition: width $transition-fast linear;
  box-shadow: 0 0 10px rgba($primary-color, 0.5);
}

.header-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: $spacing-md $spacing-lg;
  display: flex;
  align-items: center;
  gap: $spacing-lg;
  
  @include respond-to(md) {
    padding: $spacing-sm $spacing-md;
    gap: $spacing-sm;
  }
}

.header-brand {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  cursor: pointer;
  transition: all $transition-base $ease-out;
  
  &:hover {
    .logo-icon {
      transform: rotate(10deg) scale(1.1);
    }
    
    .logo-ripple {
      transform: scale(1.5);
      opacity: 0;
    }
  }
  
  .logo-wrapper {
    position: relative;
    
    .logo-icon {
      color: $primary-color;
      transition: all $transition-base $ease-out;
      position: relative;
      z-index: 1;
    }
    
    .logo-ripple {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      width: 100%;
      height: 100%;
      border-radius: 50%;
      background: radial-gradient(circle, rgba($primary-color, 0.2) 0%, transparent 70%);
      transition: all $transition-base $ease-out;
    }
  }
  
  .brand-info {
    .brand-title {
      margin: 0;
      font-size: $font-size-xl;
      font-weight: $font-weight-bold;
      line-height: 1;
      @include gradient-text($primary-gradient);
    }
    
    .brand-subtitle {
      margin: 0;
      font-size: $font-size-xs;
      color: $text-secondary;
      text-transform: uppercase;
      letter-spacing: 1px;
      transition: all $transition-base $ease-out;
    }
  }
  
  @include respond-to(md) {
    .brand-info {
      display: none;
    }
  }
}

.header-search {
  flex: 1;
  max-width: 500px;
  position: relative;
  
  .search-input {
    :deep(.el-input__wrapper) {
      @include glass-effect(0.1, 8px);
      border-radius: $radius-full;
      transition: all $transition-base $ease-out;
      
      &:hover,
      &.is-focus {
        @include glass-effect(0.15, 10px);
        box-shadow: $shadow-md;
      }
    }
  }
  
  @include respond-to(md) {
    display: none;
  }
}

.search-suggestions {
  position: absolute;
  top: calc(100% + $spacing-sm);
  left: 0;
  right: 0;
  border-radius: $radius-md;
  box-shadow: $shadow-lg;
  max-height: 300px;
  overflow-y: auto;
  @include custom-scrollbar;
  
  .suggestion-item {
    display: flex;
    align-items: center;
    gap: $spacing-sm;
    padding: $spacing-sm $spacing-md;
    cursor: pointer;
    transition: all $transition-fast $ease-out;
    
    &:hover,
    &.is-active {
      background: rgba($primary-color, 0.1);
    }
    
    .suggestion-icon {
      color: $primary-color;
    }
    
    .suggestion-name {
      flex: 1;
      font-weight: $font-weight-medium;
    }
    
    .suggestion-price {
      color: $danger-color;
      font-weight: $font-weight-bold;
    }
  }
}

.header-nav {
  display: flex;
  gap: $spacing-sm;
  
  .nav-item {
    display: flex;
    align-items: center;
    gap: $spacing-xs;
    padding: $spacing-sm $spacing-md;
    border-radius: $radius-full;
    text-decoration: none;
    color: $text-primary;
    font-weight: $font-weight-medium;
    transition: all $transition-base $ease-out;
    border: 1px solid transparent;
    
    &:hover {
      background: rgba($primary-color, 0.1);
      transform: translateY(-2px);
    }
    
    &.is-active {
      background: $primary-gradient;
      color: white;
      box-shadow: $shadow-md;
    }
  }
  
  @include respond-to(md) {
    display: none;
  }
}

.header-actions {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.cart-button-wrapper {
  position: relative;
  
  .cart-button {
    position: relative;
    width: 40px;
    height: 40px;
    border-radius: $radius-full;
    border: 1px solid rgba(255, 255, 255, 0.2);
    background: transparent;
    cursor: pointer;
    @include flex-center;
    transition: all $transition-base $ease-out;
    color: $primary-color;
    
    &:hover {
      transform: scale(1.1);
      box-shadow: $shadow-md;
    }
    
    &.has-items {
      animation: pulse 2s $ease-in-out infinite;
    }
    
    .cart-badge {
      position: absolute;
      top: -4px;
      right: -4px;
      background: $danger-color;
      color: white;
      font-size: $font-size-xs;
      font-weight: $font-weight-bold;
      padding: 2px 6px;
      border-radius: $radius-full;
      min-width: 18px;
      text-align: center;
      box-shadow: $shadow-sm;
    }
  }
  
  .cart-price-tip {
    position: absolute;
    top: calc(100% + $spacing-xs);
    right: 0;
    padding: $spacing-xs $spacing-sm;
    border-radius: $radius-sm;
    font-size: $font-size-xs;
    font-weight: $font-weight-bold;
    color: $danger-color;
    white-space: nowrap;
    box-shadow: $shadow-sm;
  }
}

.user-button-wrapper {
  .user-button {
    width: 40px;
    height: 40px;
    border-radius: $radius-full;
    border: 1px solid rgba(255, 255, 255, 0.2);
    background: transparent;
    cursor: pointer;
    @include flex-center;
    transition: all $transition-base $ease-out;
    padding: 0;
    
    &:hover {
      transform: scale(1.1);
      box-shadow: $shadow-md;
    }
  }
  
  .login-btn {
    display: flex;
    align-items: center;
    gap: $spacing-xs;
    padding: $spacing-sm $spacing-md;
    border-radius: $radius-full;
    text-decoration: none;
    color: $primary-color;
    font-weight: $font-weight-medium;
    transition: all $transition-base $ease-out;
    border: 1px solid rgba($primary-color, 0.3);
    background: transparent;
    
    &:hover {
      background: rgba($primary-color, 0.1);
      transform: translateY(-2px);
    }
    
    .login-text {
      @include respond-to(sm) {
        display: none;
      }
    }
  }
}

.mobile-menu-btn {
  display: none;
  width: 40px;
  height: 40px;
  border-radius: $radius-sm;
  border: none;
  background: transparent;
  cursor: pointer;
  @include flex-center;
  
  .menu-icon {
    width: 20px;
    height: 16px;
    position: relative;
    
    span {
      display: block;
      position: absolute;
      height: 2px;
      width: 100%;
      background: $primary-color;
      border-radius: 2px;
      transition: all $transition-base $ease-out;
      
      &:nth-child(1) {
        top: 0;
      }
      
      &:nth-child(2) {
        top: 50%;
        transform: translateY(-50%);
      }
      
      &:nth-child(3) {
        bottom: 0;
      }
    }
  }
  
  &.is-active {
    .menu-icon span {
      &:nth-child(1) {
        top: 50%;
        transform: translateY(-50%) rotate(45deg);
      }
      
      &:nth-child(2) {
        opacity: 0;
      }
      
      &:nth-child(3) {
        bottom: 50%;
        transform: translateY(50%) rotate(-45deg);
      }
    }
  }
  
  @include respond-to(md) {
    display: flex;
  }
}

.mobile-menu {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  padding: $spacing-md;
  box-shadow: $shadow-lg;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-top: 1px solid rgba(0, 0, 0, 0.1);
  z-index: 1000;
  
  .mobile-nav {
    display: flex;
    flex-direction: column;
    gap: $spacing-xs;
    
    .mobile-nav-item {
      display: flex;
      align-items: center;
      gap: $spacing-sm;
      padding: $spacing-md;
      border-radius: $radius-md;
      text-decoration: none;
      color: $text-primary;
      font-weight: $font-weight-medium;
      transition: all $transition-base $ease-out;
      cursor: pointer;
      
      &:hover {
        background: rgba($primary-color, 0.1);
      }
      
      &.is-active {
        background: $primary-gradient;
        color: white;
      }
      
      &.logout {
        color: $danger-color;
        border-top: 1px solid $border-color;
        margin-top: $spacing-sm;
        padding-top: $spacing-md;
        
        &:hover {
          background: rgba($danger-color, 0.1);
        }
      }
    }
  }
}

// 动画
.search-expand-enter-active,
.search-expand-leave-active {
  transition: all $transition-base $ease-out;
}

.search-expand-enter-from,
.search-expand-leave-to {
  opacity: 0;
  transform: scaleX(0);
}

.dropdown-fade-enter-active,
.dropdown-fade-leave-active {
  transition: all $transition-base $ease-out;
}

.dropdown-fade-enter-from,
.dropdown-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.badge-bounce-enter-active {
  animation: bounceIn 0.5s $ease-bounce;
}

.price-slide-enter-active,
.price-slide-leave-active {
  transition: all $transition-base $ease-out;
}

.price-slide-enter-from,
.price-slide-leave-to {
  opacity: 0;
  transform: translateY(-5px);
}

.mobile-menu-slide-enter-active,
.mobile-menu-slide-leave-active {
  transition: all $transition-base $ease-out;
}

.mobile-menu-slide-enter-from,
.mobile-menu-slide-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

// 移动端深度优化
@include respond-to(md) {
  .header-container {
    padding: $spacing-sm $spacing-sm;
  }
  
  .header-brand {
    gap: $spacing-sm;
    
    .logo-wrapper .logo-icon {
      font-size: 24px !important;
    }
  }
  
  .cart-button-wrapper {
    .cart-button {
      width: 36px;
      height: 36px;
      
      .cart-badge {
        font-size: 10px;
        padding: 1px 4px;
        min-width: 14px;
      }
    }
    
    .cart-price-tip {
      display: none;
    }
  }
  
  .user-button-wrapper {
    .user-button {
      width: 36px;
      height: 36px;
    }
    
    .login-btn {
      padding: $spacing-xs $spacing-sm;
      font-size: $font-size-sm;
    }
  }
  
  .mobile-menu-btn {
    width: 36px;
    height: 36px;
  }
}

// 移动端菜单优化
@include respond-to(sm) {
  .mobile-menu {
    .mobile-nav {
      .mobile-nav-item {
        padding: $spacing-md $spacing-sm;
        font-size: $font-size-base;
        
        .el-icon {
          font-size: 18px !important;
        }
      }
    }
  }
}

// 深色模式移动端适配 - 使用 CSS 选择器而非 mixin
html[data-theme='dark'] {
  .mobile-menu {
    background: rgba(45, 45, 45, 0.95);
    border-top: 1px solid #404040;
    
    .mobile-nav-item {
      color: #E0E0E0;
      
      &:hover {
        background: rgba(139, 115, 85, 0.2);
      }
      
      &.is-active {
        background: linear-gradient(135deg, #8B7355 0%, #C4A57B 100%);
        color: white;
      }
      
      &.logout {
        border-top-color: #404040;
      }
    }
  }
}
</style>
