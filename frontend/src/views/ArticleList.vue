<template>
  <div>
    <div class="mb-6 flex justify-between items-center">
      <h2 class="text-xl font-semibold text-white">Articles</h2>
      <div class="space-x-2">
        <button @click="fetchArticles" class="px-3 py-1 bg-vscode-sidebar border border-vscode-border hover:bg-vscode-bg rounded text-sm transition-colors">
          Refresh
        </button>
        <router-link v-if="isAdmin" to="/admin/create" class="px-4 py-2 bg-vscode-primary hover:bg-opacity-90 text-white rounded text-sm font-medium transition-colors">
          + New Article
        </router-link>
      </div>
    </div>

    <!-- Error/Loading States -->
    <div v-if="loading" class="text-center py-10 text-gray-500">Loading articles...</div>
    <div v-else-if="error" class="bg-red-900/20 border border-red-900 text-red-300 p-4 rounded mb-6">
      {{ error }}
    </div>

    <!-- Table -->
    <div v-else class="bg-vscode-sidebar border border-vscode-border rounded-lg overflow-hidden">
      <table class="w-full text-left border-collapse">
        <thead class="bg-[#2d2d2d]">
          <tr>
            <th class="p-3 border-b border-vscode-border font-medium text-gray-400 w-16">SN</th>
            <th class="p-3 border-b border-vscode-border font-medium text-gray-400">Title</th>
            <th class="p-3 border-b border-vscode-border font-medium text-gray-400">Category</th>
            <th class="p-3 border-b border-vscode-border font-medium text-gray-400">Tags</th>
            <th class="p-3 border-b border-vscode-border font-medium text-gray-400 w-40">Created At</th>
            <th class="p-3 border-b border-vscode-border font-medium text-gray-400 w-32 text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="article in articles" :key="article.sn" class="hover:bg-vocab-bg/50 transition-colors group">
            <td class="p-3 border-b border-vscode-border border-opacity-50 text-gray-500 font-mono text-xs">
              {{ article.sn }}
            </td>
            <td class="p-3 border-b border-vscode-border border-opacity-50 font-medium text-vscode-text">
              {{ article.title }}
              <div class="text-xs text-gray-500 truncate max-w-xs mt-1">{{ article.abstract }}</div>
            </td>
            <td class="p-3 border-b border-vscode-border border-opacity-50 text-sm">
              <span class="px-2 py-0.5 rounded-full bg-blue-900/30 text-blue-300 border border-blue-900/50 text-xs">
                {{ article.category || 'Uncategorized' }}
              </span>
            </td>
            <td class="p-3 border-b border-vscode-border border-opacity-50 text-sm">
              <div class="flex flex-wrap gap-1">
                <span v-for="tag in parseTags(article.tags)" :key="tag" 
                  class="px-2 py-0.5 rounded bg-[#3e3e42] text-gray-300 text-xs text-[10px]"
                >
                  {{ tag }}
                </span>
              </div>
            </td>
            <td class="p-3 border-b border-vscode-border border-opacity-50 text-xs text-gray-500 font-mono">
              {{ formatDate(article.create_time) }}
            </td>
            <td class="p-3 border-b border-vscode-border border-opacity-50 text-right">
              <div v-if="isAdmin" class="flex justify-end gap-2 opacity-100 sm:opacity-0 group-hover:opacity-100 transition-opacity">
                <button @click="editArticle(article.sn)" class="text-vscode-primary hover:text-white" title="Edit">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
                <button @click="deleteArticle(article.sn)" class="text-red-400 hover:text-red-200" title="Delete">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="articles.length === 0">
            <td colspan="6" class="p-8 text-center text-gray-500">
              No articles found. Create one to get started.
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { getArticles, deleteArticles } from '../api/article'
import { authStore } from '../stores/auth'

const router = useRouter()
const isAdmin = computed(() => authStore.role === 1)
const articles = ref([])
const loading = ref(false)
const error = ref(null)

const fetchArticles = async () => {
  loading.value = true
  error.value = null
  try {
    const res = await getArticles({ page: 1, size: 100 })
    if (res.data.code === 10000) {
      const d = res.data.data
      if (Array.isArray(d)) {
         articles.value = d
      } else if (d.list) {
         articles.value = d.list
      } else {
         articles.value = []
      }
    } else {
      error.value = res.data.msg || 'Failed to fetch articles'
    }
  } catch (err) {
    console.error(err)
    error.value = 'Network error or backend unreachable.'
  } finally {
    loading.value = false
  }
}

const editArticle = (sn) => {
  router.push(`/admin/edit/${sn}`)
}

const deleteArticle = async (sn) => {
  if (!confirm('Are you sure you want to delete this article?')) return
  try {
    const res = await deleteArticles([sn])
    if (res.data.code === 10000) {
      // Remove from list locally
      articles.value = articles.value.filter(a => a.sn !== sn)
    } else {
      alert('Delete failed: ' + res.data.msg)
    }
  } catch (err) {
    console.error(err)
    alert('Delete error')
  }
}

const parseTags = (tags) => {
  if (Array.isArray(tags)) return tags
  // Usually ctype.Array is just []string in JSON
  return []
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

onMounted(() => {
  fetchArticles()
})
</script>
