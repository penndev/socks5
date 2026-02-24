<template>
  <a-card class="proxy-panel" title="代理面板">
    <div class="proxy-current-server">
      <span class="proxy-label">当前节点</span>
      <span class="proxy-value">{{ selectedServer?.host ?? "未选择节点" }}</span>
      <a-button
        v-if="selectedServer"
        type="link"
        size="small"
        danger
        class="remove-btn"
        @click="serverStore.selectedServer = null"
      >
        移除
      </a-button>
    </div>

    <div class="proxy-mode-tip" v-if="!selectedServer">
      请先在下方列表中选择节点
    </div>

    <a-radio-group v-else v-model:value="proxyMode" class="proxy-mode-group">
      <a-radio-button value="manual">手动模式</a-radio-button>
      <a-radio-button value="tun">TUN 模式</a-radio-button>
      <a-radio-button value="system">系统代理</a-radio-button>
    </a-radio-group>

    <div class="proxy-mode-desc" v-if="selectedServer">
      <span>{{ modeMessage }}</span>
    </div>
  </a-card>
</template>

<script setup>
import { ref, watch, computed } from "vue";
import { message } from "ant-design-vue";
import { useServerStore } from "../stores/server";
import { Start, Stop } from "@bindings/socks5-desktop/proxy";

const serverStore = useServerStore();
const selectedServer = computed(() => serverStore.selectedServer);

const proxyMode = ref("manual");
const modeMessage = ref("启动中");

watch(
  selectedServer,
  async (newServer, oldServer) => {
    if (oldServer) {
      await Stop();
    }
    if (newServer) {
      await Start(newServer.host, newServer.username, newServer.password);
    }
  },
  { immediate: true }
);

watch(proxyMode, (newValue, oldValue) => {
  console.log(newValue, oldValue);
});

</script>

<style lang="scss" scoped>
.proxy-panel {
  flex-shrink: 0;
  border-radius: 14px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.04);

  .proxy-current-server {
    margin-bottom: 12px;
    padding: 4px 6px;
    background: #f5f5f5;
    border-radius: 8px;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;

    .proxy-label {
      font-size: 12px;
      color: #6b7280;
      margin-right: 0;
    }

    .proxy-value {
      font-size: 14px;
      font-weight: 500;
      color: #1677ff;
      flex: 1;
    }

    .remove-btn {
      padding: 0 4px;
      font-size: 12px;
    }
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

    :deep(.ant-radio-button-wrapper) {
      flex: 1;
      text-align: center;
    }
  }

  .proxy-mode-desc {
    font-size: 12px;
    color: #6b7280;
  }
}
</style>
