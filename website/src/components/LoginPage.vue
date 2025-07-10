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
                  <input maxlength="128" minlength="1" id="serverInput" type="text" class="form-control" v-model="serverId"
                         placeholder="" required>
                  <label for="serverInput">Server</label>
                </div>
                <div class="form-floating mb-3">
                  <input id="userNameInput" type="text" class="form-control" v-model="userName"
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
import { ref, computed, watch, onMounted } from 'vue'
import axios from 'axios'

const props = defineProps({
  userNameProp: String,
  serverIdProp: String
})
const emit = defineEmits(['connect', 'alert'])

const userName = ref('')
const serverId = ref('')
const formClass = ref('needs-validation')

onMounted(() => {
  userName.value = props.userNameProp || ''
  serverId.value = props.serverIdProp || ''
})

const blockConnection = computed(() =>
  userName.value.trim().length < 2 || serverId.value.trim().length === 0
)

watch([userName, serverId], () => {
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
  let url = `/api/chat/canConnect/${serverId.value}?User=${userName.value}`
  if (!await checkApiIsAlive()) {
    emit('alert', 'Connection to server failed')
    return
  }
  try {
    const response = await axios.get(url)
    if (response.status === 200) {
      emit('connect', userName.value, serverId.value)
    }
  } catch (error) {
    handleAxiosErrors(error)
  }
}
</script> 