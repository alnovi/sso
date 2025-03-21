<script setup>
import {config, validMsg, validStatus} from "../../../services/utils.js";
import {Password} from "@vicons/carbon";
import {ref} from "vue";
import {useApi} from "../../../services/api.js";
import {notifyError, notifyInfo} from "../../../services/notify.js";
import {useNotification} from "naive-ui";

const api = useApi(config('VITE_API_HOST', '/'))
const notification = useNotification()

const formData = ref({
  old_password: "",
  new_password: '',
  confirm_password: '',
})

const formErr = ref({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

const formDataReset = () => {
  formData.value.old_password = ''
  formData.value.new_password = ''
  formData.value.confirm_password = ''
}

const formErrReset = () => {
  formErr.value.old_password = ''
  formErr.value.new_password = ''
  formErr.value.confirm_password = ''
}

const formIsEmpty = () => {
  return formData.value.old_password.length < 5 ||
    formData.value.new_password.length < 5 ||
    formData.value.confirm_password.length < 5
}

const submitForm = () => {
  formErrReset()

  if (formData.value.new_password !== formData.value.confirm_password) {
    formErr.value.new_password = "пароли не совпадают"
    formErr.value.confirm_password = "пароли не совпадают"
    return false
  }

  const postData = {
    old_password: formData.value.old_password,
    new_password: formData.value.new_password,
  }

  api.put(`/profile/password`, postData)
    .then(() => {
      formDataReset()
      notification.success(notifyInfo('Пароль успешно изменен'))
    })
    .catch((err) => {
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
</script>

<template>
  <div class="block-layout">
    <div class="block-layout-header">
      <div class="block-layout-header__title">Безопасность</div>
      <div class="separator"></div>
      <div class="block-layout-header__description">Изменение пароля доступа.</div>
    </div>
    <div class="block-layout-content">
      <n-form>
        <n-form-item label="Текущий пароль" path="password" required :feedback="validMsg(formErr.old_password, 'old_password', 'пароль')" :validation-status="validStatus(formErr.old_password)">
          <n-input size="large" v-model:value="formData.old_password" type="password" show-password-on="mousedown" placeholder="Текущий пароль">
            <template #prefix>
              <n-icon :component="Password"/>
            </template>
          </n-input>
        </n-form-item>
        <n-form-item label="Новый пароль" path="new_password" required :feedback="validMsg(formErr.new_password, 'new_password', 'пароль')" :validation-status="validStatus(formErr.new_password)">
          <n-input size="large" v-model:value="formData.new_password" type="password" show-password-on="mousedown" placeholder="Новый пароль">
            <template #prefix>
              <n-icon :component="Password"/>
            </template>
          </n-input>
        </n-form-item>
        <n-form-item label="Повторите пароль" path="confirm_password" required :feedback="validMsg(formErr.confirm_password, 'confirm_password', 'пароль')" :validation-status="validStatus(formErr.confirm_password)">
          <n-input size="large" v-model:value="formData.confirm_password" type="password" show-password-on="mousedown" placeholder="Повторите пароль">
            <template #prefix>
              <n-icon :component="Password"/>
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
