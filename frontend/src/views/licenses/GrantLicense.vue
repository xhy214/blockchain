<template>
  <div class="page-container">
    <div class="page-inner">
      <div class="page-header-bar">
        <div class="page-header-title">
          <span class="header-icon"><el-icon :size="24"><Tickets /></el-icon></span>
          发放授权
        </div>
      </div>
      <div class="page-header-desc">将您的音乐作品授权给其他用户使用，可设置授权类型、有效期和使用次数</div>

      <div class="panel-row">
        <div class="panel-main">
          <div class="panel">
            <div class="panel-header">
              <div class="panel-title">
                <el-icon :size="22"><Key /></el-icon>
                授权信息
              </div>
            </div>
            <div class="panel-body">
              <el-form ref="formRef" :model="form" :rules="rules" size="large" label-width="120px">
                <el-form-item label="选择作品" prop="workID">
                  <el-select v-model="form.workID" placeholder="请选择您的作品" filterable style="width: 100%;">
                    <el-option v-for="w in myWorks" :key="w.workID" :label="`${w.title} - ${w.artist}`" :value="w.workID" />
                  </el-select>
                </el-form-item>

                <el-form-item label="授权类型" prop="licenseType">
                  <el-select v-model="form.licenseType" placeholder="选择授权类型" style="width: 240px;">
                    <el-option label="商业使用" value="COMMERCIAL" />
                    <el-option label="非商业使用" value="NON_COMMERCIAL" />
                    <el-option label="独家授权" value="EXCLUSIVE" />
                  </el-select>
                </el-form-item>

                <el-form-item label="被授权人" prop="licenseeID">
                  <el-input v-model="form.licenseeID" placeholder="被授权人的用户 ID" />
                </el-form-item>

                <el-form-item label="有效期" prop="dateRange">
                  <el-date-picker
                    v-model="form.dateRange"
                    type="daterange"
                    range-separator="至"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    value-format="YYYY-MM-DD"
                    style="width: 100%;"
                  />
                </el-form-item>

                <el-form-item label="使用次数" prop="maxUsage">
                  <el-input-number v-model="form.maxUsage" :min="0" :max="9999" />
                  <span style="margin-left: 12px; color: var(--text-secondary); font-size: 14px;">次（0 表示不限次）</span>
                </el-form-item>

                <el-form-item>
                  <el-button type="primary" class="btn-gradient" :loading="submitting" @click="handleSubmit">
                    <el-icon><Tickets /></el-icon>&nbsp;发放授权
                  </el-button>
                </el-form-item>
              </el-form>
            </div>
          </div>
        </div>

        <div class="panel-side">
          <div class="info-card">
            <div class="info-card-title">
              <span class="icon-badge icon-badge-2"><el-icon :size="24"><InfoFilled /></el-icon></span>
              授权类型说明
            </div>
            <el-descriptions :column="1" border>
              <el-descriptions-item label="商业使用">
                允许用于商业目的，需支付版权费
              </el-descriptions-item>
              <el-descriptions-item label="非商业使用">
                仅限个人学习、研究等非盈利用途
              </el-descriptions-item>
              <el-descriptions-item label="独家授权">
                仅授权给单一主体，不可再授权给第三方
              </el-descriptions-item>
            </el-descriptions>
          </div>

          <div class="info-card">
            <div class="info-card-title">
              <span class="icon-badge icon-badge-4"><el-icon :size="24"><Warning /></el-icon></span>
              注意事项
            </div>
            <ul class="tips-list">
              <li>只有作品的权利人才能发放授权</li>
              <li>授权到期后需重新发放</li>
              <li>独家授权下不可向他人再授权</li>
              <li>使用次数用完后自动失效</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import api from '@/api'

const userStore = useUserStore()
const formRef = ref(null)
const submitting = ref(false)
const myWorks = ref([])

const form = reactive({
  workID: '',
  licenseType: 'COMMERCIAL',
  licenseeID: '',
  dateRange: [],
  maxUsage: 0
})

const rules = {
  workID: [{ required: true, message: '请选择作品', trigger: 'change' }],
  licenseType: [{ required: true, message: '请选择授权类型', trigger: 'change' }],
  licenseeID: [{ required: true, message: '请输入被授权人 ID', trigger: 'blur' }],
  dateRange: [{ required: true, message: '请选择有效期', trigger: 'change' }]
}

onMounted(async () => {
  try {
    const res = await api.get('/copyright/my/list')
    myWorks.value = res.data || []
  } catch (e) { /* ignore */ }
})

async function handleSubmit() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    await api.post('/license/grant', {
      workID: form.workID,
      licenseeID: form.licenseeID,
      licenseType: form.licenseType,
      startDate: form.dateRange[0] + 'T00:00:00Z',
      endDate: form.dateRange[1] + 'T23:59:59Z',
      maxUsage: form.maxUsage
    })
    ElMessage.success('授权发放成功！')
    formRef.value?.resetFields()
    form.maxUsage = 0
  } catch (e) { /* handled */ }
  finally { submitting.value = false }
}
</script>

<style lang="scss" scoped>
.tips-list {
  padding-left: 20px;
  line-height: 1.8;
  font-size: 14px;
  color: var(--text-secondary);
}
</style>