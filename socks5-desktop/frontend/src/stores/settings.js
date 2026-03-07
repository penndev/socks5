import { defineStore } from "pinia";
import { Get, Set } from "@bindings/socks5-desktop/storage";
import { notification } from "ant-design-vue";
import { debounce } from "@/utils";
import { setLocale, defaultLocale, t } from "@/i18n";

// 存储的键名
const KEY = "settings";

export const useSettingsStore = defineStore(KEY, {
  state: () => ({
    proxy: {
      host: "127.0.0.1",
      port: 1080,
      username: "",
      password: "",
    },
    system: {
      language: defaultLocale,
      themeMode: "system",
      startupOnBoot: false,
      enableLogRecording: false,
    },
  }),

  actions: {
    /** 将 state 持久化到存储 */
    async save() {
      this.initSystem();
      try {
        await Set(KEY, {
          proxy: { ...this.proxy },
          system: { ...this.system },
        });
        notification.success({
          message: t("settings.saveSuccess"),
          placement: "topRight",
        });
      } catch (error) {
        notification.error({
          message: t("settings.saveError"),
          placement: "topRight",
        });
      }
    },

    initSystem() {
      // 设置系统语言
      setLocale(this.system.language);
    },
    /** 从存储加载并合并到 state */
    async init() {
      try {
        const storedSettings = await Get(KEY);
        if (storedSettings) {
          Object.assign(this, storedSettings);
        }
        this.$subscribe(() => debounce(this.save(), 800));
      } finally {
        // 存储加载失败时保持默认配置
        this.initSystem();
      }
    },
  },
});
