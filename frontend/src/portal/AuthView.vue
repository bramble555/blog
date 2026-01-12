<template>
  <div class="min-h-screen flex items-center justify-center bg-vscode-bg p-4">
    <div class="w-full max-w-md bg-vscode-sidebar border border-vscode-border rounded-3xl p-8 shadow-2xl">
      <div class="text-center mb-10">
         <h2 class="text-3xl font-bold text-white mb-2">{{ isRegister ? 'Create Account' : 'Welcome Back' }}</h2>
         <p class="text-gray-400">{{ isRegister ? 'Join our community today' : 'Please enter your details' }}</p>
      </div>

      <el-form :model="form" label-position="top" size="large">
        <el-form-item label="Username">
          <el-input v-model="form.username" placeholder="Enter username" />
        </el-form-item>
        
        <el-form-item v-if="isRegister" label="Email">
           <el-input v-model="form.email" placeholder="you@example.com" />
        </el-form-item>

        <el-form-item label="Password">
          <el-input 
            v-model="form.password" 
            type="password" 
            placeholder="••••••••" 
            show-password 
          />
        </el-form-item>

        <div class="mt-8">
           <el-button type="primary" class="w-full" :loading="loading" @click="handleSubmit">
             {{ isRegister ? 'Register' : 'Login' }}
           </el-button>
        </div>

        <div class="mt-6 text-center text-sm">
           <span class="text-gray-500">
             {{ isRegister ? 'Already have an account?' : "Don't have an account?" }}
           </span>
           <el-button link type="primary" @click="isRegister = !isRegister">
             {{ isRegister ? 'Sign in' : 'Register' }}
           </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '../api/user'
import { ElMessage } from 'element-plus'
import { authStore } from '../stores/auth'

const router = useRouter()
const isRegister = ref(false)
const loading = ref(false)

const form = reactive({
  username: '',
  password: '',
  email: ''
})

const handleSubmit = async () => {
  if (!form.username || !form.password) return
  loading.value = true
  try {
    if (isRegister.value) {
       // TODO: Implement Register API in user.js
       ElMessage.success('Register logic not fully implemented in API yet.')
    } else {
       const res = await login({ username: form.username, password: form.password })
       if (res.data.code === 10000) {
          // In a real app, the backend should return user info (role, etc.)
          // For now, we assume role 2 for 'admin' and 1 for others
          const role = form.username === 'admin' ? 2 : 1
          authStore.setAuth(res.data.data, form.username, role) 
          ElMessage.success('Welcome back!')
          router.push('/')
       } else {
          ElMessage.error(res.data.msg)
       }
    }
  } catch (e) {
    ElMessage.error('Auth error')
  } finally {
    loading.value = false
  }
}
</script>
