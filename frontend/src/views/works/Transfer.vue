<template>
  <div class="page-container">
    <div class="page-inner">
      <div class="page-header-bar">
        <el-button @click="$router.back()" circle plain>
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <div class="page-header-title">
          <div class="header-icon">
            <el-icon :size="24"><Share /></el-icon>
          </div>
          版权转让
        </div>
      </div>
      <div class="page-header-desc">将作品版权转让给其他用户</div>

      <div class="panel-row">
        <div class="panel-main">
          <div class="panel">
            <div class="panel-body">
              <el-alert type="warning" :closable="false" show-icon style="margin-bottom: 20px;">
                <template #title>
                  版权转让不可撤销，请谨慎操作。转让后，作品的所有权将转移给新的权利人。
                </template>
              </el-alert>

              <el-form ref="formRef" :model="form" :rules="rules" size="large" label-width="120px">
                <el-form-item label="作品 ID">
                  <el-input v-model="form.workID" disabled />
                </el-form-item>

                <el-form-item label="作品标题">
                  <el-input :model-value="work?.title" disabled />
                </el-form-item>

                <el-form-item label="当前权利人">
                  <el-input :model-value="work?.ownerID" disabled />
                </el-form-item>

                <el-form-item label="新权利人" prop="toID">
                  <el-input v-model="form.toID" placeholder="请输入新权利人的用户 ID" />
                  <div class="form-tip">新权利人必须是已注册的平台用户</div>
                </el-form-item>

                <el-form-item label="转让说明" prop="description">
                  <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请简要说明转让原因（选填）" />
                </el-form-item>

                <el-form-item>
                  <el-button class="btn-gradient" :loading="submitting" @click="handleSubmit">
                    <el-icon><Share /></el-icon>&nbsp;确认转让
                  </el-button>
                </el-form-item>
              </el-form>
            </div>
          </div>
        </div>

        <div class="panel-side">
          <div class="info-card">
            <div class="info-card-title">
              <div class="icon-badge icon-badge-2">
                <el-icon :size="24"><InfoFilled /></el-icon>
              </div>
              转让规则
            </div>
            <ul class="rules-list">
              <li>只有作品的当前权利人才能发起转让</li>
              <li>转让需要提供接收方的有效用户 ID</li>
              <li>转让完成后，作品状态变为"已转让"</li>
              <li>转让操作将被永久记录在区块链上</li>
              <li>转让后原权利人将无法再发放新授权</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

const route = useRoute()
const router = useRouter()
const formRef = ref(null)
const submitting = ref(false)
const work = ref(null)

const form = reactive({
  workID: route.params.workID || '',
  toID: '',
  description: ''
})

const rules = {
  toID: [{ required: true, message: '请输入新权利人 ID', trigger: 'blur' }]
}

onMounted(async () => {
  try {
    const res = await api.get(`/copyright/${route.params.workID}`)
    work.value = res.data
  } catch (e) { /* ignore */ }
})

async function handleSubmit() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    await ElMessageBox.confirm(
      `确定将作品"${work.value?.title || ''}"的版权转让给用户 ${form.toID} 吗？此操作不可撤销。`,
      '最终确认',
      { type: 'warning', confirmButtonText: '确认转让', cancelButtonText: '取消' }
    )
  } catch { return }

  submitting.value = true
  try {
    const res = await api.post('/copyright/transfer', {
      workID: form.workID,
      toID: form.toID
    })
    ElMessage.success('版权转让成功！')
    router.push(`/works/${form.workID}`)
  } catch (e) { /* handled */ }
  finally { submitting.value = false }
}
</script>

<style lang="scss" scoped>
.form-tip { font-size: 12px; color: var(--text-light); margin-top: 4px; }
.rules-list { padding-left: 20px; line-height: 1.8; font-size: 14px; color: var(--text-secondary); }
</style>