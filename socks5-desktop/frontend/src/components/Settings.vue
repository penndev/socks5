<template>
  <div class="settings-panel">
      <a-form layout="vertical">
        <div class="section-title">本地代理</div>
        <a-form-item label="IP 地址">
          <a-select v-model:value="form.proxy.host" placeholder="选择IP地址" style="width: 100%">
            <a-select-option value="127.0.0.1">127.0.0.1</a-select-option>
            <a-select-option value="0.0.0.0">0.0.0.0</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="端口">
          <a-input-number v-model:value="form.proxy.port" :min="1" :max="65535" placeholder="1080" style="width: 100%" />
        </a-form-item>
        <a-form-item label="用户名">
          <a-input v-model:value="form.proxy.username" placeholder="可选" allow-clear />
        </a-form-item>
        <a-form-item label="密码">
          <a-input-password v-model:value="form.proxy.password" placeholder="可选" allow-clear />
        </a-form-item>

        <div class="section-title">系统</div>
        <a-form-item label="系统语言">
          <a-select v-model:value="form.system.language" placeholder="选择语言" style="width: 100%">
            <a-select-option value="zh-CN">简体中文</a-select-option>
            <a-select-option value="en">English</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="开机启动">
          <a-switch v-model:checked="form.system.startupOnBoot" />
          <span class="setting-desc">启动时自动运行本应用</span>
        </a-form-item>
      </a-form>
  </div>
</template>

<script setup>
import { reactive, onMounted, watch } from "vue";
import { Get, Set } from "@bindings/socks5-desktop/storage";
import { message } from "ant-design-vue";
import { debounce } from "@/utils";

const KEY = "settings";

async function save() {
  try {
    await Set(KEY, { 
      proxy: { ...form.proxy }, 
      system: { ...form.system }
    });
    message.success("设置已保存");
  } catch (e) {
    message.error("保存失败");
  }
}

const form = reactive({
  proxy: { host: "127.0.0.1", port: 1080, username: "", password: "" },
  system: { language: "zh-CN", startupOnBoot:false },
});

onMounted(async() => {
  try {
    const d = await Get(KEY);
    if (d?.proxy) Object.assign(form.proxy, d.proxy);
    if (d?.system) Object.assign(form.system, d.system);
  } catch (e) {
    message.error("加载设置失败");
  }
  watch(form, debounce(save), { deep: true });
});

</script>

<style lang="scss" scoped>
.settings-panel {
  height: 100%;
  overflow-y: auto;

    .section-title {
      font-size: 14px;
      font-weight: 500;
      color: rgba(0, 0, 0, 0.88);
      margin: 16px 0 12px;
      padding: 0;

      &:first-child {
        margin-top: 0;
      }
    }

    .setting-desc {
      margin-left: 8px;
      font-size: 12px;
      color: #6b7280;
    }
  }
</style>
