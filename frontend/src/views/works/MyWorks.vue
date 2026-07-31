<template>
  <div class="page-container">
    <div class="page-header">
      <h1>我的作品</h1>
      <p>管理您已存证的音乐作品</p>
    </div>

    <div class="toolbar card">
      <div class="toolbar-left">
        <el-input
          v-model="keyword"
          placeholder="搜索作品名称或艺术家"
          :prefix-icon="Search"
          clearable
          style="width: 300px;"
          @keyup.enter="handleSearch"
        />
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
      <div class="toolbar-right">
        <router-link to="/works/register">
          <el-button type="primary">
            <el-icon><Plus /></el-icon>&nbsp;新建存证
          </el-button>
        </router-link>
      </div>
    </div>

    <div class="card" style="margin-top: 20px;" v-loading="loading">
      <el-table :data="works" style="width: 100%" empty-text="暂无作品，快去存证你的第一个作品吧！">
        <el-table-column prop="title" label="作品名称" min-width="140">
          <template #default="{ row }">
            <router-link :to="`/works/${row.workID}`" class="work-link">{{ row.title }}</router-link>
          </template>
        </el-table-column>
        <el-table-column prop="artist" label="艺术家" width="120" />
        <el-table-column prop="genre" label="类型" width="80">
          <template #default="{ row }">
            <el-tag size="small" v-if="row.genre">{{ row.genre }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="registerAt" label="存证时间" width="160">
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
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'

const router = useRouter()
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

function statusTagType(status) {
  return { ACTIVE: 'success', TRANSFERRED: 'info', DISPUTED: 'warning' }[status] || 'info'
}
function statusLabel(status) {
  return { ACTIVE: '正常', TRANSFERRED: '已转让', DISPUTED: '争议中' }[status] || status
}
function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}
</script>

<style lang="scss" scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;

  .toolbar-left { display: flex; gap: 12px; }
}

.work-link {
  color: var(--primary);
  font-weight: 500;
  &:hover { text-decoration: underline; }
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
  display: flex;
}
</style>
