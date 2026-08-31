<template>
  <div class="page-container">
    <div class="page-inner">
      <div class="page-header-bar">
        <el-button @click="$router.back()" circle plain>
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <div class="page-header-title">
          <div class="header-icon">
            <el-icon :size="24"><Warning /></el-icon>
          </div>
          版权争议存证
        </div>
      </div>
      <div class="page-header-desc">当发现版权被侵犯时，可提交争议存证作为维权依据</div>

      <div class="panel-row">
        <div class="panel-main">
          <div class="panel">
            <div class="panel-body">
              <el-form ref="formRef" :model="form" :rules="rules" size="large" label-width="120px">
                <el-form-item label="作品 ID">
                  <el-input v-model="form.workID" disabled />
                </el-form-item>

                <el-form-item label="作品标题">
                  <el-input :model-value="work?.title" disabled />
                </el-form-item>

                <el-form-item label="争议类型" prop="disputeType">
                  <el-select v-model="form.disputeType" placeholder="选择争议类型" style="width: 100%;">
                    <el-option label="未经授权使用" value="UNAUTHORIZED_USE" />
                    <el-option label="署名争议" value="ATTRIBUTION_DISPUTE" />
                    <el-option label="抄袭/侵权" value="PLAGIARISM" />
                    <el-option label="其他" value="OTHER" />
                  </el-select>
                </el-form-item>

                <el-form-item label="侵权方" prop="infringerID">
                  <el-input v-model="form.infringerID" placeholder="侵权方的用户 ID 或名称" />
                </el-form-item>

                <el-form-item label="证据说明" prop="evidence">
                  <el-input
                    v-model="form.evidence"
                    type="textarea"
                    :rows="5"
                    placeholder="详细描述侵权情况、时间、证据链接等"
                  />
                  <div class="form-tip">建议附上截图、链接等可验证的证据信息</div>
                </el-form-item>

                <el-form-item>
                  <el-button class="btn-gradient" :loading="submitting" @click="handleSubmit">
                    <el-icon><Warning /></el-icon>&nbsp;提交争议存证
                  </el-button>
                </el-form-item>
              </el-form>
            </div>
          </div>
        </div>

        <div class="panel-side">
          <div class="info-card">
            <div class="info-card-title">
              <div class="icon-badge icon-badge-4">
                <el-icon :size="24"><InfoFilled /></el-icon>
              </div>
              争议存证说明
            </div>
            <el-alert type="info" :closable="false" show-icon style="margin-bottom: 16px;">
              <template #title>
                争议存证将永久记录在区块链上，作为后续维权的重要凭据。
              </template>
            </el-alert>
            <ol class="tips-list">
              <li>争议记录包含存证时间、证据说明等关键信息</li>
              <li>司法机关可通过区块链核验争议记录真实性</li>
              <li>建议同时保留原始证据（截图、文件等）</li>
              <li>提交后可在"争议记录"页面查看状态</li>
            </ol>
          </div>

          <div class="info-card" v-if="disputes.length">
            <div class="info-card-title">
              <div class="icon-badge icon-badge-2">
                <el-icon :size="24"><Document /></el-icon>
              </div>
              该作品的争议记录
            </div>
            <div v-for="d in disputes" :key="d.disputeID" class="dispute-item">
              <div class="dispute-header">
                <el-tag :type="d.status === 'PENDING' ? 'warning' : 'success'" size="small">
                  {{ d.status === 'PENDING' ? '处理中' : '已解决' }}
                </el-tag>
                <span class="dispute-time">{{ formatTime(d.filedAt) }}</span>
              </div>
              <p class="dispute-evidence">{{ d.evidence }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '@/api'

const route = useRoute()
const formRef = ref(null)
const submitting = ref(false)
const work = ref(null)
const disputes = ref([])

const form = reactive({
  workID: route.params.workID || '',
  disputeType: '',
  infringerID: '',
  evidence: ''
})

const rules = {
  disputeType: [{ required: true, message: '请选择争议类型', trigger: 'change' }],
  evidence: [{ required: true, message: '请填写证据说明', trigger: 'blur' }]
}

onMounted(async () => {
  try {
    const res = await api.get(`/copyright/${route.params.workID}`)
    work.value = res.data
  } catch (e) { /* ignore */ }

  try {
    const res = await api.get(`/dispute/${route.params.workID}`)
    disputes.value = Array.isArray(res.data) ? res.data : []
  } catch (e) { /* ignore */ }
})

async function handleSubmit() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    await api.post('/dispute/file', {
      workID: form.workID,
      evidence: `[${form.disputeType}] 侵权方: ${form.infringerID || '未知'} | ${form.evidence}`
    })
    ElMessage.success('争议存证提交成功')
    const res = await api.get(`/dispute/${route.params.workID}`)
    disputes.value = Array.isArray(res.data) ? res.data : []
    form.evidence = ''
    form.infringerID = ''
  } catch (e) { /* handled */ }
  finally { submitting.value = false }
}

function formatTime(t) { if (!t) return '-'; return new Date(t).toLocaleString('zh-CN') }
</script>

<style lang="scss" scoped>
.form-tip { font-size: 12px; color: var(--text-light); margin-top: 4px; }
.tips-list { padding-left: 20px; line-height: 1.8; font-size: 14px; color: var(--text-secondary); }

.dispute-item {
  padding: 14px;
  border: 1px solid var(--border);
  border-radius: 12px;
  margin-bottom: 12px;
  background: var(--surface-2);
}

.dispute-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.dispute-time { font-size: 12px; color: var(--text-light); }
.dispute-evidence { font-size: 13px; color: var(--text-primary); line-height: 1.5; }
</style>