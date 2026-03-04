<template>
  <a-card class="socks5-server-card" :title="t('serverList.title')">
    <div class="server-list-scroll" v-if="servers.length > 0">
      <a-list :data-source="servers" bordered>
        <template #renderItem="{ item }">
          <a-list-item :class="{ active: selectedServer?.id === item.id }" @click="serverStore.selectedServer = item">
            <template #actions>
              <a-button type="text" size="small" @click.stop="edit.open(item)">
                <EditOutlined />
              </a-button>
              <a-button type="text" danger size="small" @click.stop="deleteModal(item)">
                <DeleteOutlined />
              </a-button>
            </template>

            <a-list-item-meta>
              <template #title>
                <span class="server-host" :title="item.remark || item.host">
                  <CheckCircleFilled v-if="selectedServer?.id === item.id" class="selected-icon" />
                  {{ item.remark || item.host }}
                </span>
              </template>
              <template #description>
                <span
                  class="server-meta"
                  :title="`${item.protocol} | ${
                    item.username || t('serverList.noAuth')
                  }`"
                >
                  {{ item.protocol }} | {{ item.username || t("serverList.noAuth") }}
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
            <a-select-option value="Socks5">Socks5</a-select-option>
            <a-select-option value="Socks5OverTLS">Socks5OverTLS</a-select-option>
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
import { ref, reactive, computed, onMounted } from "vue";
import {
  CheckCircleFilled,
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
} from "@ant-design/icons-vue";
import { Modal, message } from "ant-design-vue";
import { Get, Set } from "@bindings/socks5-desktop/storage";
import { useServerStore } from "../stores/server";
import { useI18n } from "@/i18n";

const serverStore = useServerStore();
const selectedServer = computed(() => serverStore.selectedServer);
const { t } = useI18n();

// 所有节点
const servers = ref([]);

const editRef = ref();
const edit = reactive({
  visible: false,
  loading: false,
  title: "",
  id: 0,
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
      { pattern: /^[^:]+:\d{1,5}$/, message: t("serverList.validateHostFormat") },
    ],
    protocol: [{ required: true, message: t("serverList.validateProtocolRequired") }],
  },

  open(server = null) {
    edit.id = server?.id ?? null;
    edit.title = edit.id ? t("serverList.editTitle") : t("serverList.addTitle");
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

      if (edit.id) {
        const idx = servers.value.findIndex((s) => s.id === edit.id);
        if (idx >= 0)
          servers.value[idx] = { ...servers.value[idx], ...payload };

        if (selectedServer?.id === edit.id)
          serverStore.selectedServer = { ...selectedServer, ...payload };

        message.success(t("serverList.updateSuccess"));
      } else {
        servers.value.push({ id: Date.now().toString(), ...payload });
        message.success(t("serverList.addSuccess"));
      }

      await Set(STORAGE_KEY, servers.value);
      edit.visible = false;
    } catch (e) {
      if (!e?.errorFields)
        message.error(e?.message || t("serverList.operationFailed"));
    } finally {
      edit.loading = false;
    }
  },
});

const STORAGE_KEY = "servers";

function deleteModal(item) {
  Modal.confirm({
    title: t("serverList.deleteTitle"),
    content: `${t("serverList.deleteContentPrefix")}${
      item.remark || item.host
    }${t("serverList.deleteContentSuffix")}`,
    okType: "danger",
    okText: t("serverList.deleteOkText"),
    cancelText: t("serverList.deleteCancelText"),
    async onOk() {
      servers.value = servers.value.filter((s) => s.id !== item.id);
      if (selectedServer?.id === item.id) serverStore.selectedServer = null;
      await Set(STORAGE_KEY, servers.value);
      message.success(t("serverList.deleteSuccess"));
    },
  });
}

onMounted(async () => {
  try {
    servers.value = await Get(STORAGE_KEY);
  } catch {
    message.error(t("serverList.loadFailed"));
  }
});
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
      background: #e6f4ff;
    }

    .server-host {
      font-weight: 500;
      display: flex;
      align-items: center;
      gap: 6px;
      min-width: 0;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .selected-icon {
      color: #1677ff;
      font-size: 14px;
    }

    .server-meta {
      font-size: 12px;
      color: var(--socks-text-secondary);
      display: block;
      min-width: 0;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}
</style>