<script setup>
import {defineProps, onActivated, onBeforeMount, onDeactivated, ref} from "vue"
import {NButton, NFlex, NIcon, NImage, NSelect, useNotification} from "naive-ui";
import {useRouter} from "vue-router";
import {useApi} from "../../../services/api.js";
import {config, validMsg, validStatus} from "../../../services/utils.js";
import {notifyError, notifyInfo} from "../../../services/notify.js";
import {Checkmark, Close} from "@vicons/carbon";

const props = defineProps({id: String})

const api = useApi(config('VITE_API_HOST', '/'))
const notification = useNotification()
const router = useRouter()

const user = ref({})
const formRef = ref(null);
const formData = ref({})
const formErr = ref({})
const clients = ref([])
const roleOptions = [
  {
    label: "Гость",
    value: "guest"
  }, {
    label: "Пользователь",
    value: "user"
  }, {
    label: "Менеджер",
    value: "manager"
  }, {
    label: "Администратор",
    value: "admin"
  }
]
const columns = [
  {
    title: "",
    key: "status",
    align: "center",
    width: 50,
    render: (row) => {
      return row.deleted_at
        ? h(NIcon, {size: 20, color: 'red'}, {default: () => h(Close)})
        : h(NIcon, {size: 20, color: 'green'}, {default: () => h(Checkmark)})
    }
  },
  {
    title: "",
    key: "icon",
    width: 50,
    render: (row) => h(NImage, {src: row.icon, width: 30, height: 30}),
  },
  {
    title: "Приложение",
    key: "name",
    minWidth: 200,
  }, {
    title: "Роль",
    key: "role",
    width: 300,
    render: (row) => {
      return h(
        NSelect,
        {
          options: roleOptions,
          placeholder: "Нет доступа",
          clearable: true,
          value: row.role,
          onUpdateValue: (v) => {
            row.role = v
            updateRole(props.id, row.id, row.role)
          }
        }, {
          default: () => null
        })
    }
  }
]

const loadUser = async () => {
  api.get(`/api/users/${props.id}`)
    .then(res => {
      user.value = res.data
      formData.value = {
        name: res.data.name,
        email: res.data.email,
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

const loadClients = async () => {
  api.get(`/api/users/${props.id}/clients`)
    .then(res => {
      clients.value = res.data
    })
    .catch(err => {
      if (!!err.response.data && !!err.response.data.error) {
        notification.error(notifyError(err.response.data.error))
      }
    })
}

const updateRole = async (userId, clientId, role) => {
  const postData = {
    role: role ? role : null,
  }
  api.post(`/api/users/${userId}/clients/${clientId}`, postData)
    .then(res => {
      notification.success(notifyInfo('Роль обновлена'))
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

const submitForm = () => {
  formErr.value = {}

  const postData = {
    name: formData.value.name,
    email: formData.value.email,
    password: formData.value.password ? formData.value.password : null,
  }

  api.put(`/api/users/${props.id}`, postData)
    .then(res => {
      user.value = res.data
      formData.value.password = ''
      notification.success(notifyInfo('Данные обновлены'))
    })
    .catch(err => {
      if (err.code === 'ERR_NETWORK') {
        notification.error(notifyError('Сервер не доступен'))
      }
      if (!!err.response.data && !!err.response.data.error) {
        notification.error(notifyError(err.response.data.error))
      }
      if (err.response.status === 422) {
        formErr.value = err.response.data.validate
      }
    })
}

const restoreUser = async () => {
  api.post(`/api/users/${props.id}/restore`, null)
    .then(res => {
      user.value = res.data
      notification.success(notifyInfo('Пользователь восстановлен'))
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

onActivated(() => {
  loadUser()
  loadClients()
})

onDeactivated(() => {
  user.value = {};
  formData.value = {}
  formErr.value = {}
})

onBeforeMount(() => {
  loadUser()
  loadClients()
})
</script>

<template>
  <n-breadcrumb style="margin-bottom: 24px">
    <n-breadcrumb-item @click="router.push({name: 'home'})">Главная</n-breadcrumb-item>
    <n-breadcrumb-item @click="router.push({name: 'users'})">Пользователи</n-breadcrumb-item>
    <n-breadcrumb-item>{{ user.name }}</n-breadcrumb-item>
  </n-breadcrumb>
  <n-card bordered :segmented="{content: true, footer: 'soft'}">
    <n-tabs type="line" default-value="user" animated>
      <n-tab-pane name="user" tab="Пользователь">
        <n-form :ref="formRef" :label-width="80" :model="formData">
          <n-form-item label="Имя" path="name" required :feedback="validMsg(formErr.name, 'name', 'имя')"
                       :validation-status="validStatus(formErr.name)">
            <n-input size="large" maxlength="100" show-count clearable v-model:value="formData.name" type="text"
                     placeholder="Полное имя"></n-input>
          </n-form-item>
          <n-form-item label="Email" path="email" required :feedback="validMsg(formErr.email, 'email', 'email')"
                       :validation-status="validStatus(formErr.email)">
            <n-input size="large" maxlength="100" show-count clearable v-model:value="formData.email" type="text"
                     placeholder="Email"></n-input>
          </n-form-item>
          <n-form-item label="Пароль" path="password" :feedback="validMsg(formErr.password, 'password', 'пароль')"
                       :validation-status="validStatus(formErr.password)">
            <n-input size="large" v-model:value="formData.password" type="password" show-password-on="mousedown"
                     placeholder="Новый пароль"></n-input>
          </n-form-item>
        </n-form>
        <n-flex class="tab-actions" justify="space-between">
          <div>
            <n-button v-if="user.deleted_at" tertiary @click="restoreUser">
              Восстановить
            </n-button>
          </div>
          <div>
            <n-button tertiary style="width: 100px; margin-right: 10px" @click="router.push({name: 'users'})">
              Отмена
            </n-button>
            <n-button strong secondary type="primary" style="width: 100px" @click="submitForm">
              Сохранить
            </n-button>
          </div>
        </n-flex>
      </n-tab-pane>
      <n-tab-pane name="secure" tab="Разрешения">
        <n-data-table
          :columns="columns"
          :data="clients"
          :pagination="false"
          :bordered="false"
        />
      </n-tab-pane>
    </n-tabs>
  </n-card>
</template>

<style scoped lang="scss">
.tab-actions {
  border-top: 1px solid var(--n-border-color);
  padding-top: var(--n-padding-bottom);
}
</style>
