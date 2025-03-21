<script setup>
import {ref, watch, onBeforeMount} from "vue"
import {NIcon, darkTheme, dateRuRU, lightTheme, ruRU} from "naive-ui";
import {UserAvatar, User, Logout, Menu, Moon, Sun, Home, Apps, UserMultiple, Devices} from "@vicons/carbon";
import {RouterLink, useRoute} from "vue-router";
import {useApi} from "../../services/api.js";
import {config} from "../../services/utils.js";

const api = useApi(config('VITE_API_HOST', '/'))
const route = useRoute()

const renderIcon = (icon) => {
  return () => {
    return h(NIcon, null, {default: () => h(icon)});
  };
}

const collapsed = ref(false);
const isDark = ref(false)
const theme = ref(null)
const menuSelectedKey = ref(null)
const profileOptions= ref([
  {
    label: "Профиль",
    key: "profile",
    icon: renderIcon(User)
  },
  {
    label: "Выйти",
    key: "logout",
    icon: renderIcon(Logout)
  }
])
const menuOptions = ref([
  {
    label: () => h(RouterLink, {to: {name: "home"}}, {default: () => "Главная"}),
    key: "home",
    icon: renderIcon(Home)
  }, {
    label: () => h(RouterLink, {to: {name: "clients"}}, {default: () => "Приложения"}),
    key: "clients",
    icon: renderIcon(Apps)
  }, {
    label: () => h(RouterLink, {to: {name: "users"}}, {default: () => "Пользователи"}),
    key: "users",
    icon: renderIcon(UserMultiple)
  }, {
    label: () => h(RouterLink, {to: {name: "sessions"}}, {default: () => "Устройства"}),
    key: "sessions",
    icon: renderIcon(Devices)
  }
])

const onSelectProfileOption = (option) => {
  switch (option) {
    case "profile":
      window.open(`/profile`, '_blank')
      break
    case "logout":
      api.post("/admin/logout", null).finally(() => {
        window.location.reload();
      })
      break
  }
}

const toggleSider = () => {
  collapsed.value = !collapsed.value
  localStorage.setItem("sider-is-collapsed", collapsed.value ? 'true' : 'false')
}

const changeTheme = () => {
  isDark.value = !isDark.value
  localStorage.setItem("theme-is-dark", isDark.value ? 'true' : 'false')
  applyTheme()
}

const applyTheme = () => {
  isDark.value ? theme.value = darkTheme : theme.value = lightTheme
}

watch(() => route.meta, (meta) => {
  let sider = menuOptions.value[0].key
  if (!!meta && !!meta['sider']) {
    sider = meta['sider']
  }
  menuSelectedKey.value = sider
})

onBeforeMount(() => {
  collapsed.value = localStorage.getItem("sider-is-collapsed") === "true";
  isDark.value = localStorage.getItem("theme-is-dark") === "true";
  applyTheme()
})
</script>

<template>
  <n-config-provider :theme="theme" :locale="ruRU" :date-locale="dateRuRU">
    <n-dialog-provider>
      <n-notification-provider>
        <n-layout style="height: 100vh" has-sider>
          <n-layout-sider content-style="padding: 150px 0 24px 0" width="240" :native-scrollbar="false" :collapsed="collapsed" :collapsed-width="0" bordered>
            <svg data-v-6488e27f="" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 275" class="gradient"><defs><linearGradient id="a" x1="13.74" x2="303.96" y1="183.7" y2="45.59" gradientUnits="userSpaceOnUse"><stop offset="0" stop-color="#25636c"></stop><stop offset=".6" stop-color="#3b956f"></stop><stop offset="1" stop-color="#14a058"></stop></linearGradient></defs><path fill="#14a058" d="M0 187.5v25s0 37.5 50 50S300 225 300 225v-37.5Z" opacity=".49"></path><path fill="#14a058" d="M300 237.5S287.5 275 250 275s-128.95-37.5-188.6-75 134.21 0 134.21 0Z" opacity=".49"></path><path fill="#14a058" d="M0 200v12.5a241.47 241.47 0 0 0 112.5 50c73.6 11.69 130.61-14.86 150-25L300 200Z" opacity=".38"></path><path fill="url(#a)" d="M0 0v212.5s62.5-12.5 150 25 150 0 150 0V0Z"></path></svg>
            <div class="gradient-title">
              <h1>SSO - ADMIN</h1>
              <p>Управление клиентами и пользователями.</p>
            </div>
            <n-menu :value="menuSelectedKey" :options="menuOptions" />
          </n-layout-sider>
          <n-layout-content>
            <n-layout position="absolute" :native-scrollbar="false">
              <n-layout-header style="padding: 20px" bordered>
                <n-flex justify="space-between">
                  <n-flex align="center">
                    <n-button quaternary circle size="medium" @click="toggleSider()">
                      <template #icon>
                        <n-icon size="24"><Menu /></n-icon>
                      </template>
                    </n-button>
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
                    <n-dropdown :options="profileOptions" @select="onSelectProfileOption">
                      <n-button tertiary circle>
                        <template #icon>
                          <n-icon size="36" :component="UserAvatar"></n-icon>
                        </template>
                      </n-button>
                    </n-dropdown>
                  </n-flex>
                </n-flex>
              </n-layout-header>
              <n-layout-content style="padding: 24px;">
                <router-view v-slot="{ Component }">
                  <keep-alive>
                    <component :is="Component" />
                  </keep-alive>
                </router-view>
              </n-layout-content>
            </n-layout>
          </n-layout-content>
        </n-layout>
      </n-notification-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<style scoped lang="scss">
.n-layout-sider {
  .gradient {
    position: absolute;
    top: -70px;
    left: 0;
    right: 0;
  }
  .gradient-title {
    color: rgba(255, 255, 255, 0.82);
    text-align: center;
    z-index: 1;
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    padding: 10px;
    h1 {
      font-weight: 700;
      font-size: 28px;
      margin: 0;
    }
    p {
      font-size: 16px;
      margin: 5px 0 0 0;
      line-height: 1.2;
    }
  }
}
</style>
