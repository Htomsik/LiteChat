<template>
  <div class="centerContainer rowContainer full">

    <!-- Chat info -->
    <div class="card chatContainer" style="width: 200px">

      <div class="colContainer" style="flex-shrink: 0;text-align: center">
        {{ appSettings.serverId }}
        <div style="border: 1px solid #6c757d; border-radius: 5px"></div>
      </div>

      <div class="colContainer overflow-auto" style=" flex-grow: 1">
        <div v-for="[role, users] in usersToRoleUsers(users)" :key="role">
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
        <button @click="disconnect" id="disconnectButton" style="width: 100%" class="btn btn-secondary" type="button">Disconnect</button>
      </div>

    </div>

    <!-- messages -->
    <div style="flex-grow: 1" class="card chatContainer">
      <div style="display: flex; flex-direction: column; height: 100%">
        <div style="flex-grow: 1" class="overflow-auto">
          <div v-for="item in messages" :key="item.dateTime + item.user" :class="{ messageBoxLeft: item.user === appSettings.userName }" class="test">
            <div class="message bg-secondary">
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
                  class="btn btn-secondary">
            <i class="bi bi-send-fill"></i>
          </button>
        </div>
        </div>
    </div>
  </div>

</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import {AppSettingsStore} from "../stores/appSettingsStore.js";
const appSettings = AppSettingsStore()

const emit = defineEmits(['alert'])

onMounted(() => {
  connect()
})

const chatSocket = ref(null)
const users = ref([])
const messages = ref([])
const currentMessage = ref('')

const messageType = Object.freeze({
  message: 'Message',
  userList: 'UsersList',
  UserNameChanged: 'UserNameChanged',
})

const blockSendMessage = computed(() => currentMessage.value.length === 0)

function sendMessage() {
  if (chatSocket.value) {
    chatSocket.value.send(currentMessage.value)
    currentMessage.value = ''
  }
}

function disconnect() {
  if (chatSocket.value) {
    chatSocket.value.close()
  }
  emit('alert', 'You have been disconnected')
  appSettings.disconnectFromChat()
}

function connect() {
  if (chatSocket.value)
      chatSocket.value.close()

  console.log(appSettings.userName)

  let url = `/api/chat/${appSettings.serverId}?User=${appSettings.userName}`

  chatSocket.value = new WebSocket(url)
  chatSocket.value.onopen = socketOnOpen
  chatSocket.value.onclose = socketOnClose
  chatSocket.value.onmessage = socketOnMessage
}

function socketOnOpen(evt) {}

function socketOnClose(evt) {
  disconnect()
}

function socketOnMessage(evt) {
  let messageObj = JSON.parse(evt.data)
  switch (messageObj.type) {
    case messageType.message:
      messages.value.push(messageObj)
      break
    case messageType.userList:
      for (let user of messageObj.message) {
        user.Color = getRandomHexColorByUserName(user.Name)
      }
      users.value = messageObj.message
      break
    case messageType.UserNameChanged:
      appSettings.userName = messageObj.message
      break
  }
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

function getRandomHexColorByUserName(username) {
  let hash = 0
  for (let i = 0; i < username.length; i++) {
    hash = username.charCodeAt(i) + ((hash << 10) - hash)
  }
  let hex = (hash & 0x00FFFFFF).toString(16).toUpperCase()
  let r = parseInt(hex.substring(0, 2), 16)
  let g = parseInt(hex.substring(2, 4), 16)
  let b = parseInt(hex.substring(4, 6), 16)
  let factor = 0.2
  r = Math.round(r + (255 - r) * factor)
  g = Math.round(g + (255 - g) * factor)
  b = Math.round(b + (255 - b) * factor)
  return `#${Math.round(r).toString(16)}${Math.round(g).toString(16)}${Math.round(b).toString(16)}`
}
</script> 