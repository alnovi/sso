<script setup>
import {onDeactivated, ref} from "vue"
import {useNotification} from "naive-ui";
import {useRouter} from "vue-router";
import {useApi} from "../../../services/api.js";
import {config, validMsg, validStatus} from "../../../services/utils.js";
import {notifyError, notifyInfo} from "../../../services/notify.js";

const api = useApi(config('VITE_API_HOST', '/'))
const notification = useNotification()
const router = useRouter()

const formRef = ref(null);
const formData = ref({
  name: '',
  email: '',
  password: '',
})
const formErr = ref({})

const submitForm = () => {
  formErr.value = {}

  api.post(`/api/users`, formData.value)
    .then(res => {
      formData.value = res.data
      notification.success(notifyInfo('Пользователь добавлен'))
      router.push({name: 'edit-user', params: {id: res.data.id}})
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

onDeactivated(() => {
  formErr.value = {}
  formData.value = {
    name: '',
    email: '',
    password: '',
  }
})
</script>

<template>
  <n-breadcrumb style="margin-bottom: 24px">
    <n-breadcrumb-item @click="router.push({name: 'home'})">Главная</n-breadcrumb-item>
    <n-breadcrumb-item @click="router.push({name: 'users'})">Пользователи</n-breadcrumb-item>
    <n-breadcrumb-item>Новый пользователь</n-breadcrumb-item>
  </n-breadcrumb>
  <n-card bordered :segmented="{content: true, footer: 'soft'}">
    <n-form :ref="formRef" :label-width="80" :model="formData">
      <n-form-item label="Имя" path="name" required :feedback="validMsg(formErr.name, 'name', 'имя')" :validation-status="validStatus(formErr.name)">
        <n-input size="large" maxlength="100" show-count clearable v-model:value="formData.name" type="text" placeholder="Полное имя"></n-input>
      </n-form-item>
      <n-form-item label="Email" path="email" required :feedback="validMsg(formErr.email, 'email', 'email')" :validation-status="validStatus(formErr.email)">
        <n-input size="large" maxlength="100" show-count clearable v-model:value="formData.email" type="text" placeholder="Email"></n-input>
      </n-form-item>
      <n-form-item label="Пароль" path="password" required :feedback="validMsg(formErr.password, 'password', 'пароль')" :validation-status="validStatus(formErr.password)">
        <n-input size="large" v-model:value="formData.password" type="password" show-password-on="mousedown" placeholder="Пароль"></n-input>
      </n-form-item>
    </n-form>
    <template #footer>
      <n-flex justify="end">
        <n-button tertiary style="width: 100px" @click="router.push({name: 'users'})">
          Отмена
        </n-button>
        <n-button strong secondary type="primary" style="width: 100px" @click="submitForm">
          Сохранить
        </n-button>
      </n-flex>
    </template>
  </n-card>
</template>
