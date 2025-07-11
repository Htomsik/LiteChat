<template>

  <div class="full">
    <div class="centerContainer full">
      <div class="colContainer">
        <div class="centerContainer marginAll">
          <div class="card text-center">
            <h1 style="font-weight: bold;">ChatHub</h1>
            <form id="ConnectInfo" :class="formClass" @submit.prevent>
              <div>
                <div class="form-floating mb-3">
                  <input maxlength="128" minlength="1" id="serverInput" type="text" class="form-control" v-model="appSettings.serverId"
                         placeholder="" required>
                  <label for="serverInput">Server</label>
                </div>
                <div class="form-floating mb-3">
                  <input id="userNameInput" type="text" class="form-control" v-model="appSettings.userName"
                         maxlength="20" minlength="2" placeholder="" required>
                  <label for="userNameInput">Username</label>
                </div>
              </div>
              <div>
                <button :disabled="blockConnection" type="button" @click="tryConnect"
                        class="btn btn-primary">Connect
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>

</template>

<script setup>
import { ref, computed, watch } from 'vue'
import axios from 'axios'

import {AppSettingsStore} from "../stores/appSettingsStore.js";
const appSettings = AppSettingsStore()

const emit = defineEmits(['connect', 'alert'])
const formClass = ref('needs-validation')


const blockConnection = computed(() =>
    appSettings.userName.toString().trim().length < 2 || appSettings.userName.toString().trim().length === 0
)

watch([appSettings.userName, appSettings.serverId], () => {
  formClass.value = 'was-validated'
})

function handleAxiosErrors(error) {
  if (error.response) {
    if (error.response.status === 422) {
      emit('alert', error.response.data.error)
    }
    emit('alert', 'Connection to server failed')
  }
}

async function checkApiIsAlive() {
  let ret
  let url = `/api/isAlive`
  try {
    const response = await axios.get(url)
    ret = response.status === 200
  } catch (error) {
    handleAxiosErrors(error)
    ret = false
  }
  return ret
}

async function tryConnect() {
  let url = `/api/chat/canConnect/${appSettings.serverId}?User=${appSettings.userName}`
  if (!await checkApiIsAlive()) {
    emit('alert', 'Connection to server failed')
    return
  }
  try {
    const response = await axios.get(url)
    if (response.status === 200) {
      appSettings.connectToChat()
    }
  } catch (error) {
    handleAxiosErrors(error)
  }
}
</script> 