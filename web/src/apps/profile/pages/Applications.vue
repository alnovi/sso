<script setup>
import {onMounted, ref} from "vue";
import {useApi} from "../../../services/api.js";
import {config} from "../../../services/utils.js";

const api = useApi(config('VITE_API_HOST', '/'))

const applications = ref([])

const appOpen = (app) => {
  window.open(`/oauth/authorize?response_type=code&client_id=${app.id}&redirect_uri=${app.callback}`, '_blank')
}

onMounted(() => {
  api.get(`/profile/clients`).then(res => {
    applications.value = res.data
  })
})
</script>

<template>
  <div class="block-layout">
    <div class="block-layout-header">
      <div class="block-layout-header__title">Приложения</div>
      <div class="separator"></div>
      <div class="block-layout-header__description">Список доступных приложений.</div>
    </div>
    <div class="block-layout-content">
      <div class="app-list">
        <div class="app-list-container" v-for="(app, i) in applications" :key="i">
          <n-el class="app-list-item" @click="appOpen(app)">
            <div class="app-list-item__logo">
              <img :src="app.icon || '/public/app.png'" width="50" height="50" alt="">
            </div>
            <div class="app-list-item__title">
              {{ app.name }}
            </div>
          </n-el>
        </div>
      </div>
    </div>
  </div>
</template>
