import "@wailsio/runtime";
import { createApp } from "vue";
import { createPinia } from "pinia";
import Antd from "ant-design-vue";
import App from "./App.vue";
import "ant-design-vue/dist/reset.css";

createApp(App).use(createPinia()).use(Antd).mount("#app");
