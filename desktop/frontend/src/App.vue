<template>
  <a-config-provider :theme="antdThemeConfig">
    <div class="layout">
      <div class="main">
        <div class="app" :style="{ width: appWidth + 'px' }">
          <header class="app-hd">
            <span class="app-title">{{ t("app.title") }}</span>
            <a-switch v-model:checked="extensionVisible" size="small" />
          </header>
          <action-panel />
          <serve-panel />
        </div>
        <setting-panel
          v-if="extensionVisible"
          :handleDividerMove="function(e){appWidth = Math.min(APP_MAX_WIDTH, Math.max(APP_MIN_WIDTH, e.clientX));}"
        />
      </div>
      <div v-if="settingsStore.system.enableLogRecording" class="bottom">
        <bottom-bar />
      </div>
    </div>
  </a-config-provider>
</template>

<script setup>
import {
  ref,
  onMounted,
  watch,
  computed,
} from "vue";
import { Window } from "@wailsio/runtime";
import { theme } from "ant-design-vue";
import { useSettingsStore } from "@/stores/settings";
import { t } from "@/locale";

import ActionPanel from "./components/ActionPanel.vue";
import ServePanel from "./components/ServePanel.vue";
import SettingPanel from "./components/SettingPanel.vue";
import BottomBar from "./components/BottomBar.vue";

const settingsStore = useSettingsStore();

const { token } = theme.useToken();
const colorScheme = window.matchMedia("(prefers-color-scheme: dark)")
const prefersColor = ref(colorScheme.matches);
const antdThemeConfig = computed(() => ({
  algorithm: [
    settingsStore.system.themeMode == "dark" ||
    prefersColor.value ?
    theme.darkAlgorithm :
    theme.defaultAlgorithm,
  ],
  components: { Button: { primaryShadow: "none" } },
}));

const APP_MIN_WIDTH = 400;
const APP_MAX_WIDTH = 600;
const EXTENSION_PANEL_WIDTH = 400;
const appWidth = ref(APP_MIN_WIDTH);

const extensionVisible = ref(true);
watch(extensionVisible, async (isVisible) => {
  const { height } = await Window.Size();
  const targetWindowWidth = isVisible
    ? APP_MAX_WIDTH + EXTENSION_PANEL_WIDTH
    : APP_MIN_WIDTH;
  await Window.SetSize(targetWindowWidth, height);
  appWidth.value = APP_MIN_WIDTH;
});

onMounted(async () => {
  colorScheme.addEventListener("change", (e) => {
    prefersColor.value = e.matches;
  });

  window.addEventListener("resize", () => {
    if (window.innerWidth < APP_MAX_WIDTH) {
      extensionVisible.value = false;
      appWidth.value = window.innerWidth;
      return;
    }
    extensionVisible.value = true;
  });
  await settingsStore.init();
});
</script>

<style lang="scss" scoped>
.layout {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  background: v-bind("token.colorBgLayout");

  .main {
    flex: 1;
    display: flex;
    min-height: 0;
  }

  .app {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 0 10px 8px;
    overflow: hidden;
    font-size: 14px;
    color: v-bind("token.colorText");
    background: v-bind("token.colorBgContainer");

    .app-hd {
      flex-shrink: 0;
      height: 48px;
      margin: 0 -10px;
      padding: 0 12px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      background: v-bind("token.colorBgElevated");
      border-bottom: 1px solid v-bind("token.colorBorderSecondary");

      .app-title {
        font-size: 16px;
        font-weight: 600;
        color: v-bind("token.colorText");
      }
    }
  }

  .bottom {
    flex-shrink: 0;
    border-top: 1px solid v-bind("token.colorBorderSecondary");
    background: v-bind("token.colorBgElevated");
  }
}
</style>
