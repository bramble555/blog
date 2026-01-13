<template>
  <div class="pt-4">

    <!-- Article Grid -->
    <div v-loading="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
      <div 
        v-for="article in articles" 
        :key="article.sn"
        class="group bg-vscode-sidebar rounded-2xl border border-vscode-border overflow-hidden hover:scale-[1.02] transition-all cursor-pointer shadow-lg hover:shadow-vscode-primary/10"
        @click="$router.push(`/article/${article.sn}`)"
      >
        <div class="aspect-video overflow-hidden">
          <el-image 
            :src="article.banner_url" 
            fit="cover"
            class="w-full h-full group-hover:scale-110 transition-transform duration-500"
          >
            <template #placeholder>
              <div class="w-full h-full bg-vscode-border animate-pulse flex items-center justify-center">
                <el-icon class="text-3xl text-gray-600"><Picture /></el-icon>
              </div>
            </template>
            <template #error>
              <div class="w-full h-full bg-vscode-border flex flex-col items-center justify-center text-gray-500">
                <el-icon class="text-3xl mb-1"><PictureFilled /></el-icon>
                <span class="text-[10px]">No Preview</span>
              </div>
            </template>
          </el-image>
        </div>
        <div class="p-6">
          <div class="flex items-center gap-2 mb-3">
             <el-tag size="small" effect="plain">{{ article.category }}</el-tag>
             <span class="text-xs text-gray-500">{{ formatDate(article.create_time) }}</span>
          </div>
          <h3 class="text-xl font-bold mb-3 group-hover:text-vscode-primary transition-colors line-clamp-2">
            {{ article.title }}
          </h3>
          <p class="text-gray-300 text-sm line-clamp-3 mb-4 leading-relaxed group-hover:text-white transition-colors">
            {{ article.abstract }}
          </p>
          <div class="flex items-center justify-between text-xs text-gray-400">
             <div class="flex items-center gap-3">
                <span class="flex items-center gap-1 hover:text-vscode-primary cursor-pointer"><el-icon><View /></el-icon> {{ article.look_count }}</span>
                <span class="flex items-center gap-1 hover:text-vscode-primary cursor-pointer"><el-icon><ChatLineSquare /></el-icon> {{ article.comment_count }}</span>
             </div>
             <span class="font-medium text-vscode-primary/80">By {{ article.username || 'Anonymous' }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Empty State -->
    <el-empty v-if="!loading && articles.length === 0" description="No articles found." />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getArticles } from '../api/article'
import { View, ChatLineSquare, Picture, PictureFilled } from '@element-plus/icons-vue'

const route = useRoute()
const articles = ref([])
const loading = ref(false)

const fetchData = async () => {
  loading.value = true
  try {
    const params = { 
      page: 1, 
      size: 20,
      title: route.query.title || ''
    }
    const res = await getArticles(params)
    if (res.data.code === 10000) {
       const d = res.data.data
       articles.value = Array.isArray(d) ? d : (d.list || [])
    }
  } catch (e) {
    console.error('Portal Home load failed')
  } finally {
    loading.value = false
  }
}

watch(() => route.query.title, () => fetchData())

const formatDate = (dateStr) => {
   if (!dateStr) return ''
   return new Date(dateStr).toLocaleDateString()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
