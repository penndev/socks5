<template>
  <a-card class="socks5-server-card" title="节点列表">
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
                <span class="server-host">
                  <CheckCircleFilled v-if="selectedServer?.id === item.id" class="selected-icon" />
                  {{ item.host }}
                </span>
              </template>
              <template #description>
                <span class="server-meta">
                  {{ item.protocol }} | {{ item.username || "无认证" }}
                </span>
              </template>
            </a-list-item-meta>
          </a-list-item>
        </template>
      </a-list>
    </div>

    <div class="server-list-empty" v-else>
      <a-empty description="暂无节点，点击下方按钮添加" />
    </div>

    <a-button type="primary" block class="add-server-btn" @click="edit.open()">
      <PlusOutlined />
      添加节点
    </a-button>

    <a-modal
      v-model:open="edit.visible"
      :title="edit.title"
      :confirm-loading="edit.loading"
      @ok="edit.submit"
      @cancel="edit.visible = false"
    >
      <a-form ref="editRef" :model="edit.form" :rules="edit.rules" layout="vertical">
        <a-form-item label="地址" name="host">
          <a-input v-model:value="edit.form.host" placeholder="127.0.0.1:1080" allow-clear />
        </a-form-item>

        <a-form-item label="协议" name="protocol">
          <a-select v-model:value="edit.form.protocol" placeholder="选择协议">
            <a-select-option value="Socks5">Socks5</a-select-option>
            <a-select-option value="Socks5OverTls">Socks5OverTls</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="用户名" name="username">
          <a-input v-model:value="edit.form.username" placeholder="可选" allow-clear />
        </a-form-item>

        <a-form-item label="密码" name="password">
          <a-input-password v-model:value="edit.form.password" placeholder="可选" allow-clear />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from "vue";
import { CheckCircleFilled, DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons-vue";
import { Modal, message } from "ant-design-vue";
import { Get, Set } from "@bindings/socks5-desktop/storage";
import { useServerStore } from "../stores/server";

const serverStore = useServerStore();
const selectedServer = computed(() => serverStore.selectedServer);

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
    username: "",
    password: "",
    protocol: "Socks5",
  },
  rules: {
    host: [
      { required: true, message: "请输入地址" },
      { pattern: /^[^:]+:\d{1,5}$/, message: "格式为 host:port" },
    ],
    protocol: [{ required: true, message: "请选择协议" }],
  },

  open(server = null) {
    edit.id = server?.id ?? null;
    edit.title = edit.id ? "编辑节点" : "添加节点";
    edit.form.host = server?.host ?? "";
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
        username: edit.form.username?.trim() ?? "",
        password: edit.form.password ?? "",
        protocol: edit.form.protocol,
      };

      if (edit.id) {
        const idx = servers.value.findIndex((s) => s.id === edit.id);
        if (idx >= 0) servers.value[idx] = { ...servers.value[idx], ...payload };

        if (selectedServer?.id === edit.id)
          serverStore.selectedServer = { ...selectedServer, ...payload };

        message.success("修改成功");
      } else {
        servers.value.push({ id: Date.now().toString(), ...payload });
        message.success("添加成功");
      }

      await Set(STORAGE_KEY, servers.value);
      edit.visible = false;
    } catch (e) {
      if (!e?.errorFields) message.error(e?.message || "操作失败");
    } finally {
      edit.loading = false;
    }
  },
});

const STORAGE_KEY = "servers";

function deleteModal(item) {
  Modal.confirm({
    title: "删除节点",
    content: `确定删除节点 ${item.host} 吗？`,
    okType: "danger",
    okText: "删除",
    cancelText: "取消",
    async onOk() {
      servers.value = servers.value.filter((s) => s.id !== item.id);
      if (selectedServer?.id === item.id) serverStore.selectedServer = null;
      await Set(STORAGE_KEY, servers.value);
      message.success("已删除");
    },
  });
}

onMounted(async () => {
  try {
    servers.value = await Get(STORAGE_KEY);
  } catch {
    message.error("加载节点失败");
  }
});
</script>

<style scoped lang="scss">
.socks5-server-card {
  min-height: 400px;
  flex: 1;
  display: flex;
  flex-direction: column;
  border-radius: 14px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.04);

  :deep(.ant-card-body) {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .server-list-scroll {
    flex: 1;
    margin-bottom: 12px;
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
    }

    .selected-icon {
      color: #1677ff;
      font-size: 14px;
    }

    .server-meta {
      font-size: 12px;
      color: #6b7280;
    }
  }
}
</style>