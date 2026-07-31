<template>
  <div class="page-container">
    <div class="page-header">
      <h1>我的授权</h1>
      <p>查看您持有的所有授权记录</p>
    </div>

    <div class="card" v-loading="loading">
      <el-table :data="licenses" style="width: 100%" empty-text="暂无授权记录">
        <el-table-column prop="licenseID" label="授权 ID" min-width="140">
          <template #default="{ row }"><code>{{ row.licenseID }}</code></template>
        </el-table-column>
        <el-table-column prop="workID" label="作品 ID" min-width="140">
          <template #default="{ row }"><code>{{ row.workID }}</code></template>
        </el-table-column>
        <el-table-column prop="licenseType" label="授权类型" width="120">
          <template #default="{ row }">
            <el-tag>{{ typeLabel(row.licenseType) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'ACTIVE' ? 'success' : 'info'" size="small">
              {{ row.status === 'ACTIVE' ? '有效' : '已撤销' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="startDate" label="有效期" width="200">
          <template #default="{ row }">
            {{ formatDate(row.startDate) }} ~ {{ formatDate(row.endDate) }}
          </template>
        </el-table-column>
        <el-table-column prop="usedCount" label="使用次数" width="100">
          <template #default="{ row }">{{ row.usedCount || 0 }} / {{ row.maxUsage || '∞' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'ACTIVE' && row.maxUsage > 0 && row.usedCount < row.maxUsage"
              link type="primary"
              @click="recordUsage(row)"
            >
              记录使用
            </el-button>
            <el-button
              v-if="row.status === 'ACTIVE'"
              link type="danger"
              @click="revoke(row)"
            >
              撤销
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

const loading = ref(false)
const licenses = ref([])

onMounted(fetchList)

async function fetchList() {
  loading.value = true
  try {
    const res = await api.get('/license/my')
    licenses.value = res.data || []
  } catch (e) { licenses.value = [] }
  finally { loading.value = false }
}

async function recordUsage(lic) {
  try {
    await ElMessageBox.confirm(
      `确定为授权 ${lic.licenseID} 记录一次使用吗？`,
      '确认', { type: 'info' }
    )
  } catch { return }

  try {
    await api.post('/license/record-usage', { licenseID: lic.licenseID })
    ElMessage.success('已记录使用')
    fetchList()
  } catch (e) { /* handled */ }
}

async function revoke(lic) {
  try {
    await ElMessageBox.confirm(
      `确定要撤销授权 ${lic.licenseID} 吗？撤销后将无法恢复。`,
      '撤销授权', { type: 'warning' }
    )
  } catch { return }

  try {
    await api.post('/license/revoke', { licenseID: lic.licenseID })
    ElMessage.success('授权已撤销')
    fetchList()
  } catch (e) { /* handled */ }
}

function typeLabel(t) {
  return { COMMERCIAL: '商业', NON_COMMERCIAL: '非商业', EXCLUSIVE: '独家' }[t] || t
}

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('zh-CN')
}
</script>
