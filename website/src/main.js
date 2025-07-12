import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './routes/router.js'


const app = createApp(App);
app.use(createPinia()); // Global reactive like stores
app.use(router)  //  Global view router (view locator)
app.mount('#app')
