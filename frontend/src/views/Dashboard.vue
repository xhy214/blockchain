<template>
  <div class="dashboard">
    <div class="dashboard-inner">

      <!-- ============ 四个统计卡片（255x170，间距 120） ============ -->
      <div class="stat-cards-row">
        <div class="stat-card" @click="$router.push('/works')">
          <div class="stat-icon-wrap icon-1">
            <el-icon :size="36"><Folder /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.works }}</div>
            <div class="stat-label">我的作品</div>
          </div>
        </div>

        <div class="stat-card" @click="$router.push('/licenses/my')">
          <div class="stat-icon-wrap icon-2">
            <el-icon :size="36"><Key /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.licenses }}</div>
            <div class="stat-label">持有授权</div>
          </div>
        </div>

        <div class="stat-card" @click="$router.push('/works/search')">
          <div class="stat-icon-wrap icon-3">
            <el-icon :size="36"><Search /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.searchCount }}</div>
            <div class="stat-label">作品搜索</div>
          </div>
        </div>

        <div class="stat-card" @click="$router.push('/disputes')">
          <div class="stat-icon-wrap icon-4">
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
              <el-icon :size="24"><Files /></el-icon>
              最近存证的作品
            </div>
            <div class="panel-actions">
              <el-button type="primary" round @click="$router.push('/works/register')">
                <el-icon><DocumentAdd /></el-icon>&nbsp;新建存证
              </el-button>
              <el-button round @click="$router.push('/works')">查看全部</el-button>
            </div>
          </div>

          <div class="panel-body" v-loading="loadingWorks">
            <div v-if="recentWorks.length === 0 && !loadingWorks" class="empty-state">
              <el-icon :size="48" color="#cbd5e1"><Folder /></el-icon>
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

          <!-- 使用提示 500x400 -->
          <div class="side-panel tips-panel">
            <div class="panel-header">
              <div class="panel-title">
                <el-icon :size="22"><InfoFilled /></el-icon>
                使用提示
              </div>
            </div>
            <div class="panel-body">
              <div class="tip-list">
                <div class="tip-item">
                  <div class="tip-step step-1">步骤 1</div>
                  <div class="tip-text">上传音频文件进行版权存证，获取区块链唯一存证编号</div>
                </div>
                <div class="tip-item">
                  <div class="tip-step step-2">步骤 2</div>
                  <div class="tip-text">向他人发放授权，设置授权类型、有效期和使用次数</div>
                </div>
                <div class="tip-item">
                  <div class="tip-step step-3">步骤 3</div>
                  <div class="tip-text">播放前核验授权有效性，保障版权合规使用</div>
                </div>
                <div class="tip-item">
                  <div class="tip-step step-4">步骤 4</div>
                  <div class="tip-text">随时生成存证证书，作为版权归属的法律凭证</div>
                </div>
              </div>
            </div>
          </div>

          <!-- 快速操作 500 宽 -->
          <div class="side-panel quick-panel">
            <div class="panel-header">
              <div class="panel-title">
                <el-icon :size="22"><Promotion /></el-icon>
                快速操作
              </div>
            </div>
            <div class="panel-body">
              <div class="quick-grid">
                <div class="quick-btn" @click="$router.push('/works/register')">
                  <div class="quick-icon q-icon-1">
                    <el-icon :size="28"><DocumentAdd /></el-icon>
                  </div>
                  <div class="quick-label">版权存证</div>
                </div>
                <div class="quick-btn" @click="$router.push('/licenses/grant')">
                  <div class="quick-icon q-icon-2">
                    <el-icon :size="28"><Tickets /></el-icon>
                  </div>
                  <div class="quick-label">发放授权</div>
                </div>
                <div class="quick-btn" @click="$router.push('/licenses/verify')">
                  <div class="quick-icon q-icon-3">
                    <el-icon :size="28"><CircleCheck /></el-icon>
                  </div>
                  <div class="quick-label">授权核验</div>
                </div>
                <div class="quick-btn" @click="$router.push('/works/search')">
                  <div class="quick-icon q-icon-4">
                    <el-icon :size="28"><Search /></el-icon>
                  </div>
                  <div class="quick-label">作品搜索</div>
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
    '古典': 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    '流行': 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    '摇滚': 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
    '电子': 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    '民谣': 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    '爵士': 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)'
  }
  return map[genre] || 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)'
}
</script>

<style lang="scss" scoped>
.dashboard {
  height: calc(100vh - 80px);
  width: 100%;
  background: #f5f7fb;
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
}

/* ============ 统计卡片行 ============ */
.stat-cards-row {
  display: flex;
  gap: 120px;
  margin-bottom: 40px;
  flex-shrink: 0;
}

.stat-card {
  width: 255px;
  height: 170px;
  background: #ffffff;
  border-radius: 20px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  cursor: pointer;
  transition: transform 0.25s ease, box-shadow 0.25s ease;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 32px rgba(102, 126, 234, 0.18);
  }
}

.stat-icon-wrap {
  width: 72px;
  height: 72px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;

  &.icon-1 { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
  &.icon-2 { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
  &.icon-3 { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
  &.icon-4 { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }
}

.stat-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.stat-value {
  font-size: 44px;
  font-weight: 700;
  color: #1e293b;
  line-height: 1;
}

.stat-label {
  font-size: 15px;
  color: #64748b;
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
  background: #ffffff;
  border-radius: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
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
  border-bottom: 1px solid #f1f5f9;
  flex-shrink: 0;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 700;
  color: #1e293b;
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
  color: #94a3b8;

  p {
    font-size: 15px;
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
  background: #f8fafc;
  border-radius: 20px;
  cursor: pointer;
  transition: background 0.2s;

  &:hover {
    background: #f1f5f9;
  }
}

.work-cover {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.work-meta {
  flex: 1;
  min-width: 0;
}

.work-title {
  font-size: 17px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.work-sub {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #94a3b8;

  .dot {
    display: inline-block;
    width: 4px;
    height: 4px;
    border-radius: 50%;
    background: #cbd5e1;
  }
}

.work-status {
  flex-shrink: 0;
}

.work-time {
  font-size: 13px;
  color: #94a3b8;
  flex-shrink: 0;
  min-width: 140px;
  text-align: right;
}

.work-arrow {
  color: #cbd5e1;
  flex-shrink: 0;
}

/* ============ 右侧面板组 ============ */
.right-panels {
  width: 500px;
  height: 100%;
  margin-top: -210px;
  margin-right: 60px;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.side-panel {
  background: #ffffff;
  border-radius: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  flex-shrink: 0;
}

/* 使用提示 */
.tips-panel {
  height: 500px;
  margin-bottom: 30px;

  .panel-body {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .panel-title {
    font-size: 17px;
  }
}

.tip-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding: 12px 4px;
  flex: 1;
}

.tip-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px 18px;
  border-left: 3px solid #e2e8f0;
  border-radius: 0 12px 12px 0;
  background: #f8fafc;
  flex: 1;

  &:nth-child(1) { border-left-color: #667eea; }
  &:nth-child(2) { border-left-color: #43e97b; }
  &:nth-child(3) { border-left-color: #f5576c; }
  &:nth-child(4) { border-left-color: #4facfe; }
}

.tip-step {
  font-size: 13px;
  font-weight: 700;

  &.step-1 { color: #667eea; }
  &.step-2 { color: #43e97b; }
  &.step-3 { color: #f5576c; }
  &.step-4 { color: #4facfe; }
}

.tip-text {
  font-size: 14px;
  color: #475569;
  line-height: 1.5;
}

/* 快速操作 */
.quick-panel {
  width: 500px;
  flex: 1;
  min-height: 0;
  overflow: hidden;

  .panel-body {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .panel-title {
    font-size: 17px;
  }
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  grid-template-rows: repeat(2, 1fr);
  gap: 16px;
  flex: 1;
  min-height: 0;
}

.quick-btn {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 20px;
  background: #f8fafc;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #f1f5f9;
    transform: translateY(-2px);
  }
}

.quick-icon {
  width: 52px;
  height: 52px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;

  &.q-icon-1 { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
  &.q-icon-2 { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }
  &.q-icon-3 { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
  &.q-icon-4 { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
}

.quick-label {
  font-size: 15px;
  font-weight: 600;
  color: #1e293b;
}
</style>
