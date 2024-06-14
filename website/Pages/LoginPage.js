export default {

    emits: ['connect', 'alert'],
    props:['userNameProp', 'serverIdProp'],

    data() {
        return {
            userName: "",
            serverId: "",
        };
    },
    mounted(){
        this.userName = this.userNameProp;
        this.serverId = this.serverIdProp;
    },
    computed: {
        blockConnection: function () {
            return this.userName.trim().length < 2 || this.serverId.trim().length === 0;
        },
    },
    watch: {
        userName: {
            handler(newValue, oldValue) {
                this.changeValidation()
            },
            once: true
        },

        serverId: {
            handler(newValue, oldValue) {
                this.changeValidation()
            },
            once: true
        }
    },
    methods: {
        changeValidation: function (enable = true) {
            let removeClass = enable ? 'needs-validation' : 'was-validated'
            let setClass = enable ? 'was-validated' : 'needs-validation'

            let forms = document.getElementById("ConnectInfo")

            forms.classList.remove(removeClass)
            forms.classList.add(setClass)
        },

        handleAxiosErrors: function (error) {
            if (error.response) {
                if (error.response.status === 422) {
                    this.$emit("alert", error.response.data.error)
                }
                this.$emit("alert","Connection to server failed")
            }
        },

        checkApiIsAlive: async function () {
            let ret
            let url = `/api/isAlive`

            await axios
                .get(url)
                .then(response => {
                    ret = response.status === 200
                })
                .catch((error) =>{
                    this.handleAxiosErrors(error)
                    ret = false
                });

            return ret
        },

        tryConnect: async function () {
            let url = `/api/chat/canConnect/${this.serverId}?User=${this.userName}`
            if (!await this.checkApiIsAlive()) {
                this.$emit("alert","Connection to server failed")
                return;
            }

            await axios
                .get(url)
                .then(response => {
                    if (response.status === 200) {
                        this.$emit("connect",this.userName, this.serverId);
                    }
                })
                .catch((error) =>{
                    this.handleAxiosErrors(error)
                });
        },

    },
    template: `
      <div class="full">
        <div class="centerContainer full">

          <!--    Connection window   -->
          <div class="colContainer">

            <!--    Upper data-->
            <div class="centerContainer marginAll">
              <div class="card text-center">
                <h1 style="font-weight: bold;">ChatHub</h1>
                <form id="ConnectInfo" class="needs-validation">

                  <!--    Data group  -->
                  <div>
                    <div class="form-floating mb-3">
                      <input maxlength="128" minlength="1" id="serverInput" type="text" class="form-control" v-model="serverId"
                             placeholder="" required>
                      <label for="serverInput">Server</label>
                    </div>

                    <div class="form-floating mb-3">
                      <input id="userNameInput" type="text" class="form-control" v-model="userName"
                             maxlength="20" minlength="2" placeholder="" required>
                      <label for="userNameInput">Username</label>
                    </div>
                  </div>

                  <!--    Button group    -->
                  <div>
                    <button :disabled="blockConnection" type="button" v-on:click="tryConnect"
                            class="btn btn-primary">Connect
                    </button>
                  </div>

                </form>
              </div>
            </div>

          </div>

        </div>
        
      </div>
    `
}