<template>
  <div class="page-container">
    <div class="page-inner">
      <div class="page-header-bar">
        <el-button @click="$router.back()" circle plain>
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <div class="page-header-title">
          <div class="header-icon">
            <el-icon :size="24"><Lock /></el-icon>
          </div>
          哈希验真
        </div>
      </div>
      <div class="page-header-desc">上传原始文件，验证其哈希与链上存证是否一致</div>

      <div class="panel-row">
        <div class="panel-main">
          <div class="panel">
            <div class="panel-body">
              <el-form ref="formRef" :model="form" :rules="rules" size="large">
                <el-form-item label="作品 ID" prop="workID">
                  <el-input v-model="form.workID" placeholder="输入或粘贴作品 ID" />
                  <div class="form-tip">提示：可在作品详情页查看作品 ID</div>
                </el-form-item>

                <el-form-item label="原始文件" prop="file">
                  <el-upload
                    :auto-upload="false"
                    :on-change="handleFileChange"
                    :show-file-list="false"
                    accept="audio/*"
                    class="upload-area"
                  >
                    <div v-if="!fileInfo" class="upload-placeholder">
                      <el-icon :size="48" class="upload-icon"><UploadFilled /></el-icon>
                      <p class="upload-text">选择原始音频文件</p>
                      <p class="upload-hint">上传文件用于本地计算哈希，不会发送到服务器</p>
                    </div>
                    <div v-else class="upload-selected">
                      <el-icon :size="32" color="#667eea"><Headset /></el-icon>
                      <div class="file-info">
                        <div class="file-name">{{ fileInfo.name }}</div>
                        <div class="file-size">{{ formatFileSize(fileInfo.size) }}</div>
                      </div>
                      <el-button type="danger" size="small" text @click="clearFile">
                        <el-icon><Delete /></el-icon>
                      </el-button>
                    </div>
                  </el-upload>
                  <div v-if="fileInfo && fileHash" class="hash-display">
                    <el-icon><Lock /></el-icon>
                    <span>本地 SHA-256：</span>
                    <code>{{ fileHash }}</code>
                  </div>
                </el-form-item>

                <el-form-item>
                  <el-button class="btn-gradient" :loading="verifying" @click="handleVerify">
                    <el-icon><Check /></el-icon>&nbsp;开始验真
                  </el-button>
                </el-form-item>
              </el-form>

              <div v-if="result" class="result-section" :class="{ 'result-match': result.match, 'result-mismatch': !result.match }">
                <div class="result-icon">
                  <el-icon :size="48" :color="result.match ? '#10b981' : '#ef4444'">
                    <CircleCheck v-if="result.match" />
                    <CircleClose v-else />
                  </el-icon>
                </div>
                <div class="result-content">
                  <h3>{{ result.match ? '✅ 文件哈希匹配 - 作品真实有效' : '❌ 文件哈希不匹配 - 文件可能被篡改' }}</h3>
                  <div class="result-hashes">
                    <div class="hash-row">
                      <span class="hash-label">链上哈希：</span>
                      <code>{{ result.onChainHash }}</code>
                    </div>
                    <div class="hash-row">
                      <span class="hash-label">本地上传：</span>
                      <code>{{ result.uploadedHash }}</code>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="panel-side">
          <div class="info-card">
            <div class="info-card-title">
              <div class="icon-badge icon-badge-1">
                <el-icon :size="24"><InfoFilled /></el-icon>
              </div>
              验真说明
            </div>
            <ol class="tips-list">
              <li>区块链上仅存储作品的 SHA-256 哈希值，不存储文件本身</li>
              <li>将原始文件上传到此处，系统会在本地重新计算哈希</li>
              <li>若两次哈希完全一致，证明文件未被篡改</li>
              <li>任何微小修改（哪怕是一个字节）都会导致哈希变化</li>
            </ol>
            <el-alert type="warning" :closable="false" show-icon style="margin-top: 16px;">
              <template #title>
                请使用存证时上传的原始文件进行验真，任何编辑、转码都将导致哈希不匹配。
              </template>
            </el-alert>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '@/api'

const route = useRoute()
const formRef = ref(null)
const verifying = ref(false)
const fileInfo = ref(null)
const fileHash = ref('')
const result = ref(null)

const form = reactive({
  workID: route.params.workID || '',
  file: null
})

const rules = {
  workID: [{ required: true, message: '请输入作品 ID', trigger: 'blur' }],
  file: [{ required: true, message: '请选择文件', trigger: 'change' }]
}

async function handleFileChange(uploadFile) {
  const raw = uploadFile.raw
  fileInfo.value = { name: raw.name, size: raw.size }
  form.file = raw
  try {
    const buf = await raw.arrayBuffer()
    const hashBuffer = await crypto.subtle.digest('SHA-256', buf)
    const hashArray = Array.from(new Uint8Array(hashBuffer))
    fileHash.value = hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
  } catch (e) {
    ElMessage.error('哈希计算失败')
  }
}

function clearFile() {
  fileInfo.value = null
  fileHash.value = ''
  form.file = null
  result.value = null
}

async function handleVerify() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  if (!fileHash.value) {
    ElMessage.error('请先选择文件')
    return
  }

  verifying.value = true
  try {
    const fd = new FormData()
    fd.append('workID', form.workID)
    fd.append('file', form.file)

    const res = await api.post('/copyright/verify-hash', fd, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    result.value = res.data
  } catch (e) {
    // error shown by interceptor
  } finally {
    verifying.value = false
  }
}

function formatFileSize(bytes) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}
</script>

<style lang="scss" scoped>
.form-tip { font-size: 12px; color: var(--text-light); margin-top: 4px; }

.upload-area { width: 100%; }

.upload-placeholder {
  border: 2px dashed var(--border);
  border-radius: 20px;
  padding: 48px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;

  &:hover {
    border-color: var(--primary);
    background: linear-gradient(135deg, rgba(102,126,234,0.06) 0%, rgba(118,75,162,0.06) 100%);
  }

  .upload-icon { color: var(--primary); margin-bottom: 12px; }
  .upload-text { font-size: 16px; font-weight: 600; margin-bottom: 6px; }
  .upload-hint { font-size: 12px; color: var(--text-light); }
}

.upload-selected {
  display: flex;
  align-items: center;
  gap: 14px;
  background: linear-gradient(135deg, rgba(102,126,234,0.08) 0%, rgba(118,75,162,0.08) 100%);
  border-radius: 20px;
  padding: 20px;
  border: 1px solid var(--border);
}

.file-info { flex: 1; min-width: 0; }
.file-name { font-weight: 500; }
.file-size { font-size: 12px; color: var(--text-secondary); }

.hash-display {
  margin-top: 14px;
  padding: 12px 16px;
  background: #f1f5f9;
  border-radius: 12px;
  font-size: 12px;
  color: #475569;
  display: flex;
  align-items: center;
  gap: 10px;
  word-break: break-all;

  code { font-family: monospace; color: #667eea; }
}

.result-section {
  margin-top: 24px;
  padding: 24px;
  border-radius: 16px;
  display: flex;
  gap: 16px;
  align-items: flex-start;

  &.result-match {
    background: linear-gradient(135deg, rgba(16,185,129,0.1) 0%, rgba(16,185,129,0.05) 100%);
    border: 1px solid rgba(16,185,129,0.3);
  }

  &.result-mismatch {
    background: linear-gradient(135deg, rgba(239,68,68,0.1) 0%, rgba(239,68,68,0.05) 100%);
    border: 1px solid rgba(239,68,68,0.3);
  }
}

.result-content h3 { font-size: 16px; margin-bottom: 12px; }

.result-hashes {
  .hash-row {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    font-size: 12px;
    margin-bottom: 6px;
  }
  .hash-label { color: var(--text-secondary); flex-shrink: 0; }
  code { font-family: monospace; word-break: break-all; flex: 1; }
}

.tips-list {
  padding-left: 20px;
  line-height: 1.8;
  font-size: 14px;
  color: var(--text-secondary);
}
</style>