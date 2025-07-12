import { ref } from 'vue'

export const Users = ref([])
export const Messages = ref([])

let socket = null

const messageType = Object.freeze({
    message: 'Message',
    userList: 'UsersList',
    userNameChanged: 'UserNameChanged',
})

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

    socket.onmessage = OnMessage
}

export function Disconnect(){
    if (socket){
        socket.close()
        socket = null
    }

    Users.value = []
    Messages.value = []
}

export function sendMessage(message) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(message)
    }
}

function OnMessage(evt) {
    let messageObj = JSON.parse(evt.data)
    switch (messageObj.type) {

        case messageType.message:
            Messages.value.push(messageObj)
            break

        case messageType.userList:
            for (let user of messageObj.message) {
                user.Color = getRandomHexColorByUserName(user.Name)
            }
            Users.value = [...messageObj.message]
            break

        case messageType.userNameChanged:
            emit('userNameChanged', messageObj.message)
            break
        default:
            console.warn('Unknown message type:', messageObj.type, messageObj)
            break
    }
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