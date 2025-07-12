<script setup>
// imports
import { ref } from 'vue'
import 'pinia'
import {AlertStore} from "./stores/alertStore.js"

// store, router
const alertStore = AlertStore()

// emits
const emit = defineEmits([''])

// store + routers
// ref, computed
// watch
// live cycle
// functions

</script>

<template>
  <div id="app" class="full">

    <div v-if="alertStore.show" id="alert" class="marginAll alert alert-danger alert-dismissible fade show" role="alert">
      <strong>{{ alertStore.message }}</strong>
      <button type="button" @click="alertStore.close()" class="btn-close" aria-label="Close"></button>
    </div>

    <RouterView />

  </div>
</template>

<style scoped>

</style>
