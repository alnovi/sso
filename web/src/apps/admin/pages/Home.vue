<script setup>
import {onActivated, ref} from "vue"
import {Apps, Devices, UserMultiple} from "@vicons/carbon";
import {notifyError} from "../../../services/notify.js";
import {useApi} from "../../../services/api.js";
import {config} from "../../../services/utils.js";

const api = useApi(config('VITE_API_HOST', '/'))

const stats = ref({
  users: 0,
  clients: 0,
  sessions: 0
})

const loadStats = async () => {
  api.get("/api/stats")
    .then(res => {
      stats.value = res.data
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
  loadStats()
})
</script>

<template>
  <n-breadcrumb style="padding-bottom: 24px">
    <n-breadcrumb-item>Главная</n-breadcrumb-item>
  </n-breadcrumb>
  <n-grid x-gap="12" :cols="3">
    <n-gi>
      <router-link :to="{name: 'clients'}">
        <n-card>
          <div class="card">
            <n-icon class="card-icon" size="24">
              <Apps/>
            </n-icon>
            <div class="card-title">Приложения</div>
            <div class="card-counter">{{ stats.clients }}</div>
          </div>
        </n-card>
      </router-link>
    </n-gi>
    <n-gi>
      <router-link :to="{name: 'users'}">
        <n-card>
          <div class="card">
            <n-icon class="card-icon" size="24">
              <UserMultiple/>
            </n-icon>
            <div class="card-title">Пользователи</div>
            <div class="card-counter">{{ stats.users }}</div>
          </div>
        </n-card>
      </router-link>
    </n-gi>
    <n-gi>
      <router-link :to="{name: 'sessions'}">
        <n-card>
          <div class="card">
            <n-icon class="card-icon" size="24">
              <Devices/>
            </n-icon>
            <div class="card-title">Устройства</div>
            <div class="card-counter">{{ stats.sessions }}</div>
          </div>
        </n-card>
      </router-link>
    </n-gi>
  </n-grid>
</template>

<style scoped lang="scss">
.n-card {
  &:hover {
    border-color: var(--n-color-target);
  }
}

.card {
  display: flex;
  flex-direction: row;
  align-items: center;

  .card-title {
    font-weight: 400;
    font-size: 18px;
    margin: 0 10px;
  }

  .card-counter {
    font-weight: 800;
    font-size: 18px;
    text-align: right;
    color: var(--n-color-target);
    flex-grow: 1;
    opacity: .65;
  }
}
</style>
