<script setup>
import {ref} from "vue";
import {useRouter} from "vue-router";
import {useNotification} from "naive-ui";
import {useApi} from "../../../services/api.js";
import {config, meta, validMsg, validStatus} from "../../../services/utils.js";
import {notifyError, notifyInfo} from "../../../services/notify.js";
import {User} from "@vicons/carbon";

const api = useApi(config('VITE_API_HOST', '/'))
const query = meta('auth-query', config('VITE_AUTH_QUERY'))
const router = useRouter()
const notification = useNotification()

const formRef = ref(null);

const formValue = ref({
  login: '',
})

const formError = ref({
  login: null,
})

const formIsEmpty = () => {
  return formValue.value.login.length < 5
}

async function forgot() {
  formError.value.login = null

  api.post(`/oauth/forgot-password?${query}`, formValue.value)
    .then(res => {
      formValue.value.login = ""
      notification.success(notifyInfo('Ссылка для смены пароля отправлена на почту'))
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
  <n-card title="Восстановление доступа" bordered :segmented="{content: true, footer: 'soft'}">
    <n-form :ref="formRef" :label-width="80" :model="formValue">
      <n-form-item label="Логин" path="login" required :feedback="validMsg(formError.login, 'login', 'логин')"
                   :validation-status="validStatus(formError.login)">
        <n-input size="large" v-model:value="formValue.login" type="text" placeholder="Логин">
          <template #prefix>
            <n-icon :component="User"/>
          </template>
        </n-input>
      </n-form-item>
    </n-form>
    <template #footer>
      <n-flex justify="space-between">
        <n-button text @click="router.push(`/oauth/authorize?${query}`)">У меня есть пароль</n-button>
        <n-button @click="forgot" :disabled="formIsEmpty()" size="large" type="primary" style="width: 150px">
          Восстановить
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
