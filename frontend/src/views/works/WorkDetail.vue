<template>
  <div class="page-container" v-loading="loading">
    <div class="page-inner">
      <div class="page-header-bar">
        <el-button @click="$router.back()" circle plain>
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <div class="page-header-title">
          <div class="header-icon">
            <el-icon :size="24"><Headset /></el-icon>
          </div>
          作品详情
        </div>
      </div>
      <div class="page-header-desc">查看作品存证信息、链上历史、授权记录等</div>

      <template v-if="work">
        <div class="panel-row">
          <div class="panel-main">
            <div class="panel">
              <div class="panel-body">
                <div class="info-header">
                  <div class="work-cover">
                    <el-icon :size="48" color="#221A0D"><Headset /></el-icon>
                  </div>
                  <div class="work-info">
                    <h2>{{ work.title }}</h2>
                    <p class="artist">{{ work.artist }}</p>
                    <div class="meta-tags">
                      <el-tag v-if="work.genre" type="info">{{ work.genre }}</el-tag>
                      <el-tag :type="statusTagType(work.status)">{{ statusLabel(work.status) }}</el-tag>
                    </div>
                  </div>
                  <div class="action-group">
                    <el-button type="primary" @click="$router.push(`/works/${work.workID}/certificate`)">
                      <el-icon><Medal /></el-icon>&nbsp;存证证书
                    </el-button>
                    <el-button @click="$router.push(`/works/${work.workID}/verify-hash`)">
                      <el-icon><Check /></el-icon>&nbsp;哈希验真
                    </el-button>
                    <el-button @click="$router.push(`/works/${work.workID}/transfer`)">
                      <el-icon><Share /></el-icon>&nbsp;版权转让
                    </el-button>
                    <el-button @click="$router.push(`/works/${work.workID}/dispute`)">
                      <el-icon><Warning /></el-icon>&nbsp;争议存证
                    </el-button>
                  </div>
                </div>

                <el-divider />

                <el-descriptions :column="2" border>
                  <el-descriptions-item label="作品 ID">
                    <code>{{ work.workID }}</code>
                  </el-descriptions-item>
                  <el-descriptions-item label="交易 ID">
                    <code v-if="work.txID">{{ work.txID }}</code>
                    <span v-else>-</span>
                  </el-descriptions-item>
                  <el-descriptions-item label="SHA-256 哈希" :span="2">
                    <code class="hash">{{ work.fileHash }}</code>
                    <el-button link type="primary" size="small" @click="copyHash">复制</el-button>
                  </el-descriptions-item>
                  <el-descriptions-item label="存证时间" :span="2">{{ formatTime(work.registerAt) }}</el-descriptions-item>
                  <el-descriptions-item label="作品描述" :span="2">
                    <span v-if="work.description">{{ work.description }}</span>
                    <span v-else class="text-muted">暂无描述</span>
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </div>

            <div class="panel">
              <div class="panel-header">
                <div class="panel-title">
                  <el-icon><Clock /></el-icon>
                  链上历史记录
                </div>
              </div>
              <div class="panel-body">
                <el-timeline v-if="history.length">
                  <el-timeline-item
                    v-for="(record, idx) in history"
                    :key="idx"
                    :timestamp="formatTime(record.timestamp)"
                    :type="record.type === 'RegisterWork' ? 'success' : record.type === 'TransferCopyright' ? 'warning' : 'info'"
                  >
                    <strong>{{ typeLabel(record.type) }}</strong>
                    <pre class="history-data">{{ JSON.stringify(record, null, 2) }}</pre>
                  </el-timeline-item>
                </el-timeline>
                <div v-else class="empty-state">
                  <el-icon :size="48"><Clock /></el-icon>
                  <p>暂无历史记录</p>
                </div>
              </div>
            </div>
          </div>

          <div class="panel-side">
            <div class="info-card">
              <div class="info-card-title">
                <div class="icon-badge icon-badge-1">
                  <el-icon :size="24"><Tickets /></el-icon>
                </div>
                相关授权
              </div>
              <div v-loading="loadingLicenses">
                <div v-if="licenses.length">
                  <div v-for="lic in licenses" :key="lic.licenseID" class="license-item">
                    <div class="license-header">
                      <el-tag :type="lic.status === 'ACTIVE' ? 'success' : 'info'" size="small">
                        {{ lic.status === 'ACTIVE' ? '有效' : '已撤销' }}
                      </el-tag>
                      <span class="license-type">{{ lic.licenseType }}</span>
                    </div>
                    <div class="license-detail">
                      <p><strong>被授权人：</strong>{{ lic.licenseeID }}</p>
                      <p><strong>有效期：</strong>{{ lic.startDate }} ~ {{ lic.endDate }}</p>
                      <p><strong>使用：</strong>{{ lic.usedCount }} / {{ lic.maxUsage || '不限次' }}</p>
                    </div>
                  </div>
                </div>
                <div v-else class="empty-state" style="padding: 40px 0;">
                  <el-icon :size="36"><Tickets /></el-icon>
                  <p>暂无授权记录</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <div v-else-if="!loading" class="empty-state">
        <el-icon :size="48"><Search /></el-icon>
        <p>未找到该作品</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '@/api'

const route = useRoute()
const loading = ref(false)
const loadingLicenses = ref(false)
const work = ref(null)
const history = ref([])
const licenses = ref([])

onMounted(() => {
  fetchWork()
  fetchHistory()
  fetchLicenses()
})

async function fetchWork() {
  loading.value = true
  try {
    const res = await api.get(`/copyright/${route.params.workID}`)
    work.value = res.data
  } catch (e) {
    work.value = null
  } finally {
    loading.value = false
  }
}

async function fetchHistory() {
  try {
    const res = await api.get(`/copyright/${route.params.workID}/history`)
    history.value = res.data || []
  } catch (e) { /* ignore */ }
}

async function fetchLicenses() {
  loadingLicenses.value = true
  try {
    const all = await api.get('/license/my')
    licenses.value = (all.data || []).filter(l => l.workID === route.params.workID)
  } catch (e) { /* ignore */ }
  loadingLicenses.value = false
}

function statusTagType(s) { return { ACTIVE: 'success', TRANSFERRED: 'info', DISPUTED: 'warning' }[s] || 'info' }
function statusLabel(s) { return { ACTIVE: '正常', TRANSFERRED: '已转让', DISPUTED: '争议中' }[s] || s }
function typeLabel(t) {
  const m = {
    RegisterWork: '版权存证',
    TransferCopyright: '版权转让',
    GrantLicense: '发放授权',
    RevokeLicense: '撤销授权',
    RecordUsage: '记录使用',
    FileDispute: '争议存证'
  }
  return m[t] || t
}
function formatTime(t) { if (!t) return '-'; return new Date(t).toLocaleString('zh-CN') }
function copyHash() {
  if (!work.value) return
  navigator.clipboard.writeText(work.value.fileHash).then(() => {
    ElMessage.success('哈希已复制到剪贴板')
  })
}
</script>

<style lang="scss" scoped>
.info-header {
  display: flex;
  align-items: flex-start;
  gap: 20px;
}

.work-cover {
  width: 100px;
  height: 100px;
  border-radius: 16px;
  background: var(--bg-gradient);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.work-info { flex: 1; }
.work-info h2 { font-size: 22px; font-weight: 700; margin-bottom: 6px; }
.work-info .artist { color: var(--text-secondary); margin-bottom: 10px; }
.meta-tags { display: flex; gap: 8px; }

.action-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex-shrink: 0;
}

.hash {
  font-family: monospace;
  word-break: break-all;
  font-size: 13px;
  background: var(--surface-2);
  padding: 6px 10px;
  border-radius: 8px;
  display: inline-block;
}

.history-data {
  background: var(--surface-2);
  padding: 8px 12px;
  border-radius: 10px;
  font-size: 12px;
  margin-top: 8px;
  overflow-x: auto;
}

.text-muted { color: var(--text-light); }

.license-item {
  padding: 14px;
  border: 1px solid var(--border);
  border-radius: 12px;
  margin-bottom: 12px;
  background: var(--surface-2);
}

.license-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.license-type { font-weight: 600; }
.license-detail p { font-size: 13px; margin: 4px 0; color: var(--text-secondary); }
</style>