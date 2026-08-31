<template>
  <div class="dashboard">
    <div class="dashboard-inner">

      <!-- ============ 四个统计卡片 ============ -->
      <div class="stat-cards-row">
        <div class="stat-card" @click="$router.push('/works')">
          <div class="stat-icon-wrap">
            <el-icon :size="36"><Folder /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.works }}</div>
            <div class="stat-label">我的作品</div>
          </div>
        </div>

        <div class="stat-card" @click="$router.push('/licenses/my')">
          <div class="stat-icon-wrap">
            <el-icon :size="36"><Key /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.licenses }}</div>
            <div class="stat-label">持有授权</div>
          </div>
        </div>

        <div class="stat-card" @click="$router.push('/works/search')">
          <div class="stat-icon-wrap">
            <el-icon :size="36"><Search /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.searchCount }}</div>
            <div class="stat-label">作品搜索</div>
          </div>
        </div>

        <div class="stat-card" @click="$router.push('/disputes')">
          <div class="stat-icon-wrap">
            <el-icon :size="36"><Warning /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.disputes }}</div>
            <div class="stat-label">争议存证</div>
          </div>
        </div>
      </div>

      <!-- ============ 下方：主版块(1500) + 右侧面板(500) ============ -->
      <div class="content-row">

        <!-- 主版块：最近存证的作品 -->
        <div class="main-panel">
          <div class="panel-header">
            <div class="panel-title">
              <div class="mini-vinyl vinyl"></div>
              最近存证的作品
            </div>
            <div class="panel-actions">
              <el-button type="primary" round size="large" @click="$router.push('/works/register')">
                <el-icon><DocumentAdd /></el-icon>&nbsp;新建存证
              </el-button>
              <el-button round size="large" @click="$router.push('/works')">查看全部</el-button>
            </div>
          </div>

          <div class="panel-body" v-loading="loadingWorks">
            <div v-if="recentWorks.length === 0 && !loadingWorks" class="empty-state">
              <el-icon :size="48" color="#6E7889"><Folder /></el-icon>
              <p>暂无存证作品</p>
              <el-button type="primary" @click="$router.push('/works/register')">立即存证</el-button>
            </div>

            <template v-else>
              <div class="works-list">
                <div
                  v-for="work in recentWorks"
                  :key="work.workID"
                  class="work-item"
                  @click="$router.push(`/works/${work.workID}`)"
                >
                  <div class="work-cover" :style="{ background: getCoverGradient(work.genre) }">
                    <el-icon :size="28"><Headset /></el-icon>
                  </div>
                  <div class="work-meta">
                    <div class="work-title">{{ work.title }}</div>
                    <div class="work-sub">
                      <span>{{ work.artist }}</span>
                      <span class="dot">·</span>
                      <span>{{ work.genre }}</span>
                    </div>
                  </div>
                  <div class="work-status">
                    <el-tag :type="statusTagType(work.status)" size="default" round>
                      {{ statusLabel(work.status) }}
                    </el-tag>
                  </div>
                  <div class="work-time">
                    {{ formatTime(work.registerAt) }}
                  </div>
                  <el-icon class="work-arrow" :size="18"><ArrowRight /></el-icon>
                </div>
              </div>
            </template>
          </div>
        </div>

        <!-- 右侧面板组 -->
        <div class="right-panels">

          <!-- 使用提示 -->
          <div class="side-panel tips-panel">
            <div class="panel-header">
              <div class="panel-title">
                <el-icon :size="26"><InfoFilled /></el-icon>
                使用提示
              </div>
            </div>
            <div class="panel-body">
              <div class="tip-list">
                <div class="tip-item">
                  <div class="tip-step">步骤 1</div>
                  <div class="tip-text">上传音频文件进行版权存证，获取区块链唯一存证编号</div>
                </div>
                <div class="tip-item">
                  <div class="tip-step">步骤 2</div>
                  <div class="tip-text">向他人发放授权，设置授权类型、有效期和使用次数</div>
                </div>
                <div class="tip-item">
                  <div class="tip-step">步骤 3</div>
                  <div class="tip-text">播放前核验授权有效性，保障版权合规使用</div>
                </div>
                <div class="tip-item">
                  <div class="tip-step">步骤 4</div>
                  <div class="tip-text">随时生成存证证书，作为版权归属的法律凭证</div>
                </div>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import api from '@/api'

const userStore = useUserStore()
const loadingWorks = ref(false)
const recentWorks = ref([])

const stats = reactive({
  works: 0,
  licenses: 0,
  searchCount: 0,
  disputes: 0
})

onMounted(async () => {
  loadingWorks.value = true
  try {
    const res = await api.get('/copyright/my/list')
    recentWorks.value = (res.data || []).slice(0, 5)
    stats.works = recentWorks.value.length
  } catch (e) {
    // 用演示数据
    recentWorks.value = [
      { workID: 'demo-1', title: '静夜思', artist: '李白', genre: '古典', status: 'ACTIVE', registerAt: Date.now() - 86400000 },
      { workID: 'demo-2', title: '将进酒', artist: '李白', genre: '古典', status: 'ACTIVE', registerAt: Date.now() - 172800000 },
      { workID: 'demo-3', title: '琵琶行', artist: '白居易', genre: '古典', status: 'DISPUTED', registerAt: Date.now() - 259200000 },
      { workID: 'demo-4', title: '春晓', artist: '孟浩然', genre: '古典', status: 'TRANSFERRED', registerAt: Date.now() - 345600000 },
      { workID: 'demo-5', title: '登鹳雀楼', artist: '王之涣', genre: '古典', status: 'ACTIVE', registerAt: Date.now() - 432000000 }
    ]
    stats.works = recentWorks.value.length
  } finally {
    loadingWorks.value = false
  }

  try {
    const res = await api.get('/license/my')
    stats.licenses = (res.data || []).length
  } catch (e) {
    stats.licenses = 3
  }

  try {
    const res = await api.get('/copyright/my/list')
    stats.searchCount = (res.data || []).length
  } catch (e) {
    stats.searchCount = 0
  }
})

function statusTagType(status) {
  const map = { ACTIVE: 'success', TRANSFERRED: 'info', DISPUTED: 'warning' }
  return map[status] || 'info'
}

function statusLabel(status) {
  const map = { ACTIVE: '正常', TRANSFERRED: '已转让', DISPUTED: '争议中' }
  return map[status] || status
}

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

function getCoverGradient(genre) {
  const map = {
    '古典': 'linear-gradient(135deg, #C9A86A 0%, #8A6A33 100%)',
    '流行': 'linear-gradient(135deg, #8A5E68 0%, #4A3840 100%)',
    '摇滚': 'linear-gradient(135deg, #8A4A4A 0%, #4E2A2E 100%)',
    '电子': 'linear-gradient(135deg, #3E6E8A 0%, #24404E 100%)',
    '民谣': 'linear-gradient(135deg, #5E7E58 0%, #334632 100%)',
    '爵士': 'linear-gradient(135deg, #8A6A4A 0%, #4A3826 100%)'
  }
  return map[genre] || 'linear-gradient(135deg, #4A5A78 0%, #28324A 100%)'
}
</script>

<style lang="scss" scoped>
/* 核心思路：保持原始设计像素值不变，用 zoom 等比缩放整个容器 */
/* 设计基准：2180px 宽度（不含侧边栏） */
/* 缩放因子 = (视口宽度 - 侧边栏260px) / 2180 */

.dashboard {
  height: calc(100vh - 80px);
  width: 100%;
  background: var(--bg);
  overflow: hidden;
}

.dashboard-inner {
  width: 2180px;
  height: 100%;
  margin: 0 auto;
  padding: 40px 0 40px 120px;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  /* 关键：等比缩放，完美保留原始设计比例 */
  zoom: calc((100vw - 260px) / 2180px);
}

/* ============ 统计卡片行 ============ */
.stat-cards-row {
  display: flex;
  /* 4×255 + 3×140 = 1440，与主面板右缘(1560)对齐 */
  gap: 140px;
  margin-bottom: 40px;
  flex-shrink: 0;
}

.stat-card {
  width: 255px;
  height: 170px;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  cursor: pointer;
  transition: transform 0.25s ease, box-shadow 0.25s ease, border-color 0.25s ease;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.28);

  &:hover {
    transform: translateY(-4px);
    border-color: rgba(201, 168, 106, 0.4);
    box-shadow: 0 16px 40px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(201, 168, 106, 0.12);
  }
}

.stat-icon-wrap {
  width: 72px;
  height: 72px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: var(--accent-soft);
  border: 1px solid rgba(201, 168, 106, 0.28);
  color: var(--accent-bright);
}

.stat-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.stat-value {
  /* 数字用无衬线 + tabular-nums：衬线展示体的 "0" 与字母 O 无法区分，辨识度差 */
  font-family: var(--font-body);
  font-variant-numeric: tabular-nums;
  font-size: 44px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
}

.stat-label {
  font-size: 16px;
  color: var(--text-light);
  font-weight: 500;
}

/* ============ 下方内容行 ============ */
.content-row {
  display: flex;
  gap: 60px;
  align-items: stretch;
  flex: 1;
  min-height: 0;
}

/* 主版块 */
.main-panel {
  width: 1500px;
  height: 100%;
  margin-left: -60px;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.28);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  flex-shrink: 0;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 28px;
  border-bottom: 1px solid var(--line);
  flex-shrink: 0;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.03em;
  color: var(--text-primary);
}

.mini-vinyl {
  width: 26px;
  height: 26px;
  flex-shrink: 0;
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.06);
}

.panel-actions {
  display: flex;
  gap: 12px;
}

.panel-body {
  padding: 16px 28px 24px;
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  gap: 12px;
  color: var(--text-light);

  p {
    font-size: 16px;
    margin: 0;
  }
}

/* 作品列表 */
.works-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.work-item {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 16px 20px;
  background: var(--surface-2);
  border: 1px solid var(--line);
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;

  &:hover {
    background: #1D2739;
    border-color: var(--line-strong);
  }
}

.work-cover {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.9);
  flex-shrink: 0;
}

.work-meta {
  flex: 1;
  min-width: 0;
}

.work-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.work-sub {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-light);

  .dot {
    display: inline-block;
    width: 4px;
    height: 4px;
    border-radius: 50%;
    background: var(--line-strong);
  }
}

.work-status {
  flex-shrink: 0;
}

.work-time {
  font-size: 14px;
  color: var(--text-light);
  flex-shrink: 0;
  min-width: 140px;
  text-align: right;
}

.work-arrow {
  color: var(--text-light);
  flex-shrink: 0;
}

/* ============ 右侧面板组 ============ */
.right-panels {
  width: 500px;
  /* 上移 210px 与统计卡并排，高度补回 210px，使底部与主面板对齐 */
  height: calc(100% + 210px);
  margin-top: -210px;
  margin-right: 60px;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.side-panel {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.28);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  flex-shrink: 0;
}

/* 使用提示 */
.tips-panel {
  flex: 1;
  min-height: 520px;

  .panel-body {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .panel-title {
    font-size: 20px;
  }
}

.tip-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 12px 4px;
  flex: 1;
}

.tip-item {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
  padding: 12px 18px;
  border-left: 3px solid var(--line-strong);
  border-radius: 0 8px 8px 0;
  background: var(--surface-2);
  flex: 1;

  &:nth-child(1) { border-left-color: rgba(201, 168, 106, 0.9); }
  &:nth-child(2) { border-left-color: rgba(201, 168, 106, 0.65); }
  &:nth-child(3) { border-left-color: rgba(201, 168, 106, 0.45); }
  &:nth-child(4) { border-left-color: rgba(201, 168, 106, 0.3); }
}

.tip-step {
  font-size: 18px;
  font-weight: 700;
  color: var(--accent-bright);
  letter-spacing: 0.08em;
}

.tip-text {
  font-size: 20px;
  color: var(--text-secondary);
  line-height: 1.7;
}
</style>
