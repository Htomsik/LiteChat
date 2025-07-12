import {defineStore} from "pinia";

// Global app settings
export const AppSettingsStore = defineStore("appSettingsStore",{
    state: () => ({
        userName: '',
        serverId: ''
    }),
    actions: {},
})