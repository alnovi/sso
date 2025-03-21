<script setup>
import {ref, onBeforeMount} from "vue"
import {darkTheme, dateRuRU, lightTheme, ruRU} from "naive-ui"; //lightTheme darkTheme
import Profile from "./pages/Profile.vue";
import Applications from "./pages/Applications.vue";
import Sessions from "./pages/Sessions.vue";
import {config} from "../../services/utils.js";
import {useApi} from "../../services/api.js";
import {Logout, Moon, Sun} from "@vicons/carbon";
import Secure from "./pages/Secure.vue";

const api = useApi(config('VITE_API_HOST', '/'))

const isDark = ref(false)
const theme = ref(null)

const changeTheme = () => {
  isDark.value = !isDark.value
  localStorage.setItem("theme-is-dark", isDark.value ? 'true' : 'false')
  applyTheme()
}

const applyTheme = () => {
  isDark.value ? theme.value = darkTheme : theme.value = lightTheme
}

const logout = () => {
  api.post("/profile/logout", null).finally(() => {
    window.location.reload();
  })
}

onBeforeMount(() => {
  isDark.value = localStorage.getItem("theme-is-dark") === "true";
  applyTheme()
})

</script>

<template>
  <n-config-provider :theme="theme" :locale="ruRU" :date-locale="dateRuRU">
    <n-dialog-provider>
      <n-notification-provider>
        <n-layout style="height: 100vh" :native-scrollbar="false">
          <n-layout-header bordered>
            <div class="manu-layout">
              <n-flex justify="space-between" align="center">
                <n-flex align="center">
                  <img src="/public/sso.png" width="40px" height="40px" alt=""/>
                  <h1>Единый вход</h1>
                </n-flex>
                <n-flex align="center">
                  <n-button v-if="isDark" tertiary circle @click="changeTheme">
                    <template #icon>
                      <n-icon size="22" :component="Sun"></n-icon>
                    </template>
                  </n-button>
                  <n-button v-else tertiary circle @click="changeTheme">
                    <template #icon>
                      <n-icon size="22" :component="Moon"></n-icon>
                    </template>
                  </n-button>
                  <n-button tertiary circle @click="logout">
                    <template #icon>
                      <n-icon size="22" :component="Logout"></n-icon>
                    </template>
                  </n-button>
                </n-flex>
              </n-flex>
            </div>
          </n-layout-header>
          <n-layout-content style="padding: 24px">
            <profile />
            <applications />
            <sessions />
            <secure />
          </n-layout-content>
        </n-layout>
      </n-notification-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>
