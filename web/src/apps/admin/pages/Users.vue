<script setup>
import {onActivated, onBeforeMount, ref} from "vue"
import {NIcon, NFlex, NButton, useNotification, useDialog} from "naive-ui";
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

const users = ref([])
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
    title: "Имя",
    key: "name",
    resizable: true,
    minWidth: 200,
  }, {
    title: "Email",
    key: "email",
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
      h(NButton, {type: "success", strong: true, secondary: true, circle: true, onClick: () => {router.push({name: 'edit-user', params: {id: row.id}})}}, () => h(NIcon, {component: Pen})),
      h(NButton, {type: "error", strong: true, secondary: true, circle: true, disabled: row.is_system, onClick: () => deleteUser(row)}, () => h(NIcon, {component: Delete}))
    ]),
  }
]

const loadUsers = async () => {
  api.get("/api/users")
    .then(res => {
      users.value = Array.from(res.data || []).map((user) => {
        return {
          "id": user.id,
          "name": user.name,
          "email": user.email,
          "created_at": moment(user.created_at).format("DD.MM.YYYY HH:mm"),
          "updated_at": moment(user.updated_at).format("DD.MM.YYYY HH:mm"),
          "deleted_at": user.deleted_at ? moment(user.deleted_at).format("DD.MM.YYYY HH:mm") : null,
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

const deleteUser = async (user) => {
  dialog.warning({
    title: "Внимание",
    content: `Удалить пользователя ${user.name}?`,
    positiveText: "Удалить",
    negativeText: "Отмена",
    draggable: false,
    onPositiveClick: () => {
      api.delete(`/api/users/${user.id}`)
        .then(() => {
          loadUsers()
          notification.info(notifyInfo(`Пользователь ${user.name} удален.`))
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
  loadUsers()
})

onBeforeMount(() => {
  //loadUsers()
})
</script>

<template>
  <n-flex align="center" justify="space-between" style="margin-bottom: 24px">
    <n-breadcrumb>
      <n-breadcrumb-item @click="router.push({name: 'home'})">Главная</n-breadcrumb-item>
      <n-breadcrumb-item>Пользователи</n-breadcrumb-item>
    </n-breadcrumb>
    <n-flex align="center" justify="center">
      <n-button tertiary @click="router.push({name: 'create-user'})">
        Новый пользователь
      </n-button>
    </n-flex>
  </n-flex>
  <n-data-table
    :columns="columns"
    :data="users"
    :pagination="false"
    :bordered="true"
  />
</template>
