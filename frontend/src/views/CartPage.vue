<template>
  <div class="cart-page-enhanced">
    <!-- 页面头部 -->
    <div class="page-header scroll-reveal">
      <h1 class="page-title gradient-text">
        <el-icon><ShoppingCart /></el-icon>
        购物车
      </h1>
      <p class="page-subtitle">确认您选择的商品</p>
    </div>

    <!-- 购物车内容 -->
    <div v-if="cartStore.items.length > 0" class="cart-content">
      <!-- 购物车列表 -->
      <div class="cart-items-section glass-panel">
        <div class="section-header">
          <h2>商品列表 ({{ cartStore.totalItems }})</h2>
          <el-button
            type="danger"
            :icon="Delete"
            @click="clearCart"
            class="clear-btn"
          >
            清空购物车
          </el-button>
        </div>

        <transition-group name="cart-item" tag="div" class="cart-items-list">
          <div
            v-for="item in cartStore.items"
            :key="item.id"
            class="cart-item stagger-item"
            draggable="true"
            @dragstart="handleDragStart($event, item)"
            @dragover.prevent
            @drop="handleDrop($event, item)"
          >
            <!-- 拖拽手柄 -->
            <div class="drag-handle">
              <el-icon><Rank /></el-icon>
            </div>

            <!-- 商品图片 -->
            <div class="item-image">
              <img :src="item.image_url || '/default-coffee.jpg'" :alt="item.name" />
            </div>

            <!-- 商品信息 -->
            <div class="item-info">
              <h3 class="item-name">{{ item.name }}</h3>
              <p class="item-description">{{ item.description }}</p>

              <!-- 数量控制 -->
              <div class="item-controls">
                <el-input-number
                  v-model="item.quantity"
                  :min="1"
                  :max="99"
                  @change="updateQuantity(item)"
                  class="quantity-input"
                />
              </div>
            </div>

            <!-- 价格和操作 -->
            <div class="item-actions">
              <div class="item-price">
                <span class="price">¥{{ (item.price * item.quantity).toFixed(2) }}</span>
                <span class="unit-price">单价: ¥{{ item.price }}</span>
              </div>

              <el-button
                type="danger"
                :icon="Delete"
                circle
                @click="removeItem(item)"
                class="delete-btn"
              />
            </div>
          </div>
        </transition-group>
      </div>

      <!-- 结算区域 -->
      <div class="checkout-section glass-panel scroll-reveal">
        <h2>订单汇总</h2>

        <div class="total-info">
          <div class="total-row">
            <span>商品总数：</span>
            <span class="count-up" :data-count="cartStore.totalItems">0</span>
          </div>
          <div class="total-row">
            <span>商品金额：</span>
            <span>¥{{ cartStore.totalPrice.toFixed(2) }}</span>
          </div>
        </div>

        <!-- 积分抵扣区域 -->
        <div v-if="authStore.isAuthenticated && !authStore.isAdmin && pointsStore.availablePoints > 0" class="points-section">
          <div class="points-header">
            <el-switch v-model="usePoints" active-text="使用积分抵扣" />
            <span class="points-balance">
              <el-icon><Coin /></el-icon>
              可用: {{ pointsStore.availablePoints }}
            </span>
          </div>
          
          <div v-if="usePoints" class="points-usage">
            <el-slider
              v-model="pointsToUse"
              :min="0"
              :max="maxUsablePoints"
              :step="100"
              :disabled="calculatingPoints"
              show-input
              input-size="small"
            />
            <div class="points-info">
              <span>抵扣: <strong class="discount">-¥{{ pointsDiscount.toFixed(2) }}</strong></span>
              <span>预计获得: <strong class="earn">+{{ estimatedPointsEarned }}积分</strong></span>
            </div>
          </div>
        </div>

        <!-- 未登录提示 -->
        <div v-else-if="!authStore.isAuthenticated" class="login-tip">
          <router-link to="/auth/login" class="login-link">
            <el-icon><Coin /></el-icon>
            登录后可使用积分抵扣
          </router-link>
        </div>

        <div class="total-info">
          <div class="total-row main-total">
            <span>实付金额：</span>
            <span class="total-price">
              <span class="currency">¥</span>
              <span class="price-value">{{ finalTotal.toFixed(2) }}</span>
            </span>
          </div>
        </div>

        <el-button
          type="primary"
          size="large"
          :loading="creatingOrder"
          :disabled="creatingOrder"
          @click="handleCheckout"
          class="checkout-btn ripple-button"
        >
          <el-icon v-if="!creatingOrder"><ShoppingCart /></el-icon>
          {{ creatingOrder ? '创建订单中...' : '去结算' }}
        </el-button>

        <!-- 取餐说明 -->
        <div class="pickup-info">
          <div class="info-item">
            <el-icon><Clock /></el-icon>
            <span>预计5-10分钟</span>
          </div>
          <div class="info-item">
            <el-icon><MapLocation /></el-icon>
            <span>凭取餐码到柜台取餐</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-cart scroll-reveal">
      <el-icon :size="80" class="empty-icon"><ShoppingCart /></el-icon>
      <p class="empty-text">购物车是空的</p>
      <el-button type="primary" @click="goToMenu" class="browse-btn">
        去逛逛
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ShoppingCart, Delete, Clock, MapLocation, Rank, Coin } from '@element-plus/icons-vue'
import { useCartStore } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import { usePointsStore } from '@/store/points'
import { useScrollReveal, useStaggerAnimation, useCountUp } from '@/composables/useScrollReveal'
import orderAPI from '@/api/order'
import pointsAPI from '@/api/points'

const router = useRouter()
const cartStore = useCartStore()
const authStore = useAuthStore()
const pointsStore = usePointsStore()

// 订单创建状态
const creatingOrder = ref(false)

// 积分相关状态
const usePoints = ref(true) // 默认开启积分抵扣
const pointsToUse = ref(0)
const pointsCalculation = ref(null)
const calculatingPoints = ref(false)

// 使用动画
useScrollReveal({ threshold: 0.1, once: true })
useStaggerAnimation('.stagger-item', 50)
useCountUp('.count-up', 1000)

// 拖拽相关
const draggedItem = ref(null)

const handleDragStart = (event, item) => {
  draggedItem.value = item
  event.dataTransfer.effectAllowed = 'move'
}

const handleDrop = (event, targetItem) => {
  if (!draggedItem.value || draggedItem.value.id === targetItem.id) return

  const items = [...cartStore.items]
  const draggedIndex = items.findIndex(i => i.id === draggedItem.value.id)
  const targetIndex = items.findIndex(i => i.id === targetItem.id)

  items.splice(draggedIndex, 1)
  items.splice(targetIndex, 0, draggedItem.value)

  // 更新购物车顺序（需要在store中添加此方法）
  cartStore.items = items
  cartStore.saveCartToStorage()

  draggedItem.value = null
}

const updateQuantity = (item) => {
  cartStore.updateQuantity(item.id, item.quantity)
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
  } catch {
    // 用户取消
  }
}

const handleCheckout = async () => {
  if (creatingOrder.value) return
  
  if (cartStore.items.length === 0) {
    ElMessage.warning('购物车是空的')
    return
  }

  creatingOrder.value = true

  try {
    // 准备订单数据
    const orderData = {
      items: cartStore.items.map(item => ({
        menu_id: item.id,
        quantity: item.quantity,
        unit_price: item.price
      })),
      total_price: cartStore.totalPrice
    }

    // 如果使用积分抵扣
    if (usePoints.value && pointsToUse.value > 0 && authStore.isAuthenticated) {
      orderData.use_points = true
      orderData.points_to_use = pointsToUse.value
    }

    // 调用后端API创建订单
    const response = await orderAPI.createOrder(orderData)
    
    // 后端返回格式：{ success: true, message: '', order: {...} }
    if (response.success && response.order) {
      const order = response.order

      // 保存商品信息用于结算成功页面显示
      const orderItemsForDisplay = cartStore.items.map(item => ({
        id: item.id,
        name: item.name,
        price: item.price,
        quantity: item.quantity
      }))

      // 清空购物车和积分状态（在跳转前清空）
      cartStore.clearCart()
      usePoints.value = false
      pointsToUse.value = 0
      pointsCalculation.value = null

      // 跳转到结算成功页面，传递订单信息和商品列表
      router.push({
        name: 'CheckoutSuccess',
        query: {
          orderId: order.id,
          orderNumber: order.order_number,
          pickupCode: order.pickup_code,
          totalPrice: order.final_payment_amount || order.total_price,
          originalTotal: order.original_total_price || order.total_price,
          pointsDeduction: order.points_deduction_amount || 0,
          pointsUsed: order.points_used || order.customer_points_used || 0,
          pointsEarned: order.estimated_points_earned || order.points_earned || 0,
          items: JSON.stringify(orderItemsForDisplay)
        }
      })
      
      // 刷新积分信息
      if (authStore.isAuthenticated) {
        pointsStore.fetchPointsInfo()
      }
      
      ElMessage.success('订单创建成功！')
    } else {
      throw new Error(response.message || response.errors?.[0] || '订单创建失败')
    }
  } catch (error) {
    console.error('创建订单失败:', error)
    ElMessage.error(error.message || '订单创建失败，请稍后重试')
  } finally {
    creatingOrder.value = false
  }
}

const goToMenu = () => {
  router.push('/menu')
}

// 积分相关计算
const maxUsablePoints = computed(() => {
  if (!pointsCalculation.value) return 0
  return Math.min(pointsStore.availablePoints, pointsCalculation.value.max_usable_points || 0)
})

const pointsDiscount = computed(() => {
  if (!usePoints.value || !pointsCalculation.value) return 0
  return pointsCalculation.value.points_value || 0
})

const finalTotal = computed(() => {
  if (!usePoints.value || !pointsCalculation.value) return cartStore.totalPrice
  return pointsCalculation.value.final_total || cartStore.totalPrice
})

const estimatedPointsEarned = computed(() => {
  if (!pointsCalculation.value) return Math.floor(cartStore.totalPrice)
  return pointsCalculation.value.estimated_points_earned || 0
})

// 计算积分抵扣
const calculatePointsDiscount = async () => {
  if (!authStore.isAuthenticated || !usePoints.value || cartStore.items.length === 0) {
    pointsCalculation.value = null
    return
  }
  
  try {
    calculatingPoints.value = true
    const response = await pointsAPI.calculateDiscount({
      items: cartStore.items.map(item => ({
        menu_id: item.id,
        quantity: item.quantity,
        unit_price: item.price
      })),
      points_to_use: pointsToUse.value
    })
    
    if (response.success) {
      pointsCalculation.value = response.data
      // 自动调整积分使用数量
      if (pointsToUse.value > response.data.max_usable_points) {
        pointsToUse.value = response.data.max_usable_points
      }
    }
  } catch (error) {
    console.error('计算积分抵扣失败:', error)
  } finally {
    calculatingPoints.value = false
  }
}

// 监听积分使用变化
watch([usePoints, pointsToUse], () => {
  if (usePoints.value) {
    calculatePointsDiscount()
  } else {
    pointsCalculation.value = null
  }
}, { immediate: false })

// 监听购物车变化
watch(() => cartStore.totalPrice, () => {
  if (usePoints.value) {
    calculatePointsDiscount()
  }
})

// 初始化
onMounted(async () => {
  if (authStore.isAuthenticated && !authStore.isAdmin) {
    await pointsStore.fetchPointsInfo()
    // 如果有积分且购物车有商品，自动计算积分抵扣
    if (pointsStore.availablePoints > 0 && cartStore.items.length > 0) {
      await calculatePointsDiscount()
      // 自动设置使用最大可用积分
      if (pointsCalculation.value) {
        pointsToUse.value = pointsCalculation.value.max_usable_points || 0
      }
    }
  }
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';
@import '@/styles/mixins.scss';

.cart-page-enhanced {
  min-height: 100vh;
  padding: $spacing-2xl $spacing-lg;
  max-width: 1200px;
  margin: 0 auto;

  @include respond-to(md) {
    padding: $spacing-lg $spacing-md;
  }
}

.page-header {
  text-align: center;
  margin-bottom: $spacing-2xl;

  .page-title {
    font-size: $font-size-4xl;
    font-weight: $font-weight-bold;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: $spacing-md;
    margin-bottom: $spacing-sm;
  }

  .page-subtitle {
    font-size: $font-size-lg;
    color: $text-secondary;
  }
}

.cart-content {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: $spacing-lg;

  @include respond-to(lg) {
    grid-template-columns: 1fr;
  }
}

.cart-items-section {
  padding: $spacing-lg;
  border-radius: $radius-lg;

  .section-header {
    @include flex-between;
    margin-bottom: $spacing-lg;
    padding-bottom: $spacing-md;
    border-bottom: 2px solid rgba(0, 0, 0, 0.1);

    h2 {
      margin: 0;
      font-size: $font-size-xl;
      font-weight: $font-weight-bold;
      color: $primary-color;
    }

    .clear-btn {
      border-radius: $radius-full;
    }
  }
}

.cart-items-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
}

.cart-item {
  display: flex;
  gap: $spacing-md;
  padding: $spacing-md;
  border-radius: $radius-md;
  @include glass-effect(0.05, 8px);
  cursor: move;
  transition: all $transition-base $ease-out;

  &:hover {
    @include glass-effect(0.1, 10px);
    transform: translateX(-5px);
  }

  .drag-handle {
    @include flex-center;
    width: 24px;
    color: $text-light;
    cursor: grab;

    &:active {
      cursor: grabbing;
    }
  }

  .item-image {
    width: 100px;
    height: 100px;
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

    .item-name {
      margin: 0 0 $spacing-xs 0;
      font-size: $font-size-lg;
      font-weight: $font-weight-semibold;
    }

    .item-description {
      margin: 0 0 $spacing-md 0;
      font-size: $font-size-sm;
      color: $text-secondary;
      @include text-ellipsis(2);
    }

    .item-controls {
      .quantity-input {
        width: 120px;
      }
    }
  }

  .item-actions {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    justify-content: space-between;

    .item-price {
      text-align: right;

      .price {
        display: block;
        font-size: $font-size-xl;
        font-weight: $font-weight-bold;
        color: $danger-color;
        margin-bottom: $spacing-xs;
      }

      .unit-price {
        font-size: $font-size-sm;
        color: $text-light;
      }
    }

    .delete-btn {
      @include glass-effect(0.1, 6px);
    }
  }
}

.checkout-section {
  padding: $spacing-lg;
  border-radius: $radius-lg;
  position: sticky;
  top: 100px;

  h2 {
    margin: 0 0 $spacing-lg 0;
    font-size: $font-size-xl;
    font-weight: $font-weight-bold;
    color: $primary-color;
  }

  .total-info {
    margin-bottom: $spacing-md;

    .total-row {
      @include flex-between;
      margin-bottom: $spacing-sm;
      font-size: $font-size-base;

      &.main-total {
        padding-top: $spacing-md;
        border-top: 2px solid rgba(0, 0, 0, 0.1);
        font-size: $font-size-lg;
        font-weight: $font-weight-bold;

        .total-price {
          display: flex;
          align-items: baseline;

          .currency {
            font-size: $font-size-lg;
            color: $danger-color;
          }

          .price-value {
            font-size: $font-size-3xl;
            color: $danger-color;
            margin-left: $spacing-xs;
          }
        }
      }
    }
  }

  .points-section {
    padding: $spacing-md;
    margin-bottom: $spacing-md;
    border-radius: $radius-md;
    background: rgba($warning-color, 0.1);
    border: 1px solid rgba($warning-color, 0.2);

    .points-header {
      @include flex-between;
      margin-bottom: $spacing-sm;

      .points-balance {
        display: flex;
        align-items: center;
        gap: $spacing-xs;
        font-size: $font-size-sm;
        color: $warning-color;
        font-weight: $font-weight-medium;
      }
    }

    .points-usage {
      .points-info {
        @include flex-between;
        margin-top: $spacing-sm;
        font-size: $font-size-sm;

        .discount {
          color: $danger-color;
        }

        .earn {
          color: $success-color;
        }
      }
    }
  }

  .login-tip {
    padding: $spacing-md;
    margin-bottom: $spacing-md;
    border-radius: $radius-md;
    background: rgba($info-color, 0.1);
    text-align: center;

    .login-link {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: $spacing-xs;
      color: $info-color;
      text-decoration: none;
      font-size: $font-size-sm;

      &:hover {
        text-decoration: underline;
      }
    }
  }

  .checkout-btn {
    width: 100%;
    border-radius: $radius-full;
    font-size: $font-size-lg;
    font-weight: $font-weight-semibold;
    padding: $spacing-md;
    background: $primary-gradient;
    border: none;
    @include ripple-effect;
    margin-bottom: $spacing-lg;

    &:hover {
      transform: translateY(-2px);
      box-shadow: $shadow-lg;
    }
  }

  .pickup-info {
    display: flex;
    flex-direction: column;
    gap: $spacing-sm;

    .info-item {
      display: flex;
      align-items: center;
      gap: $spacing-sm;
      font-size: $font-size-sm;
      color: $text-secondary;

      .el-icon {
        color: $primary-color;
      }
    }
  }
}

.empty-cart {
  @include flex-center;
  flex-direction: column;
  gap: $spacing-lg;
  padding: $spacing-2xl;
  min-height: 60vh;

  .empty-icon {
    color: $text-light;
    opacity: 0.5;
  }

  .empty-text {
    font-size: $font-size-xl;
    color: $text-secondary;
  }

  .browse-btn {
    border-radius: $radius-full;
    padding: $spacing-md $spacing-xl;
    font-size: $font-size-lg;
  }
}

// 动画
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

// 响应式 - 移动端适配
@include respond-to(lg) {
  .cart-content {
    grid-template-columns: 1fr;
  }
  
  .checkout-section {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    top: auto;
    border-radius: $radius-lg $radius-lg 0 0;
    padding: $spacing-md;
    z-index: 1000;
    box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.15);
    background: white;
    max-height: 40vh;
    overflow-y: auto;
    
    h2 {
      font-size: $font-size-lg;
      margin-bottom: $spacing-sm;
      display: none; // 移动端隐藏标题节省空间
    }
    
    .total-info {
      margin-bottom: $spacing-xs;
      
      .total-row {
        font-size: $font-size-sm;
        margin-bottom: $spacing-xs;
        
        &:not(.main-total) {
          display: none; // 隐藏非主要信息
        }
        
        &.main-total {
          padding-top: 0;
          border-top: none;
          
          .total-price .price-value {
            font-size: $font-size-xl;
          }
        }
      }
    }
    
    .points-section {
      padding: $spacing-xs $spacing-sm;
      margin-bottom: $spacing-xs;
      
      .points-header {
        flex-wrap: wrap;
        gap: $spacing-xs;
        font-size: $font-size-xs;
      }
      
      .points-usage {
        .el-slider {
          margin: $spacing-xs 0;
        }
        
        .points-info {
          font-size: $font-size-xs;
        }
      }
    }
    
    .login-tip {
      padding: $spacing-xs $spacing-sm;
      margin-bottom: $spacing-xs;
      font-size: $font-size-xs;
    }
    
    .checkout-btn {
      padding: $spacing-sm;
      font-size: $font-size-base;
      margin-bottom: $spacing-xs;
    }
    
    .pickup-info {
      flex-direction: row;
      justify-content: center;
      gap: $spacing-md;
      
      .info-item {
        font-size: $font-size-xs;
      }
    }
  }
  
  .cart-items-section {
    margin-bottom: 200px; // 为固定底部结算区留出空间
  }
}

// 深色模式下的固定结算区
html[data-theme='dark'] {
  @media (max-width: 1024px) {
    .checkout-section {
      background: #2D2D2D;
      border-top: 1px solid #404040;
    }
  }
}

@include respond-to(md) {
  .cart-page-enhanced {
    padding: $spacing-lg $spacing-sm;
  }
  
  .page-header {
    margin-bottom: $spacing-lg;
    
    .page-title {
      font-size: $font-size-2xl;
      gap: $spacing-sm;
    }
    
    .page-subtitle {
      font-size: $font-size-base;
    }
  }
  
  .cart-items-section {
    padding: $spacing-md;
    
    .section-header {
      flex-direction: column;
      align-items: flex-start;
      gap: $spacing-sm;
      
      h2 {
        font-size: $font-size-lg;
      }
      
      .clear-btn {
        width: 100%;
      }
    }
  }
  
  .cart-item {
    flex-wrap: wrap;
    padding: $spacing-sm;

    .drag-handle {
      display: none;
    }

    .item-image {
      width: 70px;
      height: 70px;
    }
    
    .item-info {
      flex: 1;
      min-width: 0;
      
      .item-name {
        font-size: $font-size-base;
      }
      
      .item-description {
        font-size: $font-size-xs;
        -webkit-line-clamp: 1;
      }
      
      .item-controls .quantity-input {
        width: 100px;
      }
    }

    .item-actions {
      width: 100%;
      flex-direction: row;
      justify-content: space-between;
      align-items: center;
      margin-top: $spacing-sm;
      padding-top: $spacing-sm;
      border-top: 1px solid rgba(0, 0, 0, 0.05);
      
      .item-price {
        text-align: left;
        
        .price {
          font-size: $font-size-lg;
        }
        
        .unit-price {
          font-size: $font-size-xs;
        }
      }
    }
  }
  
  .empty-cart {
    padding: $spacing-xl $spacing-md;
    min-height: 50vh;
    
    .empty-icon {
      font-size: 60px !important;
    }
    
    .empty-text {
      font-size: $font-size-lg;
    }
    
    .browse-btn {
      width: 100%;
      max-width: 200px;
    }
  }
}

// 超小屏幕适配
@media (max-width: 480px) {
  .cart-page-enhanced {
    padding: $spacing-md $spacing-xs;
  }
  
  .cart-item {
    .item-image {
      width: 60px;
      height: 60px;
    }
    
    .item-info {
      .item-name {
        font-size: $font-size-sm;
      }
      
      .item-controls .quantity-input {
        width: 90px;
      }
    }
  }
  
  .checkout-section {
    .pickup-info {
      flex-direction: column;
      gap: $spacing-xs;
    }
  }
}
</style>
