<template>
  <div class="settings-panel">
    <a-card title="设置" class="settings-card">
      <a-form layout="vertical" :model="formState">
        <a-divider orientation="left">本地代理</a-divider>
        <a-form-item label="IP 地址">
          <a-input v-model:value="formState.localProxy.host" placeholder="127.0.0.1" allow-clear />
        </a-form-item>
        <a-form-item label="端口">
          <a-input-number
            v-model:value="formState.localProxy.port"
            :min="1"
            :max="65535"
            placeholder="1080"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item label="用户名">
          <a-input v-model:value="formState.localProxy.username" placeholder="可选" allow-clear />
        </a-form-item>
        <a-form-item label="密码">
          <a-input-password v-model:value="formState.localProxy.password" placeholder="可选" allow-clear />
        </a-form-item>

        <a-divider orientation="left">系统</a-divider>
        <a-form-item label="系统语言">
          <a-select v-model:value="formState.language" placeholder="选择语言" style="width: 100%">
            <a-select-option value="zh-CN">简体中文</a-select-option>
            <a-select-option value="en">English</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="开机启动">
          <a-switch v-model:checked="formState.startupOnBoot" @change="onStartupChange" />
          <span class="setting-desc">启动时自动运行本应用</span>
        </a-form-item>
      </a-form>

      <a-button type="primary" block @click="handleSave" :loading="saveLoading">保存设置</a-button>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { Get, Set } from "../../bindings/socks5-desktop/storage.js";
import { IsLaunchAtStartup, SetLaunchAtStartup } from "../../bindings/socks5-desktop/app.js";
import { message } from "ant-design-vue";

const STORAGE_KEY = "settings";

const defaultSettings = () => ({
  localProxy: {
    host: "127.0.0.1",
    port: 1080,
    username: "",
    password: "",
  },
  language: "zh-CN",
  startupOnBoot: false,
});

const formState = reactive(defaultSettings());
const saveLoading = ref(false);

const loadSettings = async () => {
  try {
    const data = await Get(STORAGE_KEY);
    if (data && typeof data === "object") {
      Object.assign(formState, {
        localProxy: { ...defaultSettings().localProxy, ...data.localProxy },
        language: data.language || "zh-CN",
        startupOnBoot: data.startupOnBoot ?? false,
      });
    }
    try {
      const enabled = await IsLaunchAtStartup();
      formState.startupOnBoot = enabled;
    } catch (_) {}
  } catch (e) {
    message.error("加载设置失败");
  }
};

const handleSave = async () => {
  try {
    saveLoading.value = true;
    const toSave = {
      localProxy: { ...formState.localProxy },
      language: formState.language,
      startupOnBoot: formState.startupOnBoot,
    };
    await Set(STORAGE_KEY, toSave);
    await SetLaunchAtStartup(formState.startupOnBoot);
    message.success("设置已保存");
  } catch (e) {
    message.error("保存失败");
  } finally {
    saveLoading.value = false;
  }
};

const onStartupChange = async () => {
  try {
    await SetLaunchAtStartup(formState.startupOnBoot);
  } catch (e) {
    message.error("设置开机启动失败");
  }
};

onMounted(() => {
  loadSettings();
});
</script>

<style scoped>
.settings-panel {
  height: 100%;
  overflow-y: auto;
}

.settings-card {
  border-radius: 14px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.04);
}

.settings-card :deep(.ant-card-body) {
  padding-bottom: 24px;
}

.setting-desc {
  margin-left: 8px;
  font-size: 12px;
  color: #6b7280;
}

.a-divider {
  margin: 16px 0 12px;
}
</style>
