<template>
  <transition name="slide">
    <div v-show="show" class="socks5-settings-panel">
      <div class="socks5-settings-title">
        本地代理设置
      </div>

      <a-form layout="vertical">
        <a-form-item label="监听端口">
          <a-input-number
            v-model:value="localConfigProxy.listenPort"
            placeholder="1080"
            :min="1"
            :max="65535"
            style="width: 100%"
          />
        </a-form-item>

        <a-form-item label="用户名">
          <a-input
            v-model:value="localConfigProxy.username"
            placeholder="可选"
          />
        </a-form-item>

        <a-form-item label="密码">
          <a-input-password
            v-model:value="localConfigProxy.password"
            placeholder="可选"
          />
        </a-form-item>
      </a-form>

      <div class="socks5-settings-row">
        <span>启动时自动连接</span>
        <a-switch v-model:checked="localConfigProxy.autoConnect" />
      </div>

      <div class="socks5-settings-tip">
        本地 SOCKS5 地址：127.0.0.1:{{ localConfigProxy.listenPort }}
      </div>
    </div>
  </transition>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: true,
  },
  localConfig: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['update:show', 'update:localConfig'])

const localConfigProxy = computed({
  get() {
    return props.localConfig
  },
  set(val) {
    emit('update:localConfig', val)
  },
})
</script>

<style scoped>
/* ================= settings panel ================= */
.socks5-settings-panel {
  width: 320px;
  background: #ffffff;
  border-left: 1px solid #e5e7eb;
  padding: 16px;
  overflow-y: auto;
}

.socks5-settings-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
}

.socks5-settings-row {
  margin-top: 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
}

.socks5-settings-tip {
  margin-top: 14px;
  font-size: 12px;
  color: #6b7280;
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.2s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(100%);
  opacity: 0;
}
</style>

