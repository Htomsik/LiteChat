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

          <div class="userList-usersContainer" v-for="user in users" :key="user.Id">

            <div class="userList-itemContainer" >
              <i class="bi bi-chat-fill userList-permission"
                 :class="user.havePermission(PermissionType.sendMessage) ? 'userList-permission' : 'userList-permissionDenied'"/>
            </div>

            <div class="userList-itemContainer userList-Role" :style="{'background': user.Color}">
              {{ user.Role.Name[0] }}
            </div>

            <div class="userList-itemContainer userList-user">
              {{ user.Name}}
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

          <div v-for="item in ChatService.Messages.value" :key="item.DateTime + item.User" :class="{ messageBoxLeft: isCurrentUserMessage(item)}">

            <div class="message">
              <span class="message-user">{{ item.User }}</span>
              <div class="message-text">
                {{ item.MessageData }}
              </div>
              <span class="message-dateTime">{{ formatMessageDateTime(item.DateTime) }}</span>
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
import {PermissionType} from "../models/chatModels.js";

// emits
const emit = defineEmits([''])

// store + routers
const appSettings = AppSettingsStore()
const alertStore = AlertStore()
const appRouter = router

// ref, computed
const currentMessage = ref('')

const blockSendMessage = computed(() => {
  if (currentMessage.value.length === 0) return true
  if (!ChatService.CurrentUser.value) return true
  return !ChatService.CurrentUser.value.havePermission(PermissionType.sendMessage)
})

// watch
watch(ChatService.Messages, (val) => { console.log('Messages changed:', val) })
watch(ChatService.Users, (val) => { console.log('Users changed:', val) })

// live cycle
onMounted(() => {
  ChatService.Connect(appSettings.serverId, appSettings.userName)

  // Subscribes
  ChatService.On("connect", onConnect)
  ChatService.On("disconnect", onDisconnect)
  ChatService.On("userNameChanged", onUserNameChanged)
})

// Functions
function sendMessage() {
  ChatService.SendMessage(currentMessage.value)
  currentMessage.value = ''
}

function formatMessageDateTime(dateTimeString) {
  const dateTime = new Date(dateTimeString)
  const hours = dateTime.getHours()
  const minutes = dateTime.getMinutes()
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`
}

/**
 * @param {Message} message
 * @returns {boolean}
 */
function isCurrentUserMessage(message) {
  if (!ChatService.CurrentUser.value || !ChatService.CurrentUser.value.Name) {
    return false
  }
  return message.User === ChatService.CurrentUser.value.Name
}

/**
 * @param {User[]} usersArr
 * @returns {[string, User[]][]}
 */
function usersToRoleUsers(usersArr) {
  let roleUsers = new Map()
  for (let i = 0; i < usersArr.length; i++) {
    let user = usersArr[i]
    if (!roleUsers.has(user.Role.Name)) {
      roleUsers.set(user.Role.Name, [])
    }
    let usersPerRole = roleUsers.get(user.Role.Name)
    usersPerRole.push(user)
  }
  return Array.from(roleUsers.entries())
}


// Functions (Subscribes)
function onConnect(){
  alertStore.close()
}

function onDisconnect(){
  alertStore.open('You have been disconnected')
  appRouter.push("login")
}

function onUserNameChanged(){

}

</script> 