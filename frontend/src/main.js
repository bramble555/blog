import { createApp } from 'vue'
import './style.css'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'

// Crash Recovery Logic
const CRASH_KEY = 'gvb_crash_count'
const MAX_CRASHES = 3

function handleCrash(error) {
    console.error('Critical Error:', error)
    const count = parseInt(localStorage.getItem(CRASH_KEY) || '0') + 1
    localStorage.setItem(CRASH_KEY, count.toString())

    if (count >= MAX_CRASHES) {
        console.warn('Crash loop detected. Clearing storage and resetting.')
        
        // Backup Auth Data to prevent logout
        const token = localStorage.getItem('token')
        const username = localStorage.getItem('username')
        const role = localStorage.getItem('role')
        const sn = localStorage.getItem('sn')
        const avatar = localStorage.getItem('avatar')

        localStorage.clear()
        
        // Restore Auth Data
        if (token) localStorage.setItem('token', token)
        if (username) localStorage.setItem('username', username)
        if (role) localStorage.setItem('role', role)
        if (sn) localStorage.setItem('sn', sn)
        if (avatar) localStorage.setItem('avatar', avatar)

        localStorage.setItem(CRASH_KEY, '0')
        window.location.reload()
    }
}

// Reset crash count if app runs successfully for 5 seconds
setTimeout(() => {
    localStorage.setItem(CRASH_KEY, '0')
}, 5000)

const app = createApp(App)

// Global Error Handler
app.config.errorHandler = (err, instance, info) => {
    handleCrash(err)
}

window.onerror = function (message, source, lineno, colno, error) {
    handleCrash(error || message)
}

const pinia = createPinia()

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

app.use(ElementPlus, { size: 'small', zIndex: 3000 })
app.use(router)
app.use(pinia)
try {
    app.mount('#app')
} catch (e) {
    handleCrash(e)
}
