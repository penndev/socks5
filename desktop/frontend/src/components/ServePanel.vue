<template>
  <a-card class="socks5-server-card" :title="t('serverList.title')">
    <template #extra>
      <div class="card-extra-actions">
        <a-tooltip :title="t('serverList.pingAll')">
          <a-button
            type="text"
            size="small"
            :loading="pingingAll"
            :disabled="servers.length === 0"
            @click="pingAllServers"
          >
            <ThunderboltOutlined />
          </a-button>
        </a-tooltip>
        <a-tooltip title="打开订阅节点批量编辑页面">
          <a-button type="text" size="small" @click="openSubscribeEditor">
            <EditOutlined />
          </a-button>
        </a-tooltip>
      </div>
    </template>
    <div class="server-list-scroll" v-if="servers.length > 0">
      <a-list :data-source="servers" row-key="__id" bordered>
        <template #renderItem="{ item }">
          <a-list-item :class="{ active: isServerSelected(item) }" @click="selectedServer = item">
            <template #actions>
              <a-button type="text" size="small" @click.stop="edit.open(item)">
                <EditOutlined />
              </a-button>
              <a-button
                type="text"
                danger
                size="small"
                @click.stop="deleteModal(item)"
              >
                <DeleteOutlined />
              </a-button>
            </template>

            <a-list-item-meta>
              <template #title>
                <span class="server-title-row">
                  <CheckCircleFilled v-if="isServerSelected(item)" class="selected-icon" />
                  <span class="server-host" :title="item.remark || item.host">
                    {{ item.remark || item.host }}
                  </span>
                  <span
                    v-if="latencyById[item.__id] !== undefined"
                    class="latency-inline"
                  >
                    <span
                      v-if="latencyById[item.__id] >= 0"
                      :class="getLatencyClass(latencyById[item.__id])"
                    >
                      {{ latencyById[item.__id] }}ms
                    </span>
                    <span v-else class="latency-error">{{
                      t("serverList.pingFailed")
                    }}</span>
                  </span>
                </span>
              </template>
              <template #description>
                <span
                  class="server-meta"
                  :title="`${item.protocol} | ${
                    item.username || t('serverList.noAuth')
                  }`"
                >
                  <span class="server-protocol">{{ item.protocol }}</span>
                  <span class="server-separator">|</span>
                  <span class="server-username">{{
                    item.username || t("serverList.noAuth")
                  }}</span>
                </span>
              </template>
            </a-list-item-meta>
          </a-list-item>
        </template>
      </a-list>
    </div>

    <div class="server-list-empty" v-else>
      <a-empty :description="t('serverList.emptyDescription')" />
    </div>

    <a-button type="primary" block class="add-server-btn" @click="edit.open()">
      <PlusOutlined />
      {{ t("serverList.addButton") }}
    </a-button>

    <a-modal
      v-model:open="edit.visible"
      :title="edit.title"
      :confirm-loading="edit.loading"
      @ok="edit.submit"
      @cancel="edit.visible = false"
    >
      <a-form
        ref="editRef"
        :model="edit.form"
        :rules="edit.rules"
        layout="vertical"
      >
        <a-form-item :label="t('serverList.hostLabel')" name="host">
          <a-input
            v-model:value="edit.form.host"
            :placeholder="t('serverList.hostPlaceholder')"
            allow-clear
          />
        </a-form-item>

        <a-form-item :label="t('serverList.remarkLabel')" name="remark">
          <a-input
            v-model:value="edit.form.remark"
            :placeholder="t('serverList.remarkPlaceholder')"
            allow-clear
          />
        </a-form-item>

        <a-form-item :label="t('serverList.protocolLabel')" name="protocol">
          <a-select
            v-model:value="edit.form.protocol"
            :placeholder="t('serverList.selectProtocol')"
          >
            <a-select-option
              v-for="scheme in proxySchemes"
              :key="scheme"
              :value="scheme"
              >{{ scheme }}</a-select-option
            >
          </a-select>
        </a-form-item>

        <a-form-item :label="t('serverList.usernameLabel')" name="username">
          <a-input
            v-model:value="edit.form.username"
            :placeholder="t('serverList.usernamePlaceholder')"
            allow-clear
          />
        </a-form-item>

        <a-form-item :label="t('serverList.passwordLabel')" name="password">
          <a-input-password
            v-model:value="edit.form.password"
            :placeholder="t('serverList.passwordPlaceholder')"
            allow-clear
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import {
  CheckCircleFilled,
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
  ThunderboltOutlined,
} from "@ant-design/icons-vue";
import { Modal, message } from "ant-design-vue";
import { Storage } from "@bindings/desktop/storage";
import { TestServer } from "@bindings/desktop/proxy/proxyping";
import { AppConfig, ProxyScheme } from "@bindings/desktop/internal/appconst";
import { OpenExternalURL } from "@bindings/desktop/internal/appconst";
import { useServerStore } from "../stores/server";
import { useSettingsStore } from "../stores/settings";
import { t } from "@/i18n";
import { storeToRefs } from "pinia";
import { Events } from "@wailsio/runtime";

import { theme } from "ant-design-vue";
const { token } = theme.useToken();


const serverStore = useServerStore();
const { selectedServer } = storeToRefs(serverStore);
const settingsStore = useSettingsStore();

// 所有节点
const servers = ref([]);
/** 仅内存：测速结果按列表运行时 id 存，不写入 storage */
const latencyById = ref({});
const pingingAll = ref(false);

function withRuntimeIds(serverList) {
  const list = Array.isArray(serverList) ? serverList : [];
  return list.map((server) => {
    const __id = `${server.host}|${server.protocol}|${server.username}|${server.password}`;
    return { ...server, __id };
  });
}

function isServerSelected(server) {
  if (!selectedServer.value) return false;
  return selectedServer.value.__id === server?.__id;
}

const editRef = ref();
// 编辑态集中管理：弹窗状态、表单、校验与提交动作
const edit = reactive({
  visible: false,
  loading: false,
  title: "",
  key: "",
  form: {
    host: "",
    remark: "",
    username: "",
    password: "",
    protocol: "Socks5",
  },
  rules: {
    host: [
      { required: true, message: t("serverList.validateHostRequired") },
      {
        pattern: /^[^:]+:\d{1,5}$/,
        message: t("serverList.validateHostFormat"),
      },
    ],
    protocol: [
      { required: true, message: t("serverList.validateProtocolRequired") },
    ],
  },

  open(server = null) {
    edit.key = server ? server.__id : "";
    edit.title = edit.key ? t("serverList.editTitle") : t("serverList.addTitle");
    edit.form.host = server?.host ?? "";
    edit.form.remark = server?.remark ?? "";
    edit.form.username = server?.username ?? "";
    edit.form.password = server?.password ?? "";
    edit.form.protocol = server?.protocol ?? "Socks5";
    edit.visible = true;
  },

  async submit() {
    try {
      await editRef.value.validate();
      edit.loading = true;

      const payload = {
        host: edit.form.host.trim(),
        remark: edit.form.remark?.trim() ?? "",
        username: edit.form.username?.trim() ?? "",
        password: edit.form.password ?? "",
        protocol: edit.form.protocol,
      };

      const selectedWasEdited =
        !!edit.key &&
        !!selectedServer.value &&
        selectedServer.value.__id === edit.key;
      let editedIdx = -1;

      if (edit.key) {
        editedIdx = servers.value.findIndex((s) => s.__id === edit.key);
        if (editedIdx >= 0)
          servers.value[editedIdx] = {
            ...servers.value[editedIdx],
            ...payload,
          };

        delete latencyById.value[edit.key];
        message.success(t("serverList.updateSuccess"));
      } else {
        servers.value.push(payload);
        message.success(t("serverList.addSuccess"));
      }

      servers.value = withRuntimeIds(servers.value);
      if (selectedWasEdited && editedIdx >= 0)
        selectedServer.value = servers.value[editedIdx];
      await Storage.SetServers(
        servers.value.map((row) => {
          const { id: _i, __id: _x, ...rest } = /** @type {any} */ (row);
          return rest;
        }),
      );
      edit.visible = false;
    } catch (e) {
      if (!e?.errorFields)
        message.error(e?.message || t("serverList.operationFailed"));
    } finally {
      edit.loading = false;
    }
  },
});

function deleteModal(item) {
  const key = item.__id;
  Modal.confirm({
    title: t("serverList.deleteTitle"),
    content: `${t("serverList.deleteContentPrefix")}${
      item.remark || item.host
    }${t("serverList.deleteContentSuffix")}`,
    okType: "danger",
    okText: t("serverList.deleteOkText"),
    cancelText: t("serverList.deleteCancelText"),
    async onOk() {
      servers.value = servers.value.filter((s) => s.__id !== key);
      servers.value = withRuntimeIds(servers.value);
      delete latencyById.value[key];
      if (selectedServer.value && selectedServer.value.__id === key)
        selectedServer.value = null;
      await Storage.SetServers(
        servers.value.map((row) => {
          const { id: _i, __id: _x, ...rest } = /** @type {any} */ (row);
          return rest;
        }),
      );
      message.success(t("serverList.deleteSuccess"));
    },
  });
}

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

async function openSubscribeEditor() {
  const rawHost = (settingsStore.proxy.host || "").trim();
  const host = rawHost === "0.0.0.0" || rawHost === "" ? "127.0.0.1" : rawHost;
  const port = Number(settingsStore.proxy.port);
  if (!Number.isInteger(port) || port < 1 || port > 65535) {
    message.warning("请先配置正确的本地代理端口");
    return;
  }
  const url = `http://${host}:${port}/subscribe/`;
  try {
    await OpenExternalURL(url);
  } catch (e) {
    message.error(e?.message || "打开浏览器失败");
  }
}

// 获取延迟样式类
function getLatencyClass(latency) {
  if (latency < 100) return "latency-good";
  if (latency < 300) return "latency-medium";
  return "latency-bad";
}

const proxySchemes = ref([]);

onMounted(async () => {
  const schemes = await ProxyScheme();
  proxySchemes.value = schemes;
  await loadServers();
  const appConfig = await AppConfig();
  Events.On(appConfig.EventNameServersChanged, async () => {
    await loadServers();
  });
});

async function loadServers() {
  try {
    const raw = await Storage.GetServers();
    servers.value = withRuntimeIds(
      (Array.isArray(raw) ? raw : []).map((row) => {
        const { id: _i, __id: _x, ...rest } = /** @type {any} */ (row);
        return rest;
      }),
    );
  } catch {
    message.error(t("serverList.loadFailed"));
  }
}
</script>

<style scoped lang="scss">
.socks5-server-card {
  min-height: 200px;
  flex: 1;
  display: flex;
  flex-direction: column;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);

  :deep(.ant-card-body) {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 12px;
  }

  :deep(.ant-list-item-meta-title) {
    overflow: hidden;
    text-overflow: ellipsis;
  }

  :deep(.ant-list-item-meta-description) {
    min-width: 0;
  }

  .card-extra-actions {
    display: inline-flex;
    align-items: center;
    gap: 4px;
  }

  .server-list-scroll {
    flex: 1;
    margin-bottom: 8px;
    min-height: 0;
  }

  .server-list-empty {
    flex: 1;
    min-height: 120px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .add-server-btn {
    flex-shrink: 0;
  }

  :deep(.ant-list-item) {
    cursor: pointer;

    &.active {
      background: v-bind("token.colorPrimaryBgHover");
    }

    .server-title-row {
      display: flex;
      align-items: center;
      gap: 6px;
      min-width: 0;
    }

    .server-host {
      font-weight: 500;
      flex: 1;
      min-width: 0;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .latency-inline {
      flex-shrink: 0;
      margin-left: 6px;
      font-weight: 500;
      font-size: 12px;

      .latency-good {
        color: #52c41a;
      }

      .latency-medium {
        color: #faad14;
      }

      .latency-bad {
        color: #ff4d4f;
      }

      .latency-error {
        color: #ff4d4f;
        font-size: 11px;
      }
    }

    .selected-icon {
      color: #1677ff;
      font-size: 14px;
    }

    .server-meta {
      font-size: 12px;
      color: v-bind("token.colorTextSecondary");
      display: flex;
      align-items: center;
      gap: 4px;
      min-width: 0;
      overflow: hidden;
    }

    .server-protocol,
    .server-separator {
      flex-shrink: 0;
    }

    .server-username {
      min-width: 0;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      display: inline-block;
      max-width: 100%;
    }

  }
}
</style>
