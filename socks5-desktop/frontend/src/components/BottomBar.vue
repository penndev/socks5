<template>
  <div class="bottom-area">
    <StatusPanel
      v-if="activePanel === PANEL_NAMES.STATUS"
      :panelHeightPx="panelHeightPx"
      :panelTitle="panelTitle"
      :panelText="panelText"
      :onClear="clearStatus"
      :onClose="closePanel"
      :startResize="startResize"
    />
    <LogPanel
      v-else-if="activePanel === PANEL_NAMES.LOG"
      :panelHeightPx="panelHeightPx"
      :panelTitle="panelTitle"
      :panelText="panelText"
      :onClear="clearLogs"
      :onClose="closePanel"
      :startResize="startResize"
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
        <span v-if="count > 0" class="badge">{{ count }}</span>
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from "vue";
import { Events } from "@wailsio/runtime";
import { storeToRefs } from "pinia";
import { useSettingsStore } from "@/stores/settings";
import { useI18n } from "@/i18n";
import { theme } from "ant-design-vue";
import StatusPanel from "./bottombar/StatusPanel.vue";
import LogPanel from "./bottombar/LogPanel.vue";

// Constants
const MAX_LOG_LINES = 1000;
const PANEL_HEIGHT_MIN = 80;
const PANEL_HEIGHT_MAX = 480;
const PANEL_NAMES = {
  STATUS: "status",
  LOG: "log",
};
const EVENTS = {
  STATUS: "logServerStatus",
  CONNECTION: "logProxyList",
};

// Reactive data
const activePanel = ref(null);
const panelHeight = ref(160);
const lines = ref([]);
const statusText = ref("");

// Composables
const { system } = storeToRefs(useSettingsStore());
const { t } = useI18n();
const { token } = theme.useToken();

// Computed properties
const panelHeightPx = computed(() => `${panelHeight.value}px`);
const panelTitle = computed(() =>
  activePanel.value === PANEL_NAMES.STATUS
    ? t("log.statusTitle")
    : t("log.connectionTitle")
);
const panelText = computed(() => {
  if (activePanel.value === PANEL_NAMES.STATUS) {
    return statusText.value || t("log.statusEmpty");
  }
  return connectionText.value || t("log.connectionEmpty");
});
const connectionText = computed(() => lines.value.join("\n"));
const count = computed(() => lines.value.length);

// Functions
function toEventMessage(eventPayload) {
  return eventPayload?.data != null
    ? String(eventPayload.data)
    : String(eventPayload);
}

function startResize(e) {
  e.preventDefault();
  const startY = e.clientY;
  const startH = panelHeight.value;

  function onMove(ev) {
    const delta = startY - ev.clientY;
    const next = Math.round(startH + delta);
    panelHeight.value = Math.min(PANEL_HEIGHT_MAX, Math.max(PANEL_HEIGHT_MIN, next));
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

function clearStatus() {
  statusText.value = "";
}

function clearLogs() {
  lines.value = [];
}

// Event listeners
let statusEventOff = null;
let connectionEventOff = null;

onMounted(() => {
  statusEventOff = Events.On(EVENTS.STATUS, (eventPayload) => {
    statusText.value += toEventMessage(eventPayload);
  });

  connectionEventOff = Events.On(EVENTS.CONNECTION, (eventPayload) => {
    if (!system.value.enableLogRecording) return;
    lines.value.push(toEventMessage(eventPayload));
    if (lines.value.length > MAX_LOG_LINES) {
      lines.value.shift();
    }
  });
});

onBeforeUnmount(() => {
  if (typeof statusEventOff === "function") {
    statusEventOff();
  }
  if (typeof connectionEventOff === "function") {
    connectionEventOff();
  }
});
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
  color: v-bind('token.colorTextSecondary');
  background: v-bind('token.colorBgContainer');
  border-top: 1px solid v-bind('token.colorBorderSecondary');

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
