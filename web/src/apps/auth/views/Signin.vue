<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from "../store/auth.js";

const store = useAuthStore()

const isLoading = ref(false)

const isDisableButton = computed(() => {
  return !store.credential.login.trim() || !store.credential.password.trim() || isLoading.value
})

const onSubmit = () => {
  isLoading.value = true
  store.signIn().finally(() => (isLoading.value = false))
}
</script>

<template>
  <div>
    <h5>Войти</h5>
    <div class="grid p-fluid mt-3">
      <div class="col-12">
                <span class="p-float-label">
                    <InputText v-model="store.credential.login" type="text" :class="{'p-invalid': !!store.errors.login}" id="login" />
                    <label for="login">Логин или email</label>
                </span>
        <small class="p-error">{{ store.errors.login || '' }}</small>
      </div>
      <div class="col-12">
                <span class="p-float-label">
                    <Password v-model="store.credential.password" :class="{'p-invalid': !!store.errors.password}" id="password" :feedback="false" toggleMask />
                    <label for="password">Пароль</label>
                </span>
        <small class="p-error">{{ store.errors.password || '' }}</small>
      </div>
      <div class="col-12 py-2 flex align-items-center justify-content-between">
        <div class="flex align-items-center">
          <Checkbox v-model="store.credential.remember" inputId="remember" :binary="true" />
          <label for="remember" class="ml-2">Запомнить меня</label>
        </div>
        <router-link class="link" :to="{name: 'forgot-password'}">Забыли пароль?</router-link>
      </div>
      <div class="col-12">
        <Button label="Войти" :disabled="isDisableButton" @click="onSubmit"/>
      </div>
    </div>
  </div>
</template>
