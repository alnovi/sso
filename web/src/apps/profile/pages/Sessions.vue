<script setup>
import {onMounted, ref} from "vue";
import {useDialog} from "naive-ui";
import {Laptop, TrashCan} from "@vicons/carbon";
import {useApi} from "../../../services/api.js";
import {config} from "../../../services/utils.js";
import moment from "moment";

const dialog = useDialog();
const api = useApi(config('VITE_API_HOST', '/'))

const sessions = ref([])

const sessionDeleteConfirm = (session) => {
  dialog.warning({
    title: "Внимание",
    content: `Удалить сессию ${session.app} от ${moment(session.created_at).format('DD.MM.YYYY HH:mm:ss')}?`,
    positiveText: "Удалить",
    negativeText: "Отмена",
    draggable: true,
    onPositiveClick: () => {
      api.delete(`/profile/sessions/${session.id}`).then(() => {
        if (session.is_current) {
          window.location.reload()
        } else {
          sessions.value = sessions.value.filter((el) => el.id !== session.id);
        }
      })
    },
  });
}

onMounted(() => {
  api.get(`/profile/sessions`)
    .then(res => {
      sessions.value = res.data
    })
})
</script>

<template>
  <div class="block-layout">
    <div class="block-layout-header">
      <div class="block-layout-header__title">Устройства</div>
      <div class="separator"></div>
      <div class="block-layout-header__description">Список активных сеансов.</div>
    </div>
    <div class="block-layout-content">
      <div class="session-list">
        <n-el class="session-list-item" v-for="(session, i) in sessions" :key="i">
          <div class="session-list-item__logo">
            <n-icon-wrapper :size="40" :border-radius="20">
              <n-icon :size="30">
                <Laptop/>
              </n-icon>
            </n-icon-wrapper>
          </div>
          <div class="session-list-item__content">
            <div class="session-list-item__title">
              {{ session.app }} ({{ session.os }})
              <n-tag v-if="session.is_current" type="success" size="small">это устройство</n-tag>
            </div>
            <div class="session-list-item__description">
              IP: {{ session.ip }} | Дата: {{ moment(session.created_at).format('DD.MM.YYYY HH:mm:ss') }}
            </div>
          </div>
          <div class="session-list-item__action">
            <n-button strong secondary circle type="error" @click="sessionDeleteConfirm(session)">
              <template #icon>
                <n-icon>
                  <TrashCan/>
                </n-icon>
              </template>
            </n-button>
          </div>
        </n-el>
      </div>
    </div>
  </div>
</template>
