import "@wailsio/runtime";
import { createApp } from "vue";
import { createPinia } from "pinia";
import Antd from "ant-design-vue";
import App from "./App.vue";
import { useSettingsStore } from "./stores/settings";
import "ant-design-vue/dist/reset.css";

const app = createApp(App);
app.use(createPinia()).use(Antd);

(async () => {
  const settingsStore = useSettingsStore();
  await settingsStore.load();
  settingsStore.initSync();
  app.mount("#app");
})();
