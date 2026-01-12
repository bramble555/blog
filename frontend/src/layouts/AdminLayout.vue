<template>
  <div class="flex h-screen bg-vscode-bg text-vscode-text overflow-hidden">
    <!-- Sidebar -->
    <aside class="w-64 bg-vscode-sidebar border-r border-vscode-border flex flex-col">
      <div class="p-4 text-xl font-bold text-vscode-primary uppercase tracking-wider">
        Blog Admin
      </div>
      <nav class="flex-1 overflow-y-auto py-4">
        <div class="px-4 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 mt-2">Content</div>
        <router-link to="/admin/articles" class="nav-item" active-class="nav-active">
          <el-icon class="mr-3"><Document /></el-icon> Articles
        </router-link>
        <router-link to="/admin/comments" class="nav-item" active-class="nav-active">
          <el-icon class="mr-3"><ChatLineSquare /></el-icon> Comments
        </router-link>
        <router-link to="/admin/tags" class="nav-item" active-class="nav-active">
          <el-icon class="mr-3"><CollectionTag /></el-icon> Tags
        </router-link>

        <div class="px-4 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 mt-4">Messages</div>
        <router-link to="/admin/messages" class="nav-item" active-class="nav-active">
           <el-icon class="mr-3"><Bell /></el-icon> Messages
        </router-link>
         <router-link to="/admin/chat" class="nav-item" active-class="nav-active">
           <el-icon class="mr-3"><ChatDotRound /></el-icon> Chat
        </router-link>

        <div class="px-4 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 mt-4">Marketing</div>
        <router-link to="/admin/banners" class="nav-item" active-class="nav-active">
          <el-icon class="mr-3"><Picture /></el-icon> Banners
        </router-link>
        <router-link to="/admin/adverts" class="nav-item" active-class="nav-active">
          <el-icon class="mr-3"><DataBoard /></el-icon> Adverts
        </router-link>

        <div class="px-4 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 mt-4">System</div>
        <router-link to="/admin/users" class="nav-item" active-class="nav-active">
           <el-icon class="mr-3"><User /></el-icon> Users
        </router-link>

      </nav>
      <div class="p-4 border-t border-vscode-border text-xs text-gray-500 text-center">
        GVB Admin v1.0
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col min-w-0 overflow-hidden bg-vscode-bg">
      <header class="h-14 bg-vscode-bg border-b border-vscode-border flex items-center px-6 justify-between">
        <h1 class="text-lg font-medium text-white">{{ currentPageTitle }}</h1>
        <div class="flex items-center gap-6">
            <router-link to="/" class="text-xs text-vscode-text/60 hover:text-white flex items-center gap-1 transition-colors">
              <el-icon><Back /></el-icon> Exit to Portal
            </router-link>
            <div class="h-4 w-px bg-vscode-border"></div>
            <div class="flex items-center gap-3">
               <span class="text-xs font-semibold text-vscode-primary">ADMIN</span>
               <el-avatar :size="28" src="" icon="UserFilled" class="border border-vscode-border" />
            </div>
        </div>
      </header>
      
      <div class="flex-1 overflow-auto p-6">
        <router-view></router-view>
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { 
  Document, ChatLineSquare, CollectionTag, Picture, 
  DataBoard, User, Bell, ChatDotRound, Back, UserFilled 
} from '@element-plus/icons-vue'

const route = useRoute()
const currentPageTitle = computed(() => route.meta?.title || route.name || 'Dashboard')
</script>

<style scoped>
.nav-item {
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  transition-property: color, background-color, border-color, text-decoration-color, fill, stroke;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 150ms;
  border-left-width: 2px;
  border-color: transparent;
  color: #9ca3af; /* text-gray-400 */
}
.nav-item:hover {
  background-color: var(--color-vscode-bg);
  color: white;
}
.nav-active {
  background-color: var(--color-vscode-bg);
  border-color: var(--color-vscode-primary);
  color: white;
}
</style>
