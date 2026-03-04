<template>
  <a-card class="proxy-panel" title="代理面板">
    <div class="proxy-current-server">
      <span class="proxy-label">当前节点</span>
      <span class="proxy-value" :title="selectedServer?.remark || selectedServer?.host || ''">
        {{ selectedServer?.remark || selectedServer?.host || "未选择节点" }}
      </span>
      <a-button v-if="selectedServer" type="link" size="small" danger class="remove-btn"
        @click="serverStore.selectedServer = null">
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

    <div class="proxy-mode-desc">
      <pre>{{ modeMessage }}</pre>
    </div>
  </a-card>
</template>

<script setup>
import { ref, watch, computed } from "vue";
import { useServerStore } from "../stores/server";
import { useSettingsStore } from "@/stores/settings";
import { Start, Stop, SetRemote } from "@bindings/socks5-desktop/proxy";
import { Events } from "@wailsio/runtime";

const serverStore = useServerStore();
const settingsStore = useSettingsStore();
const selectedServer = computed(() => serverStore.selectedServer);

const proxyMode = ref("manual");
const modeMessage = ref("");

// 选择节点时，启动或停止本地 socks5
watch(selectedServer, async (newServer, oldServer) => {
    // 第一次选择节点时，按设置中的本地代理配置启动本地 socks5
    if (!oldServer && newServer) {
      const { host, port, username, password } = settingsStore.proxy;
      await Start(`${host || "127.0.0.1"}:${port || 1080}`, username || "", password || "");
    }

    // 不管是否第一次选择节点，都更新远程节点信息
    if (newServer) {
      const { host, username, password, protocol } = newServer;
      await SetRemote(host, username ?? "", password ?? "", protocol ?? "Socks5");
    } else if (oldServer) {
      // 取消选择节点时，停止本地服务
      await Stop();
      modeMessage.value = "已停止";
    }
  },
  { immediate: true }
);

</script>

<style lang="scss" scoped>
.proxy-panel {
  flex-shrink: 0;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);

  .proxy-current-server {
    margin-bottom: 8px;
    padding: 4px 8px;
    background: #f5f5f5;
    border-radius: 8px;
    display: flex;
    align-items: center;
    flex-wrap: nowrap;
    gap: 6px;

    .proxy-label {
      font-size: 12px;
      color: #6b7280;
      margin-right: 0;
    }

    .proxy-value {
      font-size: 13px;
      font-weight: 500;
      color: #1677ff;
      flex: 1;
      min-width: 0;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .remove-btn {
      padding: 0 4px;
      font-size: 12px;
    }
  }

  .proxy-mode-tip {
    font-size: 13px;
    color: #8c8c8c;
    padding: 6px 0;
  }

  .proxy-mode-group {
    width: 100%;
    margin-bottom: 6px;
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
