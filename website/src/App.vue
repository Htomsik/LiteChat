<script setup>
// imports
import { ref } from 'vue'
import 'pinia'

// store, router
// emits
// store + routers

// ref, computed
const errorAlert = ref('')
const showErrorAlert = ref(false)

// watch
// live cycle

// functions
function openAlert(message) {
  closeAlert()
  showErrorAlert.value = true
  errorAlert.value = message
}

function closeAlert() {
  showErrorAlert.value = false
  errorAlert.value = ''
}
</script>

<template>
  <div id="app" class="full">

    <div v-if="showErrorAlert" id="alert" class="marginAll alert alert-danger alert-dismissible fade show" role="alert">
      <strong>{{ errorAlert }}</strong>
      <button type="button" @click="closeAlert" class="btn-close" aria-label="Close"></button>
    </div>

    <RouterView />

  </div>
</template>

<style scoped>

</style>
