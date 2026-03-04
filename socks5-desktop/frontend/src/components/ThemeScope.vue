<template>
  <div :style="themeCssVars">
    <slot />
  </div>
</template>

<script setup>
import { computed } from "vue";
import { theme as antdTheme } from "ant-design-vue";

const { token } = antdTheme.useToken();

// 将 antd token 映射为业务布局变量，避免写死颜色
const themeCssVars = computed(() => ({
  "--socks-bg": token.value.colorBgLayout,
  "--socks-card-bg": token.value.colorBgContainer,
  "--socks-header-bg": token.value.colorBgElevated,
  "--socks-card-border": token.value.colorBorderSecondary,
  "--socks-text-primary": token.value.colorText,
  "--socks-text-secondary": token.value.colorTextSecondary,
  "--socks-muted-bg": token.value.colorFillAlter,
  "--socks-bottom-bg": token.value.colorBgElevated,
}));
</script>

