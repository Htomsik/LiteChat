<script setup>
import { ref } from 'vue'
import 'vue-router'
import 'pinia'

import LoginView from './components/LoginView.vue'
import ChatView from './components/ChatView.vue'

import {AppSettingsStore} from "./stores/appSettingsStore.js";


const appSettingsStore = AppSettingsStore()

const errorAlert = ref('')
const showErrorAlert = ref(false)

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

    <LoginView v-if="!appSettingsStore.isConnected" @alert="openAlert"/>
    <ChatView v-if="appSettingsStore.isConnected" @alert="openAlert"/>
  </div>
</template>

<style scoped>

</style>
