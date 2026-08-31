<template>
  <el-container class="layout-container">
    <el-aside width="260px" class="sidebar">
      <div class="logo-area" @click="$router.push('/')">
        <div class="logo-vinyl vinyl"></div>
        <span class="logo-text">版权存证</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="transparent"
        text-color="#A9B1C0"
        active-text-color="#DCC189"
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
  background: var(--bg-deep);
  border-right: 1px solid var(--line);
  overflow: hidden;
}

.logo-area {
  height: 80px;
  display: flex;
  align-items: center;
  padding: 0 28px;
  gap: 16px;
  cursor: pointer;
  border-bottom: 1px solid var(--line);
}

.logo-vinyl {
  width: 44px;
  height: 44px;
  flex-shrink: 0;
}

.logo-text {
  font-family: var(--font-display);
  font-size: 21px;
  font-weight: 700;
  letter-spacing: 0.12em;
  color: var(--text-primary);
  white-space: nowrap;
}

.sidebar-menu {
  border-right: none;
  padding: 16px;

  :deep(.el-menu-item),
  :deep(.el-sub-menu__title) {
    border-radius: var(--radius-control);
    margin-bottom: 6px;
    height: 48px;
    line-height: 48px;
    font-size: 15px;

    &:hover {
      background: var(--sidebar-hover);
      color: var(--accent-bright);
    }
  }

  :deep(.el-menu-item.is-active) {
    background: var(--accent-soft);
    color: var(--accent-bright);
    font-weight: 600;
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
  background: var(--surface);
  border-bottom: 1px solid var(--line);
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
  font-size: 15px;
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
    background: var(--surface-2);
  }
}

.user-avatar {
  background: var(--accent-soft);
  border: 1px solid rgba(201, 168, 106, 0.35);
  color: var(--accent-bright);
  font-weight: 600;
}

.username {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-primary);
}

.main-content {
  background: var(--bg);
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
    background: var(--accent-soft);
    border: 1px solid rgba(201, 168, 106, 0.35);
    color: var(--accent-bright);
    font-family: var(--font-display);
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
    border-bottom: 1px dashed var(--line-strong);

    &:last-child {
      border-bottom: none;
    }

    .label {
      color: var(--text-secondary);
      font-size: 14px;
    }

    .value {
      font-weight: 600;
      font-size: 14px;
      color: var(--text-primary);
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

<!-- 全局样式：侧边栏弹出菜单（暗色） -->
<style lang="scss">
.sidebar-menu-popup {
  --el-menu-bg-color:         #1B2536 !important;
  --el-menu-text-color:       #C7CCD6 !important;
  --el-menu-hover-bg-color:   rgba(201, 168, 106, 0.12) !important;
  --el-menu-hover-text-color: #DCC189 !important;
  --el-menu-active-color:     #DCC189 !important;
  --el-menu-border-color:     #2A3548 !important;

  background-color: #1B2536 !important;
  background: #1B2536 !important;
  border: 1px solid #2A3548 !important;
  border-radius: 12px !important;
  padding: 10px !important;
  min-width: 180px !important;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.5) !important;

  .el-menu,
  > .el-menu,
  > div > .el-menu {
    background: #1B2536 !important;
    border: none !important;
    --el-menu-bg-color:   #1B2536 !important;
    --el-menu-text-color: #C7CCD6 !important;
  }

  .el-menu-item,
  .el-sub-menu__title {
    background: transparent !important;
    color: #C7CCD6 !important;
    font-size: 14px !important;
    height: 40px !important;
    line-height: 40px !important;
    border-radius: 8px !important;
    padding: 0 16px !important;
    margin-bottom: 4px !important;
    text-indent: 0 !important;
    display: block !important;

    .el-icon {
      color: #8B93A3 !important;
      margin-right: 10px !important;
      font-size: 16px !important;
    }

    &:hover,
    &:focus-visible,
    &.is-hover {
      background: rgba(201, 168, 106, 0.12) !important;
      color: #DCC189 !important;
      --el-menu-text-color: #DCC189 !important;
      .el-icon { color: #DCC189 !important; }
    }

    &.is-active,
    &[aria-current='true'] {
      background: rgba(201, 168, 106, 0.18) !important;
      color: #DCC189 !important;
      --el-menu-text-color: #DCC189 !important;
      --el-menu-active-color: #DCC189 !important;
      font-weight: 600 !important;
      .el-icon { color: #DCC189 !important; }
      &:hover {
        background: rgba(201, 168, 106, 0.18) !important;
        color: #DCC189 !important;
      }
    }
  }

  .el-popper__arrow::before,
  [class*='arrow']::before {
    background-color: #1B2536 !important;
    background: #1B2536 !important;
    border-color: #2A3548 !important;
  }
}
</style>
