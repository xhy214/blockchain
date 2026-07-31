<template>
  <div class="auth-page">
    <div class="auth-bg">
      <div class="bg-shape shape1"></div>
      <div class="bg-shape shape2"></div>
    </div>

    <div class="auth-container">
      <div class="brand-side">
        <div class="brand-logo">
          <el-icon :size="48"><Headset /></el-icon>
        </div>
        <h1 class="brand-title">加入我们<br/>开始存证</h1>
        <p class="brand-desc">注册账号后，即可上传原创音乐作品进行区块链存证，并享受完整的版权保护服务</p>
      </div>

      <div class="form-side">
        <div class="form-header">
          <h2>创建账号</h2>
          <p>注册您的账户，开启版权保护之旅</p>
        </div>

        <el-form ref="formRef" :model="form" :rules="rules" size="large" @keyup.enter="handleRegister">
          <el-form-item prop="realName">
            <el-input v-model="form.realName" placeholder="请输入真实姓名" :prefix-icon="User" />
          </el-form-item>

          <el-form-item prop="username">
            <el-input v-model="form.username" placeholder="请输入用户名（登录账号）" :prefix-icon="Postcard" />
          </el-form-item>

          <el-form-item prop="password">
            <el-input v-model="form.password" type="password" placeholder="请输入密码（至少6位）" :prefix-icon="Lock" show-password />
          </el-form-item>

          <el-form-item prop="confirmPassword">
            <el-input v-model="form.confirmPassword" type="password" placeholder="请再次输入密码" :prefix-icon="Lock" show-password />
          </el-form-item>

          <el-button type="primary" class="submit-btn" :loading="loading" @click="handleRegister">
            注册账号
          </el-button>
        </el-form>

        <div class="form-footer">
          <span>已有账号？</span>
          <router-link to="/login" class="link">返回登录</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { User, Lock, Postcard } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  realName: '',
  username: '',
  password: '',
  confirmPassword: ''
})

const rules = {
  realName: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 32, message: '用户名长度在 3 到 32 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少 6 位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== form.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

async function handleRegister() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await userStore.register({
      username: form.username,
      password: form.password,
      realName: form.realName
    })
    ElMessage.success('注册成功，请登录')
    router.push('/login')
  } catch (e) {
    // error already shown
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.auth-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.auth-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255,255,255,0.08);
}

.shape1 { width: 350px; height: 350px; top: -80px; right: -80px; }
.shape2 { width: 250px; height: 250px; bottom: -40px; left: -40px; }

.auth-container {
  position: relative;
  display: flex;
  background: #fff;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0,0,0,0.2);
  max-width: 860px;
  width: 100%;
  min-height: 560px;
}

.brand-side {
  flex: 1;
  background: linear-gradient(135deg, #1e1b4b 0%, #4338ca 100%);
  padding: 48px 40px;
  color: #fff;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.brand-logo {
  width: 72px;
  height: 72px;
  border-radius: 16px;
  background: rgba(255,255,255,0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  margin-bottom: 28px;
}

.brand-title {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.3;
  margin-bottom: 16px;
}

.brand-desc {
  font-size: 14px;
  line-height: 1.7;
  opacity: 0.85;
}

.form-side {
  flex: 1;
  padding: 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.form-header {
  margin-bottom: 28px;

  h2 {
    font-size: 24px;
    font-weight: 700;
    margin-bottom: 8px;
  }

  p {
    color: var(--text-secondary);
    font-size: 14px;
  }
}

.submit-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  background: var(--bg-gradient);
  border: none;
  border-radius: 10px;
  margin-top: 8px;
}

.form-footer {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: var(--text-secondary);

  .link {
    color: var(--primary);
    font-weight: 500;
    margin-left: 4px;
    &:hover { text-decoration: underline; }
  }
}

@media (max-width: 768px) {
  .auth-container {
    flex-direction: column;
    min-height: auto;
  }
  .brand-side { padding: 32px 24px; }
  .form-side { padding: 32px 24px; }
}
</style>
