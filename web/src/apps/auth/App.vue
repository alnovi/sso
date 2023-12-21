<script setup>
import { watch } from "vue";
import { useToast } from "primevue/usetoast";
import { useToastStore } from "./store/toast.js";

const toast = useToast()
const store = useToastStore()

watch(store, (newMessages, oldMessages) => {
  let msg = store.getMessage()
  if (msg) {
    toast.add(msg)
  }
})

</script>

<template>
  <Toast />
  <router-view/>
</template>

<style scoped lang="scss"></style>