<template>

  <div class="full">
    <div class="centerContainer full">
      <div class="colContainer">
        <div class="centerContainer marginAll">
          <div class="card text-center">
            <h1 style="font-weight: bold;">LiteChat</h1>
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
// imports
import { ref, computed, watch } from 'vue'
import {AlertStore} from "../stores/alertStore.js";
import {AppSettingsStore} from "../stores/appSettingsStore.js";
import router from "../routes/router.js";
import * as AuthService from "../services/authorizationService.js"

// emits
const emit = defineEmits([''])

// store + routers
const appSettings = AppSettingsStore()
const alertStore = AlertStore()
const appRouter = router

// ref, computed
const formClass = ref('needs-validation')

// watch
watch([appSettings.userName, appSettings.serverId], () => {
  formClass.value = 'was-validated'
})

// live cycle

// functions
const blockConnection = computed(() =>
    appSettings.userName.toString().trim().length < 2 || appSettings.userName.toString().trim().length === 0
)

async function tryConnect() {
  const [canConnect, errorMessage] = await AuthService.canConnect(appSettings.userName, appSettings.serverId)

  if(!canConnect){
    alertStore.open(errorMessage)
    return
  }

  await appRouter.push("chat")
}
</script> 