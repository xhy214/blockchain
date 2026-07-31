<template>
  <div class="page-container">
    <div class="page-header">
      <h1>争议记录</h1>
      <p>查看所有已提交的版权争议存证</p>
    </div>

    <div class="card" v-loading="loading">
      <div v-if="disputes.length" class="dispute-list">
        <el-timeline>
          <el-timeline-item
            v-for="d in disputes"
            :key="d.disputeID"
            :timestamp="formatTime(d.filedAt)"
            :type="d.status === 'PENDING' ? 'warning' : 'success'"
            placement="top"
          >
            <div class="dispute-card">
              <div class="dispute-header">
                <el-tag :type="d.status === 'PENDING' ? 'warning' : 'success'">
                  {{ d.status === 'PENDING' ? '处理中' : '已解决' }}
                </el-tag>
                <span class="dispute-id">
                  <code>{{ d.disputeID }}</code>
                </span>
              </div>
              <el-descriptions :column="2" border size="small" style="margin-top: 12px;">
                <el-descriptions-item label="作品 ID" :span="2">
                  <code>{{ d.workID }}</code>
                </el-descriptions-item>
                <el-descriptions-item label="申请人">{{ d.claimantID }}</el-descriptions-item>
                <el-descriptions-item label="状态">{{ d.status }}</el-descriptions-item>
                <el-descriptions-item label="证据说明" :span="2">
                  <p style="line-height: 1.6;">{{ d.evidence }}</p>
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </el-timeline-item>
        </el-timeline>
      </div>
      <el-empty v-else description="暂无争议记录" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

const loading = ref(false)
const disputes = ref([])

onMounted(async () => {
  loading.value = true
  try {
    // fetch disputes for each of my works
    const worksRes = await api.get('/copyright/my/list')
    const works = worksRes.data || []
    const allDisputes = []
    for (const w of works) {
      try {
        const dRes = await api.get(`/dispute/${w.workID}`)
        if (Array.isArray(dRes.data)) {
          allDisputes.push(...dRes.data)
        }
      } catch (e) { /* skip */ }
    }
    disputes.value = allDisputes.sort((a, b) =>
      new Date(b.filedAt) - new Date(a.filedAt)
    )
  } catch (e) { /* ignore */ }
  finally { loading.value = false }
})

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}
</script>

<style lang="scss" scoped>
.dispute-card {
  background: #fff;
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 16px;
}

.dispute-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.dispute-id code {
  font-size: 12px;
  color: var(--text-secondary);
}
</style>
