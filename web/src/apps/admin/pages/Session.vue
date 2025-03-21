<script setup>
import {defineProps, onActivated, ref} from "vue";
import {useRouter} from "vue-router";
import {notifyError, notifyInfo} from "../../../services/notify.js";
import {useApi} from "../../../services/api.js";
import {config} from "../../../services/utils.js";
import {useDialog, useNotification} from "naive-ui";
import moment from "moment";

const props = defineProps({id: String})

const api = useApi(config('VITE_API_HOST', '/'))
const dialog = useDialog();
const notification = useNotification()
const router = useRouter()

const session = ref({})

const loadSession = async () => {
  api.get(`/api/sessions/${props.id}`)
    .then(res => {
      session.value = {
        id: res.data.id,
        ip: res.data.ip,
        app: res.data.app,
        os: res.data.os,
        agent: res.data.agent,
        is_current: res.data.is_current,
        created_at: moment(res.data.created_at).format("DD.MM.YYYY HH:mm"),
        updated_at: moment(res.data.updated_at).format("DD.MM.YYYY HH:mm"),
        user: {
          id: res.data.user.id,
          name: res.data.user.name,
          email: res.data.user.email,
          created_at: moment(res.data.user.created_at).format("DD.MM.YYYY HH:mm"),
          updated_at: moment(res.data.user.updated_at).format("DD.MM.YYYY HH:mm"),
          deleted_at: res.data.user.deleted_at ? moment(res.data.user.deleted_at).format("DD.MM.YYYY HH:mm") : null,
        },
      }
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

const deleteSession = async () => {
  dialog.warning({
    title: "Внимание",
    content: `Удалить сессию?`,
    positiveText: "Удалить",
    negativeText: "Отмена",
    draggable: false,
    onPositiveClick: () => {
      api.delete(`/api/sessions/${props.id}`)
        .then(() => {
          notification.info(notifyInfo(`Сессия удалена.`))
          router.push({name: 'sessions'})
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
  loadSession()
})
</script>

<template>
  <n-breadcrumb style="margin-bottom: 24px">
    <n-breadcrumb-item @click="router.push({name: 'home'})">Главная</n-breadcrumb-item>
    <n-breadcrumb-item @click="router.push({name: 'sessions'})">Устройства</n-breadcrumb-item>
    <n-breadcrumb-item>{{ session.app }} ({{ session.os }})</n-breadcrumb-item>
  </n-breadcrumb>
  <n-card bordered :segmented="{content: true, footer: 'soft'}">
    <n-flex vertical>
      <n-flex justify="start">
        <div class="session-prop">ID</div>
        <div class="session-value">{{ session.id }}</div>
      </n-flex>
      <n-flex justify="start">
        <div class="session-prop">Пользователь</div>
        <div class="session-value">{{ session.user ? session.user.name : '' }}</div>
      </n-flex>
      <n-flex justify="start">
        <div class="session-prop">Email</div>
        <div class="session-value">{{ session.user ? session.user.email : '' }}</div>
      </n-flex>
      <n-flex justify="start">
        <div class="session-prop">IP-адрес</div>
        <div class="session-value">{{ session.ip }}</div>
      </n-flex>
      <n-flex justify="start">
        <div class="session-prop">Операционная система</div>
        <div class="session-value">{{ session.os }}</div>
      </n-flex>
      <n-flex justify="start">
        <div class="session-prop">Приложение</div>
        <div class="session-value">{{ session.app }}</div>
      </n-flex>
      <n-flex justify="start">
        <div class="session-prop">Агент</div>
        <div class="session-value">{{ session.agent }}</div>
      </n-flex>
      <n-flex justify="start">
        <div class="session-prop">Дата создания</div>
        <div class="session-value">{{ session.created_at }}</div>
      </n-flex>
      <n-flex justify="start">
        <div class="session-prop">Дата обновления</div>
        <div class="session-value">{{ session.updated_at }}</div>
      </n-flex>
    </n-flex>
    <template #footer>
      <n-flex justify="space-between">
        <n-button tertiary style="width: 100px" @click="router.push({name: 'sessions'})">
          Назад
        </n-button>
        <n-button :disabled="session.is_current" strong secondary type="error" style="width: 100px" @click="deleteSession">
          Удалить
        </n-button>
      </n-flex>
    </template>
  </n-card>
</template>

<style scoped lang="scss">
.session-prop {
  width: 200px;
  text-align: right;
  font-weight: 600;
  &:after {
    content: ":";
  }
}
</style>
