import { ref, reactive, onMounted } from "vue";
import { Modal, message } from "ant-design-vue";
import { ProxyScheme } from "@bindings/desktop/internal/appconst";
import { t } from "@/locale";
import { withRuntimeIds } from "@/utils";

/**
 * 新增/编辑弹窗、删除确认；selectedServer 为 store 中当前选中项的 ref。
 */
export function useServerEditor(servers, latencyById, persistServers, selectedServer) {
  const editRef = ref();
  const proxySchemes = ref([]);

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
        await persistServers();
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
        await persistServers();
        message.success(t("serverList.deleteSuccess"));
      },
    });
  }

  onMounted(async () => {
    proxySchemes.value = await ProxyScheme();
  });

  return { edit, editRef, deleteModal, proxySchemes };
}
