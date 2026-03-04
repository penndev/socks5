<template>
  <div class="settings-panel">
      <a-form layout="vertical">
        <div class="section-title">本地代理</div>
        <a-form-item label="IP 地址">
          <a-select v-model:value="settingsStore.proxy.host" placeholder="选择IP地址" style="width: 100%">
            <a-select-option value="127.0.0.1">127.0.0.1</a-select-option>
            <a-select-option value="0.0.0.0">0.0.0.0</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="端口">
          <a-input-number v-model:value="settingsStore.proxy.port" :min="1" :max="65535" placeholder="1080" style="width: 100%" />
        </a-form-item>
        <a-form-item label="用户名">
          <a-input v-model:value="settingsStore.proxy.username" placeholder="可选" allow-clear />
        </a-form-item>
        <a-form-item label="密码">
          <a-input-password v-model:value="settingsStore.proxy.password" placeholder="可选" allow-clear />
        </a-form-item>

        <div class="section-title">系统</div>
        <a-form-item label="系统语言">
          <a-select v-model:value="settingsStore.system.language" placeholder="选择语言" style="width: 100%">
            <a-select-option value="zh-CN">简体中文</a-select-option>
            <a-select-option value="en">English</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="开机启动">
          <a-switch v-model:checked="settingsStore.system.startupOnBoot" />
          <span class="setting-desc">启动时自动运行本应用</span>
        </a-form-item>
        <a-form-item label="日志开关">
          <a-switch v-model:checked="settingsStore.system.enableLogRecording" />
          <span class="setting-desc">记录运行与连接日志，关闭后底部日志栏将隐藏</span>
        </a-form-item>
      </a-form>
  </div>
</template>

<script setup>
import { useSettingsStore } from "@/stores/settings";

const settingsStore = useSettingsStore();
</script>
