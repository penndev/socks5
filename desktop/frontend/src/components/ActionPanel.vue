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

    <template v-else>
      <a-radio-group v-model:value="proxyMode" class="proxy-mode-group">
        <a-radio-button value="manual">
          {{ t("proxy.mode.manual") }}
        </a-radio-button>
        <!-- <a-radio-button value="system">
        {{ t("proxy.mode.system") }}
      </a-radio-button> -->
        <a-radio-button value="tun">
          {{ t("proxy.mode.tun") }}
        </a-radio-button>
      </a-radio-group>

      <div class="proxy-pac-section">
        <div class="proxy-pac-line">
          <span class="proxy-pac-tag">{{ t("settings.pacTitle") }}</span>
          <span class="proxy-pac-dot">·</span>
          <a
            href="#"
            class="proxy-pac-link"
            :class="{ 'is-disabled': !webBaseURL }"
            @click.prevent="openPacEditor"
            >{{ t("settings.pacOpenEditor") }}</a
          >
          <span class="proxy-pac-dot">·</span>
          <a
            href="#"
            class="proxy-pac-link proxy-pac-js"
            :class="{ 'is-disabled': !pacScriptURL }"
            :title="pacScriptURL || undefined"
            @click.prevent="copyPacScriptURL"
            >pac.js</a
          >
        </div>
      </div>
    </template>

    <!-- <div class="proxy-mode-desc">
      <pre>{{ modeMessage }}</pre>
    </div> -->
  </a-card>
</template>

<script setup>
import { ref, watch, computed } from "vue";
import { useServerStore } from "../stores/server";
import { useSettingsStore } from "@/stores/settings";
import { SetStart, SetStop, SetRemote, SetMode } from "@bindings/desktop/proxy/proxy";
import { t } from "@/i18n";
import { theme, message } from "ant-design-vue";
import { OpenExternalURL } from "@bindings/desktop/internal/appconst";

const { token } = theme.useToken();

const settingsStore = useSettingsStore();

const webBaseURL = computed(() => {
  const rawHost = (settingsStore.proxy.host || "").trim();
  const host =
    rawHost === "0.0.0.0" || rawHost === "" ? "127.0.0.1" : rawHost;
  const port = Number(settingsStore.proxy.port);
  if (!Number.isInteger(port) || port < 1 || port > 65535) {
    return "";
  }
  return `http://${host}:${port}`;
});

const pacScriptURL = computed(() =>
  webBaseURL.value ? `${webBaseURL.value}/pac.js` : "",
);

async function openPacEditor() {
  if (!webBaseURL.value) {
    message.warning(t("settings.pacNeedPort"));
    return;
  }
  try {
    await OpenExternalURL(`${webBaseURL.value}/pac/`);
  } catch (e) {
    message.error(e?.message || t("settings.pacOpenFailed"));
  }
}

async function copyPacScriptURL() {
  if (!pacScriptURL.value) {
    message.warning(t("settings.pacNeedPort"));
    return;
  }
  try {
    await navigator.clipboard.writeText(pacScriptURL.value);
    message.success(t("settings.pacCopied"));
  } catch {
    message.error(t("settings.pacCopyFailed"));
  }
}

async function startProxy() {
  const { host, port, username, password } = settingsStore.proxy;
  if (!host || !port) {
    console.warn("[proxy] skip start: host or port is empty", { host, port });
    return;
  }
  try {
    await SetStart(`${host}:${port}`, username, password);
  } catch (error) {
    console.error("[proxy] SetStart failed", error);
  }
}


// 代理模式，默认为手动
const proxyMode = ref("manual");
watch(proxyMode, async (newMode) => {
  await SetMode(newMode);
});

// 修改远程节点
const serverStore = useServerStore();
watch(serverStore, async () => {
  const { host, username, password, protocol } = serverStore.selectedServer || {};
  if (host && protocol) {
    await SetRemote(`${protocol}://${username}:${password}@${host}`);
    await startProxy();
  }else{
    // await SetStop()
  }
});

// 启动代理，golang设置可以启动多次，会自动重启并应用新的设置
watch(
  () => settingsStore.proxy,
  async () => {
    await startProxy();
  },
  { deep: true, immediate: true },
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
    margin-bottom: 10px;
    display: flex;

    :deep(.ant-radio-button-wrapper) {
      flex: 1;
      text-align: center;
    }
  }

  .proxy-pac-section {
    margin-top: 4px;
    padding-top: 8px;
    border-top: 1px solid v-bind("token.colorBorderSecondary");
  }

  .proxy-pac-line {
    font-size: 12px;
    line-height: 1.5;
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0 2px;
  }

  .proxy-pac-tag {
    color: v-bind("token.colorTextSecondary");
  }

  .proxy-pac-dot {
    color: v-bind("token.colorTextQuaternary");
    user-select: none;
    padding: 0 3px;
  }

  .proxy-pac-link {
    color: v-bind("token.colorText");
    cursor: pointer;
    text-decoration: none;

    &:hover {
      color: #1677ff;
      text-decoration: underline;
    }

    &.is-disabled {
      color: v-bind("token.colorTextDisabled");
      cursor: not-allowed;
      pointer-events: none;
    }
  }

  a.proxy-pac-link.proxy-pac-js {
    color: #1677ff;
    font-weight: 500;

    &:hover {
      text-decoration: underline;
    }

    &.is-disabled {
      color: v-bind("token.colorTextDisabled");
    }
  }

  .proxy-mode-desc {
    font-size: 12px;
    color: v-bind("token.colorTextSecondary");
  }
}
</style>
