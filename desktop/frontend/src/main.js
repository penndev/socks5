import "ant-design-vue/dist/reset.css";
import { createApp } from "vue";
import { createPinia } from "pinia";
import Antd from "ant-design-vue";
import App from "@/App.vue";

// 创建根应用并挂载插件
const app = createApp(App);
app.use(createPinia()).use(Antd);
app.mount("#app");
