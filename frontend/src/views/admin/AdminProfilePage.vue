<template>
  <div class="admin-profile-page">
    <div class="page-header">
      <h2>个人信息</h2>
    </div>

    <el-row :gutter="20">
      <el-col :xs="24" :md="8">
        <el-card class="profile-card">
          <div class="avatar-section">
            <el-avatar :size="100" :src="userInfo.avatar_url">
              <el-icon :size="50"><User /></el-icon>
            </el-avatar>
            <h3>{{ userInfo.username }}</h3>
            <el-tag type="danger">管理员</el-tag>
            <p class="email">{{ userInfo.email }}</p>
          </div>
          <el-divider />
          <div class="info-list">
            <div class="info-item">
              <span class="label">上次登录</span>
              <span class="value">{{ formatTime(userInfo.last_login_at) }}</span>
            </div>
            <div class="info-item">
              <span class="label">注册时间</span>
              <span class="value">{{ formatTime(userInfo.created_at) }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="16">
        <el-card class="form-card">
          <template #header><h3>编辑资料</h3></template>
          <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="form.username" placeholder="请输入用户名" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" />
            </el-form-item>
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="form.phone" placeholder="请输入手机号" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSave" :loading="saving">保存修改</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card class="form-card" style="margin-top: 20px;">
          <template #header><h3>修改密码</h3></template>
          <el-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-width="100px">
            <el-form-item label="当前密码" prop="old_password">
              <el-input v-model="pwdForm.old_password" type="password" show-password placeholder="请输入当前密码" />
            </el-form-item>
            <el-form-item label="新密码" prop="new_password">
              <el-input v-model="pwdForm.new_password" type="password" show-password placeholder="请输入新密码" />
            </el-form-item>
            <el-form-item label="确认密码" prop="confirm_password">
              <el-input v-model="pwdForm.confirm_password" type="password" show-password placeholder="请再次输入新密码" />
            </el-form-item>
            <el-form-item>
              <el-button type="warning" @click="handleChangePassword" :loading="changingPwd">修改密码</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'

const authStore = useAuthStore()
const formRef = ref(null)
const pwdFormRef = ref(null)
const saving = ref(false)
const changingPwd = ref(false)

const userInfo = computed(() => authStore.user || {})

const form = reactive({
  username: '',
  email: '',
  phone: ''
})

const pwdForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

const pwdRules = {
  old_password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: (rule, value, callback) => {
      if (value !== pwdForm.new_password) callback(new Error('两次输入的密码不一致'))
      else callback()
    }, trigger: 'blur' }
  ]
}

const formatTime = (time) => time ? new Date(time).toLocaleString('zh-CN') : '-'

const initForm = () => {
  form.username = userInfo.value.username || ''
  form.email = userInfo.value.email || ''
  form.phone = userInfo.value.phone || ''
}

const handleSave = async () => {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    const result = await authStore.updateProfile(form)
    if (result.success) {
      ElMessage.success('保存成功')
    } else {
      ElMessage.error(result.errors?.[0] || '保存失败')
    }
  } finally {
    saving.value = false
  }
}

const handleChangePassword = async () => {
  if (!pwdFormRef.value) return
  const valid = await pwdFormRef.value.validate().catch(() => false)
  if (!valid) return

  changingPwd.value = true
  try {
    const result = await authStore.changePassword(pwdForm)
    if (result.success) {
      ElMessage.success('密码修改成功')
      pwdFormRef.value.resetFields()
    } else {
      ElMessage.error(result.errors?.[0] || '密码修改失败')
    }
  } finally {
    changingPwd.value = false
  }
}

onMounted(() => initForm())
</script>

<style lang="scss" scoped>
.admin-profile-page { padding: 20px; background: #f8f9fa; min-height: 100vh; }
.page-header { margin-bottom: 20px; h2 { margin: 0; color: #2c3e50; } }
.profile-card {
  text-align: center;
  .avatar-section {
    padding: 20px 0;
    h3 { margin: 16px 0 8px; color: #2c3e50; }
    .email { color: #909399; font-size: 14px; margin-top: 8px; }
  }
  .info-list {
    .info-item {
      display: flex; justify-content: space-between; padding: 12px 0;
      border-bottom: 1px solid #f0f0f0;
      &:last-child { border-bottom: none; }
      .label { color: #909399; }
      .value { color: #2c3e50; font-weight: 500; }
    }
  }
}
.form-card { h3 { margin: 0; color: #2c3e50; font-size: 16px; } }
@media (max-width: 768px) { .el-col { margin-bottom: 20px; } }
</style>
