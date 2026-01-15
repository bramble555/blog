<template>
  <div class="sidebar-ads space-y-6">
    <div v-for="ad in ads" :key="ad.sn" class="ad-item w-full aspect-square overflow-hidden rounded-xl border border-vscode-border shadow-md hover:shadow-vscode-primary/20 transition-all">
      <component
        :is="ad.href ? 'a' : 'div'"
        :href="ad.href || undefined"
        target="_blank"
        class="block w-full h-full relative group"
      >
        <el-image 
          :src="formatUrl(ad.images)" 
          fit="cover" 
          class="w-full h-full group-hover:scale-105 transition-transform duration-500" 
          loading="lazy"
        >
          <template #error>
            <div class="w-full h-full bg-vscode-sidebar flex items-center justify-center text-[#FF6600] text-xs">
               <el-icon class="mr-1"><PictureFilled /></el-icon> No Image
            </div>
          </template>
        </el-image>
        <div class="absolute bottom-0 left-0 w-full bg-gradient-to-t from-black/80 to-transparent p-3">
          <p class="text-white text-sm font-medium truncate">{{ ad.title }}</p>
          <span class="text-2xs text-gray-300 uppercase tracking-wider">Advertisement</span>
        </div>
      </component>
    </div>
  </div>
</template>

<script setup>
/**
 * SidebarAds.vue
 * 
 * @description 侧边栏广告组件。展示图片广告链接。
 * @author GVB Admin
 * @last_modified 2026-01-14
 */
import { ref, onMounted } from 'vue'
import { getAdverts } from '@/api/advert'
import { PictureFilled } from '@element-plus/icons-vue'
import { formatUrl } from '@/utils/url'

// 广告列表数据
const ads = ref([])

/**
 * 组件挂载时获取广告数据
 * @async
 */
onMounted(async () => {
  try {
    // 请求广告列表，默认获取前5条显示的广告
    const res = await getAdverts({ page: 1, limit: 5, is_show: true })
    if (res.data.code === 10000) {
      // 兼容处理返回数据，确保为数组
      const list = Array.isArray(res.data.data) ? res.data.data : []
      // 截取前5个展示
      ads.value = list.slice(0, 5)
    }
  } catch (e) {
    console.error('Failed to load sidebar ads', e)
  }
})

</script>

<style scoped>
.ad-item {
  background: var(--vscode-sidebar-background);
}
</style>
