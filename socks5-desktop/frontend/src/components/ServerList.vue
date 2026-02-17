<template>
  <a-card class="socks5-server-card" title="节点列表">
    <div class="server-list-scroll" v-if="servers.length > 0">
      <a-list :data-source="servers" bordered>
        <template #renderItem="{ item }">
          <a-list-item
            :class="{ active: selectedId === item.id }"
            @click="selectServer(item)"
          >
            <template #actions>
              <a-button type="text" size="small" @click.stop="editServer(item)">
                <template #icon><EditOutlined /></template>
              </a-button>
              <a-popconfirm
                title="确定要删除该节点吗？"
                ok-text="删除"
                cancel-text="取消"
                ok-type="danger"
                @confirm="removeServer(item.id)"
              >
                <a-button type="text" danger size="small" @click.stop>
                  <template #icon><DeleteOutlined /></template>
                </a-button>
              </a-popconfirm>
            </template>
            <a-list-item-meta>
              <template #title>
                <span class="server-host">
                  <CheckCircleFilled v-if="selectedId === item.id" class="selected-icon" />
                  {{ item.host }}
                </span>
              </template>
              <template #description>
                <span class="server-meta">{{ item.protocol }} | {{ item.username || "无认证" }}</span>
              </template>
            </a-list-item-meta>
          </a-list-item>
        </template>
      </a-list>
    </div>

    <div class="server-list-empty" v-else>
      <a-empty description="暂无节点，点击下方按钮添加" />
    </div>

    <a-button type="primary" block @click="openAddModal" class="add-server-btn">
      <template #icon><PlusOutlined /></template>
      添加节点
    </a-button>

    <a-modal
      v-model:open="modalVisible"
      :title="editingServer ? '编辑节点' : '添加节点'"
      :confirm-loading="submitLoading"
      @ok="handleSubmit"
      @cancel="handleCancel"
    >
      <a-form ref="formRef" :model="formState" :rules="formRules" layout="vertical">
        <a-form-item label="地址" name="host">
          <a-input v-model:value="formState.host" placeholder="127.0.0.1:1080" allow-clear />
        </a-form-item>

        <a-form-item label="协议" name="protocol">
          <a-select v-model:value="formState.protocol" placeholder="选择协议" allow-clear>
            <a-select-option value="socks5">SOCKS5</a-select-option>
            <a-select-option value="socks4">SOCKS4</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="用户名" name="username">
          <a-input
            v-model:value="formState.username"
            placeholder="可选，留空表示无认证"
            allow-clear
          />
        </a-form-item>

        <a-form-item label="密码" name="password">
          <a-input-password v-model:value="formState.password" placeholder="可选" allow-clear />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from "vue";
import { CheckCircleFilled, DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons-vue";
import { Get, Set } from "../../bindings/socks5-desktop/storage.js";
import { message } from "ant-design-vue";

const props = defineProps({
  selectedServer: { type: Object, default: null },
});

const emit = defineEmits(["select", "update:selectedServer"]);

const STORAGE_KEY = "servers";

const servers = ref([]);
const selectedId = computed(() => props.selectedServer?.id ?? null);
const modalVisible = ref(false);
const submitLoading = ref(false);
const formRef = ref(null);
const editingServer = ref(null);

const defaultFormState = () => ({
  host: "",
  username: "",
  password: "",
  protocol: "socks5",
});

const formState = reactive(defaultFormState());

const formRules = {
  host: [
    { required: true, message: "请输入地址" },
    { pattern: /^[^:]+:\d{1,5}$/, message: "格式为 host:port，如 127.0.0.1:1080" },
  ],
  protocol: [{ required: true, message: "请选择协议" }],
};

const loadServers = async () => {
  try {
    const data = await Get(STORAGE_KEY);
    servers.value = Array.isArray(data) ? data : [];
  } catch (e) {
    message.error("加载节点列表失败");
    servers.value = [];
  }
};

const saveServers = async () => {
  try {
    await Set(STORAGE_KEY, servers.value);
  } catch (e) {
    message.error("保存失败");
    throw e;
  }
};

const openAddModal = () => {
  editingServer.value = null;
  Object.assign(formState, defaultFormState());
  modalVisible.value = true;
};

const editServer = (item) => {
  editingServer.value = item;
  formState.host = item.host;
  formState.username = item.username || "";
  formState.password = item.password || "";
  formState.protocol = item.protocol || "socks5";
  modalVisible.value = true;
};

const handleSubmit = async () => {
  try {
    await formRef.value.validate();
    submitLoading.value = true;

    if (editingServer.value) {
      const idx = servers.value.findIndex((s) => s.id === editingServer.value.id);
      if (idx >= 0) {
        servers.value[idx] = {
          ...servers.value[idx],
          host: formState.host.trim(),
          username: formState.username?.trim() || "",
          password: formState.password || "",
          protocol: formState.protocol,
        };
      }
      message.success("修改成功");
    } else {
      const newServer = {
        id: Date.now().toString(),
        host: formState.host.trim(),
        username: formState.username?.trim() || "",
        password: formState.password || "",
        protocol: formState.protocol,
      };
      servers.value.push(newServer);
      message.success("添加成功");
    }

    await saveServers();
    modalVisible.value = false;
  } catch (e) {
    if (e?.errorFields) return;
    message.error(e?.message || "操作失败");
  } finally {
    submitLoading.value = false;
  }
};

const handleCancel = () => {
  modalVisible.value = false;
  formRef.value?.resetFields();
};

const removeServer = async (id) => {
  try {
    servers.value = servers.value.filter((s) => s.id !== id);
    if (props.selectedServer?.id === id) emit("update:selectedServer", null);
    await saveServers();
    message.success("已删除");
  } catch (e) {
    message.error("删除失败");
  }
};

const selectServer = (item) => {
  emit("update:selectedServer", item);
  emit("select", item);
};

onMounted(() => {
  loadServers();
});
</script>

<style scoped>
.socks5-server-card {
  min-height: 400px;
  flex: 1;
  display: flex;
  flex-direction: column;
  border-radius: 14px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.04);
}

.socks5-server-card :deep(.ant-card-body) {
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

.socks5-server-card :deep(.ant-list-item) {
  cursor: pointer;
}

.socks5-server-card :deep(.ant-list-item.active) {
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
</style>
