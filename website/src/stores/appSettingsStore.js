import {defineStore} from "pinia";

export const AppSettingsStore = defineStore("AppSettings",{
    state: () => ({
        userName: '',
        serverId: ''
    }),

    actions: {

    },
})