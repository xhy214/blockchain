<template>
  <div class="page-container">
    <div class="page-header">
      <h1>版权存证</h1>
      <p>上传您的原创音乐作品，将其 Hash 值存证到区块链上，获得不可篡改的版权证明</p>
    </div>

    <el-row :gutter="24">
      <el-col :xs="24" :lg="16">
        <div class="card">
          <el-form ref="formRef" :model="form" :rules="rules" size="large" label-width="100px">
            <el-form-item label="作品名称" prop="title">
              <el-input v-model="form.title" placeholder="如：稻香" />
            </el-form-item>

            <el-form-item label="艺术家" prop="artist">
              <el-input v-model="form.artist" placeholder="如：周杰伦" />
            </el-form-item>

            <el-form-item label="作品类型" prop="genre">
              <el-select v-model="form.genre" placeholder="选择类型" style="width: 200px;">
                <el-option label="流行" value="流行" />
                <el-option label="古典" value="古典" />
                <el-option label="摇滚" value="摇滚" />
                <el-option label="民谣" value="民谣" />
                <el-option label="电子" value="电子" />
                <el-option label="爵士" value="爵士" />
                <el-option label="R&B" value="R&B" />
                <el-option label="说唱" value="说唱" />
                <el-option label="其他" value="其他" />
              </el-select>
            </el-form-item>

            <el-form-item label="作品描述" prop="description">
              <el-input
                v-model="form.description"
                type="textarea"
                :rows="3"
                placeholder="描述作品的创作背景、灵感来源等"
              />
            </el-form-item>

            <el-form-item label="音频文件" prop="file">
              <el-upload
                :auto-upload="false"
                :on-change="handleFileChange"
                :show-file-list="false"
                :limit="1"
                accept="audio/*"
                class="upload-area"
              >
                <div v-if="!fileInfo" class="upload-placeholder">
                  <el-icon :size="48" class="upload-icon"><UploadFilled /></el-icon>
                  <p class="upload-text">点击选择音频文件</p>
                  <p class="upload-hint">支持 MP3、WAV、FLAC 等格式，文件仅用于计算 Hash 存证，不会上传到服务器</p>
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
                <span>SHA-256 哈希：</span>
                <code>{{ fileHash }}</code>
              </div>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" class="submit-btn" :loading="submitting" @click="handleSubmit">
                提交存证
              </el-button>
              <el-button @click="resetForm">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-col>

      <el-col :xs="24" :lg="8">
        <div class="card tips-card">
          <div class="card-title">
            <el-icon><InfoFilled /></el-icon>
            存证说明
          </div>
          <el-steps direction="vertical" :active="0" class="steps">
            <el-step title="上传文件" description="选择您的原创音频文件" />
            <el-step title="计算 Hash" description="本地计算 SHA-256 哈希，不上传文件" />
            <el-step title="链上存证" description="将哈希写入区块链，生成唯一编号" />
            <el-step title="获取证书" description="可随时下载版权存证证书 PDF" />
          </el-steps>
          <el-alert type="info" :closable="false" show-icon style="margin-top: 16px;">
            <template #title>
              文件始终在本地处理，仅将哈希值上链，您的作品版权得到严格保护。
            </template>
          </el-alert>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '@/api'

const router = useRouter()
const formRef = ref(null)
const submitting = ref(false)
const fileInfo = ref(null)
const fileHash = ref('')

const form = reactive({
  title: '',
  artist: '',
  genre: '',
  description: '',
  file: null
})

const rules = {
  title: [{ required: true, message: '请输入作品名称', trigger: 'blur' }],
  artist: [{ required: true, message: '请输入艺术家', trigger: 'blur' }],
  file: [{ required: true, message: '请上传音频文件', trigger: 'change' }]
}

async function handleFileChange(uploadFile) {
  const raw = uploadFile.raw
  fileInfo.value = { name: raw.name, size: raw.size }
  form.file = raw

  // compute SHA-256 locally
  try {
    const hash = await computeSHA256(raw)
    fileHash.value = hash
  } catch (e) {
    ElMessage.error('哈希计算失败')
  }
}

async function computeSHA256(file) {
  const buf = await file.arrayBuffer()
  const hashBuffer = await crypto.subtle.digest('SHA-256', buf)
  const hashArray = Array.from(new Uint8Array(hashBuffer))
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
}

function clearFile() {
  fileInfo.value = null
  fileHash.value = ''
  form.file = null
}

async function handleSubmit() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  if (!fileHash.value) {
    ElMessage.error('请先上传文件并计算哈希')
    return
  }

  submitting.value = true
  try {
    const fd = new FormData()
    fd.append('title', form.title)
    fd.append('artist', form.artist)
    fd.append('genre', form.genre)
    fd.append('description', form.description)
    fd.append('file', form.file)

    const res = await api.post('/copyright/register', fd, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    ElMessage.success('存证成功！')
    router.push(`/works/${res.data.workID}`)
  } catch (e) {
    // error handled
  } finally {
    submitting.value = false
  }
}

function resetForm() {
  formRef.value?.resetFields()
  clearFile()
}

function formatFileSize(bytes) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}
</script>

<style lang="scss" scoped>
.submit-btn {
  background: var(--bg-gradient) !important;
  border: none !important;
}

.upload-area {
  width: 100%;
}

.upload-placeholder {
  border: 2px dashed var(--border);
  border-radius: 12px;
  padding: 40px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;

  &:hover {
    border-color: var(--primary);
    background: var(--bg-gradient-light);
  }

  .upload-icon {
    color: var(--primary);
    margin-bottom: 12px;
  }

  .upload-text {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 6px;
  }

  .upload-hint {
    font-size: 12px;
    color: var(--text-light);
  }
}

.upload-selected {
  display: flex;
  align-items: center;
  gap: 12px;
  background: var(--bg-gradient-light);
  border-radius: 12px;
  padding: 16px;
  border: 1px solid var(--border);
}

.file-info { flex: 1; min-width: 0; }
.file-name { font-weight: 500; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.file-size { font-size: 12px; color: var(--text-secondary); }

.hash-display {
  margin-top: 12px;
  padding: 10px 12px;
  background: #f1f5f9;
  border-radius: 8px;
  font-size: 12px;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 8px;
  word-break: break-all;

  code {
    font-family: monospace;
    color: var(--primary);
  }
}

.steps {
  :deep(.el-step__title) {
    font-weight: 500;
  }
}
</style>
