<script setup>
import {ref} from "vue";
import {useNotification} from "naive-ui";
import {useRouter} from "vue-router"
import {Password, User} from "@vicons/carbon";
import {useApi} from "../../../services/api.js";
import {config, meta, validStatus, validMsg} from "../../../services/utils.js";
import {notifyError} from "../../../services/notify.js";

const api = useApi(config('VITE_API_HOST', '/'))
const query = meta('auth-query', config('VITE_AUTH_QUERY'))
const router = useRouter()
const notification = useNotification()

const formRef = ref(null);

const formValue = ref({
  login: '',
  password: '',
  remember: false,
})

const formError = ref({
  login: null,
  password: null,
})

const formIsEmpty = () => {
  return formValue.value.login.length < 5 || formValue.value.password.length < 5
}

async function authorize() {
  formError.value.login = null
  formError.value.password = null

  const data = {
    login: formValue.value.login,
    password: formValue.value.password,
    remember: formValue.value.remember,
  }

  api.post(`oauth/authorize?${query}`, data)
    .then(res => {
      window.location.replace(res.data.url)
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
  <n-card title="Авторизация" bordered :segmented="{content: true, footer: 'soft'}">
    <n-form :ref="formRef" :label-width="80" :model="formValue">
      <n-form-item label="Логин" path="login" required :feedback="validMsg(formError.login, 'login', 'логин')" :validation-status="validStatus(formError.login)">
        <n-input size="large" v-model:value="formValue.login" type="text" placeholder="Логин">
          <template #prefix>
            <n-icon :component="User"/>
          </template>
        </n-input>
      </n-form-item>
      <n-form-item label="Пароль" path="password" required :feedback="validMsg(formError.password, 'password', 'пароль')" :validation-status="validStatus(formError.password)">
        <n-input size="large" v-model:value="formValue.password" type="password" show-password-on="mousedown" placeholder="Пароль">
          <template #prefix>
            <n-icon :component="Password"/>
          </template>
        </n-input>
      </n-form-item>
      <div>
        <n-checkbox size="large" label="Не выходить" v-model:checked="formValue.remember" />
      </div>
    </n-form>
    <template #footer>
      <n-flex justify="space-between">
        <n-button text @click="router.push(`/v1/oauth/forgot-password?${query}`)">Забыли свой пароль?</n-button>
        <n-button @click="authorize" :disabled="formIsEmpty()" size="large" type="primary" style="width: 150px">
          Войти
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
