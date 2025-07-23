import { ref } from 'vue'
import {Message, MessageType, User} from "../models/chatModels.js";

/** @type {Vue.ref<User>} */
export const CurrentUser = ref(null)

/** @type {Vue.ref<User[]>} */
export const Users = ref([])

/** @type {Vue.ref<Message[]>} */
export const Messages = ref([])

let socket = null

const listeners = {
    connect: [],
    disconnect: [],
    userNameChanged: [],
}

// Listeners functions
export function On(event, handler) {
    if (listeners[event])
        listeners[event].push(handler)
}

function emit(event, ...args) {
    if (listeners[event])
        listeners[event].forEach(fn => fn(...args))
}

// Chat functions
export function Connect(serverId, userName){
    if (socket)
        Disconnect()

    let url = `/api/chat/${serverId}?User=${userName}`

    socket = new WebSocket(url)

    socket.onopen = (evt) => {
        emit('connect', evt)
    }

    socket.onclose = (evt) => {
        emit('disconnect', evt)
    }

    socket.onmessage = onMessage
}

export function Disconnect(){
    if (socket){
        socket.close()
        socket = null
    }

    CurrentUser.value = null
    Users.value = []
    Messages.value = []
}

/**
 * @param {Message} message
 */
export function SendMessage(message) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(message)
    }
}

function onMessage(evt) {
    const message = new Message(JSON.parse(evt.data))
    switch (message.Type) {
        case MessageType.message:
            Messages.value.push(message)
            break

        case MessageType.userList:
            onMessageUserList(message.MessageData)
            break

        case MessageType.userNameChanged:
            emit('userNameChanged', message.MessageData)
            break

        default:
            console.warn('Unknown message type:', message.Type, message)
            break
    }
}

/**
 * @param {User[]} usersData
 */
function onMessageUserList(usersData){
    if (!Array.isArray(usersData)) {
        console.error('Users data is not an array:', usersData)
        return
    }

    // Создаем экземпляры User класса
    const users = usersData.map(userData => {
        const user = new User(userData)
        user.Color = getHexColorByUserRole(user.Role)
        return user
    })
    Users.value = [...users]

    // TODO Add UserChanged event on api later. Its time
    // Change only if user doesn't exist
    if(CurrentUser.value && CurrentUser.value.id !== '')
    {
        CurrentUser.value = users.find(user => user.Id === CurrentUser.value.id)
        return
    }

    // For Admin
    if(users.length === 1)
    {
        CurrentUser.value = users[0]
        return
    }

    CurrentUser.value = users.find(user => user.Id && user.Id !== '')
}

/**
 * @param {UserRole} role
 * @return {string} hexCode
 */
function getHexColorByUserRole(role) {
    return role.IsAdmin ? "#cf6679" : "#4caf50"
}