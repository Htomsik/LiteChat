<script setup>
import { ref } from 'vue'
import LoginPage from './components/LoginPage.vue'
import ChatPage from './components/ChatPage.vue'

const isConnected = ref(false)
const userName = ref('')
const serverId = ref('')
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

function connectToChat(name, id) {
  userName.value = name
  serverId.value = id
  closeAlert()
  isConnected.value = true
}

function disconnectFromChat(name, id) {
  userName.value = name
  serverId.value = id
  isConnected.value = false
}
</script>

<template>
  <div id="app" class="full">
    <div v-if="showErrorAlert" id="alert" class="marginAll alert alert-danger alert-dismissible fade show" role="alert">
      <strong>{{ errorAlert }}</strong>
      <button type="button" @click="closeAlert" class="btn-close" aria-label="Close"></button>
    </div>
    <LoginPage v-if="!isConnected" @alert="openAlert" @connect="connectToChat" :userNameProp="userName" :serverIdProp="serverId" />
    <ChatPage v-if="isConnected" @alert="openAlert" @disconnect="disconnectFromChat" :userNameProp="userName" :serverIdProp="serverId" />
  </div>
</template>

<style scoped>

</style>
