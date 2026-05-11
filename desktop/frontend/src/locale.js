import { ref } from "vue";
import { Bundle, CurrentLocale } from "@bindings/desktop/lang/lang";
import { Events } from "@wailsio/runtime";
import { AppConfig } from "@bindings/desktop/internal/appconst";


// 语言文件对象
export const languageLocale = ref({});

// 监听语言改变事件
export async function subscribeLocaleEvents() {
  const appConst = await AppConfig();
  Events.On(appConst.EventNameLocaleChanged, async (ev) => {
    languageLocale.value = await Bundle(ev.data);
  });
}

/**
 * 
 * @param {string} key 
 * @returns {string}
 */
export function t(key) {
  return languageLocale.value[key] ?? key;
}

