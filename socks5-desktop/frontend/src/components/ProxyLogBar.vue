<template>
  <div class="bottom-area">
    <!-- 单一面板，根据 activePanel 切换内容，避免两个 Transition 切换时闪烁 -->
    <Transition name="fade">
      <div v-show="activePanel" class="panel" :style="{ height: panelHeightPx }">
        <div
          class="panel-resize-handle"
          :title="t('log.dragToResize')"
          @mousedown="startResize"
        />
        <template v-if="activePanel === 'status'">
          <div class="panel-header">
            <span>{{ t("log.statusTitle") }}</span>
            <div class="panel-actions">
              <a-button type="text" size="small" @click="clearStatus">
                {{ t("log.clear") }}
              </a-button>
              <a-button type="text" size="small" @click="activePanel = null">
                <CloseOutlined />
              </a-button>
            </div>
          </div>
          <pre class="panel-content">
{{ statusText || t("log.statusEmpty") }}</pre
          >
        </template>
        <template v-else-if="activePanel === 'log'">
          <div class="panel-header">
            <span>{{ t("log.connectionTitle") }}</span>
            <div class="panel-actions">
              <a-button type="text" size="small" @click="clearLogs">
                {{ t("log.clear") }}
              </a-button>
              <a-button type="text" size="small" @click="activePanel = null">
                <CloseOutlined />
              </a-button>
            </div>
          </div>
          <pre class="panel-content">
{{ connectionText || t("log.connectionEmpty") }}</pre
          >
        </template>
      </div>
    </Transition>

    <!-- 底部状态栏：状态日志 | 连接日志 -->
    <div class="status-bar">
      <span
        class="status-item"
        :class="{ active: activePanel === 'status' }"
        @click="togglePanel('status')"
      >
        {{ t("log.statusTitle") }}
      </span>
      <span
        class="status-item"
        :class="{ active: activePanel === 'log' }"
        @click="togglePanel('log')"
      >
        {{ t("log.connectionTitle") }}
        <span v-if="count > 0" class="badge">{{ count }}</span>
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from "vue";
import { CloseOutlined } from "@ant-design/icons-vue";
import { Events } from "@wailsio/runtime";
import { storeToRefs } from "pinia";
import { useSettingsStore } from "@/stores/settings";
import { useI18n } from "@/i18n";

const MAX_LOG_LINES = 1000;

const activePanel = ref(null);
const lines = ref([]);
const statusText = ref("");
const { system } = storeToRefs(useSettingsStore());
const { t } = useI18n();

const PANEL_HEIGHT_MIN = 80;
const PANEL_HEIGHT_MAX = 480;
const panelHeight = ref(160);
const panelHeightPx = computed(() => `${panelHeight.value}px`);

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

let off = null;
let offStatus = null;

function togglePanel(name) {
  activePanel.value = activePanel.value === name ? null : name;
}

function clearStatus() {
  statusText.value = "";
}

function clearLogs() {
  lines.value = [];
}

onMounted(() => {
  offStatus = Events.On("logServerStatus", (ev) => {
    const msg = ev?.data != null ? String(ev.data) : String(ev);
    statusText.value += msg;
  });

  off = Events.On("logProxyList", (ev) => {
    if (!system.value.enableLogRecording) return;
    const msg = ev?.data != null ? String(ev.data) : String(ev);
    lines.value.push(msg);
    if (lines.value.length > MAX_LOG_LINES) {
      lines.value.shift();
    }
  });
});

onBeforeUnmount(() => {
  if (typeof off === "function") {
    off();
  }
  if (typeof offStatus === "function") {
    offStatus();
  }
});

const connectionText = computed(() => lines.value.join("\n"));
const count = computed(() => lines.value.length);
</script>

<style scoped lang="scss">
.bottom-area {
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding-top: 6px;
  border-top: 1px solid #e5e7eb;
  background: #f9fafb;
  position: relative;
  box-sizing: border-box;

  .panel-resize-handle {
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    height: 4px;
    cursor: ns-resize;
    z-index: 1;
    background: transparent;
    transition: background 0.15s;
    &:hover {
      background: rgba(66, 133, 244, 0.2);
    }
  }

  .panel-header {
    flex-shrink: 0;
    height: 32px;
    padding: 0 10px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .panel-actions {
      display: flex;
      align-items: center;
      gap: 4px;
    }
    font-size: 12px;
    font-weight: 500;
    color: #374151;
    border-bottom: 1px solid #e5e7eb;
    background: #fff;
  }

  .panel-content {
    flex: 1;
    margin: 0;
    padding: 8px 10px;
    overflow: auto;
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas,
      "Liberation Mono", "Courier New", monospace;
    font-size: 12px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-all;
    color: #374151;
    background: #f9fafb;
  }
}

.status-bar {
  height: 28px;
  padding: 0 10px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  color: #6b7280;
  background: #fff;
  border-top: 1px solid #e5e7eb;

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

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
