<template>
  <div class="centerContainer rowContainer full">

    <!-- Chat info -->
    <div class="cardSC chatContainer" style="width: 200px">

      <div class="colContainer" style="flex-shrink: 0;text-align: center">
        <h4 style="font-weight: bold;">{{appSettings.serverId}}</h4>
        <div style="border: 1px solid var(--bs-border-color); border-radius: 5px"></div>
      </div>

      <div class="colContainer overflow-auto" style=" flex-grow: 1">
        <div v-for="[role, users] in usersToRoleUsers(ChatService.Users.value)" :key="role">
          <div class="userList-Role">
            {{ role }}
          </div>

          <div class="rowContainer userList-usersContainer" v-for="user in users" :key="user.Id">
            <div class="userList-userAvatar centerContainer" :style="{'background': user.Color}">
              {{ user.Name[0] }}
            </div>
            <div class="userList-user">
              {{ user.Name }}
            </div>

          </div>
        </div>
      </div>

      <div class="rowContainer" style="flex-shrink: 0">
        <button @click="ChatService.Disconnect" id="disconnectButton" style="width: 100%" class="btn btn-primary" type="button">Disconnect</button>
      </div>

    </div>

    <!-- messages -->
    <div style="flex-grow: 1" class="cardSC chatContainer">

      <div style="display: flex; flex-direction: column; height: 100%">
        <div style="flex-grow: 1" class="overflow-auto">

          <div v-for="item in ChatService.Messages.value" :key="item.dateTime + item.user" :class="{ messageBoxLeft: item.user === appSettings.userName }" class="test">

            <div class="message">
              <span class="message-user">{{ item.user }}</span>
              <div class="message-text">
                {{ item.message }}
              </div>
              <span class="message-dateTime">{{ formatMessageDateTime(item.dateTime) }}</span>
            </div>

          </div>

        </div>
        <div class="rowContainer centerContainer">
          <textarea maxlength="256" v-model="currentMessage" class="marginAll"
                    style="flex-grow: 1; height: 35px; max-height: 70px; min-height: 35px"></textarea>
          <button :disabled="blockSendMessage" @click="sendMessage" type="button"
                  class="btn btn-primary">
            <i class="bi bi-send-fill"></i>
          </button>
        </div>
        </div>
    </div>

  </div>

</template>

<script setup>
// imports
import {ref, computed, onMounted, watch} from 'vue'
import {AlertStore} from "../stores/alertStore.js";
import {AppSettingsStore} from "../stores/appSettingsStore.js";
import router from "../routes/router.js";
import * as ChatService from "../services/chatService.js"

// emits
const emit = defineEmits([''])

// store + routers
const appSettings = AppSettingsStore()
const alertStore = AlertStore()
const appRouter = router

// ref, computed
const currentMessage = ref('')

const blockSendMessage = computed(() => currentMessage.value.length === 0)

const messageType = Object.freeze({
  message: 'Message',
  userList: 'UsersList',
  UserNameChanged: 'UserNameChanged',
})

// watch
watch(ChatService.Messages, (val) => { console.log('Messages changed:', val) })
watch(ChatService.Users, (val) => { console.log('Users changed:', val) })

// live cycle
onMounted(() => {
  ChatService.Connect(appSettings.userName, appSettings.serverId)

  // Subscribes
  ChatService.On("connect", onConnect)
  ChatService.On("disconnect", onDisconnect)
  ChatService.On("userNameChanged", onUserNameChanged)
})

// Functions
function sendMessage() {
  ChatService.sendMessage(currentMessage.value)
  currentMessage.value = ''
}

function formatMessageDateTime(dateTimeString) {
  const dateTime = new Date(dateTimeString)
  const hours = dateTime.getHours()
  const minutes = dateTime.getMinutes()
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`
}

function usersToRoleUsers(usersArr) {
  let roleUsers = new Map()
  for (let i = 0; i < usersArr.length; i++) {
    let user = usersArr[i]
    if (!roleUsers.has(user.Role)) {
      roleUsers.set(user.Role, [])
    }
    let usersPerRole = roleUsers.get(user.Role)
    usersPerRole.push(user)
  }
  return Array.from(roleUsers.entries())
}


// Functions (Subscribes)
function onConnect(){

}

function onDisconnect(){
  alertStore.open('You have been disconnected')
  appRouter.push("login")
}

function onUserNameChanged(){

}

</script> 