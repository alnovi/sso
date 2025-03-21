<script setup>
import {Apps, UserMultiple, Devices} from "@vicons/carbon";
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
              <n-icon class="card-icon" size="24"><Apps/></n-icon>
              <div class="card-title">Приложения</div>
              <div class="card-counter">0</div>
            </div>
          </n-card>
      </router-link>
    </n-gi>
    <n-gi>
      <router-link :to="{name: 'clients'}">
        <n-card>
          <div class="card">
            <n-icon class="card-icon" size="24"><UserMultiple/></n-icon>
            <div class="card-title">Пользователи</div>
            <div class="card-counter">0</div>
          </div>
        </n-card>
      </router-link>
    </n-gi>
    <n-gi>
      <router-link :to="{name: 'clients'}">
        <n-card>
          <div class="card">
            <n-icon class="card-icon" size="24"><Devices/></n-icon>
            <div class="card-title">Устройства</div>
            <div class="card-counter">0</div>
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
    opacity: 0;
  }
}
</style>
