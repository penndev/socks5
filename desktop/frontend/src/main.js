import "ant-design-vue/dist/reset.css";
import { createApp } from "vue";
import { createPinia } from "pinia";
import Antd from "ant-design-vue";
import App from "@/App.vue";
import { bootI18n } from "@/i18n";

await bootI18n();

const app = createApp(App);
app.use(createPinia()).use(Antd);
app.mount("#app");
