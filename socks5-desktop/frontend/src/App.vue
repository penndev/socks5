<template>
  <a-config-provider :theme="antdThemeConfig">
    <div class="socks5-layout">
      <div class="socks5-main">
        <div class="socks5-app" :style="{ width: appWidth + 'px' }">
          <div class="socks5-app-header">
            <div class="socks5-app-title">{{ t("app.title") }}</div>
            <a-switch v-model:checked="extensionVisible" size="small" />
          </div>
          <div class="socks5-app-body">
            <action-panel />
            <serve-panel />
          </div>
        </div>
        <setting-panel
          v-if="extensionVisible"
          :handleDividerMove="handleDividerMove"
        />
      </div>
      <!-- 底部连接日志状态栏组件（设置开启时显示） -->
      <div v-if="settingsStore.system.enableLogRecording" class="socks5-bottom">
        <bottom-bar />
      </div>
    </div>
  </a-config-provider>
</template>

<script setup>
import { ref, onMounted, watch, computed, nextTick } from "vue";
import { Window } from "@wailsio/runtime";
import { theme } from "ant-design-vue";
import { useSettingsStore } from "@/stores/settings";
import { setLocale, t } from "@/i18n";

import ActionPanel from "./components/ActionPanel.vue";
import ServePanel from "./components/ServePanel.vue";
import SettingPanel from "./components/SettingPanel.vue";
import BottomBar from "./components/BottomBar.vue";
import { Start } from "@bindings/socks5-desktop/proxy";

const settingsStore = useSettingsStore();

// (async () => {
//   console.log("我准备初始化设置")
// })();

watch(
  () => settingsStore.system.language,
  () => setLocale(settingsStore.system.language)
);

// 将 antd token 映射为布局 CSS 变量
const { token } = theme.useToken();

// 设置系统皮肤模板
const antdThemeConfig = computed(() => {
  let baseAlgorithm = theme.defaultAlgorithm;
  let token = {};
  if (settingsStore.system.themeMode === "dark") {
    baseAlgorithm = theme.darkAlgorithm;
  }
  return {
    algorithm: [baseAlgorithm],
    token,
    components: { Button: { primaryShadow: "none" } },
  };
});

const APP_MIN_WIDTH = 400;
const APP_MAX_WIDTH = 600;
const EXTENSION_PANEL_WIDTH = 400;
const appWidth = ref(APP_MIN_WIDTH);

// 右侧设置面板显示状态
const extensionVisible = ref(true);

watch(extensionVisible, async (isVisible) => {
    const { height } = await Window.Size();
    const targetWindowWidth = isVisible
      ? APP_MAX_WIDTH + EXTENSION_PANEL_WIDTH
      : APP_MIN_WIDTH;
    await Window.SetSize(targetWindowWidth, height);
    appWidth.value = APP_MIN_WIDTH;
  }
);

const handleDividerMove = (e) => {
  appWidth.value = Math.min(APP_MAX_WIDTH, Math.max(APP_MIN_WIDTH, e.clientX));
};

onMounted(async() => {
  window.addEventListener("resize", () => {
    if (window.innerWidth < APP_MAX_WIDTH) {
      extensionVisible.value = false;
      appWidth.value = window.innerWidth;
      return;
    }
    extensionVisible.value = true;
  });
  
  await settingsStore.init();
  await nextTick()
  const { host, port, username, password } = settingsStore.proxy
  console.log(`我准备启动${host}:${port}`, username, password)
  Start(`${host}:${port}`, username, password)
});
</script>

<style lang="scss" scoped>
.socks5-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: v-bind("token.colorBgLayout");

  .socks5-main {
    flex: 1;
    display: flex;
    min-height: 0;
  }

  .socks5-app {
    display: flex;
    flex-direction: column;
    background: v-bind("token.colorBgContainer");

    .socks5-app-header {
      height: 48px;
      padding: 0 12px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      background: v-bind("token.colorBgElevated");
      border-bottom: 1px solid v-bind("token.colorBorderSecondary");

      .socks5-app-title {
        font-size: 16px;
        font-weight: 600;
        color: v-bind("token.colorText");
      }
    }

    .socks5-app-body {
      flex: 1;
      padding: 10px 12px;
      font-size: 14px;
      color: v-bind("token.colorText");
      display: flex;
      flex-direction: column;
      gap: 10px;
      overflow-y: auto;
    }
  }

  .socks5-bottom {
    flex-shrink: 0;
    border-top: 1px solid v-bind("token.colorBorderSecondary");
    background: v-bind("token.colorBgElevated");
  }
}

:global(.theme-dark)
  .socks5-layout
  .socks5-app
  .socks5-app-header
  .socks5-app-title {
  color: #ffffff;
}
</style>
