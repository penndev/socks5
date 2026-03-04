import { defineStore } from "pinia";
import { Get, Set } from "@bindings/socks5-desktop/storage";
import { message } from "ant-design-vue";
import { debounce } from "@/utils";

const KEY = "settings";

const defaultProxy = () => ({
  host: "127.0.0.1",
  port: 1080,
  username: "",
  password: "",
});

const defaultSystem = () => ({
  language: "zh-CN",
  startupOnBoot: false,
  enableLogRecording: true,
});

export const useSettingsStore = defineStore(KEY, {
  state: () => ({
    proxy: defaultProxy(),
    system: defaultSystem(),
  }),

  actions: {
    /** 从存储加载并合并到 state */
    async load() {
      try {
        const d = await Get(KEY);
        if (d?.proxy) Object.assign(this.proxy, d.proxy);
        if (d?.system) Object.assign(this.system, d.system);
      } catch (_) {}
    },

    /** 将 state 持久化到存储 */
    async save() {
      try {
        await Set(KEY, {
          proxy: { ...this.proxy },
          system: { ...this.system },
        });
        message.success("设置已保存");
      } catch (e) {
        message.error("保存失败");
      }
    },

    /** 初始化：加载后订阅 state 变化，自动防抖保存 */
    initSync() {
      const store = this;
      const debouncedSave = debounce(() => store.save(), 500);
      store.$subscribe(debouncedSave);
    },
  },
});
