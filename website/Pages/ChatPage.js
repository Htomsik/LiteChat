const messageType = Object.freeze({
    message:   "Message",
    userList:   "UsersList", // Receiving array of chat users
    UserNameChanged:   "UserNameChanged",
});


// TODO
// По нормальному разделить монокомпонента на
// 1. Блок показа пользователей
// 2. Блок отправки сообщений
// 3. Всплывающую панель информации о выбранном пользователе
export default {

    emits: ['disconnect','alert'],
    props:['userNameProp', 'serverIdProp'],

    data() {
        return {
            userName: "",
            serverId: "",

            chatSocket: null,

            users:[

            ],

            messages: [],
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

            switch (messageObj.type) {

                case messageType.message:
                    this.messages.push(messageObj)
                    break;

                case messageType.userList:

                    // Set colors by usernames
                    for (let user of messageObj.message) {
                        user.Color = this.getRandomHexColorByUserName(user.Name)
                    }
                    this.users = messageObj.message
                    console.log(messageObj.message)
                    this.usersToMapRoleUsers(this.users)

                    break;

                case messageType.UserNameChanged:
                    this.userName = messageObj.message
                    break;
            }
        },

        formatMessageDateTime: function (dateTimeString){
            const dateTime = new Date(dateTimeString);
            const hours = dateTime.getHours();
            const minutes = dateTime.getMinutes();

            return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`
        },

        /// Users to Role - Users[]
        usersToMapRoleUsers: function (users){

            let roleUsersMap = new Map()

            for(let i = 0; i < users.length; i++){
                let user = users[i];

                if (!roleUsersMap.has(user.Role)){
                    roleUsersMap.set(user.Role, [user])
                }

                roleUsersMap.set(user.Role,roleUsersMap.get(user.Role).push(user))
            }

            console.log(roleUsersMap)
        },

        getRandomHexColorByUserName: function (username){

            // Generate color
            let hash = 0;
            for (let i = 0; i < username.length; i++) {
                hash = username.charCodeAt(i) + ((hash << 10) - hash);
            }

            // Generate random color by hash
            let hex = (hash & 0x00FFFFFF).toString(16).toUpperCase();

            let r = parseInt(hex.substring(0, 2), 16);
            let g = parseInt(hex.substring(2, 4), 16);
            let b = parseInt(hex.substring(4, 6), 16);

            // add or remove brightness 1 light 0 dark 0.5 original
            let factor = 0.2
            r = Math.round(r + (255 - r) * factor);
            g = Math.round(g + (255 - g) * factor);
            b = Math.round(b + (255 - b) * factor);

            // Translate into hex
            return `#${Math.round(r).toString(16)}${Math.round(g).toString(16)}${Math.round(b).toString(16)}`
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
          <div class="colContainer overflow-auto" style=" flex-grow: 1">
                <div class="userList-user rowContainer" v-for="user in users" >
                    
                    <div class="userList-userAvatar centerContainer" :style="{'background': user.Color}">
                        {{user.Name[0]}} 
                    </div>
                    
                    <div>
                        {{user.Name}} 
                    </div>
                    
                </div>
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
                  <span class="message-dateTime">{{ formatMessageDateTime(item.dateTime) }}</span>
                  
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