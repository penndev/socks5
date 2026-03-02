<template>
  <div class="bottom-area">
    <!-- 面板区域：点击状态栏项时显示，在状态栏上方 -->
    <Transition name="fade">
      <div v-show="activePanel === 'log'" class="panel">
        <div class="panel-header">
          <span>连接日志</span>
          <a-button type="text" size="small" @click="activePanel = null">
            <CloseOutlined />
          </a-button>
        </div>
        <pre class="panel-content">{{ text || "暂无连接日志" }}</pre>
      </div>
    </Transition>

    <!-- 状态栏：固定在最底部 -->
    <div class="status-bar">
      <span
        class="status-item"
        :class="{ active: activePanel === 'log' }"
        @click="togglePanel('log')"
      >
        连接日志
        <span v-if="count > 0" class="badge">{{ count }}</span>
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from "vue";
import { CloseOutlined } from "@ant-design/icons-vue";
import { Events } from "@wailsio/runtime";

const activePanel = ref(null);
const lines = ref([]);

let off = null;

function togglePanel(name) {
  activePanel.value = activePanel.value === name ? null : name;
}

onMounted(() => {
  off = Events.On("logProxyList", (ev) => {
    const msg = ev?.data != null ? String(ev.data) : String(ev);
    lines.value.push(msg);
    if (lines.value.length > 500) {
      lines.value.shift();
    }
  });
});

onBeforeUnmount(() => {
  if (typeof off === "function") {
    off();
  }
});

const text = computed(() => lines.value.join("\n"));
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
  height: 160px;
  border-top: 1px solid #e5e7eb;
  background: #f9fafb;

  .panel-header {
    flex-shrink: 0;
    height: 32px;
    padding: 0 10px;
    display: flex;
    align-items: center;
    justify-content: space-between;
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
