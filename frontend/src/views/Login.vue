<template>
  <div class="auth-page">
    <div class="auth-bg">
      <div class="bg-glow glow-gold"></div>
      <div class="bg-glow glow-blue"></div>
    </div>

    <div class="auth-container">
      <div class="brand-side">
        <div class="vinyl brand-vinyl"></div>
        <div class="staff-divider brand-divider"></div>
        <h1 class="brand-title">音乐数字<br/>版权保护平台</h1>
        <p class="brand-desc">基于 Hyperledger Fabric 区块链技术，为音乐作品提供不可篡改的版权存证与授权管理服务</p>
        <div class="brand-features">
          <div class="feature-item">
            <el-icon :size="18"><Lock /></el-icon>
            <span>区块链存证</span>
          </div>
          <div class="feature-item">
            <el-icon :size="18"><Key /></el-icon>
            <span>授权核验</span>
          </div>
          <div class="feature-item">
            <el-icon :size="18"><Medal /></el-icon>
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
  background: var(--bg-deep);
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

.bg-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
}

.glow-gold {
  width: 560px;
  height: 560px;
  top: -180px;
  left: -120px;
  background: rgba(201, 168, 106, 0.13);
}

.glow-blue {
  width: 640px;
  height: 640px;
  bottom: -240px;
  right: -160px;
  background: rgba(62, 82, 118, 0.22);
}

.auth-container {
  position: relative;
  display: flex;
  background: var(--surface);
  border: 1px solid var(--line-strong);
  border-radius: var(--radius-panel);
  overflow: hidden;
  box-shadow: 0 32px 90px rgba(0, 0, 0, 0.55);
  max-width: 880px;
  width: 100%;
  min-height: 520px;
}

.brand-side {
  flex: 1;
  background: linear-gradient(160deg, #0C121E 0%, #101927 100%);
  border-right: 1px solid var(--line);
  padding: 48px 40px;
  color: var(--text-primary);
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.brand-vinyl {
  width: 150px;
  height: 150px;
  margin-bottom: 28px;
}

.brand-divider {
  margin-bottom: 24px;
}

.brand-title {
  font-family: var(--font-display);
  font-size: 27px;
  font-weight: 700;
  letter-spacing: 0.06em;
  line-height: 1.4;
  margin-bottom: 14px;
  color: var(--text-primary);
}

.brand-desc {
  font-size: 14px;
  line-height: 1.8;
  color: var(--text-light);
  margin-bottom: 28px;
}

.brand-features {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: var(--text-secondary);
  letter-spacing: 0.06em;

  .el-icon {
    color: var(--accent);
  }
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
    font-family: var(--font-display);
    font-size: 24px;
    font-weight: 700;
    letter-spacing: 0.05em;
    margin-bottom: 8px;
  }

  p {
    color: var(--text-light);
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
  border-radius: var(--radius-control);
  color: var(--on-accent);
  margin-top: 8px;
}

.form-footer {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: var(--text-light);

  .link {
    color: var(--accent-bright);
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
    border-right: none;
    border-bottom: 1px solid var(--line);
  }
  .brand-vinyl {
    width: 110px;
    height: 110px;
  }
  .form-side {
    padding: 32px 24px;
  }
}
</style>
