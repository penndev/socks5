<template>
  <a-card class="proxy-panel" :title="t('proxy.title')">
    <div class="proxy-current-server">
      <span class="proxy-label">{{ t("proxy.currentServerLabel") }}</span>
      <span
        class="proxy-value"
        :title="selectedServer?.remark || selectedServer?.host || ''"
      >
        {{
          selectedServer?.remark ||
          selectedServer?.host ||
          t("proxy.noSelectedServer")
        }}
      </span>
      <a-button
        v-if="selectedServer"
        type="link"
        size="small"
        danger
        class="remove-btn"
        @click="serverStore.selectedServer = null"
      >
        {{ t("proxy.removeButton") }}
      </a-button>
    </div>

    <div class="proxy-mode-tip" v-if="!selectedServer">
      {{ t("proxy.selectTip") }}
    </div>

    <a-radio-group v-else v-model:value="proxyMode" class="proxy-mode-group">
      <a-radio-button value="manual">
        {{ t("proxy.mode.manual") }}
      </a-radio-button>
      <a-radio-button value="tun">
        {{ t("proxy.mode.tun") }}
      </a-radio-button>
      <a-radio-button value="system">
        {{ t("proxy.mode.system") }}
      </a-radio-button>
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
import { useI18n } from "@/i18n";

import { theme } from "ant-design-vue";
const { token } = theme.useToken();

const serverStore = useServerStore();
const settingsStore = useSettingsStore();
const selectedServer = computed(() => serverStore.selectedServer);
const { t } = useI18n();

// 仅用于 UI 切换展示，保留现有交互行为
const proxyMode = ref("manual");
const modeMessage = ref("");

function toSafeString(value, fallback = "") {
  if (typeof value === "string") return value;
  if (value == null) return fallback;
  return String(value);
}

// 选择节点时，启动或停止本地 socks5
watch(
  selectedServer,
  async (newServer, oldServer) => {
    // 第一次选择节点时，按设置中的本地代理配置启动本地 socks5
    if (!oldServer && newServer) {
      const { host, port, username, password } = settingsStore.proxy;
      await Start(
        `${host || "127.0.0.1"}:${port || 1080}`,
        username || "",
        password || ""
      );
    }

    // 不管是否第一次选择节点，都更新远程节点信息
    if (newServer) {
      const { host, username, password, protocol } = newServer;
      await SetRemote(
        toSafeString(host),
        toSafeString(username),
        toSafeString(password),
        toSafeString(protocol, "Socks5"),
      );
    } else if (oldServer) {
      // 取消选择节点时，停止本地服务
      await Stop();
      modeMessage.value = t("proxy.stopped");
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
    background: v-bind('token.colorFillAlter');
    border-radius: 8px;
    display: flex;
    align-items: center;
    flex-wrap: nowrap;
    gap: 6px;

    .proxy-label {
      font-size: 12px;
      color: v-bind('token.colorTextSecondary');
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
    color: v-bind('token.colorTextSecondary');
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
    color: v-bind('token.colorTextSecondary');
  }
}
</style>
