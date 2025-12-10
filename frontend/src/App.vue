<template>
  <div id="app">
    <router-view />
    
    <!-- 全局购物车侧边栏 -->
    <CartSidebar ref="cartSidebarRef" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import CartSidebar from '@/components/CartSidebar.vue'

const cartSidebarRef = ref(null)

// 监听全局购物车切换事件
onMounted(() => {
  window.addEventListener('toggle-cart-sidebar', () => {
    cartSidebarRef.value?.toggle()
  })
})
</script>

<style lang="scss">
// 全局样式已在 styles/index.scss 中定义
</style>