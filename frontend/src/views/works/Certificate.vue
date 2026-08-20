<template>
  <div class="page-container">
    <div class="page-inner">
      <div class="page-header-bar">
        <el-button @click="$router.back()" :icon="ArrowLeft" circle plain />
        <div class="page-header-title">
          <div class="header-icon">
            <el-icon :size="24"><Medal /></el-icon>
          </div>
          存证证书
        </div>
      </div>
      <div class="page-header-desc">下载区块链版权存证证书 PDF，作为版权归属的法律凭证</div>

      <div class="panel-row">
        <div class="panel-main">
          <div class="panel" v-loading="loading">
            <div class="panel-body">
              <div class="certificate-preview">
                <div class="cert-header">
                  <div class="cert-logo">
                    <el-icon :size="32"><Headset /></el-icon>
                  </div>
                  <h2>音乐数字版权存证证书</h2>
                  <p class="cert-subtitle">本证书由区块链网络自动生成，具有不可篡改性</p>
                </div>

                <el-divider />

                <el-descriptions :column="2" border size="large">
                  <el-descriptions-item label="作品名称" :span="2">
                    <strong>{{ work?.title || '-' }}</strong>
                  </el-descriptions-item>
                  <el-descriptions-item label="艺术家" :span="2">{{ work?.artist || '-' }}</el-descriptions-item>
                  <el-descriptions-item label="版权人">{{ work?.ownerID || '-' }}</el-descriptions-item>
                  <el-descriptions-item label="作品类型">{{ work?.genre || '-' }}</el-descriptions-item>
                  <el-descriptions-item label="存证时间" :span="2">{{ formatTime(work?.registerAt) }}</el-descriptions-item>
                  <el-descriptions-item label="交易 ID" :span="2">
                    <code class="cert-hash">{{ work?.txID || '-' }}</code>
                  </el-descriptions-item>
                  <el-descriptions-item label="文件 SHA-256" :span="2">
                    <code class="cert-hash">{{ work?.fileHash || '-' }}</code>
                  </el-descriptions-item>
                  <el-descriptions-item label="作品 ID" :span="2">
                    <code>{{ work?.workID || '-' }}</code>
                  </el-descriptions-item>
                </el-descriptions>

                <div class="cert-footer">
                  <div class="cert-stamp">
                    <div class="stamp-ring"></div>
                    <span>区块链存证</span>
                  </div>
                  <p class="cert-notice">
                    本证书所记载的版权信息已记录于 Hyperledger Fabric 区块链，<br/>
                    可通过作品 ID 在区块链网络中验证其真实性。
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="panel-side">
          <div class="info-card" style="text-align: center;">
            <el-button class="btn-gradient" size="large" :loading="downloading" @click="handleDownload" style="width: 100%; height: 48px; font-size: 16px;">
              <el-icon :size="20"><Download /></el-icon>&nbsp;下载 PDF 证书
            </el-button>
            <p style="margin-top: 12px; font-size: 13px; color: var(--text-secondary);">
              证书包含完整的链上存证信息，可作为版权证明材料
            </p>
          </div>

          <div class="info-card">
            <div class="info-card-title">
              <div class="icon-badge icon-badge-3">
                <el-icon :size="24"><InfoFilled /></el-icon>
              </div>
              证书说明
            </div>
            <ol class="cert-tips">
              <li>PDF 证书由系统自动生成，包含区块链存证编号</li>
              <li>证书可作为版权归属的初步法律凭证</li>
              <li>任何人可通过作品 ID 在本系统验证证书真伪</li>
              <li>建议妥善保管证书文件，用于侵权维权时的举证</li>
            </ol>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '@/api'

const route = useRoute()
const loading = ref(false)
const downloading = ref(false)
const work = ref(null)

onMounted(async () => {
  loading.value = true
  try {
    const res = await api.get(`/copyright/${route.params.workID}`)
    work.value = res.data
  } catch (e) { /* ignore */ }
  finally { loading.value = false }
})

async function handleDownload() {
  downloading.value = true
  try {
    const res = await fetch(`/api/v1/copyright/${route.params.workID}/certificate`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    if (!res.ok) throw new Error('下载失败')
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `copyright_certificate_${route.params.workID.substring(0, 8)}.pdf`
    a.click()
    URL.revokeObjectURL(url)
    ElMessage.success('证书下载成功')
  } catch (e) {
    ElMessage.error('证书下载失败')
  } finally { downloading.value = false }
}

function formatTime(t) { if (!t) return '-'; return new Date(t).toLocaleString('zh-CN') }
</script>

<style lang="scss" scoped>
.certificate-preview {
  padding: 32px;
  background: linear-gradient(180deg, #fefce8 0%, #fff 30%);
  border-radius: 16px;
}

.cert-header { text-align: center; margin-bottom: 16px; }
.cert-logo {
  width: 64px;
  height: 64px;
  margin: 0 auto 12px;
  border-radius: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.cert-header h2 { font-size: 22px; font-weight: 700; margin-bottom: 4px; color: #92400e; }
.cert-subtitle { color: var(--text-secondary); font-size: 13px; }

.cert-hash {
  font-family: monospace;
  word-break: break-all;
  font-size: 12px;
  color: var(--text-primary);
  background: #f8fafc;
  padding: 4px 8px;
  border-radius: 6px;
}

.cert-footer {
  margin-top: 24px;
  text-align: center;
}

.cert-stamp {
  position: relative;
  width: 100px;
  height: 100px;
  margin: 0 auto 16px;

  .stamp-ring {
    width: 100px;
    height: 100px;
    border: 3px solid #ef4444;
    border-radius: 50%;
    opacity: 0.6;
    position: absolute;
    top: 0;
    left: 0;
  }

  span {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100px;
    width: 100%;
    font-weight: 700;
    color: #ef4444;
    font-size: 13px;
    z-index: 1;
  }
}

.cert-notice {
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.8;
}

.cert-tips { padding-left: 20px; line-height: 1.8; font-size: 14px; color: var(--text-secondary); }
</style>