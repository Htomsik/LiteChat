import {defineStore} from "pinia";

// Global error alert
export const AlertStore = defineStore("alertStore",{
    state: () => ({
        message : '',
        show : false
    }),

    actions: {
        open(message){
            this.close()
            this.message = message
            this.show = true
        },

        close(){
            this.message = ""
            this.show = false
        }
    },
})