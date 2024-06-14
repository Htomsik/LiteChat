export default {

    emits: ['disconnect','alert'],
    props:['userNameProp', 'serverIdProp'],

    data() {
        return {
            userName: "",
            serverId: "",

            chatSocket: null,
            messages: [
            ],
            currentMessage: "",
        }
    },
    mounted(){
        this.userName = this.userNameProp;
        this.serverId = this.serverIdProp;

        this.connect()
    },
    computed:{
        blockSendMessage: function () {
            return this.currentMessage.length === 0
        }
    },
    methods:{

        sendMessage: function () {
            this.chatSocket.send(this.currentMessage)
            this.currentMessage = "";
        },

        // Close connection
        disconnect: function () {
            if (this.chatSocket) {
                this.chatSocket.close()
            }

            this.$emit('alert', 'You have been disconnected')
            this.$emit('disconnect', this.userName, this.serverId)
        },

        // Connect to websocket chat
        connect: function () {

            if (this.chatSocket != null)
                this.chatSocket.close()


            let url = `/api/chat/${this.serverId}?User=${this.userName}`

            this.chatSocket = new WebSocket(url);

            this.chatSocket.onopen = this.socketOnOpen;
            this.chatSocket.onclose = this.socketOnClose;
            this.chatSocket.onmessage = this.socketOnMessage;
        },

        socketOnOpen: function (evt) {

        },

        socketOnClose: function (evt) {
            this.disconnect()
        },

        socketOnMessage:function (evt) {
            let messageObj = JSON.parse(evt.data)
            this.messages.push(messageObj)
        }
    },
    template:`
      
      <div class="centerContainer rowContainer full">
        
        <!--    Chat info   -->
        <div class="card chatContainer" style="width: 200px">

          <!--    Chat name   -->
          <div class="colContainer" style="flex-shrink: 0;text-align: center">
            {{serverId}}
            <div style="border: 1px solid #6c757d; border-radius: 5px"></div>
          </div>

          <!--    Users   -->
          <div class="colContainer" style=" flex-grow: 1">

          </div>

          <!--    Additional button -->
          <div class="rowContainer"  style="flex-shrink: 0">
            <button v-on:click="disconnect" id="disconnectButton" style="width: 100%" class="btn btn-secondary" type="button">Disconnect</button>
          </div>

        </div>

        <!--   messages    -->
        <div style="flex-grow: 1" class="card chatContainer">

          <!--    Messages    -->
          <div style="display: flex; flex-direction: column; height: 100%">

            <div style="flex-grow: 1" class="overflow-auto">
              <div v-for="item in messages" :class="{ messageBoxLeft: item.user === userName}" class="test">

                <div class="message bg-secondary">
                  <span class="message-user">{{ item.user }}</span>
                  <div class="message-text">
                    {{ item.message }}
                  </div>
                </div>

              </div>
            </div>

            <div class=" rowContainer centerContainer">
                    <textarea maxlength="256" v-model="currentMessage" class="marginAll"
                              style="flex-grow: 1; height: 35px; max-height: 70px; min-height: 35px"></textarea>
              <button :disabled="blockSendMessage" v-on:click="sendMessage" type="button"
                      class="btn btn-secondary">
                <i class="bi bi-send-fill"></i>
              </button>
            </div>

          </div>
        </div>
      </div>
    `
}