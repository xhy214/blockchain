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
  min-height: 560px;
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
  width: 130px;
  height: 130px;
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
  font-size: 13px;
  line-height: 1.8;
  color: var(--text-light);
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
    width: 100px;
    height: 100px;
  }
  .form-side { padding: 32px 24px; }
}
</style>
