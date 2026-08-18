<template>
  <div class="page-container">
    <div class="page-header">
      <h1>工作台</h1>
      <p>欢迎回来，{{ userStore.userInfo?.realName || userStore.userInfo?.username }}！以下是您的版权数据概览</p>
    </div>

    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="12" :md="6">
        <div class="stat-card" @click="$router.push('/works')">
          <div class="stat-icon icon-1">
            <el-icon :size="24"><Folder /></el-icon>
          </div>
          <div class="stat-value">{{ stats.works }}</div>
          <div class="stat-label">我的作品</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="12" :md="6">
        <div class="stat-card" @click="$router.push('/licenses/my')">
          <div class="stat-icon icon-2">
            <el-icon :size="24"><Key /></el-icon>
          </div>
          <div class="stat-value">{{ stats.licenses }}</div>
          <div class="stat-label">持有授权</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="12" :md="6">
        <div class="stat-card" @click="$router.push('/works/search')">
          <div class="stat-icon icon-3">
            <el-icon :size="24"><Search /></el-icon>
          </div>
          <div class="stat-value">{{ stats.searchCount }}</div>
          <div class="stat-label">作品搜索</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="12" :md="6">
        <div class="stat-card" @click="$router.push('/disputes')">
          <div class="stat-icon icon-4">
            <el-icon :size="24"><Warning /></el-icon>
          </div>
          <div class="stat-value">{{ stats.disputes }}</div>
          <div class="stat-label">争议存证</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 24px;">
      <el-col :xs="24" :md="16">
        <div class="card">
          <div class="card-title">
            <el-icon><Files /></el-icon>
            最近存证的作品
          </div>
          <div v-loading="loadingWorks">
            <el-table :data="recentWorks" style="width: 100%" empty-text="暂无作品">
              <el-table-column prop="title" label="作品名称" min-width="140">
                <template #default="{ row }">
                  <router-link :to="`/works/${row.workID}`" class="work-link">{{ row.title }}</router-link>
                </template>
              </el-table-column>
              <el-table-column prop="artist" label="艺术家" width="120" />
              <el-table-column prop="genre" label="类型" width="80" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="statusTagType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="registerAt" label="存证时间" width="160">
                <template #default="{ row }">
                  {{ formatTime(row.registerAt) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="100" fixed="right">
                <template #default="{ row }">
                  <router-link :to="`/works/${row.workID}`"><el-icon><View /></el-icon></router-link>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :md="8">
        <div class="card" style="margin-bottom: 20px;">
          <div class="card-title">
            <el-icon><DocumentAdd /></el-icon>
            快速操作
          </div>
          <div class="quick-actions">
            <el-button type="primary" size="large" @click="$router.push('/works/register')">
              <el-icon><DocumentAdd /></el-icon>&nbsp;版权存证
            </el-button>
            <el-button size="large" @click="$router.push('/licenses/grant')">
              <el-icon><Tickets /></el-icon>&nbsp;发放授权
            </el-button>
            <el-button size="large" @click="$router.push('/licenses/verify')">
              <el-icon><CircleCheck /></el-icon>&nbsp;授权核验
            </el-button>
          </div>
        </div>

        <div class="card">
          <div class="card-title">
            <el-icon><InfoFilled /></el-icon>
            使用提示
          </div>
          <el-timeline>
            <el-timeline-item type="primary">
              <div class="step-label step-primary">步骤 1</div>
              上传音频文件进行版权存证，获取区块链上唯一存证编号
            </el-timeline-item>
            <el-timeline-item type="success">
              <div class="step-label step-success">步骤 2</div>
              向他人发放授权，设置授权类型、有效期和使用次数
            </el-timeline-item>
            <el-timeline-item type="warning">
              <div class="step-label step-warning">步骤 3</div>
              播放前核验授权有效性，保障版权合规使用
            </el-timeline-item>
            <el-timeline-item type="info">
              <div class="step-label step-info">步骤 4</div>
              随时生成存证证书，作为版权归属的法律凭证
            </el-timeline-item>
          </el-timeline>
        </div>
      </el-col>
    </el-row>
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
    // ignore
  } finally {
    loadingWorks.value = false
  }

  try {
    const res = await api.get('/license/my')
    stats.licenses = (res.data || []).length
  } catch (e) { /* ignore */ }
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
</script>

<style lang="scss" scoped>
.stats-row {
  .stat-card {
    cursor: pointer;
    height: 100%;
  }
}

.stat-icon {
  margin-bottom: 12px;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-1 { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.icon-2 { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.icon-3 { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
.icon-4 { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }

.work-link {
  color: var(--primary);
  font-weight: 500;
  &:hover { text-decoration: underline; }
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;

  .el-button {
    justify-content: flex-start;
    padding-left: 20px;
  }
}

:deep(.el-timeline-item__timestamp) {
  color: var(--text-light);
  font-weight: 600;
}

.step-label {
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 4px;
  line-height: 1.4;
}
.step-primary { color: #409eff; }
.step-success { color: #67c23a; }
.step-warning { color: #e6a23c; }
.step-info    { color: #909399; }
</style>
