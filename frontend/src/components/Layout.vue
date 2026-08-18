<template>
  <el-container class="layout-container">
    <el-aside :width="isCollapsed ? '64px' : '220px'" class="sidebar" :class="{ 'is-collapsed': isCollapsed }">
      <div class="logo-area" @click="$router.push('/')">
        <div class="logo-icon">
          <el-icon :size="28"><Headset /></el-icon>
        </div>
        <transition name="fade">
          <span v-show="!isCollapsed" class="logo-text">版权存证</span>
        </transition>
      </div>
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapsed"
        :collapse-transition="false"
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
          <el-icon class="collapse-btn" :size="20" @click="toggleCollapse">
            <Fold v-if="!isCollapsed" />
            <Expand v-else />
          </el-icon>
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
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const isCollapsed = ref(false)

const activeMenu = computed(() => route.path)
const currentTitle = computed(() => route.meta.title)

function toggleCollapse() {
  isCollapsed.value = !isCollapsed.value
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
    // Could navigate to profile page
  }
}
</script>

<style lang="scss" scoped>
.layout-container {
  height: 100vh;
}

.sidebar {
  background: linear-gradient(180deg, #1e1b4b 0%, #312e81 100%);
  transition: width 0.3s;
  overflow: hidden;
}

.logo-area {
  height: 64px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  gap: 12px;
  cursor: pointer;
  border-bottom: 1px solid rgba(255,255,255,0.1);
}

.logo-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  white-space: nowrap;
}

.sidebar-menu {
  border-right: none;
  padding: 12px;

  :deep(.el-menu-item),
  :deep(.el-sub-menu__title) {
    border-radius: 8px;
    margin-bottom: 4px;
    height: 44px;
    line-height: 44px;

    &:hover {
      background: rgba(255,255,255,0.1);
    }
  }

  :deep(.el-menu-item.is-active) {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: #fff;
  }

  :deep(.el-sub-menu .el-menu-item) {
    padding-left: 52px !important;
    min-width: calc(100% - 24px);
    margin-left: 12px;
  }

  /* =========================================================================
   * 收起态对齐：回退到你说的「第一次改法」（能让 2/5 居中的那一版）
   *   margin 左右对称 8px + flex justify-content center
   *   （这版对 padding-left / text-indent 的内部 offset 不敏感，4 个 icon 能一起居中）
   *   仅单独加一条隐藏 > 箭头（不让 3/4 上面叠 >）
   *   ⚠️ 图标若还有轻微偏差，改 margin-left/right 数值（保持左右相等）即可
   * ========================================================================= */
  &:deep(.el-menu--collapse) {
    padding: 12px 0 !important;

    .el-menu-item,
    .el-sub-menu__title {
      padding: 0 !important;
      /* ↓↓↓ 唯一可调参数：左右 margin 对称，控制整个色块的水平位置
             22 = 8+8 对称宽度 + 48 内容 = 64。若整体偏右 → 改大到 10/10；偏左 → 改小到 6/6 ↓↓↓ */
      margin: 0 8px 4px 8px !important;
      width: calc(100% - 16px) !important;
      height: 44px !important;
      line-height: 44px !important;
      min-width: 0 !important;
      text-indent: 0 !important;
      display: flex !important;
      flex-direction: row !important;
      align-items: center !important;
      justify-content: center !important;   /* 让 icon 盒子在 flex 主轴正中 */
      text-align: center;
      overflow: hidden;
    }

    /* 图标盒子：固定 20x20，不占多余空间（去掉 > 直接子选择器，避免包装节点导致匹配不到） */
    .el-menu-item .el-icon,
    .el-sub-menu__title .el-icon {
      margin: 0 !important;
      padding: 0 !important;
      width: 20px !important;
      height: 20px !important;
      line-height: 1;
      flex: 0 0 20px !important;        /* 防止被 flex 压扁 */
      display: block !important;
    }

    /* 隐藏文字（不参与 flex 布局，避免占空间推偏 icon） */
    .el-menu-item span:not(.el-icon),
    .el-sub-menu__title span:not(.el-icon) {
      display: none !important;
      width: 0 !important;
      height: 0 !important;
      padding: 0 !important;
      margin: 0 !important;
      flex: 0 0 0;
      overflow: hidden;
    }

    /* 子菜单内部项（展开态才用得到，收起态同步对齐） */
    .el-sub-menu .el-menu-item {
      padding: 0 !important;
      margin: 0 8px 4px 8px !important;
      width: calc(100% - 16px) !important;
    }

    /* ══════════════════════════════════════
       唯一针对「3/4 上面叠 >」的补丁：
       彻底干掉子菜单展开箭头（包括宽度占位）
       ══════════════════════════════════════ */
    .el-sub-menu__icon-arrow {
      display: none !important;
      width: 0 !important;
      height: 0 !important;
      flex: 0 0 0 !important;
      margin: 0 !important;
      padding: 0 !important;
      overflow: hidden !important;
      opacity: 0 !important;
      visibility: hidden !important;
    }
  }
}

/* 收起态 logo 色块：图标保持左对齐不跳（x 坐标与展开态连续，不做先居中再左移） */
.sidebar.is-collapsed {
  .logo-area {
    padding: 0 12px;
    justify-content: flex-start;
    gap: 0;
  }
  .logo-text {
    display: none;
  }
}

.header {
  background: #fff;
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 64px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-btn {
  cursor: pointer;
  color: var(--text-secondary);
  &:hover { color: var(--primary); }
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 12px;
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
  font-weight: 500;
  color: var(--text-primary);
}

.main-content {
  background: #f5f7fb;
  overflow-y: auto;
  padding: 0;
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

<!-- 全局样式：解决浮层白字 —— 核心思路：用 popper-class="sidebar-menu-popup" 作为唯一锚点
     （类名是我们自己写在 <el-sub-menu> 上的，不会受 EP 版本内部 DOM 类名变化影响）
     任何 Element Plus 版本，挂了自定义 popper-class 后浮层都会包含这个类名，直接针对它写样式最稳妥 -->
<style lang="scss">
/* ============== 自定义浮层：sidebar-menu-popup（版权管理 / 授权管理 的展开面板） ==============
   Element Plus 浮层 DOM 结构：body > div.el-popper.sidebar-menu-popup > .el-menu
   所以用 .sidebar-menu-popup 打头能精确命中这个面板的所有层级，不需要猜 EP 的内部类名 */
.sidebar-menu-popup {
  /* ↓↓↓ 先覆盖 Element Plus 的 CSS 变量（EP 菜单颜色全靠这些变量，直接从根覆盖最稳） ↓↓↓ */
  --el-menu-bg-color:            #ffffff !important;
  --el-menu-text-color:          #000000 !important;
  --el-menu-hover-bg-color:      #dbeafe !important;
  --el-menu-hover-text-color:    #1e3a8a !important;
  --el-menu-active-color:        #ffffff !important;
  --el-menu-border-color:        #cbd5e1 !important;

  /* 面板自身外观：白底 + 灰边 + 阴影 */
  background-color: #ffffff !important;
  background: #ffffff !important;
  border: 1px solid #cbd5e1 !important;
  border-radius: 8px !important;
  padding: 6px !important;
  min-width: 150px !important;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2) !important;

  /* ↓↓↓ 终极兜底：无论内部节点多少层，全纯黑字（包括伪元素），通配 + !important
       特异性 0,1,0 + !important，EP 任何内部 color: var(--xxx) 写法都会被干掉 ↓↓↓ */
  *, *::before, *::after {
    color: #000000 !important;
  }

  /* 子菜单外层容器 .el-menu 再次确认 */
  .el-menu,
  > .el-menu,
  > div > .el-menu {
    background: #ffffff !important;
    border: none !important;
    color: #000000 !important;
    --el-menu-bg-color:     #ffffff !important;
    --el-menu-text-color:   #000000 !important;
  }

  /* —— 每个菜单项（我的作品 / 版权存证 / 作品搜索 / ...） —— */
  .el-menu-item,
  .el-sub-menu__title {
    background: #ffffff !important;
    color: #000000 !important;
    font-weight: 700 !important;
    font-size: 14px !important;
    height: 36px !important;
    line-height: 36px !important;
    border-radius: 6px !important;
    padding: 0 14px !important;
    margin-bottom: 2px !important;
    text-indent: 0 !important;
    display: block !important;

    /* 图标：纯黑（和字同步） */
    .el-icon {
      color: #000000 !important;
      margin-right: 8px !important;
      font-size: 16px !important;
    }

    /* 悬停：淡蓝底 + 深蓝字 */
    &:hover,
    &:focus-visible,
    &.is-hover {
      background: #dbeafe !important;
      color: #1e3a8a !important;
      --el-menu-text-color: #1e3a8a !important;
      .el-icon { color: #1e3a8a !important; }
    }

    /* 当前激活路由：紫色渐变 + 纯白字（和主站 brand 一致） */
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

  /* 浮层箭头：白底 + 灰边 */
  .el-popper__arrow::before,
  [class*='arrow']::before {
    background-color: #ffffff !important;
    background: #ffffff !important;
    border-color: #cbd5e1 !important;
  }
}

/* ============== 单行 tooltip 浮层（工作台 / 争议存证 等 el-menu-item 悬停提示） ==============
   Element Plus tooltip 在 collapse 菜单下自动生成，无自定义 popper-class 钩子。
   这次用 ID 选择器最高特异性 + 直接写 style 的 100% 兜底方案：
   只要是挂在 body 直接子节点上的 tooltip/popper，不管类名是什么，都强制白底黑字 */
html > body > div[role='tooltip'],
html > body > div.el-tooltip__popper,
html > body > div[class*='el-tooltip'],
html > body > div[class*='el-popper'] {
  --el-bg-color-overlay:      #ffffff !important;
  --el-text-color-primary:    #000000 !important;
  --el-color-white:           #000000 !important;   /* 兜底：防止用 color:white 直接写字 */
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
  border-radius: 6px !important;
  border: 1px solid #cbd5e1 !important;
  box-shadow: none !important;
}
</style>
