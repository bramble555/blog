<template>
  <!-- 
    主页布局容器 
    - 修改背景为暗色 (bg-vscode-bg)
    - 去除间距 (gap-0)
  -->
  <div class="pt-0 flex flex-col lg:flex-row gap-0 bg-vscode-bg text-vscode-text min-h-screen">
    <!-- 侧边栏 (广告) -->
    <!-- 保持原有逻辑，但在移动端隐藏 -->
    <aside class="hidden lg:block w-1/4 flex-shrink-0 sticky top-16 h-fit pl-0 border-r border-vscode-border">
      <SidebarAds />
    </aside>

    <!-- 文章列表区域 -->
    <!-- 
      修改说明：
      1. md:grid-cols-2: 桌面端双列显示，实现"每行两篇"
      2. gap-0: 去除文章之间的间距
      3. items-stretch: 确保卡片高度一致，避免背景空隙
    -->
    <div v-loading="loading" class="flex-1 grid grid-cols-1 md:grid-cols-2 gap-0 p-0 items-stretch content-start">
      <!-- 单个文章卡片 -->
      <div 
        v-for="article in articles" 
        :key="article.sn"
        class="group bg-vscode-bg border-b border-vscode-border md:odd:border-r overflow-hidden hover:bg-vscode-sidebar transition-all cursor-pointer flex flex-col"
        @click="$router.push(`/article/${article.sn}`)"
      >
        <!-- 文章封面图 -->
        <!-- 双列布局下改为上图下文 -->
        <div class="h-48 w-full flex-shrink-0 overflow-hidden relative group">
          <el-image 
            :src="formatUrl(article.banner_url)" 
            fit="contain" 
            class="w-full h-full transition-transform duration-700 group-hover:scale-105 bg-black/20"
            loading="lazy"
            :alt="article.title"
            :title="article.title"
          >
            <template #placeholder>
              <div class="w-full h-full bg-vscode-sidebar animate-pulse flex items-center justify-center">
                <el-icon class="text-3xl text-gray-500"><Picture /></el-icon>
              </div>
            </template>
            <template #error>
              <div class="w-full h-full bg-vscode-sidebar flex flex-col items-center justify-center text-gray-500">
                <el-icon class="text-3xl mb-1"><PictureFilled /></el-icon>
                <span class="text-2xs">No Preview</span>
              </div>
            </template>
          </el-image>
        </div>
        
        <!-- 文章内容区域 -->
        <div class="p-4 flex-1 flex flex-col justify-between">
          <div>
            <div class="flex items-center gap-2 mb-2">
               <span class="text-xs text-gray-400">{{ formatDate(article.create_time) }}</span>
            </div>
            <!-- 标题：使用浅色字体 -->
            <h3 class="text-lg font-bold mb-2 text-gray-200 group-hover:text-blue-400 transition-colors line-clamp-1">
              {{ article.title }}
            </h3>
            <!-- 摘要：使用灰色字体 -->
            <p class="text-gray-400 text-sm line-clamp-2 mb-3 leading-relaxed">
              {{ article.abstract }}
            </p>
          </div>
          
          <!-- 底部元数据 -->
          <div class="flex items-center justify-between text-xs text-gray-500 mt-2">
             <div class="flex items-center gap-3">
                <span class="flex items-center gap-1 hover:text-blue-400 cursor-pointer"><el-icon><View /></el-icon> {{ article.look_count }}</span>
                <span class="flex items-center gap-1 hover:text-blue-400 cursor-pointer"><el-icon><ChatLineSquare /></el-icon> {{ article.comment_count }}</span>
             </div>
             <span class="font-medium text-gray-500">By {{ article.username || 'Anonymous' }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 空状态 -->
    <el-empty v-if="!loading && articles.length === 0" description="暂无文章" />
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
import { useRoute } from 'vue-router'
import { getArticles } from '../api/article'
import { View, ChatLineSquare, Picture, PictureFilled } from '@element-plus/icons-vue'
import SidebarAds from '@/components/SidebarAds.vue'
import { formatDate } from '@/utils/date'
import { formatUrl } from '@/utils/url'

const route = useRoute()
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
      size: 20,
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
