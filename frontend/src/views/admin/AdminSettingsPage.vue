<template>
  <div class="admin-settings-page">
    <div class="page-header">
      <h2>系统设置</h2>
    </div>

    <el-row :gutter="20">
      <el-col :xs="24" :md="12">
        <el-card class="settings-card">
          <template #header><h3><el-icon><Shop /></el-icon> 店铺设置</h3></template>
          <el-form :model="shopSettings" label-width="120px">
            <el-form-item label="店铺名称">
              <el-input v-model="shopSettings.name" placeholder="咖啡点餐系统" />
            </el-form-item>
            <el-form-item label="联系电话">
              <el-input v-model="shopSettings.phone" placeholder="请输入联系电话" />
            </el-form-item>
            <el-form-item label="店铺地址">
              <el-input v-model="shopSettings.address" type="textarea" :rows="2" placeholder="请输入店铺地址" />
            </el-form-item>
            <el-form-item label="营业时间">
              <el-time-picker v-model="shopSettings.openTime" placeholder="开始时间" format="HH:mm" style="width: 45%;" />
              <span style="margin: 0 10px;">至</span>
              <el-time-picker v-model="shopSettings.closeTime" placeholder="结束时间" format="HH:mm" style="width: 45%;" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveShopSettings">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="settings-card">
          <template #header><h3><el-icon><Setting /></el-icon> 订单设置</h3></template>
          <el-form :model="orderSettings" label-width="120px">
            <el-form-item label="取餐码前缀">
              <el-input v-model="orderSettings.pickupPrefix" placeholder="A" style="width: 100px;" />
            </el-form-item>
            <el-form-item label="自动完成时间">
              <el-input-number v-model="orderSettings.autoCompleteMinutes" :min="10" :max="120" />
              <span style="margin-left: 10px;">分钟</span>
            </el-form-item>
            <el-form-item label="订单超时提醒">
              <el-switch v-model="orderSettings.timeoutAlert" />
            </el-form-item>
            <el-form-item label="允许取消订单">
              <el-switch v-model="orderSettings.allowCancel" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveOrderSettings">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="settings-card">
          <template #header><h3><el-icon><Coin /></el-icon> 积分设置</h3></template>
          <el-form :model="pointsSettings" label-width="120px">
            <el-form-item label="消费积分比例">
              <el-input-number v-model="pointsSettings.earnRate" :min="1" :max="100" />
              <span style="margin-left: 10px;">元 = 1积分</span>
            </el-form-item>
            <el-form-item label="积分抵扣比例">
              <el-input-number v-model="pointsSettings.redeemRate" :min="1" :max="1000" />
              <span style="margin-left: 10px;">积分 = 1元</span>
            </el-form-item>
            <el-form-item label="注册赠送积分">
              <el-input-number v-model="pointsSettings.signupBonus" :min="0" :max="1000" />
            </el-form-item>
            <el-form-item label="积分有效期">
              <el-input-number v-model="pointsSettings.expiryMonths" :min="1" :max="36" />
              <span style="margin-left: 10px;">个月</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="savePointsSettings">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="settings-card">
          <template #header><h3><el-icon><Monitor /></el-icon> 系统信息</h3></template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="系统版本">v3.0.0</el-descriptions-item>
            <el-descriptions-item label="前端框架">Vue 3.3.4 + Element Plus</el-descriptions-item>
            <el-descriptions-item label="后端框架">Go Gin 1.9.1</el-descriptions-item>
            <el-descriptions-item label="数据库">MySQL 8.0</el-descriptions-item>
            <el-descriptions-item label="服务器时间">{{ serverTime }}</el-descriptions-item>
          </el-descriptions>
          <div style="margin-top: 20px; text-align: center;">
            <el-button type="info" @click="clearCache">清除缓存</el-button>
            <el-button type="warning" @click="exportData">导出数据</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Shop, Setting, Coin, Monitor } from '@element-plus/icons-vue'

const serverTime = ref(new Date().toLocaleString('zh-CN'))
let timer = null

const shopSettings = reactive({
  name: '咖啡点餐系统',
  phone: '',
  address: '',
  openTime: null,
  closeTime: null
})

const orderSettings = reactive({
  pickupPrefix: 'A',
  autoCompleteMinutes: 30,
  timeoutAlert: true,
  allowCancel: true
})

const pointsSettings = reactive({
  earnRate: 1,
  redeemRate: 100,
  signupBonus: 50,
  expiryMonths: 12
})

const saveShopSettings = () => {
  localStorage.setItem('shopSettings', JSON.stringify(shopSettings))
  ElMessage.success('店铺设置已保存')
}

const saveOrderSettings = () => {
  localStorage.setItem('orderSettings', JSON.stringify(orderSettings))
  ElMessage.success('订单设置已保存')
}

const savePointsSettings = () => {
  localStorage.setItem('pointsSettings', JSON.stringify(pointsSettings))
  ElMessage.success('积分设置已保存')
}

const clearCache = () => {
  ElMessage.success('缓存已清除')
}

const exportData = () => {
  ElMessage.info('数据导出功能开发中...')
}

const loadSettings = () => {
  const shop = localStorage.getItem('shopSettings')
  const order = localStorage.getItem('orderSettings')
  const points = localStorage.getItem('pointsSettings')
  if (shop) Object.assign(shopSettings, JSON.parse(shop))
  if (order) Object.assign(orderSettings, JSON.parse(order))
  if (points) Object.assign(pointsSettings, JSON.parse(points))
}

onMounted(() => {
  loadSettings()
  timer = setInterval(() => { serverTime.value = new Date().toLocaleString('zh-CN') }, 1000)
})

onUnmounted(() => { if (timer) clearInterval(timer) })
</script>

<style lang="scss" scoped>
.admin-settings-page { padding: 20px; background: #f8f9fa; min-height: 100vh; }
.page-header { margin-bottom: 20px; h2 { margin: 0; color: #2c3e50; } }
.settings-card {
  margin-bottom: 20px;
  h3 { margin: 0; color: #2c3e50; font-size: 16px; display: flex; align-items: center; gap: 8px; }
}
</style>
