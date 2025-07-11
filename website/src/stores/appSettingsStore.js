import {defineStore} from "pinia";

export const AppSettingsStore = defineStore("AppSettings",{
    state: () => ({
        userName: '',
        serverId: '',
        isConnected: false
    }),

    // TODO переделать на полноценный сервис
    actions: {
        connectToChat(){
            this.isConnected = true
        },
        disconnectFromChat(){
            this.isConnected = false
        },
    },
})