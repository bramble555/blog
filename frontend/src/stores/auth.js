import { reactive, watch } from 'vue'

export const authStore = reactive({
    token: localStorage.getItem('token') || '',
    username: localStorage.getItem('username') || 'Guest',
    sn: localStorage.getItem('sn') || '0',  // Renamed from id to sn
    role: (() => {
        const r = localStorage.getItem('role');
        if (r === '管理员' || r === 'admin') return 1;
        if (r === '用户' || r === 'user') return 2;
        const parsed = parseInt(r);
        return isNaN(parsed) ? 2 : parsed;
    })(),
    isLoggedIn: !!localStorage.getItem('token'),

    setAuth(token, username, role, sn) {
        this.token = token
        this.username = username
        this.role = role
        this.sn = sn
        this.isLoggedIn = true
        localStorage.setItem('token', token)
        localStorage.setItem('username', username)
        localStorage.setItem('role', role.toString())
        localStorage.setItem('sn', sn)
    },

    clearAuth() {
        this.token = ''
        this.username = 'Guest'
        this.sn = '0'
        this.role = 2
        this.isLoggedIn = false
        localStorage.removeItem('token')
        localStorage.removeItem('username')
        localStorage.removeItem('role')
        localStorage.removeItem('sn')
    }
})
