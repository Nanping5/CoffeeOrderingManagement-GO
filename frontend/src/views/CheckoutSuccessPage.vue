<template>
  <div class="checkout-success-page">
    <div class="success-container animate__animated animate__zoomIn">
      <!-- 成功图标 -->
      <div class="success-icon">
        <el-icon :size="100" color="#27ae60">
          <CircleCheck />
        </el-icon>
      </div>

      <!-- 成功标题 -->
      <h1 class="success-title">下单成功！</h1>
      <p class="success-subtitle">请凭取餐码到柜台取餐</p>

      <!-- 取餐码 -->
      <div class="pickup-code-section">
        <div class="pickup-code-label">您的取餐码</div>
        <div class="pickup-code">{{ pickupCode }}</div>
        <div class="pickup-code-tip">请记住您的取餐码</div>
      </div>

      <!-- 订单信息 -->
      <el-card class="order-info-card">
        <template #header>
          <div class="card-header">
            <h3>订单信息</h3>
            <el-tag type="success">已支付</el-tag>
          </div>
        </template>

        <div class="order-details">
          <div class="detail-row">
            <span class="label">订单号：</span>
            <span class="value">{{ orderNumber }}</span>
          </div>
          <div class="detail-row">
            <span class="label">下单时间：</span>
            <span class="value">{{ orderTime }}</span>
          </div>
          <div class="detail-row">
            <span class="label">商品数量：</span>
            <span class="value">{{ itemCount }} 件</span>
          </div>
          <div v-if="pointsUsed > 0" class="detail-row points">
            <span class="label">积分抵扣：</span>
            <span class="value discount">-{{ pointsUsed }}积分</span>
          </div>
          <div class="detail-row total">
            <span class="label">实付金额：</span>
            <span class="value">¥{{ totalAmount }}</span>
          </div>
          <div v-if="pointsEarned > 0" class="detail-row earned">
            <span class="label">获得积分：</span>
            <span class="value success">+{{ pointsEarned }}积分</span>
          </div>
        </div>

        <!-- 商品列表 -->
        <el-divider />
        <div class="order-items">
          <h4>商品明细</h4>
          <div
            v-for="item in orderItems"
            :key="item.id"
            class="order-item"
          >
            <div class="item-info">
              <span class="item-name">{{ item.name }}</span>
              <span class="item-quantity">x{{ item.quantity }}</span>
            </div>
            <span class="item-price">¥{{ (item.price * item.quantity).toFixed(2) }}</span>
          </div>
        </div>
      </el-card>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <el-button
          type="primary"
          size="large"
          @click="continueShopping"
          class="continue-btn"
        >
          <el-icon><ShoppingCart /></el-icon>
          继续点餐
        </el-button>
        <el-button
          size="large"
          @click="printOrder"
          class="print-btn"
        >
          <el-icon><Printer /></el-icon>
          打印小票
        </el-button>
      </div>

      <!-- 温馨提示 -->
      <el-alert
        title="温馨提示"
        type="info"
        :closable="false"
        show-icon
        class="tips-alert"
      >
        <p>• 请凭取餐码到柜台取餐</p>
        <p>• 预计制作时间：5-10分钟</p>
        <p>• 如有问题请联系工作人员</p>
      </el-alert>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  CircleCheck,
  ShoppingCart,
  Printer
} from '@element-plus/icons-vue'
import { useCartStore } from '@/store/cart'

const router = useRouter()
const route = useRoute()
const cartStore = useCartStore()

// 响应式数据
const pickupCode = ref('')
const orderNumber = ref('')
const orderTime = ref('')
const totalAmount = ref(0)
const orderItems = ref([])
const pointsUsed = ref(0)
const pointsEarned = ref(0)

// 计算属性
const itemCount = computed(() => {
  return orderItems.value.reduce((sum, item) => sum + item.quantity, 0)
})

// 生成取餐码
const generatePickupCode = () => {
  // 格式：A001, A002, ..., Z999
  const date = new Date()
  const hour = date.getHours()
  const minute = date.getMinutes()
  const random = Math.floor(Math.random() * 1000)
  
  const letter = String.fromCharCode(65 + (hour % 26))
  const number = String(random).padStart(3, '0')
  
  return `${letter}${number}`
}

// 生成订单号
const generateOrderNumber = () => {
  const date = new Date()
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const random = Math.floor(Math.random() * 10000)
  
  return `CO${year}${month}${day}${String(random).padStart(4, '0')}`
}

// 格式化时间
const formatTime = () => {
  return new Date().toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 初始化订单信息
const initOrderInfo = () => {
  // 优先从路由参数获取真实订单信息
  if (route.query.pickupCode) {
    // 从后端返回的订单信息
    pickupCode.value = route.query.pickupCode
    orderNumber.value = route.query.orderNumber || generateOrderNumber()
    totalAmount.value = parseFloat(route.query.totalPrice) || 0
    orderTime.value = formatTime()
    pointsUsed.value = parseInt(route.query.pointsUsed) || 0
    pointsEarned.value = parseInt(route.query.pointsEarned) || 0
    
    // 优先从路由参数获取商品信息
    if (route.query.items) {
      try {
        const items = JSON.parse(route.query.items)
        orderItems.value = items.map(item => ({
          id: item.id,
          name: item.name,
          price: item.price || item.unit_price,
          quantity: item.quantity
        }))
      } catch (e) {
        console.error('解析商品信息失败:', e)
        orderItems.value = []
      }
    } else {
      // 兼容：从购物车获取（如果路由参数没有商品信息）
      orderItems.value = cartStore.items.map(item => ({
        id: item.id,
        name: item.name,
        price: item.price || item.unit_price,
        quantity: item.quantity
      }))
    }
  } else {
    // 兼容旧逻辑：从购物车生成（如果没有后端订单信息）
    const items = route.query.items ? JSON.parse(route.query.items) : cartStore.items
    const amount = route.query.amount || cartStore.totalPrice
    
    pickupCode.value = generatePickupCode()
    orderNumber.value = generateOrderNumber()
    orderTime.value = formatTime()
    totalAmount.value = amount
    orderItems.value = items.map(item => ({
      id: item.id,
      name: item.name,
      price: item.price || item.unit_price,
      quantity: item.quantity
    }))
  }
}

// 继续点餐
const continueShopping = () => {
  router.push('/menu')
}

// 打印订单
const printOrder = () => {
  // 创建打印内容
  const printContent = `
    <div style="text-align: center; padding: 20px; font-family: Arial, sans-serif;">
      <h2>咖啡点餐系统</h2>
      <h1 style="font-size: 48px; margin: 20px 0;">取餐码: ${pickupCode.value}</h1>
      <hr style="margin: 20px 0;">
      <div style="text-align: left; margin: 20px 0;">
        <p><strong>订单号：</strong>${orderNumber.value}</p>
        <p><strong>下单时间：</strong>${orderTime.value}</p>
        <p><strong>商品数量：</strong>${itemCount.value} 件</p>
      </div>
      <hr style="margin: 20px 0;">
      <div style="text-align: left;">
        <h3>商品明细</h3>
        ${orderItems.value.map(item => `
          <p>${item.name} x${item.quantity} - ¥${(item.price * item.quantity).toFixed(2)}</p>
        `).join('')}
      </div>
      <hr style="margin: 20px 0;">
      <h2 style="margin: 20px 0;">总计: ¥${totalAmount.value}</h2>
      <p style="margin-top: 30px;">感谢您的光临！</p>
    </div>
  `
  
  // 打开打印窗口
  const printWindow = window.open('', '_blank')
  printWindow.document.write(printContent)
  printWindow.document.close()
  printWindow.print()
  
  ElMessage.success('正在打印...')
}

// 组件挂载
onMounted(() => {
  initOrderInfo()
})
</script>

<style lang="scss" scoped>
.checkout-success-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 40px 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.success-container {
  max-width: 600px;
  width: 100%;
}

.success-icon {
  text-align: center;
  margin-bottom: 20px;
}

.success-title {
  text-align: center;
  font-size: 2.5rem;
  color: #27ae60;
  margin: 0 0 10px 0;
  font-weight: bold;
}

.success-subtitle {
  text-align: center;
  font-size: 1.2rem;
  color: #7f8c8d;
  margin: 0 0 30px 0;
}

.pickup-code-section {
  background: linear-gradient(135deg, #8b4513 0%, #a0522d 100%);
  border-radius: 16px;
  padding: 40px;
  text-align: center;
  margin-bottom: 30px;
  box-shadow: 0 8px 24px rgba(139, 69, 19, 0.3);

  .pickup-code-label {
    color: rgba(255, 255, 255, 0.9);
    font-size: 1.2rem;
    margin-bottom: 15px;
  }

  .pickup-code {
    font-size: 4rem;
    font-weight: bold;
    color: white;
    letter-spacing: 8px;
    margin-bottom: 15px;
    text-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  }

  .pickup-code-tip {
    color: rgba(255, 255, 255, 0.8);
    font-size: 1rem;
  }
}

.order-info-card {
  margin-bottom: 30px;
  border-radius: 12px;
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    h3 {
      margin: 0;
      color: #2c3e50;
      font-size: 1.3rem;
    }
  }

  .order-details {
    .detail-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 0;
      border-bottom: 1px solid #f0f0f0;

      &:last-child {
        border-bottom: none;
      }

      &.points {
        .value.discount {
          color: #f39c12;
          font-weight: 500;
        }
      }

      &.total {
        margin-top: 10px;
        padding-top: 15px;
        border-top: 2px solid #e0e0e0;
        font-size: 1.2rem;
        font-weight: bold;

        .value {
          color: #e74c3c;
          font-size: 1.5rem;
        }
      }

      &.earned {
        .value.success {
          color: #27ae60;
          font-weight: 500;
        }
      }

      .label {
        color: #7f8c8d;
        font-size: 1rem;
      }

      .value {
        color: #2c3e50;
        font-size: 1rem;
        font-weight: 500;
      }
    }
  }

  .order-items {
    h4 {
      margin: 0 0 15px 0;
      color: #2c3e50;
      font-size: 1.1rem;
    }

    .order-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 10px 0;
      border-bottom: 1px solid #f0f0f0;

      &:last-child {
        border-bottom: none;
      }

      .item-info {
        display: flex;
        align-items: center;
        gap: 10px;

        .item-name {
          color: #2c3e50;
          font-size: 1rem;
        }

        .item-quantity {
          color: #7f8c8d;
          font-size: 0.9rem;
        }
      }

      .item-price {
        color: #e74c3c;
        font-weight: 500;
        font-size: 1rem;
      }
    }
  }
}

.action-buttons {
  display: flex;
  gap: 20px;
  margin-bottom: 30px;

  .continue-btn,
  .print-btn {
    flex: 1;
    height: 50px;
    font-size: 1.1rem;
    border-radius: 25px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
    }
  }

  .continue-btn {
    background: linear-gradient(135deg, #8b4513, #a0522d);
    border: none;
  }
}

.tips-alert {
  border-radius: 12px;

  :deep(.el-alert__content) {
    p {
      margin: 5px 0;
      line-height: 1.6;
    }
  }
}

// 响应式设计 - 移动端适配
@media (max-width: 768px) {
  .checkout-success-page {
    padding: 20px 12px;
    align-items: flex-start;
    padding-top: 30px;
  }

  .success-icon {
    .el-icon {
      font-size: 70px !important;
    }
  }

  .success-title {
    font-size: 1.8rem;
  }

  .success-subtitle {
    font-size: 1rem;
    margin-bottom: 20px;
  }

  .pickup-code-section {
    padding: 24px 16px;
    border-radius: 12px;
    margin-bottom: 20px;

    .pickup-code-label {
      font-size: 1rem;
    }

    .pickup-code {
      font-size: 2.5rem;
      letter-spacing: 4px;
    }

    .pickup-code-tip {
      font-size: 0.85rem;
    }
  }

  .order-info-card {
    margin-bottom: 20px;
    border-radius: 10px;

    :deep(.el-card__header) {
      padding: 12px 16px;
    }

    :deep(.el-card__body) {
      padding: 16px;
    }

    .card-header h3 {
      font-size: 1.1rem;
    }

    .order-details {
      .detail-row {
        padding: 10px 0;
        font-size: 0.9rem;

        &.total {
          .value {
            font-size: 1.3rem;
          }
        }
      }
    }

    .order-items {
      h4 {
        font-size: 1rem;
        margin-bottom: 10px;
      }

      .order-item {
        padding: 8px 0;

        .item-info {
          .item-name {
            font-size: 0.9rem;
          }

          .item-quantity {
            font-size: 0.8rem;
          }
        }

        .item-price {
          font-size: 0.9rem;
        }
      }
    }
  }

  .action-buttons {
    flex-direction: column;
    gap: 12px;
    margin-bottom: 20px;

    .continue-btn,
    .print-btn {
      width: 100%;
      height: 46px;
      font-size: 1rem;
    }
  }

  .tips-alert {
    border-radius: 10px;

    :deep(.el-alert__content) {
      p {
        font-size: 0.85rem;
        margin: 4px 0;
      }
    }
  }
}

@media (max-width: 480px) {
  .checkout-success-page {
    padding: 16px 10px;
    padding-top: 20px;
  }

  .success-icon .el-icon {
    font-size: 60px !important;
  }

  .success-title {
    font-size: 1.5rem;
  }

  .success-subtitle {
    font-size: 0.9rem;
  }

  .pickup-code-section {
    padding: 20px 12px;

    .pickup-code {
      font-size: 2rem;
      letter-spacing: 3px;
    }
  }

  .order-info-card {
    .order-details .detail-row {
      flex-wrap: wrap;
      
      .label {
        flex: 0 0 100%;
        margin-bottom: 4px;
        font-size: 0.8rem;
      }
      
      .value {
        font-size: 0.9rem;
      }
      
      &.total {
        flex-wrap: nowrap;
        
        .label {
          flex: 0 0 auto;
          margin-bottom: 0;
        }
      }
    }
  }
}
</style>
