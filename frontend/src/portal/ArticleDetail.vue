<template>
  <div v-loading="loading" class="max-w-4xl mx-auto">
    <template v-if="article">
      <!-- Article Header -->
      <header class="mb-10 text-center">
        <h1 class="text-4xl md:text-5xl font-extrabold mb-6 leading-tight">{{ article.title }}</h1>
        <div class="flex items-center justify-center gap-6 text-sm text-[#FF6600]">
           <div class="flex items-center gap-2">
              <el-avatar :size="32" :src="article.user_avatar">{{ article.username?.[0] }}</el-avatar>
              <span>{{ article.username }}</span>
           </div>
           <span>•</span>
           <span>{{ formatDate(article.create_time) }}</span>
           <span>•</span>
           <span class="flex items-center gap-1"><el-icon><View /></el-icon> {{ article.look_count }} views</span>
           <span>•</span>
           <DiggButton 
              :sn="article.sn" 
              :count="article.digg_count" 
              :isDigg="article.is_digg"
              type="article"
              @update:count="val => article.digg_count = val"
              @update:isDigg="val => article.is_digg = val"
           />
        </div>
      </header>

      <!-- Banner -->
      <div v-if="article.banner_url" class="rounded-3xl overflow-hidden mb-12 border border-vscode-border shadow-2xl">
         <el-image :src="formatUrl(article.banner_url)" class="w-full bg-black/20">
            <template #error>
               <div class="w-full h-64 bg-vscode-sidebar flex flex-col items-center justify-center text-[#FF6600]">
                 <el-icon class="text-4xl mb-2"><PictureFilled /></el-icon>
                 <span>No Preview</span>
               </div>
             </template>
         </el-image>
      </div>

      <!-- Content -->
      <div class="mb-16">
         <MarkdownRenderer 
            :parsed="article.parsed_content" 
         />
      </div>

      <!-- Tags -->
      <div v-if="tags.length" class="flex gap-2 mb-16">
         <el-tag v-for="t in tags" :key="t" size="small" type="info"># {{ t }}</el-tag>
      </div>

      <el-divider />
      
      <AdCarousel class="mb-10" />

      <!-- Comments Section -->
      <section class="mt-16">
         <h3 class="text-2xl font-bold mb-8 flex items-center gap-2">
            Comments <span class="text-sm font-normal text-[#FF6600]">({{ article.comment_count }})</span>
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
         <div class="space-y-4">
            <CommentItem 
               v-for="c in comments" 
               :key="c.sn" 
               :comment="c" 
               @reply-success="fetchComments"
            />
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
/**
 * ArticleDetail.vue
 * 
 * @description 门户文章详情页。展示文章内容、侧边栏广告以及评论区。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires vue, vue-router, ../api/article, ../api/comment, ../components/CommentItem, ../components/SidebarAds, @/utils/url
 */
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getArticle } from '../api/article'
import { getComments, createComment } from '../api/comment'
import { View, PictureFilled } from '@element-plus/icons-vue'
import { authStore } from '../stores/auth'
import { ElMessage } from 'element-plus'
import DiggButton from '@/components/DiggButton.vue'
import MarkdownRenderer from '@/components/MarkdownRenderer.vue'
import CommentItem from '@/components/CommentItem.vue'
import { formatUrl } from '@/utils/url'
import { formatDate } from '@/utils/date'
// ... (imports)

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
   // 如果文章详情未加载或无标签，返回空数组
   if (!article.value?.tags) return []
   
   // 情况1：标签已经是数组格式
   if (Array.isArray(article.value.tags)) {
      return article.value.tags
   }
   
   // 情况2：标签是逗号分隔的字符串
   if (typeof article.value.tags === 'string') {
      // 分割字符串，去除首尾空格，并过滤空项
      return article.value.tags.split(',').map(t => t.trim()).filter(t => t)
   }
   
   // 其他情况返回空数组
   return []
})

/**
 * 异步函数：获取文章详情数据
 * @description 根据路由参数中的 SN 获取文章详细信息，并同步收藏状态
 * @async
 * @returns {Promise<void>}
 */
const fetchData = async () => {
   // 设置加载状态为 true，显示 loading 效果
   loading.value = true
   try {
      // 调用 API 获取文章详情
      const res = await getArticle(route.params.sn)
      
      // 检查 API 响应码
      if (res.data.code === 10000) {
         // 更新文章数据
         article.value = res.data.data
         
         // 获取该文章的评论列表
         fetchComments()
      }
   } catch (e) {
      // 捕获并记录错误
      console.error('Detail load failed', e)
   } finally {
      // 无论成功失败，结束加载状态
      loading.value = false
   }
}

/**
 * 异步函数：获取评论列表
 * @description 分页获取当前文章的评论数据
 * @async
 * @returns {Promise<void>}
 */
const fetchComments = async () => {
   try {
      // 调用 API 获取评论，默认获取第一页，100条
      const res = await getComments(route.params.sn, { page: 1, size: 100 })
      
      // 检查响应成功
      if (res.data.code === 10000) {
         // 更新评论列表数据
         comments.value = res.data.data.list || []
      }
   } catch (e) {
      // 忽略错误或按需处理
      console.error('Fetch comments error', e)
   }
}

/**
 * 提交评论函数
 * @description 校验输入内容，向后端发送创建评论请求，处理成功或失败的响应
 * @async
 * @returns {Promise<void>} 无返回值
 */
const postComment = async () => {
   // 1. 校验评论内容是否为空
   if (!commentText.value) return
   
   try {
      // 2. 发送创建评论请求
      // 参数包含：文章序列号(sn) 和 评论内容(content)
      const res = await createComment({
         article_sn: route.params.sn,
         content: commentText.value
      })
      
      // 3. 处理响应结果
      if (res.data.code === 10000) {
         // 成功：显示提示消息
         ElMessage.success('Comment posted successfully')
         // 清空输入框
         commentText.value = ''
         // 刷新评论列表以显示新评论
         fetchComments() 
      } else {
         // 失败：显示后端返回的错误消息
         ElMessage.error(res.data.msg)
      }
   } catch (e) {
      // 4. 捕获网络或其他未知错误
      console.error('Post comment error:', e)
      ElMessage.error('Failed to post comment')
   }
}

/**
 * 生命周期钩子：组件挂载完成
 * @description 组件加载时触发，开始获取文章详情数据
 */
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
   color: #FF6600; 
   margin-top: 1.5rem; 
   margin-bottom: 1.5rem;
}
.prose hr {
   display: none;
   border: none;
}

/* Limit content image size */
:deep(.prose img) {
   max-height: 600px;
   width: auto;
   margin: 1rem auto;
   border-radius: 0.5rem;
   max-width: 100%;
   object-fit: contain;
   background-color: rgba(0,0,0,0.2);
}
</style>
