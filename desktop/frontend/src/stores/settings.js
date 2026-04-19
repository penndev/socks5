import { defineStore } from "pinia";
import { Storage } from "@bindings/desktop/storage";
import { notification } from "ant-design-vue";
import { debounce } from "@/utils";
import { defaultLocale, setLocale, t } from "@/i18n";

// 存储的键名
const KEY = "settings";

function defaultProxyState() {
  return {
    host: "127.0.0.1",
    port: 1080,
    latencyTestHost: "",
    username: "",
    password: "",
  };
}

function defaultSystemState() {
  return {
    language: defaultLocale,
    themeMode: "system",
    startupOnBoot: false,
    enableLogRecording: false,
  };
}

export const useSettingsStore = defineStore(KEY, {
  state: () => ({
    proxy: defaultProxyState(),
    system: defaultSystemState(),
  }),

  actions: {
    /** 将 state 持久化到存储 */
    async save() {
      try {
        await Storage.SetSettings({
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


    /** 从存储加载并合并到 state */
    async init() {
      try {
        const storedSettings = await Storage.GetSettings();
        if (storedSettings) {
          this.proxy = {
            ...defaultProxyState(),
            ...(storedSettings.proxy ?? {}),
          };
          this.system = {
            ...defaultSystemState(),
            ...(storedSettings.system ?? {}),
          };
        }
        const save = debounce(this.save, 800)
        this.$subscribe(save);
      } catch (error) {
        notification.error({
          message: t("settings.loadError"),
          placement: "topRight",
        });
      }
    },
  },
});
