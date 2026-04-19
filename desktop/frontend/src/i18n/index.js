import { ref, computed } from "vue";
import * as Lang from "@bindings/desktop/lang/lang";
import { Events } from "@wailsio/runtime";
import { AppConfig } from "@bindings/desktop/internal/appconst";

export var defaultLocale = "zh-CN";

if (typeof navigator !== "undefined" && navigator.language) {
  const language = navigator.languages?.[0] || navigator.language;
  const lower = language.toLowerCase();
  if (lower.startsWith("zh")) defaultLocale = "zh-CN";
  if (lower.startsWith("en")) defaultLocale = "en";
}

const currentLocale = ref(defaultLocale);
const messages = ref({});

export const locale = computed(() => currentLocale.value);

export function t(key) {
  const m = messages.value;
  return m[key] ?? key;
}

/**
 * 加载文案并可选择是否同步 Go（Go 已通过其它路径切换时不要 notifyGo）。
 */
async function applyLocale(locale, notifyGo) {
  if (!locale) return;
  try {
    const bundle = await Lang.Bundle(locale);
    messages.value =
      bundle && typeof bundle === "object" ? { ...bundle } : {};
    currentLocale.value = locale;
    if (notifyGo) {
      await Lang.SetLocale(locale);
    }
  } catch (e) {
    console.error("[i18n] applyLocale failed:", e);
  }
}

/** 用户在界面选择语言或由 Pinia 写入后调用：更新文案并通知 Go（派发 localeChanged）。 */
export async function setLocale(localeId) {
  await applyLocale(localeId, true);
  return localeId;
}

/** 首次启动时注入 Bundle（不调用 Go SetLocale，避免与持久化语言重复同步）。 */
export async function bootI18n() {
  await applyLocale(defaultLocale, false);
}

let localeListenerOff = null;

/**
 * 监听 Go 端语言切换，可选写回 Pinia（避免环：仅在值变化时写入）。
 */
export async function subscribeLocaleEvents(syncPiniaLanguage) {
  try {
    const appConst = await AppConfig();
    localeListenerOff = Events.On(
      appConst.EventNameLocaleChanged,
      (eventPayload) => {
        const loc = String(eventPayload?.data ?? "");
        if (!loc || loc === currentLocale.value) return;
        void (async () => {
          await applyLocale(loc, false);
          if (typeof syncPiniaLanguage === "function") {
            syncPiniaLanguage(loc);
          }
        })();
      },
    );
  } catch (e) {
    console.error("[i18n] subscribeLocaleEvents failed:", e);
  }
}

export function unsubscribeLocaleEvents() {
  if (localeListenerOff) {
    localeListenerOff();
    localeListenerOff = null;
  }
}
