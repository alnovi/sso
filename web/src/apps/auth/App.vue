<script setup>
import {ref, onBeforeMount} from "vue";
import {dateRuRU, lightTheme, darkTheme, ruRU} from "naive-ui"; //lightTheme darkTheme
import {Moon, Sun} from "@vicons/carbon";
import {meta} from "./../../services/utils";

const client = ref({
  icon: meta("client-icon", "/public/app.png"),
  name: meta("client-name", "Приложение"),
})
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

onBeforeMount(() => {
  isDark.value = localStorage.getItem("theme-is-dark") === "true";
  applyTheme()
})
</script>

<template>
  <n-config-provider :theme="theme" :locale="ruRU" :date-locale="dateRuRU">
    <n-notification-provider>
      <n-layout position="absolute" class="main-layout" :class="isDark ? 'dark' : 'light'">
        <n-layout-header>
          <n-flex justify="space-between" align="center">
            <n-flex align="center">
              <img :src="client.icon" width="30px" height="30px" alt=""/>
              <h2>{{ client.name }}</h2>
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
            </n-flex>
          </n-flex>
        </n-layout-header>
        <n-layout-content style="padding: 24px">
          <n-flex align="center" justify="center" style="height: 100%">
            <router-view/>
          </n-flex>
        </n-layout-content>
      </n-layout>
    </n-notification-provider>
  </n-config-provider>
</template>

<style scoped lang="scss">
.main-layout {
  background-position: center center;
  background-attachment: fixed;
  background-size: cover;
  background-repeat: no-repeat;
  &.light {
    background-image: url("./../../public/light.jpg");
  }
  &.dark {
    background-image: url("./../../public/dark.jpg");
  }
}

.n-layout-header {
  padding: 20px;
  background-color: rgba(255, 255, 255, 0);
  h2 {
    margin: 0;
    font-weight: normal;
  }
}

.n-layout-content {
  background-color: rgba(255, 255, 255, 0);
  height: 80%;
}
</style>
