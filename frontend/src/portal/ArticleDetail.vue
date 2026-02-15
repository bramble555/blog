<template>
  <div v-loading="loading" class="max-w-4xl mx-auto px-4 sm:px-0">
    <template v-if="article">
      <!-- Article Header -->
      <header class="mb-12 text-center">
        <h1 class="text-4xl md:text-6xl font-extrabold mb-8 leading-tight text-text-primary tracking-tight">{{ article.title }}</h1>
        <div class="flex flex-wrap items-center justify-center gap-4 md:gap-8 text-base font-medium text-text-tertiary">
           <div class="flex items-center gap-2 px-3 py-1 rounded-full bg-bg-secondary border border-border-primary/50">
              <el-avatar :size="28" :src="formatUrl(article.user_avatar)" class="ring-2 ring-border-secondary">{{ article.username?.[0] }}</el-avatar>
              <span class="text-text-secondary">{{ article.username }}</span>
           </div>
           
           <div class="flex items-center gap-6">
             <span class="flex items-center gap-2"><el-icon><Calendar /></el-icon> {{ formatDate(article.create_time) }}</span>
             <span class="flex items-center gap-2"><el-icon><View /></el-icon> {{ article.look_count }} views</span>
           </div>

           <div class="flex items-center gap-4">
             <DiggButton 
                :sn="article.sn" 
                :count="article.digg_count" 
                :isDigg="article.is_digg"
                type="article"
                @update:count="val => article.digg_count = val"
                @update:isDigg="val => article.is_digg = val"
                class="scale-90"
             />
             
             <div 
               class="flex items-center gap-2 cursor-pointer transition-all hover:scale-110 active:scale-95 px-3 py-1.5 rounded-full bg-bg-secondary hover:bg-bg-tertiary border border-transparent hover:border-accent-primary/30" 
               @click="handleCollect" 
               title="Collect"
             >
                <el-icon :size="18" :class="article.is_collect ? 'text-yellow-500' : 'text-text-tertiary'">
                   <component :is="article.is_collect ? StarFilled : Star" />
                </el-icon>
                <span :class="article.is_collect ? 'text-text-primary' : 'text-text-tertiary'">{{ article.collects_count }}</span>
             </div>
           </div>
        </div>
      </header>

      <!-- Banner -->
      <div v-if="article.banner_url" class="rounded-3xl overflow-hidden mb-16 shadow-2xl shadow-indigo-500/10 ring-1 ring-white/10 group">
         <el-image :src="formatUrl(article.banner_url)" class="w-full bg-bg-tertiary transition-transform duration-700 group-hover:scale-105" fit="cover">
            <template #error>
               <div class="w-full h-96 bg-bg-secondary flex flex-col items-center justify-center text-text-tertiary">
                 <el-icon class="text-6xl mb-4 opacity-50"><PictureFilled /></el-icon>
                 <span class="text-lg">No Preview Available</span>
               </div>
             </template>
         </el-image>
      </div>

      <!-- Content -->
      <div class="mb-20">
         <MarkdownRenderer 
            :parsed="article.parsed_content" 
            class="prose prose-lg md:prose-xl max-w-none text-text-secondary"
         />
      </div>

      <!-- Tags -->
      <div v-if="tags.length" class="flex flex-wrap gap-3 mb-16">
         <el-tag 
            v-for="t in tags" 
            :key="t" 
            effect="dark"
            round
            class="!bg-bg-secondary !border-accent-primary/30 !text-accent-primary hover:!bg-accent-primary hover:!text-white transition-all cursor-pointer !px-4 !py-1.5 !text-sm"
         >
           # {{ t }}
         </el-tag>
      </div>

      <el-divider class="!border-border-primary/50 !my-12" />
      
      <AdCarousel class="mb-16 rounded-2xl overflow-hidden shadow-lg" />

      <!-- Comments Section -->
      <section class="mt-16 max-w-3xl mx-auto">
         <h3 class="text-3xl font-bold mb-10 flex items-center gap-3 text-text-primary">
            Comments <span class="px-3 py-1 text-sm font-bold bg-accent-primary/10 text-accent-primary rounded-full">{{ article.comment_count }}</span>
         </h3>

         <!-- Post Comment -->
         <div class="bg-bg-secondary/50 backdrop-blur-sm border border-border-primary/50 rounded-2xl p-8 mb-12 shadow-sm focus-within:shadow-lg focus-within:border-accent-primary/50 transition-all">
            <el-input 
               v-model="commentText" 
               type="textarea" 
               :rows="3" 
               placeholder="Share your thoughts..." 
               class="mb-6 !text-lg bg-transparent"
               input-style="background-color: transparent; box-shadow: none; color: var(--text-primary); padding: 0;"
               resize="none"
            />
            <div class="flex justify-between items-center">
               <span class="text-xs text-text-tertiary">Markdown supported</span>
               <el-button type="primary" size="large" class="!px-8 !rounded-xl !font-bold shadow-lg shadow-indigo-500/20" @click="postComment">Post Comment</el-button>
            </div>
         </div>

         <!-- Comment List -->
         <div class="space-y-8">
            <CommentItem 
               v-for="c in comments" 
               :key="c.sn" 
               :comment="c" 
               @reply-success="fetchComments"
            />
         </div>
      </section>
    </template>
    
    <el-result v-else-if="!loading" icon="error" title="404" sub-title="Article not found" class="mt-20">
       <template #extra>
          <el-button type="primary" @click="$router.push('/')">Back to Home</el-button>
       </template>
    </el-result>
  </div>
</template>

<script setup>
/**
 * ArticleDetail.vue
 * 
 * @description 门户文章详情页。展示文章内容、侧边栏广告以及评论区。
 * @author GVB Admin
 * @last_modified 2026-02-02
 * @requires vue, vue-router, ../api/article, ../api/comment, ../components/CommentItem, ../components/SidebarAds, @/utils/url
 */
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getArticle } from '../api/article'
import { getComments, createComment } from '../api/comment'
import { View, PictureFilled, Star, StarFilled, Calendar } from '@element-plus/icons-vue'
import { authStore } from '../stores/auth'
import { ElMessage } from 'element-plus'
import DiggButton from '@/components/DiggButton.vue'
import MarkdownRenderer from '@/components/MarkdownRenderer.vue'
import CommentItem from '@/components/CommentItem.vue'
import AdCarousel from '@/components/AdCarousel.vue'
import { formatUrl } from '@/utils/url'
import { formatDate } from '@/utils/date'
import { postCollect } from '../api/article'

const route = useRoute()
const article = ref(null)
const loading = ref(false)
const comments = ref([])
const commentText = ref('')

/**
 * 计算属性：解析文章标签
 * @description 处理后端返回的标签数据，支持数组或逗号分隔的字符串格式
 * @returns {Array<string>} 标签字符串数组
 */
const tags = computed(() => {
   if (!article.value?.tags) return []
   if (Array.isArray(article.value.tags)) {
      return article.value.tags
   }
   if (typeof article.value.tags === 'string') {
      return article.value.tags.split(',').map(t => t.trim()).filter(t => t)
   }
   return []
})

/**
 * 异步函数：获取文章详情数据
 * @description 根据路由参数中的 SN 获取文章详细信息，并同步收藏状态
 */
const fetchData = async () => {
   loading.value = true
   try {
      const res = await getArticle(route.params.sn)
      if (res.data.code === 10000) {
         article.value = res.data.data
         fetchComments()
      }
   } catch (e) {
      console.error('Detail load failed', e)
   } finally {
      loading.value = false
   }
}

/**
 * 异步函数：获取评论列表
 * @description 分页获取当前文章的评论数据
 */
const fetchComments = async () => {
   try {
      const res = await getComments(route.params.sn, { page: 1, size: 100 })
      if (res.data.code === 10000) {
         comments.value = res.data.data.list || []
      }
   } catch (e) {
      console.error('Fetch comments error', e)
   }
}

/**
 * 提交评论函数
 */
const postComment = async () => {
   if (!commentText.value) return
   
   if (!authStore.isLoggedIn) {
       ElMessage.warning('Please login first')
       return
   }

   try {
      const res = await createComment({
         article_sn: route.params.sn,
         content: commentText.value
      })
      
      if (res.data.code === 10000) {
         ElMessage.success('Comment posted successfully')
         commentText.value = ''
         fetchComments() 
      } else {
         ElMessage.error(res.data.msg)
      }
   } catch (e) {
      console.error('Post comment error:', e)
      ElMessage.error('Failed to post comment')
   }
}

/**
 * 处理收藏点击
 */
const handleCollect = async () => {
   if (!article.value) return
   if (!authStore.isLoggedIn) {
      ElMessage.warning('Please login first')
      return
   }
   
   try {
      const res = await postCollect({ sn: String(article.value.sn) })
      if (res.data.code === 10000) {
         article.value.is_collect = !article.value.is_collect
         article.value.collects_count += article.value.is_collect ? 1 : -1
         ElMessage.success(article.value.is_collect ? 'Collected successfully' : 'Uncollected successfully')
      } else {
         ElMessage.error(res.data.msg)
      }
   } catch (e) {
      console.error('Collect error:', e)
      ElMessage.error('Operation failed')
   }
}

onMounted(() => {
   fetchData()
})
</script>

<style scoped>
/* Global variables are available without explicit reference */

/* Markdown styles shim with Semantic Vars */
:deep(.prose) {
  color: var(--text-secondary);
}
:deep(.prose h1), :deep(.prose h2), :deep(.prose h3), :deep(.prose h4) {
   color: var(--text-primary); 
   margin-top: 2.5rem; 
   margin-bottom: 1.25rem; 
   font-weight: 800;
   letter-spacing: -0.025em;
   line-height: 1.25;
}
:deep(.prose h2) {
   border-left-width: 4px; 
   border-color: var(--accent-primary); 
   padding-left: 1.25rem;
}
:deep(.prose p) {
   margin-bottom: 1.5rem;
   line-height: 1.8;
}
:deep(.prose strong) {
   color: var(--text-primary);
   font-weight: 700;
}
:deep(.prose a) {
   color: var(--accent-primary);
   text-decoration: underline;
   text-decoration-color: var(--accent-primary);
   text-underline-offset: 4px;
   transition: all 0.2s;
}
:deep(.prose a:hover) {
   color: var(--accent-secondary);
   text-decoration-color: var(--accent-secondary);
}
:deep(.prose blockquote) {
   border-left-width: 4px; 
   border-color: var(--border-primary); 
   padding-left: 1.5rem; 
   padding-top: 0.5rem;
   padding-bottom: 0.5rem;
   background: var(--bg-secondary);
   border-radius: 0 0.5rem 0.5rem 0;
   font-style: italic; 
   color: var(--text-tertiary); 
   margin-top: 2rem; 
   margin-bottom: 2rem;
}
:deep(.prose code) {
   color: var(--accent-secondary);
   background-color: var(--bg-secondary);
   padding: 0.2rem 0.4rem;
   border-radius: 0.25rem;
   font-size: 0.875em;
   font-family: 'Fira Code', monospace;
}
:deep(.prose pre) {
   background-color: var(--bg-secondary);
   padding: 1.25rem;
   border-radius: 0.5rem;
   overflow-x: auto;
   margin-top: 1.5rem;
   margin-bottom: 1.5rem;
   border: 1px solid var(--border-primary);
}
:deep(.prose pre code) {
   background-color: transparent;
   padding: 0;
   color: inherit;
   font-size: 0.875em;
}
:deep(.prose img) {
   max-height: 650px;
   width: auto;
   margin: 2rem auto;
   border-radius: 1rem;
   max-width: 100%;
   object-fit: contain;
   box-shadow: 0 10px 30px -10px rgba(0,0,0,0.2);
   border: 1px solid var(--border-primary);
}
</style>
