<template>
  <a-card class="proxy-panel" :title="t('proxy.title')">
    <div class="proxy-current-server">
      <span class="proxy-label">{{ t("proxy.currentServerLabel") }}</span>
      <span class="proxy-value">
        {{ serverStore.selectedServer?.remark || serverStore.selectedServer?.host || t("proxy.noSelectedServer") }}
      </span>
      <a-button v-if="serverStore.selectedServer" type="link" size="small" danger
        @click="serverStore.selectedServer = null">
        {{ t("proxy.removeButton") }}
      </a-button>
    </div>

    <div class="proxy-mode-tip" v-if="!serverStore.selectedServer">
      {{ t("proxy.selectTip") }}
    </div>

    <a-radio-group v-else v-model:value="proxyMode" class="proxy-mode-group">
      <a-radio-button value="manual">
        {{ t("proxy.mode.manual") }}
      </a-radio-button>
      <a-radio-button value="system">
        {{ t("proxy.mode.system") }}
      </a-radio-button>
      <a-radio-button value="tun">
        {{ t("proxy.mode.tun") }}
      </a-radio-button>
    </a-radio-group>

    <!-- <div class="proxy-mode-desc">
      <pre>{{ modeMessage }}</pre>
    </div> -->
  </a-card>
</template>

<script setup>
import { onMounted, ref, watch, nextTick } from "vue";
import { useServerStore } from "../stores/server";
import { useSettingsStore } from "@/stores/settings";
import { Start, SetRemote, SetMode } from "@bindings/socks5-desktop/proxy";
import { t } from "@/i18n";
import { theme } from "ant-design-vue";

const { token } = theme.useToken();


// 代理模式，默认为手动
const proxyMode = ref("manual");
watch(proxyMode, async (newMode) => {
  await SetMode(newMode);
});

// 修改远程节点
const serverStore = useServerStore();
watch(serverStore, async () => {
  const { host, username, password, protocol } = serverStore.selectedServer || {};
  await SetRemote(host, username, password, protocol);
});

// 启动代理，golang设置可以启动多次，会自动重启并应用新的设置
const settingsStore = useSettingsStore();


watch(settingsStore.proxy, async () => {
  const { host, port, username, password } = settingsStore.proxy
  console.log(`${host}:${port}`, username, password)
  await Start(`${host}:${port}`, username, password);
});

onMounted(() => {

});

</script>

<style lang="scss" scoped>
.proxy-panel {
  flex-shrink: 0;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);

  .proxy-current-server {
    margin-bottom: 8px;
    padding: 4px 8px;
    background: v-bind("token.colorFillAlter");
    border-radius: 8px;
    display: flex;
    align-items: center;
    flex-wrap: nowrap;
    gap: 6px;

    .proxy-label {
      font-size: 12px;
      color: v-bind("token.colorTextSecondary");
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

  }

  .proxy-mode-tip {
    font-size: 13px;
    color: v-bind("token.colorTextSecondary");
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
    color: v-bind("token.colorTextSecondary");
  }
}
</style>
