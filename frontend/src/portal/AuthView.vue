<template>
  <div class="min-h-screen flex items-center justify-center bg-vscode-bg p-4">
    <div class="w-full max-w-md bg-vscode-sidebar border border-vscode-border rounded-3xl p-8 shadow-2xl">
      <div class="text-center mb-10">
         <h2 class="text-3xl font-bold text-white mb-2">{{ isRegister ? 'Create Account' : 'Welcome Back' }}</h2>
         <p class="text-[#FFA500]">{{ isRegister ? 'Join our community today' : 'Please enter your details' }}</p>
      </div>

      <el-form :model="form" label-position="top" size="large">
        <el-form-item label="Username">
          <el-input v-model="form.username" placeholder="Enter username" />
        </el-form-item>
        
        <el-form-item v-if="isRegister" label="Email">
           <el-input v-model="form.email" placeholder="you@example.com" />
        </el-form-item>

		<el-form-item v-if="isRegister" label="Verification Code">
		  <div class="flex gap-2 w-full">
			<el-input v-model="form.code" placeholder="Enter code" class="flex-1" />
			<el-button
			  type="primary"
			  @click="handleSendRegisterCode"
			  :disabled="codeCountdown > 0"
			>
			  {{ codeCountdown > 0 ? codeCountdown + 's' : 'Send Code' }}
			</el-button>
		  </div>
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
           <span class="text-[#FFA500]">
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
import { ref, reactive, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { login, register, sendRegisterCode } from '../api/user'
import { ElMessage } from 'element-plus'
import { authStore } from '../stores/auth'

const router = useRouter()
const isRegister = ref(false)
const loading = ref(false)

const form = reactive({
  username: '',
  password: '',
  email: '',
	code: ''
})

const codeCountdown = ref(0)
let timer = null

const handleSendRegisterCode = async () => {
	if (!form.email) {
		ElMessage.warning('Please input email first')
		return
	}
	if (codeCountdown.value > 0) {
		return
	}
	try {
		const res = await sendRegisterCode({ email: form.email })
		if (res.data.code === 10000) {
			ElMessage.success(res.data.msg || 'Verification code sent')
			codeCountdown.value = 60
			timer = setInterval(() => {
				codeCountdown.value--
				if (codeCountdown.value <= 0) {
					clearInterval(timer)
					timer = null
				}
			}, 1000)
		} else {
			ElMessage.error(res.data.msg)
		}
	} catch (e) {
		ElMessage.error('Failed to send code')
	}
}

const handleSubmit = async () => {
  if (!form.username || !form.password || (isRegister.value && (!form.email || !form.code))) return
  loading.value = true
  try {
    if (isRegister.value) {
	   const res = await register({ username: form.username, password: form.password, email: form.email, code: form.code })
	   if (res.data.code === 10000) {
		  const { token, username, role, sn, avatar } = res.data.data
		  authStore.setAuth(token, username, role, sn, avatar)
		  ElMessage.success('Registration successful')
		  router.push('/')
	   } else {
		  ElMessage.error(res.data.msg)
	   }
    } else {
       const res = await login({ username: form.username, password: form.password })
       if (res.data.code === 10000) {
          const { token, username, role, sn, avatar } = res.data.data
          authStore.setAuth(token, username, role, sn, avatar) 
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

onUnmounted(() => {
	if (timer) {
		clearInterval(timer)
		timer = null
	}
})
</script>
