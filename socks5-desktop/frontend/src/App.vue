<template>
  <div class="socks5-layout">

    <!-- 左侧 app -->
    <div class="socks5-app" :style="{ width: appWidth + 'px' }">
      <div class="socks5-app-header">
        <div class="socks5-app-title">Socks5 App</div>
        <a-switch v-model:checked="extensionVisible" @click="hanleToggleExtension" size="small" />
      </div>

      <div class="socks5-app-body">
        <proxy-panel
          :selected-server="selectedServer"
          @tun-change="onTunChange"
          @system-proxy-change="onSystemProxyChange"
        />
        <server-list v-model:selected-server="selectedServer" />
      </div>
    </div>

    <!-- 右侧 extension -->
    <div v-if="extensionVisible" class="socks5-extension">
      <div class="socks5-divider" @mousedown="socks5Dragging = true"></div>
      <div class="socks5-extension-body">
        Extension</div>
    </div>

  </div>
</template>


<script setup>
import { ref, onMounted, watch } from 'vue'
import { WindowSetSize } from '@wails/runtime'


const appWidth = ref(400)
const selectedServer = ref(null)
const appMinWidth = 300
const appMaxWidth = 600




// 调整左右布局宽度拖拽
const socks5Dragging = ref(false)

const handleMouseMove = (e) => {
  if (!socks5Dragging.value) return
  let w = e.clientX
  if (w < appMinWidth) w = appMinWidth
  if (w > appMaxWidth) w = appMaxWidth
  appWidth.value = w
}

const handleMouseUp = () => {
  socks5Dragging.value = false
}

watch(socks5Dragging, (val) => {
  if (val) {
    document.body.style.cursor = 'e-resize'
    document.body.style.userSelect = 'none'
    window.addEventListener('mousemove', handleMouseMove)
    window.addEventListener('mouseup', handleMouseUp)
  } else {
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
    window.removeEventListener('mousemove', handleMouseMove)
    window.removeEventListener('mouseup', handleMouseUp)
  }
})

// 是否展示拓展屏
const extensionVisible = ref(true)

const hanleToggleExtension = () => {
  if (extensionVisible.value && appWidth.value < appMaxWidth) {
    WindowSetSize(appMaxWidth + 400, window.innerHeight)
  }
}

const onTunChange = (enabled) => {
  // TODO: 调用后端启动/关闭 TUN
  console.log('TUN 模式:', enabled)
}

const onSystemProxyChange = (enabled) => {
  // TODO: 调用后端设置/取消系统代理
  console.log('系统代理:', enabled)
}

watch(extensionVisible, (newVal) => {
  if (newVal) {
    appWidth.value = appMinWidth
  } else {
    appWidth.value = Math.max(window.innerWidth, appMinWidth)
  }
})

onMounted(() => {
  window.addEventListener('resize', () => {
    if (window.innerWidth < appMaxWidth) {
      extensionVisible.value = false
    } else {
      extensionVisible.value = true
    }
  })
})

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
      min-height: 0;
      padding: 16px;
      font-size: 14px;
      color: #374151;
      display: flex;
      flex-direction: column;
      gap: 16px;
      overflow: hidden;
    }
  }

  .socks5-extension {
    flex: 1;
    display: flex;
    background: #f9fafb;
    border-left: 1px solid #e5e7eb;

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