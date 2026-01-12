import { reactive, watch } from 'vue'

export const authStore = reactive({
    token: localStorage.getItem('token') || '',
    username: localStorage.getItem('username') || 'Guest',
    role: parseInt(localStorage.getItem('role') || '1'),
    isLoggedIn: !!localStorage.getItem('token'),

    setAuth(token, username, role) {
        this.token = token
        this.username = username
        this.role = role
        this.isLoggedIn = true
        localStorage.setItem('token', token)
        localStorage.setItem('username', username)
        localStorage.setItem('role', role.toString())
    },

    clearAuth() {
        this.token = ''
        this.username = 'Guest'
        this.role = 1
        this.isLoggedIn = false
        localStorage.removeItem('token')
        localStorage.removeItem('username')
        localStorage.removeItem('role')
    }
})
