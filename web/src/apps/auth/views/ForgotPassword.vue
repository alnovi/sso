<script setup>
import { computed, ref } from "vue";
import { useAuthStore } from "../store/auth.js";

const store = useAuthStore()

const isLoading = ref(false)

const onSubmit = () => {
  isLoading.value = true
  store.forgotPassword().finally(() => (isLoading.value = false))
}

const isDisabledButton = computed(() => {
  return !store.credential.login.trim() || isLoading.value
})
</script>

<template>
  <div>
    <Toast />
    <h5>Восстановление аккаунта</h5>
    <div class="grid p-fluid mt-3">
      <div class="col-12">
                <span class="p-float-label">
                    <InputText v-model="store.credential.login" :class="{'p-invalid': !!store.errors.login}" id="login" />
                    <label for="login">Логин или email</label>
                </span>
        <small class="p-error">{{ store.errors.login || '' }}</small>
      </div>
      <div class="col-12">
        <Button label="Далее" :disabled="isDisabledButton" @click="onSubmit"/>
      </div>
      <div class="col-12 py-2 flex align-items-center justify-content-center">
        <router-link class="link" :to="{name: 'signin'}">Войти с паролем</router-link>
      </div>
    </div>
  </div>
</template>
