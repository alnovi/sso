<script setup>
import {defineProps, onActivated, onBeforeMount, onDeactivated, ref} from "vue"
import {useNotification} from "naive-ui";
import {useRouter} from "vue-router";
import {useApi} from "../../../services/api.js";
import {config, validMsg, validStatus} from "../../../services/utils.js";
import {notifyError, notifyInfo} from "../../../services/notify.js";

const props = defineProps({id: String})

const api = useApi(config('VITE_API_HOST', '/'))
const notification = useNotification()
const router = useRouter()

const formRef = ref(null);
const client = ref({})
const formData = ref({})
const formErr = ref({})

const loadClient = async () => {
  api.get(`/api/clients/${props.id}`)
    .then(res => {
      client.value = res.data
      formData.value = res.data
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
    icon: formData.value.icon ? formData.value.icon : null,
    callback: formData.value.callback,
    secret: formData.value.secret,
  }

  api.put(`/api/clients/${client.value.id}`, postData)
    .then(res => {
      client.value = res.data
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

const restoreClient = async () => {
  api.post(`/api/clients/${props.id}/restore`, null)
    .then(res => {
      client.value = res.data
      formData.value = res.data
      notification.success(notifyInfo('Приложение восстановлено'))
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
  loadClient()
})

onDeactivated(() => {
  client.value = {};
  formData.value = {}
  formErr.value = {}
})

onBeforeMount(() => {
  loadClient()
})
</script>

<template>
  <n-breadcrumb style="margin-bottom: 24px">
    <n-breadcrumb-item @click="router.push({name: 'home'})">Главная</n-breadcrumb-item>
    <n-breadcrumb-item @click="router.push({name: 'clients'})">Приложения</n-breadcrumb-item>
    <n-breadcrumb-item>{{ client.name || props.id }}</n-breadcrumb-item>
  </n-breadcrumb>
  <n-card bordered :segmented="{content: true, footer: 'soft'}">
    <div class="client-form-layout">
      <div class="client-form-layout__icon">
        <n-image :src="formData.icon || '/public/app.png'" width="120" height="120" alt=""></n-image>
      </div>
      <div class="client-form-layout__form">
        <n-form :ref="formRef" :label-width="80" :model="formData">
          <n-form-item label="ClientID" path="id" required>
            <n-input size="large" maxlength="50" disabled v-model:value="formData.id" type="text"></n-input>
          </n-form-item>
          <n-form-item label="Название" path="name" required :feedback="validMsg(formErr.name, 'name', 'название')"
                       :validation-status="validStatus(formErr.name)">
            <n-input size="large" maxlength="50" show-count clearable v-model:value="formData.name" type="text"
                     placeholder="Название"></n-input>
          </n-form-item>
          <n-form-item label="Иконка" path="icon" :feedback="validMsg(formErr.icon, 'icon', 'иконка')"
                       :validation-status="validStatus(formErr.icon)">
            <n-input size="large" maxlength="250" show-count clearable v-model:value="formData.icon" type="text"
                     placeholder="Ссылка на иконку"></n-input>
          </n-form-item>
          <n-form-item label="Callback" path="callback" required
                       :feedback="validMsg(formErr.callback, 'callback', 'callback')"
                       :validation-status="validStatus(formErr.callback)">
            <n-input size="large" maxlength="250" show-count clearable v-model:value="formData.callback" type="text"
                     placeholder="Callback"></n-input>
          </n-form-item>
          <n-form-item label="Secret" path="secret" required :feedback="validMsg(formErr.secret, 'secret', 'secret')"
                       :validation-status="validStatus(formErr.secret)">
            <n-input size="large" maxlength="100" show-count clearable v-model:value="formData.secret" type="text"
                     placeholder="Secret"></n-input>
          </n-form-item>
        </n-form>
      </div>
    </div>
    <template #footer>
      <n-flex justify="space-between">
        <div>
          <n-button v-if="client.deleted_at" tertiary @click="restoreClient">
            Восстановить
          </n-button>
        </div>
        <div>
          <n-button tertiary style="width: 100px; margin-right: 10px" @click="router.push({name: 'clients'})">
            Отмена
          </n-button>
          <n-button strong secondary type="primary" style="width: 100px" @click="submitForm">
            Сохранить
          </n-button>
        </div>
      </n-flex>
    </template>
  </n-card>
</template>
