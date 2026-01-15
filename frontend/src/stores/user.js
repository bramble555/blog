import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
    const token = ref(localStorage.getItem('token') || '')
    const username = ref(localStorage.getItem('username') || 'Guest')
    const sn = ref(localStorage.getItem('sn') || '0')
    const role = ref(parseInt(localStorage.getItem('role') || '2'))

    const isLoggedIn = computed(() => !!token.value)
    const isAdmin = computed(() => role.value === 1)

    function setAuth(t, u, r, s) {
        token.value = t
        username.value = u
        role.value = r
        sn.value = s
        
        localStorage.setItem('token', t)
        localStorage.setItem('username', u)
        localStorage.setItem('role', r)
        localStorage.setItem('sn', s)
    }

    function clearAuth() {
        token.value = ''
        username.value = 'Guest'
        role.value = 2
        sn.value = '0'
        
        localStorage.removeItem('token')
        localStorage.removeItem('username')
        localStorage.removeItem('role')
        localStorage.removeItem('sn')
    }

    return { token, username, sn, role, isLoggedIn, isAdmin, setAuth, clearAuth }
})
