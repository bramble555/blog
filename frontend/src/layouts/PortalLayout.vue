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
                  to="/admin/articles" 
                  class="text-sm font-medium text-[#FF6600] hover:text-vscode-primary transition-colors"
                >
                  Dashboard
                </router-link>
                <button 
                  @click="handleLogout" 
                  class="text-sm font-medium text-[#FF6600] hover:text-red-400 transition-colors cursor-pointer"
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
    <main :class="['py-8', $route.name === 'Home' ? 'w-full px-0' : 'container mx-auto px-4 max-w-6xl']">
       <router-view :key="$route.fullPath"></router-view>
    </main>

    <!-- Footer -->
    <footer class="border-t border-vscode-border py-8 mt-auto">
      <div class="container mx-auto px-4 text-center text-sm text-[#FF6600]">
        © 2026 GVB Blog. Powered by Go & Vue 3.
      </div>
    </footer>
  </div>
</template>

<script setup>
/**
 * PortalLayout.vue
 * 
 * @description 门户前台布局组件。包含顶部导航栏、主体内容区域和底部页脚。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires vue, vue-router, @element-plus/icons-vue, ../stores/auth
 */
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Search, UserFilled } from '@element-plus/icons-vue'
import { authStore } from '../stores/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const searchQuery = ref('')

const isLoggedIn = computed(() => authStore.isLoggedIn)
const isAdmin = computed(() => authStore.role === 1)
const isUser = computed(() => authStore.role === 2)
const username = computed(() => authStore.username)

/**
 * 处理搜索跳转
 */
const handleSearch = () => {
  router.push({
    path: '/',
    query: { title: searchQuery.value }
  })
}

/**
 * 处理用户退出登录
 */
const handleLogout = () => {
  authStore.clearAuth()
  router.push('/')
}
</script>

<style>
/* Global Portal Styles if needed */
</style>
