<script setup>
import {onActivated, ref} from "vue"
import {useRouter} from "vue-router";
import {notifyError, notifyInfo} from "../../../services/notify.js";
import {NButton, NFlex, NIcon, useDialog, useNotification} from "naive-ui";
import {Delete, Search} from "@vicons/carbon";
import {useApi} from "../../../services/api.js";
import {config} from "../../../services/utils.js";
import moment from "moment/moment";

const api = useApi(config('VITE_API_HOST', '/'))
const dialog = useDialog();
const notification = useNotification()
const router = useRouter()

const sessions = ref([])
const columns = [
  {
    title: "Пользователь",
    key: "user.name",
    resizable: true,
    minWidth: 200,
    render: (row) => h('div', {innerHTML: `${row.user.name} <span class="muted">${row.user.email}</span>`})
  }, {
    title: "IP",
    key: "ip",
    width: 200,
  }, {
    title: "OS",
    key: "os",
    width: 200,
  }, {
    title: "App",
    key: "app",
    width: 200,
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
      h(NButton, {type: "success", strong: true, secondary: true, circle: true, onClick: () => {router.push({name: 'session', params: {id: row.id}})}}, () => h(NIcon, {component: Search})),
      h(NButton, {type: "error", strong: true, secondary: true, circle: true, disabled: row.is_current, onClick: () => deleteSession(row.id)}, () => h(NIcon, {component: Delete}))
    ]),
  }
]

const loadSessions = async () => {
  api.get("/api/sessions")
    .then(res => {
      sessions.value = Array.from(res.data || []).map((session) => {
        return {
          "id": session.id,
          "ip": session.ip,
          "app": session.app,
          "os": session.os,
          "agent": session.agent,
          "is_current": session.is_current,
          "created_at": moment(session.created_at).format("DD.MM.YYYY HH:mm"),
          "updated_at": moment(session.updated_at).format("DD.MM.YYYY HH:mm"),
          "user": {
            "id": session.user.id,
            "name": session.user.name,
            "email": session.user.email,
            "created_at": moment(session.user.created_at).format("DD.MM.YYYY HH:mm"),
            "updated_at": moment(session.user.updated_at).format("DD.MM.YYYY HH:mm"),
            "deleted_at": session.user.deleted_at ? moment(session.user.deleted_at).format("DD.MM.YYYY HH:mm") : null,
          },
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

const deleteSession = (id) => {
  dialog.warning({
    title: "Внимание",
    content: `Удалить сессию?`,
    positiveText: "Удалить",
    negativeText: "Отмена",
    draggable: false,
    onPositiveClick: () => {
      api.delete(`/api/sessions/${id}`)
        .then(() => {
          notification.info(notifyInfo(`Сессия удалена.`))
          loadSessions()
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
  loadSessions()
})
</script>

<template>
  <n-breadcrumb style="margin-bottom: 24px">
    <n-breadcrumb-item @click="router.push({name: 'home'})">Главная</n-breadcrumb-item>
    <n-breadcrumb-item>Устройства</n-breadcrumb-item>
  </n-breadcrumb>
  <n-data-table
    :columns="columns"
    :data="sessions"
    :pagination="false"
    :bordered="true"
  />
</template>

<style lang="scss">
.muted {
  padding: 0 5px;
  opacity: 0.6;
}
</style>
