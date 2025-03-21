<script setup>
</script>

<template>
  <n-card bordered :segmented="{content: true, footer: 'soft'}">
    <n-flex vertical align="center" justify="center" style="height: 100%">
      <div class="err_code">404</div>
      <div class="err_message">Страница не найдена</div>
    </n-flex>
  </n-card>
</template>

<style scoped>
.n-card {
  height: 40%;
  box-shadow: 0 10px 20px 0 rgba(0, 0, 0, .2);
  max-width: 600px;
  background-color: rgba(255, 255, 255, .95);
  border-radius: 12px;
}

.n-flex .err_code {
  font-size: 84px;
  font-weight: bold;
}
.n-flex .err_message {
  font-size: 24px;
  font-weight: normal;
}
</style>
