<template>
  <div class="page-container">
    <div class="page-inner">
      <div class="page-header-bar">
        <div class="page-header-title">
          <div class="header-icon"><el-icon :size="24"><Folder /></el-icon></div>
          <div>我的作品</div>
        </div>
        <el-button type="primary" round class="btn-gradient" @click="$router.push('/works/register')" style="margin-left: auto;">
          <el-icon><DocumentAdd /></el-icon>&nbsp;新建存证
        </el-button>
      </div>
      <div class="page-header-desc">管理您已存证的音乐作品</div>

      <div class="panel">
        <div class="panel-body" v-loading="loading">
          <div class="toolbar-row">
            <el-input
              v-model="keyword"
              placeholder="搜索作品名称或艺术家"
              prefix-icon="Search"
              clearable
              style="width: 360px;"
              @keyup.enter="handleSearch"
            />
            <el-button type="primary" class="btn-gradient" @click="handleSearch">搜索</el-button>
          </div>

          <template v-if="works.length === 0 && !loading">
            <div class="empty-state">
              <el-icon :size="48" color="#6E7889"><Folder /></el-icon>
              <p>暂无作品，快去存证你的第一个作品吧！</p>
              <el-button type="primary" class="btn-gradient" @click="$router.push('/works/register')">立即存证</el-button>
            </div>
          </template>

          <template v-else>
            <el-table :data="works" style="width: 100%" empty-text="暂无作品">
              <el-table-column prop="title" label="作品名称" min-width="200">
                <template #default="{ row }">
                  <div class="work-cell" @click="$router.push(`/works/${row.workID}`)">
                    <div class="work-cover-sm" :style="{ background: getCoverGradient(row.genre) }">
                      <el-icon :size="20"><Headset /></el-icon>
                    </div>
                    <div class="work-cell-info">
                      <div class="work-cell-title">{{ row.title }}</div>
                      <div class="work-cell-sub">{{ row.artist }} · {{ row.genre }}</div>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="120">
                <template #default="{ row }">
                  <el-tag :type="statusTagType(row.status)" round>{{ statusLabel(row.status) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="registerAt" label="存证时间" width="200">
                <template #default="{ row }">{{ formatTime(row.registerAt) }}</template>
              </el-table-column>
              <el-table-column label="操作" width="220" fixed="right">
                <template #default="{ row }">
                  <el-button link type="primary" @click="$router.push(`/works/${row.workID}`)">详情</el-button>
                  <el-button link type="primary" @click="$router.push(`/works/${row.workID}/certificate`)">证书</el-button>
                  <el-button link type="primary" @click="$router.push(`/works/${row.workID}/verify-hash`)">验真</el-button>
                </template>
              </el-table-column>
            </el-table>

            <el-pagination
              v-if="total > 0"
              class="pagination"
              v-model:current-page="page"
              v-model:page-size="size"
              :total="total"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="fetchList"
              @current-change="fetchList"
            />
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

const keyword = ref('')
const page = ref(1)
const size = ref(10)
const total = ref(0)
const loading = ref(false)
const works = ref([])

onMounted(() => fetchList())

async function fetchList() {
  loading.value = true
  try {
    const res = await api.get('/copyright/my/list')
    let data = res.data || []
    if (keyword.value) {
      data = data.filter(w =>
        w.title?.toLowerCase().includes(keyword.value.toLowerCase()) ||
        w.artist?.toLowerCase().includes(keyword.value.toLowerCase())
      )
    }
    total.value = data.length
    const start = (page.value - 1) * size.value
    works.value = data.slice(start, start + size.value)
  } catch (e) {
    works.value = []
  } finally {
    loading.value = false
  }
}

function handleSearch() { page.value = 1; fetchList() }
function statusTagType(s) { return { ACTIVE: 'success', TRANSFERRED: 'info', DISPUTED: 'warning' }[s] || 'info' }
function statusLabel(s) { return { ACTIVE: '正常', TRANSFERRED: '已转让', DISPUTED: '争议中' }[s] || s }
function formatTime(t) { if (!t) return '-'; return new Date(t).toLocaleString('zh-CN') }
function getCoverGradient(genre) {
  const m = { '古典': 'linear-gradient(135deg, #C9A86A, #8A6A33)', '流行': 'linear-gradient(135deg, #6E5A7E, #3E3452)', '摇滚': 'linear-gradient(135deg, #8A4A4A, #4E2A2E)', '电子': 'linear-gradient(135deg, #3E6E8A, #24404E)', '民谣': 'linear-gradient(135deg, #5E7E58, #334632)', '爵士': 'linear-gradient(135deg, #8A6A4A, #4A3826)' }
  return m[genre] || 'linear-gradient(135deg, #4A5A78, #28324A)'
}
</script>

<style lang="scss" scoped>
.toolbar-row {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.work-cell {
  display: flex;
  align-items: center;
  gap: 14px;
  cursor: pointer;

  &:hover .work-cell-title {
    color: var(--accent-bright);
  }
}

.work-cover-sm {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.work-cell-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.work-cell-sub {
  font-size: 14px;
  color: var(--text-light);
  margin-top: 2px;
}

.pagination {
  margin-top: 24px;
  justify-content: flex-end;
  display: flex;
}
</style>
