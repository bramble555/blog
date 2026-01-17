<template>
  <div>
    <div class="mb-6 flex justify-between items-center">
      <h2 class="text-xl font-semibold text-white">Articles</h2>
      <div class="space-x-2">
        <button @click="fetchArticles" class="px-3 py-1 bg-vscode-sidebar border border-vscode-border hover:bg-vscode-bg rounded text-sm transition-colors">
          Refresh
        </button>
        <button
          v-if="isAdmin"
          :disabled="selectedSNList.length === 0"
          @click="deleteSelectedArticles"
          class="px-3 py-1 bg-red-600/30 border border-red-600/60 hover:bg-red-600/40 rounded text-sm transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
        >
          Delete Selected<span v-if="selectedSNList.length"> ({{ selectedSNList.length }})</span>
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
            <th v-if="isAdmin" class="p-3 border-b border-vscode-border font-medium text-[#FFA500] w-10">
              <input
                type="checkbox"
                :checked="isAllSelected"
                @change="toggleSelectAll"
                class="accent-[#FFA500] cursor-pointer"
              />
            </th>
            <th class="p-3 border-b border-vscode-border font-medium text-[#FFA500] w-16">SN</th>
            <th class="p-3 border-b border-vscode-border font-medium text-[#FFA500]">Title</th>
            <th class="p-3 border-b border-vscode-border font-medium text-[#FFA500]">Tags</th>
            <th class="p-3 border-b border-vscode-border font-medium text-[#FFA500] w-40">Created At</th>
            <th class="p-3 border-b border-vscode-border font-medium text-[#FFA500] w-32 text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="article in articles" :key="article.sn" class="hover:bg-vocab-bg/50 transition-colors group">
            <td v-if="isAdmin" class="p-3 border-b border-vscode-border border-opacity-50">
              <input
                type="checkbox"
                :checked="isSelected(article.sn)"
                @change="toggleSelected(article.sn)"
                class="accent-[#FFA500] cursor-pointer"
              />
            </td>
            <td class="p-3 border-b border-vscode-border border-opacity-50 text-[#FFA500] font-mono text-xs">
              {{ article.sn }}
            </td>
            <td class="p-3 border-b border-vscode-border border-opacity-50 font-medium text-vscode-text">
              {{ article.title }}
              <div class="text-xs text-[#FFA500] truncate max-w-xs mt-1">{{ article.abstract }}</div>
            </td>
            <td class="p-3 border-b border-vscode-border border-opacity-50 text-sm">
              <div class="flex flex-wrap gap-1">
                <span v-for="tag in parseTags(article.tags)" :key="tag" 
                  class="px-2 py-0.5 rounded bg-[#3e3e42] text-[#FFA500] text-2xs"
                >
                  {{ tag }}
                </span>
              </div>
            </td>
            <td class="p-3 border-b border-vscode-border border-opacity-50 text-xs text-[#FFA500] font-mono">
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
            <td :colspan="isAdmin ? 6 : 5" class="p-8 text-center text-gray-500">
              No articles found. Create one to get started.
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
/**
 * ArticleList.vue
 * 
 * @description 后台文章管理页面。提供文章的列表展示、搜索、编辑和删除功能。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires element-plus, vue, ../api/article, @/utils/date, @element-plus/icons-vue
 */
import { ref, reactive, onMounted, computed } from 'vue'
import { getArticles, deleteArticles } from '../api/article'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDateTime } from '@/utils/date'
import { Edit, Delete, Search } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { authStore } from '../stores/auth'

const router = useRouter()
const articles = ref([])
const loading = ref(false)
const error = ref(null)
const searchKeyword = ref('')
const selectedSNList = ref([])
const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

/**
 * 计算属性：是否为管理员
 * @returns {boolean}
 */
const isAdmin = computed(() => authStore.role === 1)

/**
 * 格式化日期
 * @param {string} date
 * @returns {string}
 */
const formatDate = (date) => {
  return formatDateTime(date)
}

/**
 * 解析标签
 * @param {string|Array} tags
 * @returns {Array}
 */
const parseTags = (tags) => {
  if (!tags) return []
  if (Array.isArray(tags)) return tags
  if (typeof tags === 'string') {
    try {
      // 尝试解析 JSON 字符串
      return JSON.parse(tags)
    } catch (e) {
      // 如果不是 JSON，尝试逗号分隔
      return tags.split(',').map(t => t.trim()).filter(t => t)
    }
  }
  return []
}

/**
 * 获取文章数据
 */
const fetchArticles = async () => {
  loading.value = true
  error.value = null
  try {
    const res = await getArticles({ 
      page: pagination.page, 
      size: pagination.size,
      keyword: searchKeyword.value
    })
    if (res.data.code === 10000) {
       articles.value = res.data.data.list
       pagination.total = res.data.data.count
       selectedSNList.value = []
    } else {
       error.value = res.data.msg || '获取文章失败'
    }
  } catch (err) {
    console.error(err)
    error.value = '获取文章失败'
    ElMessage.error('获取文章失败')
  } finally {
    loading.value = false
  }
}

/**
 * 处理搜索
 */
const handleSearch = () => {
  pagination.page = 1
  fetchArticles()
}

/**
 * 处理分页变化
 * @param {number} page - 新页码
 */
const handlePageChange = (page) => {
  pagination.page = page
  fetchArticles()
}

/**
 * 跳转编辑页面
 * @param {string} id - 文章 ID
 */
const editArticle = (id) => {
  router.push(`/admin/edit/${id}`)
}

/**
 * 删除文章
 * @param {string} id - 文章 ID
 */
const deleteArticle = (id) => {
  ElMessageBox.confirm('确定删除该文章吗?', '提示', {
    type: 'warning'
  }).then(async () => {
    const res = await deleteArticles([id])
    if (res.data.code === 10000) {
      ElMessage.success('删除成功')
      fetchArticles()
      selectedSNList.value = selectedSNList.value.filter(sn => sn !== id)
    }
  })
}

const isSelected = (sn) => {
  return selectedSNList.value.includes(sn)
}

const toggleSelected = (sn) => {
  if (isSelected(sn)) {
    selectedSNList.value = selectedSNList.value.filter(v => v !== sn)
    return
  }
  selectedSNList.value = [...selectedSNList.value, sn]
}

const isAllSelected = computed(() => {
  return articles.value.length > 0 && selectedSNList.value.length === articles.value.length
})

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedSNList.value = []
    return
  }
  selectedSNList.value = articles.value.map(a => a.sn)
}

const deleteSelectedArticles = async () => {
  if (selectedSNList.value.length === 0) return
  try {
    await ElMessageBox.confirm(`确定删除选中的 ${selectedSNList.value.length} 篇文章吗?`, '提示', { type: 'warning' })
    const res = await deleteArticles(selectedSNList.value)
    if (res.data.code === 10000) {
      ElMessage.success('删除成功')
      selectedSNList.value = []
      fetchArticles()
      return
    }
    ElMessage.error(res.data.msg || '删除失败')
  } catch (err) {
    if (err !== 'cancel') {
      console.error(err)
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  fetchArticles()
})
</script>
