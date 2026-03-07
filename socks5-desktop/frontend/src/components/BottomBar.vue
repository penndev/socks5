<template>
  <div class="bottom-area">
    <StatusPanel
      v-show="activePanel === PANEL_NAMES.STATUS"
      :panelHeightPx="panelHeightPx"
      :onClose="closePanel"
      :startResize="startResize"
    />
    <LogPanel
      v-show="activePanel === PANEL_NAMES.LOG"
      :panelHeightPx="panelHeightPx"
      :onClose="closePanel"
      :startResize="startResize"
      @update:count="logCount = $event"
    />

    <!-- 底部状态栏：状态日志 | 连接日志 -->
    <div class="status-bar">
      <span
        class="status-item"
        :class="{ active: activePanel === PANEL_NAMES.STATUS }"
        @click="togglePanel(PANEL_NAMES.STATUS)"
      >
        {{ t("log.statusTitle") }}
      </span>
      <span
        class="status-item"
        :class="{ active: activePanel === PANEL_NAMES.LOG }"
        @click="togglePanel(PANEL_NAMES.LOG)"
      >
        {{ t("log.connectionTitle") }}
        <span v-if="logCount > 0" class="badge">{{ logCount }}</span>
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { theme } from "ant-design-vue";
import { t } from "@/i18n";
import StatusPanel from "./bottombar/StatusPanel.vue";
import LogPanel from "./bottombar/LogPanel.vue";

const PANEL_HEIGHT_MIN = 80;
const PANEL_HEIGHT_MAX = 480;
const PANEL_NAMES = {
  STATUS: "status",
  LOG: "log",
};

const activePanel = ref(null);
const panelHeight = ref(160);
const logCount = ref(0);

const { token } = theme.useToken();

const panelHeightPx = computed(() => `${panelHeight.value}px`);

function startResize(e) {
  e.preventDefault();
  const startY = e.clientY;
  const startH = panelHeight.value;

  function onMove(ev) {
    const delta = startY - ev.clientY;
    const next = Math.round(startH + delta);
    panelHeight.value = Math.min(
      PANEL_HEIGHT_MAX,
      Math.max(PANEL_HEIGHT_MIN, next),
    );
  }

  function onUp() {
    document.removeEventListener("mousemove", onMove);
    document.removeEventListener("mouseup", onUp);
    document.body.style.cursor = "";
    document.body.style.userSelect = "";
  }

  document.body.style.cursor = "ns-resize";
  document.body.style.userSelect = "none";
  document.addEventListener("mousemove", onMove);
  document.addEventListener("mouseup", onUp);
}

function togglePanel(name) {
  activePanel.value = activePanel.value === name ? null : name;
}

function closePanel() {
  activePanel.value = null;
}
</script>

<style scoped lang="scss">
.bottom-area {
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.status-bar {
  height: 28px;
  padding: 0 10px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  color: v-bind("token.colorTextSecondary");
  background: v-bind("token.colorBgContainer");
  border-top: 1px solid v-bind("token.colorBorderSecondary");

  .status-item {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    cursor: pointer;
    user-select: none;

    &:hover {
      color: #1677ff;
    }

    &.active {
      color: #1677ff;
      font-weight: 500;
    }

    .badge {
      min-width: 16px;
      padding: 0 4px;
      height: 14px;
      border-radius: 999px;
      background: #e5f0ff;
      color: #1d4ed8;
      font-size: 11px;
      display: inline-flex;
      align-items: center;
      justify-content: center;
    }
  }
}
</style>
