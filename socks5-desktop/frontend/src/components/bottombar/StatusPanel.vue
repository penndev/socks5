<template>
  <div class="panel" :style="{ height: panelHeightPx }">
    <div
      class="panel-resize-handle"
      :title="t('log.dragToResize')"
      @mousedown="startResize($event)"
    />
    <div class="panel-header">
      <span>{{ t('log.statusTitle') }}</span>
      <div class="panel-actions">
        <a-button type="text" size="small" @click="clearStatus">
          {{ t("log.clear") }}
        </a-button>
        <a-button type="text" size="small" @click="onClose">
          <CloseOutlined />
        </a-button>
      </div>
    </div>
    <pre class="panel-content">{{ displayText }}</pre>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from "vue";
import { Events } from "@wailsio/runtime";
import { CloseOutlined } from "@ant-design/icons-vue";
import { t } from "@/i18n";
import { theme } from "ant-design-vue";

const { token } = theme.useToken();

defineProps({
  panelHeightPx: String,
  onClose: Function,
  startResize: Function,
});

const statusText = ref("");
const displayText = computed(
  () => statusText.value || t("log.statusEmpty"),
);

function clearStatus() {
  statusText.value = "";
}

let statusEventOff = null;
onMounted(() => {
  statusEventOff = Events.On("logServerStatus", (eventPayload) => {
    statusText.value += String(eventPayload.data) + "\n";
  });
});

onBeforeUnmount(() => {
  if (typeof statusEventOff === "function") {
    statusEventOff();
  }
});
</script>

<style scoped lang="scss">
.panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding-top: 6px;
  border-top: 1px solid v-bind("token.colorBorderSecondary");
  background: v-bind("token.colorBgContainer");
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
    color: v-bind("token.colorText");
    border-bottom: 1px solid v-bind("token.colorBorderSecondary");
    background: v-bind("token.colorBgContainer");
  }

  .panel-content {
    flex: 1;
    margin: 0;
    padding: 8px 10px;
    overflow: auto;
    font-family:
      ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
      "Courier New", monospace;
    font-size: 12px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-all;
    color: v-bind("token.colorText");
    background: v-bind("token.colorFillAlter");
  }
}
</style>
