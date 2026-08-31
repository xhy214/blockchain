<template>
  <div class="page-container">
    <div class="page-inner">
      <div class="page-header-bar">
        <div class="page-header-title">
          <span class="header-icon"><el-icon :size="24"><CircleCheck /></el-icon></span>
          授权核验
        </div>
      </div>
      <div class="page-header-desc">模拟流媒体播放前的版权核验场景，验证用户对特定作品是否持有有效授权</div>

      <div class="panel-row">
        <div class="panel-main">
          <div class="panel">
            <div class="panel-header">
              <div class="panel-title">
                <el-icon :size="22"><Search /></el-icon>
                核验查询
              </div>
            </div>
            <div class="panel-body">
              <el-form ref="formRef" :model="form" :rules="rules" size="large" inline>
                <el-form-item label="作品 ID" prop="workID">
                  <el-input v-model="form.workID" placeholder="输入作品 ID" style="width: 300px;" />
                </el-form-item>
                <el-form-item label="用户 ID" prop="licenseeID">
                  <el-input v-model="form.licenseeID" placeholder="输入被核验用户 ID" style="width: 240px;" />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" class="btn-gradient" :loading="verifying" @click="handleVerify">
                    <el-icon><CircleCheck /></el-icon>&nbsp;核验
                  </el-button>
                </el-form-item>
              </el-form>
            </div>
          </div>

          <div v-if="result" class="panel result-panel">
            <div class="result-header" :class="{ valid: result.valid, invalid: !result.valid }">
              <el-icon :size="40">
                <CircleCheck v-if="result.valid" />
                <CircleClose v-else />
              </el-icon>
              <h3>{{ result.valid ? '授权有效' : '授权无效' }}</h3>
              <p v-if="!result.valid" class="reason">{{ result.reason }}</p>
            </div>

            <template v-if="result.valid && result.license">
              <div class="panel-body">
                <el-descriptions :column="2" border>
                  <el-descriptions-item label="授权类型">{{ result.license.licenseType }}</el-descriptions-item>
                  <el-descriptions-item label="授权 ID">
                    <code>{{ result.license.licenseID }}</code>
                  </el-descriptions-item>
                  <el-descriptions-item label="有效期">{{ formatDate(result.license.startDate) }} ~ {{ formatDate(result.license.endDate) }}</el-descriptions-item>
                  <el-descriptions-item label="使用次数">
                    {{ result.license.usedCount }} / {{ result.license.maxUsage || '不限次' }}
                  </el-descriptions-item>
                </el-descriptions>
                <el-button
                  type="success"
                  style="margin-top: 16px;"
                  :loading="recording"
                  @click="handleRecordUsage"
                >
                  <el-icon><VideoPlay /></el-icon>&nbsp;模拟播放（记录一次使用）
                </el-button>
              </div>
            </template>

            <template v-if="!result.valid && result.license">
              <div class="panel-body">
                <el-descriptions :column="2" border>
                  <el-descriptions-item label="授权状态">
                    {{ result.license.status === 'ACTIVE' ? '有效' : '已撤销/过期' }}
                  </el-descriptions-item>
                  <el-descriptions-item label="有效期">{{ formatDate(result.license.startDate) }} ~ {{ formatDate(result.license.endDate) }}</el-descriptions-item>
                  <el-descriptions-item label="使用次数" :span="2">
                    {{ result.license.usedCount }} / {{ result.license.maxUsage || '不限次' }}
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </template>
          </div>
        </div>

        <div class="panel-side">
          <div class="info-card">
            <div class="info-card-title">
              <span class="icon-badge icon-badge-3"><el-icon :size="24"><InfoFilled /></el-icon></span>
              核验逻辑
            </div>
            <el-steps direction="vertical" :active="0" class="steps">
              <el-step title="查找授权" description="查找该用户对该作品的所有授权" />
              <el-step title="检查状态" description="授权是否为 ACTIVE 状态" />
              <el-step title="检查时效" description="当前时间是否在有效期内" />
              <el-step title="检查次数" description="使用次数是否未超过上限" />
              <el-step title="返回结果" description="全部通过 → 授权有效" />
            </el-steps>
          </div>

          <div class="info-card">
            <div class="info-card-title">
              <span class="icon-badge icon-badge-1"><el-icon :size="24"><Reading /></el-icon></span>
              场景示例
            </div>
            <p class="example-text">
              音乐平台在用户点击播放按钮时，调用此接口核验该用户是否对目标歌曲持有有效授权。这是版权合规播放的核心环节，确保每一次播放都有合法授权依据。
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import api from '@/api'

const formRef = ref(null)
const verifying = ref(false)
const recording = ref(false)
const result = ref(null)

const form = reactive({
  workID: '',
  licenseeID: ''
})

const rules = {
  workID: [{ required: true, message: '请输入作品 ID', trigger: 'blur' }],
  licenseeID: [{ required: true, message: '请输入用户 ID', trigger: 'blur' }]
}

async function handleVerify() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  verifying.value = true
  result.value = null
  try {
    const res = await api.get('/license/verify', {
      params: { workID: form.workID, licenseeID: form.licenseeID }
    })
    result.value = res.data
  } catch (e) { /* handled */ }
  finally { verifying.value = false }
}

async function handleRecordUsage() {
  if (!result.value?.license?.licenseID) return
  recording.value = true
  try {
    await api.post('/license/record-usage', { licenseID: result.value.license.licenseID })
    ElMessage.success('已记录一次使用')
    handleVerify()
  } catch (e) { /* handled */ }
  finally { recording.value = false }
}

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('zh-CN')
}
</script>

<style lang="scss" scoped>
.result-panel {
  gap: 0;
  overflow: hidden;
  padding: 0;
}

.result-header {
  padding: 32px;
  text-align: center;

  &.valid {
    background: linear-gradient(135deg, rgba(16,185,129,0.1) 0%, rgba(16,185,129,0.05) 100%);
    color: #10b981;
  }

  &.invalid {
    background: linear-gradient(135deg, rgba(239,68,68,0.1) 0%, rgba(239,68,68,0.05) 100%);
    color: #ef4444;
  }

  h3 {
    margin: 12px 0 6px;
    font-size: 20px;
  }

  .reason {
    font-size: 14px;
    color: var(--text-secondary);
  }
}

.steps {
  :deep(.el-step__title) {
    font-weight: 500;
  }
}

.example-text {
  line-height: 1.8;
  font-size: 14px;
  color: var(--text-secondary);
}
</style>