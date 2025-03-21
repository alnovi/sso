<script setup>
import {onActivated, onBeforeMount, ref} from "vue"
import {NIcon, NImage, NFlex, NButton, useNotification, useDialog} from "naive-ui";
import {Checkmark, Close, Pen, Delete} from "@vicons/carbon"
import {useRouter} from "vue-router";
import {useApi} from "../../../services/api.js";
import {config} from "../../../services/utils.js";
import {notifyError, notifyInfo} from "../../../services/notify.js";
import moment from "moment";

const api = useApi(config('VITE_API_HOST', '/'))
const dialog = useDialog();
const notification = useNotification()
const router = useRouter()

const clients = ref([])
const columns = [
  {
    title: "",
    key: "status",
    align: "center",
    width: 50,
    render: (row) => {
      return row.deleted_at
        ? h(NIcon, {size: 20, color: 'red'}, { default: () => h(Close) })
        : h(NIcon, {size: 20, color: 'green'}, { default: () => h(Checkmark) })
    }
  },
  {
    title: "",
    key: "icon",
    width: 50,
    render: (row) => h(NImage, {src: row.icon, width: 30, height: 30}),
  },
  {
    title: "Название",
    key: "name",
    resizable: true,
    minWidth: 200,
  }, {
    title: "ID",
    key: "id",
    resizable: true,
    minWidth: 200,
  }, {
    title: "Secret",
    key: "secret",
    ellipsis: true,
    resizable: true,
    minWidth: 200,
  }, {
    title: "Даты",
    key: "date_at",
    width: 140,
    render: (row) => h('div', {innerHTML: `${row.created_at}<br/>${row.updated_at}`}),
  }, {
    title: "Действия",
    key: "action",
    width: 150,
    render: (row) => h(NFlex, {align: "center", justify: "center"}, () => [
      h(NButton, {type: "success", strong: true, secondary: true, circle: true, onClick: () => {router.push({name: 'edit-client', params: {id: row.id}})}}, () => h(NIcon, {component: Pen})),
      h(NButton, {type: "error", strong: true, secondary: true, circle: true, disabled: row.is_system, onClick: () => deleteClient(row)}, () => h(NIcon, {component: Delete}))
    ]),
  }
]

const loadClients = async () => {
  api.get("/api/clients")
    .then(res => {
      clients.value = Array.from(res.data || []).map((client) => {
        return {
          "id": client.id,
          "name": client.name,
          "icon": client.icon || '/public/app.png',
          "secret": client.secret,
          "callback": client.callback,
          "is_system": client.is_system,
          "created_at": moment(client.created_at).format("DD.MM.YYYY HH:mm"),
          "updated_at": moment(client.updated_at).format("DD.MM.YYYY HH:mm"),
          "deleted_at": client.deleted_at ? moment(client.deleted_at).format("DD.MM.YYYY HH:mm") : null,
        }
      })
    })
    .catch(err => {
      if (err.code === 'ERR_NETWORK') {
        notification.error(notifyError('Сервер не доступен'))
      }
      if (!!err.response.data && !!err.response.data.error) {
        notification.error(notifyError(err.response.data.error))
      }
    })
}

const deleteClient = async (client) => {
  dialog.warning({
    title: "Внимание",
    content: `Удалить приложение ${client.name} (id: ${client.id})?`,
    positiveText: "Удалить",
    negativeText: "Отмена",
    draggable: false,
    onPositiveClick: () => {
      api.delete(`/api/clients/${client.id}`)
        .then(() => {
          loadClients()
          notification.info(notifyInfo(`Приложение ${client.name} удалено.`))
        })
        .catch(err => {
          if (err.code === 'ERR_NETWORK') {
            notification.error(notifyError('Сервер не доступен'))
          }
          if (!!err.response.data && !!err.response.data.error) {
            notification.error(notifyError(err.response.data.error))
          }
        })
    },
  });
}

onActivated(() => {
  loadClients()
})

onBeforeMount(() => {
  //loadClients()
})
</script>

<template>
  <n-flex align="center" justify="space-between" style="margin-bottom: 24px">
    <n-breadcrumb>
      <n-breadcrumb-item @click="router.push({name: 'home'})">Главная</n-breadcrumb-item>
      <n-breadcrumb-item>Приложения</n-breadcrumb-item>
    </n-breadcrumb>
    <n-flex align="center" justify="center">
      <n-button tertiary @click="router.push({name: 'create-client'})">
        Новое приложение
      </n-button>
    </n-flex>
  </n-flex>
  <n-data-table
    :columns="columns"
    :data="clients"
    :pagination="false"
    :bordered="true"
  />
</template>
