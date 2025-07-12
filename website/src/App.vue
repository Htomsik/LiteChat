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

    <div v-if="alertStore.show" id="alert" class="alert alert-danger alert-dismissible fade show" role="alert" @click="alertStore.close()" >
      <strong>{{ alertStore.message }}</strong>
    </div>

    <RouterView />

  </div>
</template>

<style scoped>

</style>
