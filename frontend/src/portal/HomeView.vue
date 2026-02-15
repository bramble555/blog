<template>
  <!-- 
    主页布局容器 
    - 使用语义化背景色
    - 响应式布局优化
  -->
  <div class="flex flex-col lg:flex-row gap-8 min-h-screen">
    <!-- 侧边栏 (广告) -->
    <aside class="hidden lg:block w-80 flex-shrink-0 sticky top-24 h-fit">
      <SidebarAds />
    </aside>

    <!-- 文章列表区域 -->
    <div v-loading="loading" class="flex-1 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 xl:grid-cols-2 gap-6 items-start content-start">
      <!-- 单个文章卡片 -->
      <div 
        v-for="article in articles" 
        :key="article.sn"
        class="group relative flex flex-col bg-bg-secondary rounded-xl overflow-hidden border border-border-primary/50 hover:border-accent-primary/50 transition-all duration-300 hover:shadow-xl hover:shadow-indigo-500/10 hover:-translate-y-1 cursor-pointer h-full"
        @click="$router.push(`/article/${article.sn}`)"
      >
        <!-- 文章封面图 -->
        <div class="h-56 w-full flex-shrink-0 overflow-hidden relative">
          <!-- Image Overlay Gradient -->
          <div class="absolute inset-0 bg-gradient-to-t from-bg-secondary via-transparent to-transparent opacity-60 z-10"></div>
          
          <el-image 
            :src="formatUrl(article.banner_url)" 
            fit="cover" 
            class="w-full h-full transition-transform duration-700 group-hover:scale-110"
            loading="lazy"
            :alt="article.title"
          >
            <template #placeholder>
              <div class="w-full h-full bg-bg-tertiary animate-pulse flex items-center justify-center">
                <el-icon class="text-3xl text-text-tertiary"><Picture /></el-icon>
              </div>
            </template>
            <template #error>
              <div class="w-full h-full bg-bg-tertiary flex flex-col items-center justify-center text-text-tertiary">
                <el-icon class="text-3xl mb-1"><PictureFilled /></el-icon>
                <span class="text-xs">No Preview</span>
              </div>
            </template>
          </el-image>
          
          <!-- Floating Date Badge -->
          <div class="absolute top-3 right-3 z-20 bg-bg-primary/90 backdrop-blur-sm px-3 py-1 rounded-full text-xs font-mono text-text-secondary border border-border-primary shadow-sm">
             {{ formatDate(article.create_time).split(' ')[0] }}
          </div>
        </div>
        
        <!-- 文章内容区域 -->
        <div class="p-6 flex-1 flex flex-col">
          <!-- 标题 -->
          <h3 class="text-xl font-bold mb-3 text-text-primary group-hover:text-accent-primary transition-colors line-clamp-2 leading-tight">
            {{ article.title }}
          </h3>
          
          <!-- 摘要 -->
          <p class="text-text-secondary text-base line-clamp-3 mb-6 leading-relaxed flex-1">
            {{ article.abstract }}
          </p>
          
          <!-- 底部元数据 -->
          <div class="flex items-center justify-between text-sm text-text-tertiary pt-4 border-t border-border-primary/50 mt-auto">
             <div class="flex items-center gap-4">
                <span class="flex items-center gap-1.5 hover:text-accent-primary transition-colors"><el-icon><View /></el-icon> {{ article.look_count }}</span>
                <span class="flex items-center gap-1.5 hover:text-accent-primary transition-colors"><el-icon><ChatLineSquare /></el-icon> {{ article.comment_count }}</span>
             </div>
             
             <div class="flex items-center gap-2">
                <span class="text-xs font-medium text-text-secondary">By {{ article.username || 'Admin' }}</span>
             </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 空状态 -->
    <div v-if="!loading && articles.length === 0" class="flex-1 flex flex-col items-center justify-center min-h-[400px] text-text-tertiary">
      <el-empty description="暂无文章" :image-size="200" />
    </div>
  </div>
</template>

<script setup>
/**
 * HomeView.vue
 * 
 * @description 门户首页组件。展示文章列表和侧边栏广告。
 * 按照用户要求修改：
 * 1. 去除文章间距 (gap-0)
 * 2. 去除黑色背景 (使用 bg-white, text-gray-*)
 * 3. 每一行一篇文章 (grid-cols-1, flex-row layout)
 * 
 * @author GVB Admin
 * @last_modified 2026-01-14
 */
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticles } from '../api/article'
import { View, ChatLineSquare, Picture, PictureFilled } from '@element-plus/icons-vue'
import SidebarAds from '@/components/SidebarAds.vue'
import { formatDate } from '@/utils/date'
import { formatUrl } from '@/utils/url'

const route = useRoute()
const router = useRouter()
// 文章列表数据
const articles = ref([])
// 加载状态
const loading = ref(false)

/**
 * 获取文章列表数据
 * @description 根据 URL 查询参数获取文章列表，支持分页和标题搜索
 * @async
 */
const fetchData = async () => {
  loading.value = true
  try {
    // 构建查询参数
    const params = { 
      page: 1, 
      size: 10,
      title: route.query.title || ''
    }
    // 调用 API
    const res = await getArticles(params)
    if (res.data.code === 10000) {
       const d = res.data.data
       // 处理返回数据，兼容列表或直接数组格式
       articles.value = Array.isArray(d) ? d : (d.list || [])
    }
  } catch (e) {
    console.error('Portal Home load failed')
  } finally {
    loading.value = false
  }
}

// 监听路由参数变化（如搜索关键词），重新获取数据
watch(() => route.query.title, () => fetchData())

// 组件挂载时获取数据
onMounted(() => {
  // 前端优化：如果在搜索状态下刷新页面，直接跳转回首页
  // Refresh logic moved to router/index.js
  fetchData()
})
</script>

<style scoped>
/* 隐藏滚动条样式 */
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
