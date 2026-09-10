<template>
  <span class="copy-field">
    <el-tooltip :content="value" placement="top" :disabled="!truncated" :show-after="200">
      <code class="copy-field-text">{{ display }}</code>
    </el-tooltip>
    <el-icon v-if="value" class="copy-field-icon" title="复制" @click="copy">
      <CopyDocument />
    </el-icon>
  </span>
</template>

<script setup>
import { computed } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  value: { type: String, default: '' },
  length: { type: Number, default: 18 }
})

const truncated = computed(() => !!props.value && props.value.length > props.length)

const display = computed(() => {
  if (!props.value) return '-'
  return truncated.value ? props.value.slice(0, props.length) + '…' : props.value
})

async function copy() {
  try {
    await navigator.clipboard.writeText(props.value)
    ElMessage.success('已复制')
  } catch { /* ignore */ }
}
</script>

<style lang="scss" scoped>
.copy-field {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  max-width: 100%;
  vertical-align: middle;
}

.copy-field-text {
  font-family: monospace;
  font-size: 14px;
  background: var(--surface-2);
  padding: 4px 8px;
  border-radius: 6px;
  white-space: nowrap;
}

.copy-field-icon {
  flex-shrink: 0;
  cursor: pointer;
  color: var(--text-light);
  transition: color 0.2s;

  &:hover {
    color: var(--accent-bright);
  }
}
</style>
