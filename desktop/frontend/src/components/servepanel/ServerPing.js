import { ref } from "vue";
import { message } from "ant-design-vue";
import { TestServer } from "@bindings/desktop/proxy/proxyping";
import { useSettingsStore } from "@/stores/settings";
import { t } from "@/locale";

export function useServerPing(servers) {
  const settingsStore = useSettingsStore();
  const latencyById = ref({});
  const pingingAll = ref(false);

  async function pingAllServers() {
    if (servers.value.length === 0) return;
    const latencyHost = (settingsStore.proxy.latencyTestHost || "").trim();
    if (!latencyHost) {
      message.warning(t("settings.latencyTestHostRequired"));
      return;
    }
    pingingAll.value = true;
    try {
      await Promise.all(
        servers.value.map(async (server) => {
          const key = server.__id;
          const protocol = server.protocol.toLowerCase();
          const username = server.username || "";
          const password = server.password || "";
          const proxyURL = `${protocol}://${username}:${password}@${server.host}`;
          try {
            const result = await TestServer(proxyURL, latencyHost);
            latencyById.value[key] = result.success ? result.latency : -1;
          } catch {
            latencyById.value[key] = -1;
          }
        }),
      );
      message.success(t("serverList.pingAllDone"));
    } finally {
      pingingAll.value = false;
    }
  }

  function getLatencyClass(latency) {
    if (latency < 100) return "latency-good";
    if (latency < 300) return "latency-medium";
    return "latency-bad";
  }

  return { latencyById, pingingAll, pingAllServers, getLatencyClass };
}
