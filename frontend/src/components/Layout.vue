<template>
  <el-container class="layout-container">
    <el-aside width="260px" class="sidebar">
      <div class="logo-area" @click="$router.push('/')">
        <div class="logo-icon">
          <el-icon :size="28"><Headset /></el-icon>
        </div>
        <span class="logo-text">版权存证</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="transparent"
        text-color="rgba(255,255,255,0.7)"
        active-text-color="#ffffff"
        class="sidebar-menu"
      >
        <el-menu-item index="/">
          <el-icon><Odometer /></el-icon>
          <template #title>工作台</template>
        </el-menu-item>

        <el-sub-menu index="works" popper-class="sidebar-menu-popup">
          <template #title>
            <el-icon><Files /></el-icon>
            <span>版权管理</span>
          </template>
          <el-menu-item index="/works">
            <el-icon><Folder /></el-icon>
            <template #title>我的作品</template>
          </el-menu-item>
          <el-menu-item index="/works/register">
            <el-icon><DocumentAdd /></el-icon>
            <template #title>版权存证</template>
          </el-menu-item>
          <el-menu-item index="/works/search">
            <el-icon><Search /></el-icon>
            <template #title>作品搜索</template>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="licenses" popper-class="sidebar-menu-popup">
          <template #title>
            <el-icon><Key /></el-icon>
            <span>授权管理</span>
          </template>
          <el-menu-item index="/licenses/grant">
            <el-icon><Tickets /></el-icon>
            <template #title>发放授权</template>
          </el-menu-item>
          <el-menu-item index="/licenses/verify">
            <el-icon><CircleCheck /></el-icon>
            <template #title>授权核验</template>
          </el-menu-item>
          <el-menu-item index="/licenses/my">
            <el-icon><Collection /></el-icon>
            <template #title>我的授权</template>
          </el-menu-item>
        </el-sub-menu>

        <el-menu-item index="/disputes">
          <el-icon><Warning /></el-icon>
          <template #title>争议存证</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentTitle">{{ currentTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown trigger="click" @command="handleCommand">
            <div class="user-info">
              <el-avatar :size="36" class="user-avatar">
                {{ userStore.userInfo?.username?.[0]?.toUpperCase() || 'U' }}
              </el-avatar>
              <span class="username">{{ userStore.userInfo?.realName || userStore.userInfo?.username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>个人信息
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>

      <el-dialog v-model="profileVisible" title="个人信息" width="420px">
        <div v-loading="profileLoading" class="profile-dialog">
          <div class="profile-avatar">
            {{ profileData?.username?.[0]?.toUpperCase() || 'U' }}
          </div>
          <div class="profile-row"><span class="label">用户名</span><span class="value">{{ profileData?.username }}</span></div>
          <div class="profile-row"><span class="label">真实姓名</span><span class="value">{{ profileData?.realName || '未填写' }}</span></div>
          <div class="profile-row"><span class="label">用户ID</span><span class="value">{{ profileData?.id }}</span></div>
          <div class="profile-row"><span class="label">注册时间</span><span class="value">{{ profileData?.createdAt }}</span></div>
        </div>
      </el-dialog>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const currentTitle = computed(() => route.meta.title)

const profileVisible = ref(false)
const profileLoading = ref(false)
const profileData = ref(null)

async function showProfile() {
  profileVisible.value = true
  profileLoading.value = true
  try {
    profileData.value = await userStore.fetchProfile()
  } catch (e) {
    // interceptor 已提示错误
  } finally {
    profileLoading.value = false
  }
}

function handleCommand(command) {
  if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      userStore.logout()
      router.push('/login')
    }).catch(() => {})
  } else if (command === 'profile') {
    showProfile()
  }
}
</script>

<style lang="scss" scoped>
.layout-container {
  height: 100vh;
}

.sidebar {
  width: 260px !important;
  background: linear-gradient(180deg, #1e1b4b 0%, #312e81 100%);
  overflow: hidden;
}

.logo-area {
  height: 80px;
  display: flex;
  align-items: center;
  padding: 0 28px;
  gap: 16px;
  cursor: pointer;
  border-bottom: 1px solid rgba(255,255,255,0.1);
}

.logo-icon {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.logo-text {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  white-space: nowrap;
}

.sidebar-menu {
  border-right: none;
  padding: 16px;

  :deep(.el-menu-item),
  :deep(.el-sub-menu__title) {
    border-radius: 12px;
    margin-bottom: 6px;
    height: 48px;
    line-height: 48px;
    font-size: 15px;

    &:hover {
      background: rgba(255,255,255,0.1);
    }
  }

  :deep(.el-menu-item.is-active) {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: #fff;
  }

  :deep(.el-sub-menu .el-menu-item) {
    padding-left: 56px !important;
    min-width: calc(100% - 24px);
    margin-left: 12px;
    height: 42px;
    line-height: 42px;
  }
}

.header {
  background: #fff;
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 40px;
  height: 80px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
  font-size: 16px;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding: 8px 18px;
  border-radius: 24px;
  transition: background 0.2s;

  &:hover {
    background: #f1f5f9;
  }
}

.user-avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-weight: 600;
}

.username {
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
}

.main-content {
  background: #f5f7fb;
  overflow-y: auto;
  padding: 0;
}

.profile-dialog {
  padding: 8px 0;

  .profile-avatar {
    width: 72px;
    height: 72px;
    margin: 0 auto 20px;
    border-radius: 50%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: #fff;
    font-size: 28px;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .profile-row {
    display: flex;
    justify-content: space-between;
    padding: 12px 4px;
    border-bottom: 1px dashed var(--border);

    &:last-child {
      border-bottom: none;
    }

    .label {
      color: var(--text-secondary, #64748b);
      font-size: 14px;
    }

    .value {
      font-weight: 600;
      font-size: 14px;
      color: var(--text-primary, #1e293b);
      max-width: 220px;
      word-break: break-all;
      text-align: right;
    }
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

<!-- 全局样式：解决浮层白字 -->
<style lang="scss">
.sidebar-menu-popup {
  --el-menu-bg-color:            #ffffff !important;
  --el-menu-text-color:          #000000 !important;
  --el-menu-hover-bg-color:      #dbeafe !important;
  --el-menu-hover-text-color:    #1e3a8a !important;
  --el-menu-active-color:        #ffffff !important;
  --el-menu-border-color:        #cbd5e1 !important;

  background-color: #ffffff !important;
  background: #ffffff !important;
  border: 1px solid #cbd5e1 !important;
  border-radius: 20px !important;
  padding: 10px !important;
  min-width: 180px !important;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2) !important;

  *, *::before, *::after {
    color: #000000 !important;
  }

  .el-menu,
  > .el-menu,
  > div > .el-menu {
    background: #ffffff !important;
    border: none !important;
    color: #000000 !important;
    --el-menu-bg-color:     #ffffff !important;
    --el-menu-text-color:   #000000 !important;
  }

  .el-menu-item,
  .el-sub-menu__title {
    background: #ffffff !important;
    color: #000000 !important;
    font-weight: 700 !important;
    font-size: 14px !important;
    height: 40px !important;
    line-height: 40px !important;
    border-radius: 10px !important;
    padding: 0 16px !important;
    margin-bottom: 4px !important;
    text-indent: 0 !important;
    display: block !important;

    .el-icon {
      color: #000000 !important;
      margin-right: 10px !important;
      font-size: 16px !important;
    }

    &:hover,
    &:focus-visible,
    &.is-hover {
      background: #dbeafe !important;
      color: #1e3a8a !important;
      --el-menu-text-color: #1e3a8a !important;
      .el-icon { color: #1e3a8a !important; }
    }

    &.is-active,
    &[aria-current='true'] {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
      color: #ffffff !important;
      --el-menu-text-color: #ffffff !important;
      --el-menu-active-color: #ffffff !important;
      .el-icon { color: #ffffff !important; }
      &:hover {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
        color: #ffffff !important;
      }
    }
  }

  .el-popper__arrow::before,
  [class*='arrow']::before {
    background-color: #ffffff !important;
    background: #ffffff !important;
    border-color: #cbd5e1 !important;
  }
}

html > body > div[role='tooltip'],
html > body > div.el-tooltip__popper,
html > body > div[class*='el-tooltip'],
html > body > div[class*='el-popper'] {
  --el-bg-color-overlay:      #ffffff !important;
  --el-text-color-primary:    #000000 !important;
  --el-color-white:           #000000 !important;
  --el-text-color-regular:    #000000 !important;
  color: #000000 !important;
  background: #ffffff !important;
  border: 1px solid #cbd5e1 !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

html > body > div[role='tooltip'] *,
html > body > div.el-tooltip__popper *,
html > body > div[class*='el-tooltip'] *,
html > body > div[class*='el-popper'] * {
  color: #000000 !important;
}

html > body > div[role='tooltip'] .el-tooltip__inner,
html > body > div.el-tooltip__popper .el-tooltip__inner {
  background: #ffffff !important;
  color: #000000 !important;
  font-weight: 700 !important;
  font-size: 13px !important;
  padding: 8px 14px !important;
  border-radius: 10px !important;
  border: 1px solid #cbd5e1 !important;
  box-shadow: none !important;
}
</style>
