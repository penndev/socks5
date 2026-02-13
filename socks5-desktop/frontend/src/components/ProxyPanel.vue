<template>
  <a-card class="proxy-panel" title="代理面板">
    <div class="proxy-current-server">
      <span class="proxy-label">当前节点</span>
      <span class="proxy-value">{{ serverLabel }}</span>
    </div>

    <div class="proxy-mode-tip" v-if="!selectedServer">
      请先在下方列表中选择节点
    </div>

    <a-radio-group
      v-else
      v-model:value="proxyMode"
      class="proxy-mode-group"
      @change="onModeChange"
    >
      <a-radio-button value="off">关闭</a-radio-button>
      <a-radio-button value="tun">TUN 模式</a-radio-button>
      <a-radio-button value="system">系统代理</a-radio-button>
    </a-radio-group>
    <div class="proxy-mode-desc" v-if="selectedServer">
      <span v-if="proxyMode === 'tun'">虚拟网卡，全局流量</span>
      <span v-else-if="proxyMode === 'system'">设置系统代理</span>
      <span v-else>未启用代理</span>
    </div>
  </a-card>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  selectedServer: { type: Object, default: null },
})

const emit = defineEmits(['tun-change', 'system-proxy-change'])

const proxyMode = ref('off') // off | tun | system

const serverLabel = computed(() => {
  if (!props.selectedServer) return '未选择节点'
  return props.selectedServer.host
})

// 切换节点时重置模式
watch(() => props.selectedServer, () => {
  proxyMode.value = 'off'
  emit('tun-change', false)
  emit('system-proxy-change', false)
})

const onModeChange = () => {
  emit('tun-change', proxyMode.value === 'tun')
  emit('system-proxy-change', proxyMode.value === 'system')
}
</script>

<style scoped>
.proxy-panel {
  flex-shrink: 0;
  border-radius: 14px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.04);
}

.proxy-current-server {
  margin-bottom: 12px;
  padding: 8px 12px;
  background: #f5f5f5;
  border-radius: 8px;
}

.proxy-label {
  font-size: 12px;
  color: #6b7280;
  margin-right: 8px;
}

.proxy-value {
  font-size: 14px;
  font-weight: 500;
  color: #1677ff;
}

.proxy-mode-tip {
  font-size: 13px;
  color: #8c8c8c;
  padding: 12px 0;
}

.proxy-mode-group {
  width: 100%;
  margin-bottom: 8px;
  display: flex;
}

.proxy-mode-group :deep(.ant-radio-button-wrapper) {
  flex: 1;
  text-align: center;
}

.proxy-mode-desc {
  font-size: 12px;
  color: #6b7280;
}
</style>
