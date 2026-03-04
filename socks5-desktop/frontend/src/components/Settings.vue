<template>
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
        <a-switch v-model:checked="settingsStore.system.enableLogRecording" />
        <span class="setting-desc">
          {{ t("settings.enableLogRecordingDesc") }}
        </span>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { useSettingsStore } from "@/stores/settings";
import { useI18n } from "@/i18n";

const settingsStore = useSettingsStore();
const { t } = useI18n();
</script>
