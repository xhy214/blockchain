<template>
  <div class="page-container">
    <div class="page-header">
      <h1>区块链作品搜索</h1>
      <p>按作品名称或艺术家在区块链上检索音乐作品</p>
    </div>

    <div class="card search-card">
      <el-form :inline="true" @submit.prevent>
        <el-form-item label="关键词">
          <el-input
            v-model="keyword"
            placeholder="输入作品名或艺术家"
            :prefix-icon="Search"
            clearable
            style="width: 400px;"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSearch">
            <el-icon><Search /></el-icon>&nbsp;搜索
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="card" style="margin-top: 20px;" v-loading="loading">
      <div v-if="!loading && total === 0 && searched" class="empty">
        <el-empty description="未找到匹配的作品" />
      </div>

      <el-row v-else :gutter="16">
        <el-col
          v-for="work in works"
          :key="work.workID"
          :xs="24" :sm="12" :md="8" :lg="6"
        >
          <div class="work-card" @click="$router.push(`/works/${work.workID}`)">
            <div class="work-cover">
              <el-icon :size="40" color="#fff"><Headset /></el-icon>
            </div>
            <div class="work-info">
              <h3 class="work-title" :title="work.title">{{ work.title }}</h3>
              <p class="work-artist">{{ work.artist }}</p>
              <div class="work-meta">
                <el-tag size="small" v-if="work.genre">{{ work.genre }}</el-tag>
                <el-tag :type="statusTagType(work.status)" size="small">{{ statusLabel(work.status) }}</el-tag>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>

      <el-pagination
        v-if="total > 0"
        class="pagination"
        v-model:current-page="page"
        v-model:page-size="size"
        :total="total"
        :page-sizes="[12, 24, 48]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSearch"
        @current-change="handleSearch"
      />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '@/api'

const keyword = ref('')
const page = ref(1)
const size = ref(12)
const total = ref(0)
const loading = ref(false)
const searched = ref(false)
const works = ref([])

async function handleSearch() {
  if (!keyword.value.trim()) return
  searched.value = true
  loading.value = true
  try {
    const res = await api.get('/copyright/search', {
      params: { keyword: keyword.value, page: page.value, size: size.value }
    })
    works.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    works.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function statusTagType(s) { return { ACTIVE: 'success', TRANSFERRED: 'info', DISPUTED: 'warning' }[s] || 'info' }
function statusLabel(s) { return { ACTIVE: '正常', TRANSFERRED: '已转让', DISPUTED: '争议中' }[s] || s }
</script>

<style lang="scss" scoped>
.search-card {
  padding: 24px;
}

.work-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--border);
  transition: transform 0.2s, box-shadow 0.2s;
  margin-bottom: 16px;
  cursor: pointer;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 32px rgba(0,0,0,0.1);
  }
}

.work-cover {
  height: 120px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.work-info {
  padding: 12px 14px;
}

.work-title {
  font-size: 15px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.work-artist {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.work-meta {
  display: flex;
  gap: 6px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
  display: flex;
}
</style>
