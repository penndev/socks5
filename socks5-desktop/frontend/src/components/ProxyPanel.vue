<template>
  <a-card class="proxy-panel" title="代理面板">
    <div class="proxy-current-server">
      <span class="proxy-label">当前节点</span>
      <span class="proxy-value">{{ selectedServer?.host ?? "未选择节点" }}</span>
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

    <div class="proxy-mode-desc" v-if="selectedServer">
      <pre>{{ modeMessage }}</pre>
    </div>
  </a-card>
</template>

<script setup>
import { ref, watch, computed } from "vue";
import { message } from "ant-design-vue";
import { useServerStore } from "../stores/server";
import { Start, Stop, SetRemote } from "@bindings/socks5-desktop/proxy";
import { Get } from "@bindings/socks5-desktop/storage";
import { Events } from "@wailsio/runtime";

const serverStore = useServerStore();
const selectedServer = computed(() => serverStore.selectedServer);

const proxyMode = ref("manual");
const modeMessage = ref("");

// 后端通过事件推送一些状态时，附加到文案中（如当前使用的本地地址等）
Events.On("logServerStatus", (ev) => {
  const msg = ev?.data != null ? ev.data : ev;
  modeMessage.value +=  msg;
});

const SETTINGS_KEY = "settings";


watch(
  selectedServer,
  async (newServer, oldServer) => {
    console.log(newServer, oldServer);

    // 第一次选择节点时，按设置中的本地代理配置启动本地 socks5
    if (!oldServer && newServer) {
      const data = await Get(SETTINGS_KEY);
      const proxy = data?.proxy || {};
      const host = proxy.host || "127.0.0.1";
      const port = proxy.port || 1080;
      const username = proxy.username || "";
      const password = proxy.password || "";
      await Start(`${host}:${port}`, username, password);
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
