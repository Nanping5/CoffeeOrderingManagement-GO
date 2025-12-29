<template>
  <div class="cart-sidebar-wrapper">
    <!-- 购物车浮动按钮 -->
    <transition name="bounce">
      <div 
        v-if="!isOpen && cartStore.totalItems > 0"
        class="cart-float-button"
        @click="toggleSidebar"
      >
        <div class="cart-icon-wrapper">
          <el-icon :size="24"><ShoppingCart /></el-icon>
          <transition name="badge-pop">
            <span v-if="cartStore.totalItems > 0" class="cart-badge">
              {{ cartStore.totalItems }}
            </span>
          </transition>
        </div>
        <div class="cart-price">
          ¥{{ cartStore.totalPrice.toFixed(2) }}
        </div>
      </div>
    </transition>

    <!-- 侧边栏 -->
    <transition name="slide-fade">
      <div v-if="isOpen" class="cart-sidebar-overlay" @click="closeSidebar">
        <div
          class="cart-sidebar glass-panel"
          @click.stop
        >
          <!-- 头部 -->
          <div class="sidebar-header">
            <div class="header-title">
              <el-icon :size="20"><ShoppingCart /></el-icon>
              <h3>购物车</h3>
              <span class="item-count">({{ cartStore.totalItems }})</span>
            </div>
            <el-button
              :icon="Close"
              circle
              size="small"
              @click="closeSidebar"
              class="close-btn"
            />
          </div>

          <!-- 购物车内容 -->
          <div class="sidebar-content" v-if="cartStore.items.length > 0">
            <transition-group name="cart-item" tag="div" class="cart-items-list">
              <div
                v-for="item in cartStore.items"
                :key="item.id"
                class="cart-item"
              >
                <!-- 商品图片 -->
                <div class="item-image">
                  <img
                    :src="item.image_url || '/default-coffee.jpg'"
                    :alt="item.name"
                    @error="handleImageError"
                  />
                </div>

                <!-- 商品信息 -->
                <div class="item-info">
                  <h4 class="item-name">{{ item.name }}</h4>
                  <div class="item-price">
                    <span class="price">¥{{ item.price }}</span>
                    <span class="quantity">x{{ item.quantity }}</span>
                  </div>
                  
                  <!-- 数量控制 -->
                  <div class="quantity-controls">
                    <el-button
                      :icon="Minus"
                      size="small"
                      circle
                      @click="decreaseQuantity(item)"
                      :disabled="item.quantity <= 1"
                      class="quantity-btn"
                    />
                    <span class="quantity-display">{{ item.quantity }}</span>
                    <el-button
                      :icon="Plus"
                      size="small"
                      circle
                      @click="increaseQuantity(item)"
                      :disabled="item.quantity >= 99"
                      class="quantity-btn"
                    />
                  </div>
                </div>

                <!-- 小计 -->
                <div class="item-subtotal">
                  <div class="subtotal-price">
                    ¥{{ (item.price * item.quantity).toFixed(2) }}
                  </div>
                  <el-button
                    :icon="Delete"
                    type="danger"
                    size="small"
                    circle
                    @click="removeItem(item)"
                    class="delete-btn"
                  />
                </div>

                <!-- 滑动删除指示器 -->
                <div class="swipe-indicator">
                  <el-icon><Delete /></el-icon>
                  <span>滑动删除</span>
                </div>
              </div>
            </transition-group>
          </div>

          <!-- 空状态 -->
          <div v-else class="empty-cart">
            <el-icon :size="64" class="empty-icon"><ShoppingCart /></el-icon>
            <p class="empty-text">购物车是空的</p>
            <el-button type="primary" @click="goToMenu" class="browse-btn">
              去逛逛
            </el-button>
          </div>

          <!-- 底部操作栏 -->
          <div v-if="cartStore.items.length > 0" class="sidebar-footer">
            <!-- 总价信息 -->
            <div class="total-section">
              <div class="total-row">
                <span class="total-label">商品总数：</span>
                <span class="total-value">{{ cartStore.totalItems }} 件</span>
              </div>
              <div class="total-row main-total">
                <span class="total-label">总计：</span>
                <div class="total-price-wrapper">
                  <span class="currency">¥</span>
                  <span class="total-price">{{ cartStore.totalPrice.toFixed(2) }}</span>
                </div>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="action-buttons">
              <el-button
                size="large"
                @click="clearCart"
                class="clear-btn"
              >
                <el-icon><Delete /></el-icon>
                清空购物车
              </el-button>
              <el-button
                type="primary"
                size="large"
                @click="goToCheckout"
                class="checkout-btn ripple-button"
              >
                <el-icon><ShoppingCart /></el-icon>
                去结算
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  ShoppingCart, Close, Delete, Plus, Minus 
} from '@element-plus/icons-vue'
import { useCartStore } from '@/store/cart'

const router = useRouter()
const cartStore = useCartStore()

// 响应式数据
const isOpen = ref(false)

// 方法
const toggleSidebar = () => {
  isOpen.value = !isOpen.value
}

const closeSidebar = () => {
  isOpen.value = false
}

const handleImageError = (event) => {
  event.target.src = '/default-coffee.jpg'
}

const increaseQuantity = (item) => {
  cartStore.updateQuantity(item.id, item.quantity + 1)
}

const decreaseQuantity = (item) => {
  if (item.quantity > 1) {
    cartStore.updateQuantity(item.id, item.quantity - 1)
  }
}

const removeItem = async (item) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 ${item.name} 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    cartStore.removeItem(item.id)
    ElMessage.success('已删除')
  } catch {
    // 用户取消
  }
}

const handleSwipeDelete = (item) => {
  removeItem(item)
}

const clearCart = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清空购物车吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    cartStore.clearCart()
    ElMessage.success('购物车已清空')
    closeSidebar()
  } catch {
    // 用户取消
  }
}

const goToMenu = () => {
  closeSidebar()
  router.push('/menu')
}

const goToCheckout = () => {
  closeSidebar()
  router.push('/cart')
}

// 暴露方法供父组件调用
defineExpose({
  open: () => { isOpen.value = true },
  close: closeSidebar,
  toggle: toggleSidebar
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';
@import '@/styles/mixins.scss';

.cart-sidebar-wrapper {
  position: fixed;
  z-index: $z-fixed;
}

.cart-float-button {
  position: fixed;
  bottom: $spacing-2xl;
  right: $spacing-2xl;
  @include glass-effect(0.2, 12px);
  padding: $spacing-md $spacing-lg;
  border-radius: $radius-full;
  box-shadow: $shadow-float;
  cursor: pointer;
  transition: all $transition-base $ease-out;
  display: flex;
  align-items: center;
  gap: $spacing-md;
  z-index: $z-fixed;
  
  &:hover {
    transform: translateY(-5px) scale(1.05);
    box-shadow: $shadow-xl;
  }
  
  &:active {
    transform: translateY(-2px) scale(1.02);
  }
  
  .cart-icon-wrapper {
    position: relative;
    color: $primary-color;
    
    .cart-badge {
      position: absolute;
      top: -8px;
      right: -8px;
      background: $danger-color;
      color: white;
      font-size: $font-size-xs;
      font-weight: $font-weight-bold;
      padding: 2px 6px;
      border-radius: $radius-full;
      min-width: 20px;
      text-align: center;
      box-shadow: $shadow-sm;
    }
  }
  
  .cart-price {
    font-size: $font-size-lg;
    font-weight: $font-weight-bold;
    color: $primary-color;
  }
  
  @include respond-to(md) {
    bottom: $spacing-lg;
    right: $spacing-lg;
    padding: $spacing-sm $spacing-md;
    
    .cart-price {
      font-size: $font-size-base;
    }
  }
}

.cart-sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  z-index: $z-modal;
  display: flex;
  justify-content: flex-end;
}

.cart-sidebar {
  width: 420px;
  max-width: 90vw;
  height: 100vh;
  @include glass-effect(0.95, 20px);
  display: flex;
  flex-direction: column;
  box-shadow: -4px 0 20px rgba(0, 0, 0, 0.2);
  
  @include respond-to(md) {
    width: 100%;
    max-width: 100vw;
  }
}

.sidebar-header {
  @include flex-between;
  padding: $spacing-lg;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  
  .header-title {
    display: flex;
    align-items: center;
    gap: $spacing-sm;
    
    h3 {
      margin: 0;
      font-size: $font-size-xl;
      font-weight: $font-weight-bold;
      color: $text-primary;
    }
    
    .item-count {
      font-size: $font-size-sm;
      color: $text-secondary;
    }
  }
  
  .close-btn {
    @include glass-effect(0.1, 8px);
  }
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
  padding: $spacing-md;
  @include custom-scrollbar;
}

.cart-items-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
}

.cart-item {
  @include glass-effect(0.1, 8px);
  padding: $spacing-md;
  border-radius: $radius-md;
  display: flex;
  gap: $spacing-md;
  position: relative;
  transition: all $transition-base $ease-out;
  
  &:hover {
    transform: translateX(-5px);
    box-shadow: $shadow-md;
  }
  
  .item-image {
    width: 80px;
    height: 80px;
    border-radius: $radius-sm;
    overflow: hidden;
    flex-shrink: 0;
    
    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }
  
  .item-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: $spacing-xs;
    
    .item-name {
      margin: 0;
      font-size: $font-size-base;
      font-weight: $font-weight-semibold;
      color: $text-primary;
      @include text-ellipsis(2);
    }
    
    .item-price {
      display: flex;
      align-items: center;
      gap: $spacing-sm;
      
      .price {
        font-size: $font-size-lg;
        font-weight: $font-weight-bold;
        color: $danger-color;
      }
      
      .quantity {
        font-size: $font-size-sm;
        color: $text-secondary;
      }
    }
    
    .quantity-controls {
      display: flex;
      align-items: center;
      gap: $spacing-sm;
      
      .quantity-btn {
        width: 28px;
        height: 28px;
        @include glass-effect(0.1, 6px);
        
        &:hover:not(:disabled) {
          transform: scale(1.1);
        }
      }
      
      .quantity-display {
        min-width: 30px;
        text-align: center;
        font-weight: $font-weight-semibold;
      }
    }
  }
  
  .item-subtotal {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    justify-content: space-between;
    
    .subtotal-price {
      font-size: $font-size-lg;
      font-weight: $font-weight-bold;
      color: $primary-color;
    }
    
    .delete-btn {
      @include glass-effect(0.1, 6px);
    }
  }
  
  .swipe-indicator {
    position: absolute;
    right: 0;
    top: 0;
    bottom: 0;
    width: 80px;
    background: $danger-color;
    @include flex-center;
    flex-direction: column;
    gap: $spacing-xs;
    color: white;
    font-size: $font-size-xs;
    border-radius: 0 $radius-md $radius-md 0;
    opacity: 0;
    transform: translateX(100%);
    transition: all $transition-base $ease-out;
  }
}

.empty-cart {
  @include flex-center;
  flex-direction: column;
  gap: $spacing-lg;
  padding: $spacing-2xl;
  height: 100%;
  
  .empty-icon {
    color: $text-light;
    opacity: 0.5;
  }
  
  .empty-text {
    font-size: $font-size-lg;
    color: $text-secondary;
    margin: 0;
  }
  
  .browse-btn {
    border-radius: $radius-full;
    padding: $spacing-sm $spacing-xl;
  }
}

.sidebar-footer {
  border-top: 1px solid rgba(0, 0, 0, 0.1);
  padding: $spacing-lg;
  @include glass-effect(0.1, 12px);
  
  .total-section {
    margin-bottom: $spacing-md;
    
    .total-row {
      @include flex-between;
      margin-bottom: $spacing-sm;
      
      .total-label {
        font-size: $font-size-base;
        color: $text-secondary;
      }
      
      .total-value {
        font-size: $font-size-base;
        font-weight: $font-weight-semibold;
        color: $text-primary;
      }
      
      &.main-total {
        padding-top: $spacing-sm;
        border-top: 2px solid rgba(0, 0, 0, 0.1);
        margin-bottom: 0;
        
        .total-label {
          font-size: $font-size-lg;
          font-weight: $font-weight-bold;
          color: $text-primary;
        }
        
        .total-price-wrapper {
          display: flex;
          align-items: baseline;
          
          .currency {
            font-size: $font-size-lg;
            color: $danger-color;
            margin-right: $spacing-xs;
          }
          
          .total-price {
            font-size: $font-size-3xl;
            font-weight: $font-weight-bold;
            color: $danger-color;
          }
        }
      }
    }
  }
  
  .action-buttons {
    display: flex;
    gap: $spacing-sm;
    
    .clear-btn,
    .checkout-btn {
      flex: 1;
      border-radius: $radius-full;
      font-weight: $font-weight-semibold;
      
      &:hover {
        transform: translateY(-2px);
      }
    }
    
    .checkout-btn {
      background: $primary-gradient;
      border: none;
      @include ripple-effect;
    }
  }
}

// 动画
.bounce-enter-active {
  animation: bounceIn 0.6s $ease-bounce;
}

.badge-pop-enter-active {
  animation: scaleIn 0.3s $ease-bounce;
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all $transition-base $ease-out;
}

.slide-fade-enter-from {
  opacity: 0;
  
  .cart-sidebar {
    transform: translateX(100%);
  }
}

.slide-fade-leave-to {
  opacity: 0;
  
  .cart-sidebar {
    transform: translateX(100%);
  }
}

.cart-item-enter-active,
.cart-item-leave-active {
  transition: all $transition-base $ease-out;
}

.cart-item-enter-from {
  opacity: 0;
  transform: translateX(50px);
}

.cart-item-leave-to {
  opacity: 0;
  transform: translateX(-50px);
}

.cart-item-move {
  transition: transform $transition-base $ease-out;
}
</style>
