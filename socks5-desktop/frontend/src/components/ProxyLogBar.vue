<template>
  <div class="proxy-log-bar" @click="visible = true">
    <span class="bar-title">连接日志</span>
    <span v-if="count > 0" class="bar-count">{{ count }}</span>
    <span class="bar-hint">点击查看</span>

    <a-modal v-model:open="visible" title="连接日志" :footer="null" width="720">
      <div class="proxy-log-modal-body">
        <pre class="proxy-log-text">{{ text || "暂无连接日志" }}</pre>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from "vue";
import { Events } from "@wailsio/runtime";

const visible = ref(false);
const lines = ref([]);

let off = null;

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
.proxy-log-bar {
  height: 32px;
  padding: 0 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #6b7280;
  cursor: pointer;
  user-select: none;

  &:hover {
    background: #f3f4f6;
    color: #374151;
  }

  .bar-title {
    font-weight: 500;
  }

  .bar-count {
    min-width: 20px;
    padding: 0 6px;
    height: 18px;
    border-radius: 999px;
    background: #e5f0ff;
    color: #1d4ed8;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 11px;
  }

  .bar-hint {
    margin-left: auto;
    font-size: 11px;
    color: #9ca3af;
  }
}

.proxy-log-modal-body {
  max-height: 420px;
  overflow: auto;
}

.proxy-log-text {
  margin: 0;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
    "Courier New", monospace;
  font-size: 12px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>

