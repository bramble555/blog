<template>
  <div class="min-h-screen bg-bg-primary text-text-primary transition-colors duration-300">
    <!-- Navbar -->
    <header class="sticky top-0 z-50 w-full border-b border-border-primary/40 bg-bg-primary/80 backdrop-blur-md shadow-sm transition-colors duration-300">
      <div class="container mx-auto flex h-16 items-center justify-between px-6">
        <!-- Logo Area -->
        <div class="flex items-center gap-10">
          <router-link to="/" class="group flex items-center space-x-3">
            <div class="relative flex items-center justify-center w-8 h-8 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 shadow-lg group-hover:shadow-indigo-500/30 transition-shadow">
              <span class="text-white font-bold text-lg select-none">G</span>
            </div>
            <span class="text-xl font-bold bg-gradient-to-r from-text-primary to-text-secondary bg-clip-text text-transparent group-hover:from-accent-primary group-hover:to-accent-secondary transition-all duration-300">
              GVB BLOG
            </span>
          </router-link>
          
          <!-- Desktop Nav -->
          <nav class="hidden md:flex items-center space-x-1">
            <router-link 
              to="/" 
              class="px-4 py-2 text-base font-medium text-text-secondary hover:text-accent-primary hover:bg-bg-secondary rounded-md transition-all duration-200"
              active-class="bg-bg-secondary text-accent-primary"
            >
              Home
            </router-link>
            <!-- Extensible nav items can go here -->
          </nav>
        </div>

        <!-- Right Side Actions -->
        <div class="flex items-center gap-4">
          <!-- Search Bar -->
          <div class="relative hidden sm:block group">
             <el-input
               v-model="searchQuery"
               placeholder="Search articles..."
               class="w-64 transition-all duration-300 focus-within:w-72"
               @keyup.enter="handleSearch"
             >
               <template #prefix>
                 <el-icon class="text-text-tertiary group-focus-within:text-accent-primary transition-colors"><Search /></el-icon>
               </template>
             </el-input>
          </div>
          
          <!-- Auth Actions -->
          <template v-if="!isLoggedIn">
            <router-link to="/login">
              <el-button link class="!text-text-secondary hover:!text-text-primary text-base">Login</el-button>
            </router-link>
            <router-link to="/register">
              <el-button type="primary" class="!font-medium !px-6 !rounded-lg shadow-lg shadow-indigo-500/20 hover:shadow-indigo-500/40 transition-all">
                Get Started
              </el-button>
            </router-link>
          </template>

          <!-- User Menu -->
          <template v-else>
             <div class="flex items-center gap-4">
                <router-link 
                  to="/admin/articles" 
                  class="text-base font-medium text-accent-primary hover:text-accent-secondary transition-colors"
                >
                  Dashboard
                </router-link>
                <button 
                  @click="handleLogout" 
                  class="text-base font-medium text-text-secondary hover:text-red-400 transition-colors cursor-pointer"
                >
                  Logout
                </button>
                <div class="h-6 w-px bg-border-primary mx-1"></div>
                <div class="flex items-center gap-3 pl-2 py-1 pr-1 rounded-full hover:bg-bg-secondary transition-colors cursor-pointer">
                   <el-avatar :size="32" :src="formatUrl(avatar)" :icon="UserFilled" class="ring-2 ring-border-secondary" />
                   <span class="text-base font-medium text-text-primary hidden lg:inline pr-2">{{ username }}</span>
                </div>
             </div>
          </template>
        </div>
      </div>
    </header>

    <!-- Content -->
    <main :class="['py-8 transition-all duration-300', $route.name === 'Home' ? 'w-full px-0' : 'container mx-auto px-6 max-w-6xl']">
       <router-view v-slot="{ Component }">
          <transition name="fade-slide-up" mode="out-in">
             <component :is="Component" :key="$route.fullPath" />
          </transition>
       </router-view>
    </main>

    <!-- Footer -->
    <footer class="border-t border-border-primary/40 py-12 mt-auto bg-bg-secondary/30 backdrop-blur-sm">
      <div class="container mx-auto px-6 text-center">
        <p class="text-base text-text-tertiary">
          © 2026 GVB Blog. <span class="text-text-secondary">Powered by Go & Vue 3.</span>
        </p>
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
import { formatUrl } from '@/utils/url'

const router = useRouter()
const searchQuery = ref('')

const isLoggedIn = computed(() => authStore.isLoggedIn)
const isAdmin = computed(() => authStore.role === 1)
const isUser = computed(() => authStore.role === 2)
const username = computed(() => authStore.username)
const avatar = computed(() => authStore.avatar)

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
