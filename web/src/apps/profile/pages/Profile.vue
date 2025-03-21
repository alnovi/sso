<script setup>
import {onMounted, ref} from "vue";
import {config, validMsg, validStatus} from "../../../services/utils.js";
import {Email, User} from "@vicons/carbon";
import {notifyError, notifyInfo} from "../../../services/notify.js";
import {useApi} from "../../../services/api.js";
import {useNotification} from "naive-ui";

const api = useApi(config('VITE_API_HOST', '/'))
const notification = useNotification()

const formData = ref({
  name: '',
  email: '',
})

const formErr = ref({
  name: '',
  email: '',
})

const formErrReset = () => {
  formErr.value.name = ''
  formErr.value.email = ''
}

const formIsEmpty = () => {
  return formData.value.name.length < 3 || formData.value.email.length < 5
}

const submitForm = () => {
  formErrReset()
  api.put(`/profile/me`, formData.value)
    .then((res) => {
      notification.success(notifyInfo('Данные успешно изменены'))
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

onMounted(() => {
  api.get(`/profile/me`)
    .then(res => {
      formData.value.name = res.data.name
      formData.value.email = res.data.email
    })
    .catch(err => {
      if (err.code === 'ERR_NETWORK') {
        notification.error(notifyError('Сервер не доступен'))
      }
    })
})
</script>

<template>
  <div class="block-layout">
    <div class="block-layout-header">
      <div class="block-layout-header__title">Мой профиль</div>
      <div class="separator"></div>
      <div class="block-layout-header__description">Управление личной информацией.</div>
    </div>
    <div class="block-layout-content">
      <n-form>
        <n-form-item label="Имя" path="name" required :feedback="validMsg(formErr.name, 'name', 'имя')" :validation-status="validStatus(formErr.name)">
          <n-input size="large" v-model:value="formData.name" type="text" placeholder="Имя">
            <template #prefix>
              <n-icon :component="User"/>
            </template>
          </n-input>
        </n-form-item>
        <n-form-item label="Email" path="email" required :feedback="validMsg(formErr.email, 'email', 'email')" :validation-status="validStatus(formErr.email)">
          <n-input size="large" v-model:value="formData.email" type="text" placeholder="Email">
            <template #prefix>
              <n-icon :component="Email"/>
            </template>
          </n-input>
        </n-form-item>
      </n-form>
      <n-flex justify="end">
        <n-button @click="submitForm" :disabled="formIsEmpty()" size="large" type="primary" style="width: 150px">
          Изменить
        </n-button>
      </n-flex>
    </div>
  </div>
</template>
