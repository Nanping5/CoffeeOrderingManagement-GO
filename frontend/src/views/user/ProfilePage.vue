<template>
  <div class="profile-page">
    <div class="page-header">
      <h1 class="page-title gradient-text">
        <el-icon><User /></el-icon>
        个人中心
      </h1>
      <p class="page-subtitle">管理您的账户信息</p>
    </div>

    <div v-if="authStore.isAuthenticated" class="profile-content">
      <!-- 用户信息卡片 -->
      <el-card class="profile-card glass-panel">
        <div class="profile-header">
          <div class="avatar-section">
            <el-avatar :size="100" :src="authStore.userInfo?.avatar_url">
              <el-icon :size="50"><User /></el-icon>
            </el-avatar>
            <div class="member-badge" :class="pointsStore.memberLevel">
              {{ pointsStore.memberLevelName }}
            </div>
          </div>
          <div class="user-details">
            <h2>{{ authStore.userInfo?.username }}</h2>
            <p class="email">{{ authStore.userInfo?.email }}</p>
            <p class="points">
              <el-icon><Coin /></el-icon>
              {{ pointsStore.availablePoints }} 可用积分
            </p>
          </div>
        </div>
      </el-card>

      <!-- 个人信息表单 -->
      <el-card class="info-card glass-panel">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
            <el-button v-if="!isEditing" type="primary" text @click="startEdit">
              <el-icon><Edit /></el-icon> 编辑
            </el-button>
          </div>
        </template>

        <el-form
          ref="profileFormRef"
          :model="formData"
          :rules="rules"
          :disabled="!isEditing"
          label-width="80px"
        >
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item label="名字" prop="first_name">
                <el-input v-model="formData.first_name" placeholder="请输入名字" />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12">
              <el-form-item label="姓氏" prop="last_name">
                <el-input v-model="formData.last_name" placeholder="请输入姓氏" />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item label="手机号" prop="phone">
                <el-input v-model="formData.phone" placeholder="请输入手机号" />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12">
              <el-form-item label="性别" prop="gender">
                <el-select v-model="formData.gender" placeholder="请选择性别" style="width: 100%">
                  <el-option label="男" value="male" />
                  <el-option label="女" value="female" />
                  <el-option label="其他" value="other" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="生日" prop="birth_date">
            <el-date-picker
              v-model="formData.birth_date"
              type="date"
              placeholder="选择生日"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              style="width: 100%"
            />
          </el-form-item>

          <el-form-item v-if="isEditing">
            <el-button type="primary" @click="saveProfile" :loading="saving">
              保存修改
            </el-button>
            <el-button @click="cancelEdit">取消</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 修改密码 -->
      <el-card class="password-card glass-panel">
        <template #header>
          <span>修改密码</span>
        </template>

        <el-form
          ref="passwordFormRef"
          :model="passwordData"
          :rules="passwordRules"
          label-width="100px"
        >
          <el-form-item label="当前密码" prop="current_password">
            <el-input
              v-model="passwordData.current_password"
              type="password"
              placeholder="请输入当前密码"
              show-password
            />
          </el-form-item>

          <el-form-item label="新密码" prop="new_password">
            <el-input
              v-model="passwordData.new_password"
              type="password"
              placeholder="请输入新密码（至少8位）"
              show-password
            />
          </el-form-item>

          <el-form-item label="确认新密码" prop="confirm_password">
            <el-input
              v-model="passwordData.confirm_password"
              type="password"
              placeholder="请再次输入新密码"
              show-password
            />
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="changePassword" :loading="changingPassword">
              修改密码
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 快捷操作 -->
      <el-card class="actions-card glass-panel">
        <el-row :gutter="16">
          <el-col :xs="12" :sm="6">
            <el-button type="info" plain @click="goToOrders" class="action-btn">
              <el-icon><List /></el-icon>
              我的订单
            </el-button>
          </el-col>
          <el-col :xs="12" :sm="6">
            <el-button type="success" plain @click="goToMember" class="action-btn">
              <el-icon><Trophy /></el-icon>
              会员中心
            </el-button>
          </el-col>
          <el-col :xs="12" :sm="6">
            <el-button type="primary" plain @click="goToMenu" class="action-btn">
              <el-icon><Coffee /></el-icon>
              继续购物
            </el-button>
          </el-col>
          <el-col :xs="12" :sm="6">
            <el-button type="danger" plain @click="handleLogout" class="action-btn">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-button>
          </el-col>
        </el-row>
      </el-card>
    </div>

    <!-- 未登录状态 -->
    <div v-else class="not-logged-in">
      <el-icon :size="80" class="empty-icon"><User /></el-icon>
      <h3>请先登录</h3>
      <p>登录后即可管理您的个人信息</p>
      <el-button type="primary" @click="goToLogin" size="large">立即登录</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  User, Coin, Edit, Trophy, Coffee, ShoppingCart, SwitchButton, List 
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'
import { usePointsStore } from '@/store/points'

const router = useRouter()
const authStore = useAuthStore()
const pointsStore = usePointsStore()

const profileFormRef = ref(null)
const passwordFormRef = ref(null)
const isEditing = ref(false)
const saving = ref(false)
const changingPassword = ref(false)

const formData = reactive({
  first_name: '',
  last_name: '',
  phone: '',
  gender: '',
  birth_date: ''
})

const passwordData = reactive({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

const rules = {
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== passwordData.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  current_password: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 8, message: '密码至少8位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

// 初始化表单数据
const initFormData = () => {
  const user = authStore.userInfo
  if (user) {
    formData.first_name = user.first_name || ''
    formData.last_name = user.last_name || ''
    formData.phone = user.phone || ''
    formData.gender = user.gender || ''
    formData.birth_date = user.birth_date || ''
  }
}

const startEdit = () => {
  isEditing.value = true
}

const cancelEdit = () => {
  isEditing.value = false
  initFormData()
}

const saveProfile = async () => {
  if (!profileFormRef.value) return
  
  try {
    await profileFormRef.value.validate()
    saving.value = true
    
    const result = await authStore.updateProfile(formData)
    if (result.success) {
      ElMessage.success('保存成功')
      isEditing.value = false
    } else {
      ElMessage.error(result.errors?.[0] || '保存失败')
    }
  } catch (error) {
    if (error !== false) {
      ElMessage.error('保存失败')
    }
  } finally {
    saving.value = false
  }
}

const changePassword = async () => {
  if (!passwordFormRef.value) return
  
  try {
    await passwordFormRef.value.validate()
    changingPassword.value = true
    
    const result = await authStore.changePassword({
      current_password: passwordData.current_password,
      new_password: passwordData.new_password
    })
    
    if (result.success) {
      ElMessage.success('密码修改成功')
      passwordData.current_password = ''
      passwordData.new_password = ''
      passwordData.confirm_password = ''
      passwordFormRef.value.resetFields()
    } else {
      ElMessage.error(result.errors?.[0] || '密码修改失败')
    }
  } catch (error) {
    if (error !== false) {
      ElMessage.error('密码修改失败')
    }
  } finally {
    changingPassword.value = false
  }
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    authStore.logout()
    pointsStore.clearPointsData()
    ElMessage.success('已退出登录')
    router.push('/menu')
  } catch {
    // 用户取消
  }
}

const goToLogin = () => router.push('/auth/login')
const goToOrders = () => router.push('/user/orders')
const goToMember = () => router.push('/user/member')
const goToMenu = () => router.push('/menu')
const goToCart = () => router.push('/cart')

watch(() => authStore.userInfo, initFormData, { immediate: true })

onMounted(async () => {
  if (authStore.isAuthenticated) {
    await pointsStore.fetchPointsInfo()
  }
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';
@import '@/styles/mixins.scss';

.profile-page {
  min-height: 100vh;
  padding: $spacing-2xl $spacing-lg;
  max-width: 800px;
  margin: 0 auto;
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

.profile-card {
  margin-bottom: $spacing-lg;
  padding: $spacing-xl;
  border-radius: $radius-lg;

  .profile-header {
    display: flex;
    align-items: center;
    gap: $spacing-xl;

    @include respond-to(sm) {
      flex-direction: column;
      text-align: center;
    }
  }

  .avatar-section {
    position: relative;

    .member-badge {
      position: absolute;
      bottom: -5px;
      left: 50%;
      transform: translateX(-50%);
      padding: 2px 10px;
      border-radius: $radius-full;
      font-size: $font-size-xs;
      font-weight: $font-weight-bold;
      color: white;
      white-space: nowrap;

      &.bronze { background: linear-gradient(135deg, #CD7F32, #B8860B); }
      &.silver { background: linear-gradient(135deg, #C0C0C0, #808080); }
      &.gold { background: linear-gradient(135deg, #FFD700, #FFA500); }
      &.platinum { background: linear-gradient(135deg, #E5E4E2, #9370DB); }
    }
  }

  .user-details {
    h2 {
      margin: 0 0 $spacing-xs 0;
      font-size: $font-size-2xl;
      font-weight: $font-weight-bold;
    }

    .email {
      margin: 0 0 $spacing-sm 0;
      color: $text-secondary;
    }

    .points {
      margin: 0;
      display: flex;
      align-items: center;
      gap: $spacing-xs;
      color: $warning-color;
      font-weight: $font-weight-medium;
    }
  }
}

.info-card, .password-card {
  margin-bottom: $spacing-lg;
  border-radius: $radius-lg;

  .card-header {
    @include flex-between;
  }
}

.actions-card {
  border-radius: $radius-lg;
  padding: $spacing-lg;

  .action-btn {
    width: 100%;
    height: 60px;
    flex-direction: column;
    gap: $spacing-xs;
    border-radius: $radius-md;
    margin-bottom: $spacing-md;
  }
}

.not-logged-in {
  @include flex-center;
  flex-direction: column;
  gap: $spacing-lg;
  min-height: 50vh;
  text-align: center;

  .empty-icon { color: $text-light; opacity: 0.5; }
  h3 { margin: 0; color: $text-primary; }
  p { margin: 0; color: $text-secondary; }
}
</style>
