<template>
  <div class="socks5-extension">
    <div class="socks5-divider" @mousedown="isDividerDragging = true" />
    <div class="socks5-extension-body">
      <div class="settings-panel">
        <a-form layout="vertical">
          <div class="section-title">{{ t("settings.localProxy") }}</div>
          <a-form-item :label="t('settings.ipAddress')">
            <a-select
              v-model:value="settingsStore.proxy.host"
              :placeholder="t('settings.ipPlaceholder')"
              style="width: 100%"
            >
              <a-select-option value="127.0.0.1">127.0.0.1</a-select-option>
              <a-select-option value="0.0.0.0">0.0.0.0</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label="t('settings.port')">
            <a-input-number
              v-model:value="settingsStore.proxy.port"
              :min="1"
              :max="65535"
              :placeholder="t('settings.portPlaceholder')"
              style="width: 100%"
            />
          </a-form-item>
          <a-form-item :label="t('settings.username')">
            <a-input
              v-model:value="settingsStore.proxy.username"
              :placeholder="t('settings.usernamePlaceholder')"
              allow-clear
            />
          </a-form-item>
          <a-form-item :label="t('settings.password')">
            <a-input-password
              v-model:value="settingsStore.proxy.password"
              :placeholder="t('settings.passwordPlaceholder')"
              allow-clear
            />
          </a-form-item>

          <div class="section-title">{{ t("settings.system") }}</div>
          <a-form-item :label="t('settings.systemLanguage')">
            <a-select
              v-model:value="settingsStore.system.language"
              :placeholder="t('settings.selectLanguage')"
              style="width: 100%"
            >
              <a-select-option value="zh-CN">简体中文</a-select-option>
              <a-select-option value="en">English</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label="t('settings.theme')">
            <a-select
              v-model:value="settingsStore.system.themeMode"
              style="width: 100%"
            >
              <a-select-option value="system">
                {{ t("settings.theme.system") }}
              </a-select-option>
              <a-select-option value="light">
                {{ t("settings.theme.light") }}
              </a-select-option>
              <a-select-option value="dark">
                {{ t("settings.theme.dark") }}
              </a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label="t('settings.startupOnBoot')">
            <a-switch v-model:checked="settingsStore.system.startupOnBoot" />
            <span class="setting-desc">
              {{ t("settings.startupOnBootDesc") }}
            </span>
          </a-form-item>
          <a-form-item :label="t('settings.enableLogRecording')">
            <a-switch
              v-model:checked="settingsStore.system.enableLogRecording"
            />
            <span class="setting-desc">
              {{ t("settings.enableLogRecordingDesc") }}
            </span>
          </a-form-item>
        </a-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from "vue";
import { useSettingsStore } from "@/stores/settings";
import { t } from "@/i18n";
import { theme } from "ant-design-vue";

// 将 antd token 映射为布局 CSS 变量
const { token } = theme.useToken();

const settingsStore = useSettingsStore();

const props = defineProps({
  handleDividerMove: {
    type: Function,
    required: true,
  },
});

const isDividerDragging = ref(false);

// 拓展面板左右拖拽限制大小
const handleDividerMove = (e) => {
  if (!isDividerDragging.value) return;
  props.handleDividerMove(e);
};

watch(isDividerDragging, (isDragging) => {
  if (isDragging) {
    document.body.style.cursor = "e-resize";
    document.body.style.userSelect = "none";
    window.addEventListener("mousemove", handleDividerMove);
    window.addEventListener("mouseup", () => {
      isDividerDragging.value = false;
    });
  } else {
    document.body.style.cursor = "";
    document.body.style.userSelect = "";
    window.removeEventListener("mousemove", handleDividerMove);
  }
});
</script>

<style lang="scss" scoped>
.socks5-extension {
  flex: 1;
  display: flex;
  background: v-bind("token.colorBgContainer");
  border-left: 1px solid v-bind("token.colorBorderSecondary");

  .socks5-divider {
    width: 4px;
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
    overflow-y: auto;
  }
}
</style>
