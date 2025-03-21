<script setup>
import {onBeforeMount, ref} from 'vue'
import {darkTheme, lightTheme, ruRU, dateRuRU} from "naive-ui";
import {Moon, Sun} from "@vicons/carbon";
import {meta} from "../../services/utils.js";

const data = ref({
  code: meta('err-code', '500'),
  error: meta('err-message', 'Internal error')
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
    <n-layout position="absolute" class="main-layout" :class="isDark ? 'dark' : 'light'">
      <n-layout-header>
        <n-flex justify="end" align="center">
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
      </n-layout-header>
      <n-layout-content style="padding: 24px">
        <n-flex vertical align="center" justify="center" style="height: 100%">
          <n-card bordered :segmented="{content: true, footer: 'soft'}">
            <n-flex vertical align="center" justify="center" style="height: 100%">
              <div class="err-code">{{ data.code }}</div>
              <div class="err-message">{{ data.error }}</div>
            </n-flex>
          </n-card>
        </n-flex>
      </n-layout-content>
    </n-layout>
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
</style>
