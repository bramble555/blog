<template>
  <div class="min-h-screen bg-vscode-bg text-vscode-text">
    <!-- Navbar -->
    <header class="sticky top-0 z-50 w-full border-b border-vscode-border bg-vscode-bg/95 backdrop-blur">
      <div class="container mx-auto flex h-16 items-center justify-between px-4">
        <div class="flex items-center gap-8">
          <router-link to="/" class="flex items-center space-x-2">
            <span class="text-xl font-bold bg-gradient-to-r from-blue-500 to-purple-500 bg-clip-text text-transparent">
              GVB BLOG
            </span>
          </router-link>
          <nav class="hidden md:flex items-center space-x-6 text-sm font-medium">
            <router-link to="/" class="hover:text-vscode-primary transition-colors">Home</router-link>
          </nav>
        </div>

        <div class="flex items-center gap-4">
          <div class="relative hidden sm:block">
             <el-input
               v-model="searchQuery"
               placeholder="Search articles..."
               size="small"
               prefix-icon="Search"
               class="w-64"
               @keyup.enter="handleSearch"
             />
          </div>
          
          <template v-if="!isLoggedIn">
            <router-link to="/login">
              <el-button link class="text-vscode-text">Login</el-button>
            </router-link>
            <router-link to="/register">
              <el-button type="primary" size="small">Get Started</el-button>
            </router-link>
          </template>
          <template v-else>
             <div class="flex items-center gap-4">
                <router-link 
                  v-if="isAdmin" 
                  to="/admin/articles" 
                  class="text-sm font-medium text-gray-400 hover:text-vscode-primary transition-colors"
                >
                  Dashboard
                </router-link>
                <button 
                  @click="handleLogout" 
                  class="text-sm font-medium text-gray-400 hover:text-red-400 transition-colors cursor-pointer"
                >
                  Logout
                </button>
                <div class="h-6 w-px bg-vscode-border mx-1"></div>
                <div class="flex items-center gap-2">
                   <el-avatar :size="32" icon="UserFilled" />
                   <span class="text-sm font-medium hidden lg:inline">{{ username }}</span>
                </div>
             </div>
          </template>
        </div>
      </div>
    </header>

    <!-- Content -->
    <main class="container mx-auto px-4 py-8 max-w-6xl">
       <router-view :key="$route.fullPath"></router-view>
    </main>

    <!-- Footer -->
    <footer class="border-t border-vscode-border py-8 mt-auto">
      <div class="container mx-auto px-4 text-center text-sm text-gray-500">
        © 2026 GVB Blog. Powered by Go & Vue 3.
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Search, UserFilled } from '@element-plus/icons-vue'
import { authStore } from '../stores/auth'

const router = useRouter()
const searchQuery = ref('')

const isLoggedIn = computed(() => authStore.isLoggedIn)
const isAdmin = computed(() => authStore.role === 2)
const username = computed(() => authStore.username)

const handleSearch = () => {
  router.push({
    path: '/',
    query: { title: searchQuery.value }
  })
}

const handleLogout = () => {
  authStore.clearAuth()
  router.push('/')
}
</script>

<style>
/* Global Portal Styles if needed */
</style>
