<template>
  <div class="socks5-layout">
    <div class="socks5-app" :style="{ width: appWidth + 'px' }">
      <div class="socks5-app-header">
        <div class="socks5-app-title">Socks5 App</div>
        <a-switch v-model:checked="extensionVisible" @click="toggleExtension" size="small" />
      </div>

      <div class="socks5-app-body">
        <proxy-panel
          :selected-server="selectedServer"
        />
        <server-list v-model:selected-server="selectedServer" />
      </div>
    </div>

    <div v-if="extensionVisible" class="socks5-extension">
      <div class="socks5-divider" @mousedown="socks5Dragging = true"></div>
      <div class="socks5-extension-body">
        <settings />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import { Window } from "@wailsio/runtime";
import Settings from "./components/Settings.vue";



const selectedServer = ref(null);

const appMinWidth = 400;
const appMaxWidth = 600;


const appWidth = ref(400);


// 判断是否显示扩展窗口
const extensionVisible = ref(true);

// 用户切换扩展窗口显示状态
const toggleExtension = async () => {
  const { height } = await Window.Size();
  if (extensionVisible.value && appWidth.value < appMaxWidth) {
    await Window.SetSize(appMaxWidth + 400, height);
  } else {
    await Window.SetSize(appMinWidth, height);
  }
};

// 
watch(extensionVisible, (newVal) => {
  if (newVal < appMinWidth) {
    appWidth.value = appMinWidth;
  } else {
    appWidth.value = Math.max(window.innerWidth, appMinWidth);
  }
});

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
  height: 100vh;
  background: #f7f9fc;

  .socks5-app {
    display: flex;
    flex-direction: column;
    background: #fff;

    .socks5-app-header {
      height: 56px;
      padding: 0 16px;
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
      padding: 16px;
      font-size: 14px;
      color: #374151;
      display: flex;
      flex-direction: column;
      gap: 16px;
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
      padding: 16px;
      font-size: 14px;
      color: #374151;
      background: #ffffff;
    }
  }
}
</style>
