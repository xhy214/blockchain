<template>
  <div class="auth-page">
    <div class="auth-bg">
      <div class="bg-shape shape1"></div>
      <div class="bg-shape shape2"></div>
      <div class="bg-shape shape3"></div>
    </div>

    <div class="auth-container">
      <div class="brand-side">
        <div class="brand-logo">
          <el-icon :size="48"><Headset /></el-icon>
        </div>
        <h1 class="brand-title">音乐数字<br/>版权保护平台</h1>
        <p class="brand-desc">基于 Hyperledger Fabric 区块链技术，为音乐作品提供不可篡改的版权存证与授权管理服务</p>
        <div class="brand-features">
          <div class="feature-item">
            <el-icon :size="20"><Lock /></el-icon>
            <span>区块链存证</span>
          </div>
          <div class="feature-item">
            <el-icon :size="20"><Key /></el-icon>
            <span>授权核验</span>
          </div>
          <div class="feature-item">
            <el-icon :size="20"><Medal /></el-icon>
            <span>存证证书</span>
          </div>
        </div>
      </div>

      <div class="form-side">
        <div class="form-header">
          <h2>欢迎登录</h2>
          <p>登录账号以使用版权保护功能</p>
        </div>

        <el-form ref="formRef" :model="form" :rules="rules" size="large" @keyup.enter="handleLogin">
          <el-form-item prop="username">
            <el-input v-model="form.username" placeholder="请输入用户名" :prefix-icon="User" />
          </el-form-item>

          <el-form-item prop="password">
            <el-input v-model="form.password" type="password" placeholder="请输入密码" :prefix-icon="Lock" show-password />
          </el-form-item>

          <el-button type="primary" class="submit-btn" :loading="loading" @click="handleLogin">
            登录
          </el-button>
        </el-form>

        <div class="form-footer">
          <span>还没有账号？</span>
          <router-link to="/register" class="link">立即注册</router-link>
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
import { User, Lock } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

async function handleLogin() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await userStore.login({ ...form })
    ElMessage.success('登录成功')
    const redirect = router.currentRoute.value.query.redirect || '/'
    router.push(redirect)
  } catch (e) {
    // error already shown by interceptor
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

.shape1 { width: 400px; height: 400px; top: -100px; right: -100px; }
.shape2 { width: 300px; height: 300px; bottom: -50px; left: -50px; }
.shape3 { width: 200px; height: 200px; top: 40%; right: 10%; }

.auth-container {
  position: relative;
  display: flex;
  background: #fff;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0,0,0,0.2);
  max-width: 860px;
  width: 100%;
  min-height: 500px;
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
  margin-bottom: 32px;
}

.brand-features {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  opacity: 0.9;
}

.form-side {
  flex: 1;
  padding: 48px 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.form-header {
  margin-bottom: 32px;

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
  .brand-side {
    padding: 32px 24px;
  }
  .form-side {
    padding: 32px 24px;
  }
}
</style>
