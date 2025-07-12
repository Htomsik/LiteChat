import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../components/LoginView.vue'
import ChatView from '../components/ChatView.vue'

const routes = [
    { path: '/', name: 'login', component: LoginView },
    { path: '/chat', name: 'chat', component: ChatView },
    { path: '/:pathMatch(.*)*', redirect: '/' } // Redirect to login on 404
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router


