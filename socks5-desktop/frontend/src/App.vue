<template>
  <a-config-provider :theme="antdThemeConfig">
    <theme-scope>
      <div class="socks5-layout">
        <div class="socks5-main">
          <div class="socks5-app" :style="{ width: appWidth + 'px' }">
            <div class="socks5-app-header">
              <div class="socks5-app-title">{{ t("app.title") }}</div>
              <a-switch v-model:checked="extensionVisible" size="small" />
            </div>

            <div class="socks5-app-body">
              <proxy-panel />
              <server-list />
            </div>
          </div>

          <div v-if="extensionVisible" class="socks5-extension">
            <div class="socks5-divider" @mousedown="socks5Dragging = true"></div>
            <div class="socks5-extension-body">
              <settings />
            </div>
          </div>
        </div>

        <!-- 底部连接日志状态栏组件（设置开启时显示） -->
        <div v-if="settingsStore.system.enableLogRecording" class="socks5-bottom">
          <proxy-log-bar />
        </div>
      </div>
    </theme-scope>
  </a-config-provider>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch, computed } from "vue";
import { Window } from "@wailsio/runtime";
import { theme as antdTheme } from "ant-design-vue";
import Settings from "./components/Settings.vue";
import ProxyLogBar from "./components/ProxyLogBar.vue";
import ThemeScope from "./components/ThemeScope.vue";
import { useSettingsStore } from "@/stores/settings";
import { useI18n } from "@/i18n";
import { resolvedTheme } from "@/theme";

const settingsStore = useSettingsStore();
const { t } = useI18n();

const antdThemeConfig = computed(() => {
  const isDark = resolvedTheme.value === "dark";
  const baseAlgorithm =
    isDark ? antdTheme.darkAlgorithm : antdTheme.defaultAlgorithm;

  return {
    algorithm: [baseAlgorithm],
    ...(isDark
      ? {
          token: {
            // 仅暗色模式降低饱和度
            colorPrimary: "#6f8fb8",
            colorLink: "#6f8fb8",
          },
        }
      : {}),
    components: {
      Button: {
        primaryShadow: "none",
      },
    },
  };
});

const appMinWidth = 400;
const appMaxWidth = 600;
const appWidth = ref(400);

// 右侧设置面板显示状态
const extensionVisible = ref(true);

async function syncWindowSizeByExtensionVisible(isVisible) {
  const { height } = await Window.Size();
  if (isVisible) {
    await Window.SetSize(appMaxWidth + 400, height);
    appWidth.value = appMinWidth;
    return;
  }
  await Window.SetSize(appMinWidth, height);
  appWidth.value = appMinWidth;
}

watch(extensionVisible, syncWindowSizeByExtensionVisible, { immediate: true });

const socks5Dragging = ref(false);
const updateLayoutByWindowWidth = () => {
  if (window.innerWidth < appMaxWidth) {
    extensionVisible.value = false;
    appWidth.value = window.innerWidth;
    return;
  }
  extensionVisible.value = true;
};

const handleDividerMouseMove = (e) => {
  if (!socks5Dragging.value) return;

  let w = e.clientX;
  if (w < appMinWidth) w = appMinWidth;
  if (w > appMaxWidth) w = appMaxWidth;
  appWidth.value = w;
};

const stopDividerDragging = () => {
  socks5Dragging.value = false;
};

watch(socks5Dragging, (isDragging) => {
  if (isDragging) {
    document.body.style.cursor = "e-resize";
    document.body.style.userSelect = "none";
    window.addEventListener("mousemove", handleDividerMouseMove);
    window.addEventListener("mouseup", stopDividerDragging);
  } else {
    document.body.style.cursor = "";
    document.body.style.userSelect = "";
    window.removeEventListener("mousemove", handleDividerMouseMove);
    window.removeEventListener("mouseup", stopDividerDragging);
  }
});

onMounted(() => {
  window.addEventListener("resize", updateLayoutByWindowWidth);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", updateLayoutByWindowWidth);
  window.removeEventListener("mousemove", handleDividerMouseMove);
  window.removeEventListener("mouseup", stopDividerDragging);
});
</script>

<style lang="scss" scoped>
.socks5-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--socks-bg);

  .socks5-main {
    flex: 1;
    display: flex;
    min-height: 0;
  }

  .socks5-app {
    display: flex;
    flex-direction: column;
    background: var(--socks-card-bg);

    .socks5-app-header {
      height: 48px;
      padding: 0 12px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      background: var(--socks-header-bg);
      border-bottom: 1px solid var(--socks-card-border);

      .socks5-app-title {
        font-size: 16px;
        font-weight: 600;
        color: var(--socks-text-primary);
      }
    }

    .socks5-app-body {
      flex: 1;
      padding: 10px 12px;
      font-size: 14px;
      color: var(--socks-text-primary);
      display: flex;
      flex-direction: column;
      gap: 10px;
      overflow-y: auto;
    }
  }

  .socks5-extension {
    flex: 1;
    display: flex;
    background: var(--socks-muted-bg);
    border-left: 1px solid var(--socks-card-border);
    overflow-y: auto;

    .socks5-divider {
      width: 4px;
      height: 100%;
      cursor: e-resize;
      background: transparent;
      transition: background 0.15s;
      &:hover {
        background: rgba(66, 133, 244, 0.2);
      }
    }

    .socks5-extension-body {
      flex: 1;
      padding: 10px 12px;
      font-size: 14px;
      color: var(--socks-text-primary);
      background: var(--socks-card-bg);
    }
  }

  .socks5-bottom {
    flex-shrink: 0;
    border-top: 1px solid var(--socks-card-border);
    background: var(--socks-bottom-bg);
  }
}

:global(.theme-dark) .socks5-layout .socks5-app .socks5-app-header .socks5-app-title {
  color: #ffffff;
}
</style>
