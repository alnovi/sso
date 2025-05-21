<script setup>
import {ref} from "vue";
import {useNotification} from "naive-ui";
import {useRouter} from "vue-router"
import {Password, User} from "@vicons/carbon";
import {useApi} from "../../../services/api.js";
import {config, meta, validStatus, validMsg} from "../../../services/utils.js";
import {notifyError, notifyInfo} from "../../../services/notify.js";

const api = useApi(config('VITE_API_HOST', '/'))
const query = meta('auth-query', config('VITE_AUTH_QUERY'))
const hash = (new URLSearchParams(window.location.search)).get('hash')
const router = useRouter()
const notification = useNotification()

const formRef = ref(null);

const formValue = ref({
  password: '',
  passwordConfirmation: '',
})

const formError = ref({
  password: null,
})

const formIsValid = () => {
  return formValue.value.password.length >= 5 && formValue.value.password === formValue.value.passwordConfirmation
}

async function reset() {
  formError.value.password = null

  const data = {
    token: hash,
    password: formValue.value.password,
  }

  api.post(`oauth/reset-password`, data)
    .then(res => {
      notification.success(notifyInfo('Пароль успешно изменен'))
      router.push(`/oauth/authorize?${query}`)
    })
    .catch(error => {
      if (error.code === 'ERR_NETWORK') {
        notification.error(notifyError('Сервер не доступен'))
        return
      }
      if (!!error.response.data && !!error.response.data.error) {
        notification.error(notifyError(error.response.data.error))
      }
      if (error.response.status === 422) {
        formError.value = error.response.data.validate
      }
    })
}
</script>

<template>
  <n-card title="Смена пароля" bordered :segmented="{content: true, footer: 'soft'}">
    <n-form :ref="formRef" :label-width="80" :model="formValue">
      <n-form-item label="Новый пароль" path="password" required :feedback="validMsg(formError.password, 'password', 'пароль')" :validation-status="validStatus(formError.password)">
        <n-input size="large" v-model:value="formValue.password" type="password" show-password-on="mousedown" placeholder="Новый пароль">
          <template #prefix>
            <n-icon :component="Password"/>
          </template>
        </n-input>
      </n-form-item>
      <n-form-item label="Повторите пароль" path="password" required :feedback="validMsg(formError.password, 'password', 'пароль')" :validation-status="validStatus(formError.password)">
        <n-input size="large" v-model:value="formValue.passwordConfirmation" type="password" show-password-on="mousedown" placeholder="Повторите пароль">
          <template #prefix>
            <n-icon :component="Password"/>
          </template>
        </n-input>
      </n-form-item>
    </n-form>
    <template #footer>
      <n-flex justify="space-between">
        <n-button text @click="router.push(`/oauth/authorize?${query}`)">Войти с паролем</n-button>
        <n-button @click="reset" :disabled="!formIsValid()" size="large" type="primary" style="width: 150px">
          Сохранить
        </n-button>
      </n-flex>
    </template>
  </n-card>
</template>

<style scoped>
.n-card {
  box-shadow: 0 10px 20px 0 rgba(0, 0, 0, .2);
  max-width: 500px;
  border-radius: 12px;
}
</style>
