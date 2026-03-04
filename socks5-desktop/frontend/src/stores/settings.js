import { defineStore } from "pinia";
import { Get, Set } from "@bindings/socks5-desktop/storage";
import { notification } from "ant-design-vue";
import { debounce } from "@/utils";
import { setLocale, useI18n, detectSystemLocale } from "@/i18n";
import { applyTheme } from "@/theme";

const KEY = "settings";
const SAVE_DEBOUNCE_MS = 500;
const DEFAULT_THEME_MODE = "system";
const NOTIFICATION_PLACEMENT = "topRight";

const defaultProxy = () => ({
  host: "127.0.0.1",
  port: 1080,
  username: "",
  password: "",
});

const defaultSystem = () => ({
  language: detectSystemLocale(),
  themeMode: DEFAULT_THEME_MODE,
  startupOnBoot: false,
  enableLogRecording: true,
});

function normalizeThemeMode(mode) {
  if (mode === "dark" || mode === "light" || mode === "system") {
    return mode;
  }
  return DEFAULT_THEME_MODE;
}

function showSaveNotification(type, message) {
  notification[type]({
    message,
    placement: NOTIFICATION_PLACEMENT,
  });
}

function syncRuntimePreferences(systemSettings) {
  systemSettings.language = setLocale(systemSettings.language);
  systemSettings.themeMode = normalizeThemeMode(systemSettings.themeMode);
  applyTheme(systemSettings.themeMode);
}

export const useSettingsStore = defineStore(KEY, {
  state: () => ({
    proxy: defaultProxy(),
    system: defaultSystem(),
  }),

  actions: {
    /** 从存储加载并合并到 state */
    async load() {
      try {
        const storedSettings = await Get(KEY);
        if (storedSettings?.proxy) {
          Object.assign(this.proxy, storedSettings.proxy);
        }
        if (storedSettings?.system) {
          Object.assign(this.system, storedSettings.system);
        }
      } catch (_) {
        // 存储加载失败时保持默认配置
      }

      syncRuntimePreferences(this.system);
    },

    /** 将 state 持久化到存储 */
    async save() {
      const { t } = useI18n();

      try {
        await Set(KEY, {
          proxy: { ...this.proxy },
          system: { ...this.system },
        });
        showSaveNotification("success", t("settings.saveSuccess"));
      } catch (e) {
        showSaveNotification("error", t("settings.saveError"));
      }
    },

    /** 初始化：加载后订阅 state 变化，自动防抖保存 */
    initSync() {
      const debouncedSave = debounce(() => this.save(), SAVE_DEBOUNCE_MS);

      this.$subscribe(() => {
        syncRuntimePreferences(this.system);
        debouncedSave();
      });
    },
  },
});
