<template>
  <div v-loading="loading" class="max-w-4xl mx-auto">
    <template v-if="article">
      <!-- Article Header -->
      <header class="mb-10 text-center">
        <el-tag class="mb-4" effect="dark">{{ article.category }}</el-tag>
        <h1 class="text-4xl md:text-5xl font-extrabold mb-6 leading-tight">{{ article.title }}</h1>
        <div class="flex items-center justify-center gap-6 text-sm text-gray-400">
           <div class="flex items-center gap-2">
              <el-avatar :size="32" :src="article.user_avatar">{{ article.username?.[0] }}</el-avatar>
              <span>{{ article.username }}</span>
           </div>
           <span>•</span>
           <span>{{ formatDate(article.create_time) }}</span>
           <span>•</span>
           <span class="flex items-center gap-1"><el-icon><View /></el-icon> {{ article.look_count }} views</span>
           <span>•</span>
           <button 
             @click="handleCollect" 
             class="flex items-center gap-1 hover:text-yellow-400 transition-colors cursor-pointer"
             :class="{ 'text-yellow-500': isCollected }"
           >
             <el-icon><Star v-if="!isCollected" /><StarFilled v-else /></el-icon>
             {{ isCollected ? 'Collected' : 'Collect' }} ({{ article.collects_count }})
           </button>
        </div>
      </header>

      <!-- Banner -->
      <div v-if="article.banner_url" class="rounded-3xl overflow-hidden mb-12 border border-vscode-border shadow-2xl">
         <el-image :src="article.banner_url" fit="cover" class="w-full max-h-[500px]" />
      </div>

      <!-- Content -->
      <article class="prose prose-invert prose-blue max-w-none mb-16 text-gray-300 leading-relaxed text-lg">
         <!-- Simple Markdown Display (User asked to use textarea before, but here we render) -->
         <div v-html="renderMarkdown(article.content)"></div>
      </article>

      <!-- Tags -->
      <div v-if="tags.length" class="flex gap-2 mb-16">
         <el-tag v-for="t in tags" :key="t" size="small" type="info"># {{ t }}</el-tag>
      </div>

      <el-divider />

      <!-- Comments Section -->
      <section class="mt-16">
         <h3 class="text-2xl font-bold mb-8 flex items-center gap-2">
            Comments <span class="text-sm font-normal text-gray-500">({{ article.comment_count }})</span>
         </h3>

         <!-- Post Comment -->
         <div class="bg-vscode-sidebar border border-vscode-border rounded-2xl p-6 mb-12">
            <el-input 
               v-model="commentText" 
               type="textarea" 
               :rows="3" 
               placeholder="Write your thought..." 
               class="mb-4"
            />
            <div class="flex justify-end">
               <el-button type="primary" @click="postComment">Post Comment</el-button>
            </div>
         </div>

         <!-- Comment List -->
         <div class="space-y-8">
            <div v-for="c in comments" :key="c.sn" class="flex gap-4">
               <el-avatar :src="c.user_avatar">{{ c.username?.[0] }}</el-avatar>
               <div class="flex-1">
                  <div class="flex items-center gap-2 mb-1">
                     <span class="font-bold text-vscode-primary">{{ c.username }}</span>
                     <span class="text-xs text-gray-500">{{ formatDate(c.create_time) }}</span>
                  </div>
                  <p class="text-gray-300">{{ c.content }}</p>
                  <div class="mt-2 flex items-center gap-4 text-xs text-gray-500">
                     <span class="cursor-pointer hover:text-white flex items-center gap-1">
                        <el-icon><Pointer /></el-icon> {{ c.digg_count }}
                     </span>
                     <span class="cursor-pointer hover:text-white">Reply</span>
                  </div>
               </div>
            </div>
         </div>
      </section>
    </template>
    
    <el-result v-else-if="!loading" icon="error" title="404" sub-title="Article not found">
       <template #extra>
          <el-button type="primary" @click="$router.push('/')">Back to Home</el-button>
       </template>
    </el-result>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getArticle, collectArticle } from '../api/article'
import { getComments, createComment } from '../api/comment'
import { View, Pointer, Star, StarFilled } from '@element-plus/icons-vue'
import { authStore } from '../stores/auth'
import { collectionStore } from '../stores/collection'
import { ElMessage } from 'element-plus'

const route = useRoute()
const article = ref(null)
const loading = ref(false)
const comments = ref([])
const commentText = ref('')

const tags = computed(() => {
   if (!article.value?.tags) return []
   // Handle both array (from JSON) and string (fallback)
   if (Array.isArray(article.value.tags)) {
      return article.value.tags
   }
   if (typeof article.value.tags === 'string') {
      return article.value.tags.split(',').map(t => t.trim()).filter(t => t)
   }
   return []
})

const fetchData = async () => {
   loading.value = true
   try {
      const res = await getArticle(route.params.sn)
      if (res.data.code === 10000) {
         article.value = res.data.data
         // Initialize state from backend, sync with global store
         if (article.value.is_collect) {
            collectionStore.add(article.value.sn)
         } else {
            collectionStore.remove(article.value.sn)
         }
         fetchComments()
      }
   } catch (e) {
      console.error('Detail load failed')
   } finally {
      loading.value = false
   }
}

const fetchComments = async () => {
   try {
      const res = await getComments({ article_sn: route.params.sn })
      if (res.data.code === 10000) {
         comments.value = res.data.data
      }
   } catch (e) {}
}

const postComment = async () => {
   if (!commentText.value) return
   try {
      const res = await createComment({
         article_sn: route.params.sn,
         content: commentText.value
      })
      if (res.data.code === 10000) {
         ElMessage.success('Comment posted!')
         commentText.value = ''
         fetchComments() // Refresh list
      } else {
         ElMessage.error(res.data.msg)
      }
   } catch (e) {
      ElMessage.error('Post failed')
   }
}

const isCollected = computed(() => {
   return article.value ? collectionStore.isCollected(article.value.sn) : false
})
const collectLoading = ref(false)

const handleCollect = async () => {
   if (!authStore.isLoggedIn) {
      ElMessage.warning('Please login first')
      return
   }
   if (collectLoading.value || !article.value) return
   
   // Backup current state for potential rollback
   const previousIsCollected = isCollected.value
   const previousCollectsCount = article.value.collects_count

   // Optimistic Update: Immediately toggle store state
   if (!previousIsCollected) {
      article.value.collects_count++
      collectionStore.add(article.value.sn)
   } else {
      article.value.collects_count--
      collectionStore.remove(article.value.sn)
   }

   collectLoading.value = true
   try {
      const res = await collectArticle(article.value.sn)
      if (res.data.code === 10000) {
         // Success: Keep the optimistic changes and show feedback
         const msg = res.data.data
         ElMessage.success(msg === '收藏成功' ? 'Collected successfully' : 'Uncollected successfully')
      } else {
         // Failure: Rollback to previous state
         article.value.collects_count = previousCollectsCount
         if (previousIsCollected) {
            collectionStore.add(article.value.sn)
         } else {
            collectionStore.remove(article.value.sn)
         }
         ElMessage.error(res.data.msg || 'Action failed')
      }
   } catch (e) {
      // Network/Server error: Rollback to previous state
      article.value.collects_count = previousCollectsCount
      if (previousIsCollected) {
         collectionStore.add(article.value.sn)
      } else {
         collectionStore.remove(article.value.sn)
      }
      console.error('Collect error:', e)
      ElMessage.error('Network error or server busy')
   } finally {
      collectLoading.value = false
   }
}

const formatDate = (d) => d ? new Date(d).toLocaleString() : ''

// Simple Markdown placeholder renderer
const renderMarkdown = (content) => {
   if (!content) return ''
   // TODO: Use marked or other lib. For now, preserve line breaks.
   return content.replace(/\n/g, '<br>')
}

onMounted(() => {
   fetchData()
})
</script>

<style scoped>
@reference "../style.css";

/* Markdown styles shim */
.prose h1, .prose h2, .prose h3 {
   color: white; 
   margin-top: 2rem; 
   margin-bottom: 1rem; 
   border-left-width: 4px; 
   border-color: var(--color-vscode-primary); 
   padding-left: 1rem;
}
.prose blockquote {
   border-left-width: 4px; 
   border-color: #4b5563; 
   padding-left: 1rem; 
   font-style: italic; 
   color: #9ca3af; 
   margin-top: 1.5rem; 
   margin-bottom: 1.5rem;
}
</style>
