<template>
  <div class="socks5-layout">
    <div class="socks5-main">
      <div class="socks5-app" :style="{ width: appWidth + 'px' }">
        <div class="socks5-app-header">
          <div class="socks5-app-title">Socks5 App</div>
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
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import { Window } from "@wailsio/runtime";
import Settings from "./components/Settings.vue";
import ProxyLogBar from "./components/ProxyLogBar.vue";
import { useSettingsStore } from "@/stores/settings";

const settingsStore = useSettingsStore();

const appMinWidth = 400;
const appMaxWidth = 600;
const appWidth = ref(400);

// 判断是否显示扩展窗口
const extensionVisible = ref(true);

// 扩展面板显示状态变化时同步窗口尺寸与主区宽度（immediate: true 保证初始为 true 时也会执行）
watch(
  extensionVisible,
  async (newVal) => {
    const { height } = await Window.Size();
    if (newVal) {
      await Window.SetSize(appMaxWidth + 400, height);
      appWidth.value = appMinWidth;
    } else {
      await Window.SetSize(appMinWidth, height);
      appWidth.value = appMinWidth;
    }
  },
  { immediate: true }
);

// 拖拽窗口条大小事件
const socks5Dragging = ref(false);

const handleMouseMove = (e) => {
  if (!socks5Dragging.value) return;
  let w = e.clientX;
  if (w < appMinWidth) w = appMinWidth;
  if (w > appMaxWidth) w = appMaxWidth;
  appWidth.value = w;
};

const handleMouseUp = () => {
  socks5Dragging.value = false;
};

watch(socks5Dragging, (val) => {
  if (val) {
    document.body.style.cursor = "e-resize";
    document.body.style.userSelect = "none";
    window.addEventListener("mousemove", handleMouseMove);
    window.addEventListener("mouseup", handleMouseUp);
  } else {
    document.body.style.cursor = "";
    document.body.style.userSelect = "";
    window.removeEventListener("mousemove", handleMouseMove);
    window.removeEventListener("mouseup", handleMouseUp);
  }
});

onMounted(() => {
  window.addEventListener("resize", () => {
    if (window.innerWidth < appMaxWidth) {
      extensionVisible.value = false;
      appWidth.value = window.innerWidth;
    } else {
      extensionVisible.value = true;
    }
  });
});
</script>

<style lang="scss" scoped>
.socks5-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f7f9fc;

  .socks5-main {
    flex: 1;
    display: flex;
    min-height: 0;
  }

  .socks5-app {
    display: flex;
    flex-direction: column;
    background: #fff;

    .socks5-app-header {
      height: 48px;
      padding: 0 12px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      border-bottom: 1px solid #e5e7eb;

      .socks5-app-title {
        font-size: 16px;
        font-weight: 600;
      }
    }

    .socks5-app-body {
      flex: 1;
      padding: 10px 12px;
      font-size: 14px;
      color: #374151;
      display: flex;
      flex-direction: column;
      gap: 10px;
      overflow-y: auto;
    }
  }

  .socks5-extension {
    flex: 1;
    display: flex;
    background: #f9fafb;
    border-left: 1px solid #e5e7eb;
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
      color: #374151;
      background: #ffffff;
    }
  }

  .socks5-bottom {
    flex-shrink: 0;
    border-top: 1px solid #e5e7eb;
    background: #ffffff;
  }
}
</style>
