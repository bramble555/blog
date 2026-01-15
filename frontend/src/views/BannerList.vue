<template>
  <div>
    <!-- 标题区域 -->
    <h2 class="text-xl font-semibold mb-6 flex justify-between">
      Banner Management
      <div class="flex gap-2 items-center">
        <!-- 批量删除按钮：仅管理员可见，且有选中项时可用 -->
        <el-button 
          v-if="isAdmin" 
          type="danger" 
          size="small" 
          :disabled="selectedBanners.length === 0"
          @click="handleBatchDelete"
        >
          Delete Selected ({{ selectedBanners.length }})
        </el-button>
        <!-- 上传按钮：所有用户可见 -->
        <el-upload
          action=""
          :http-request="upload"
          :show-file-list="false"
          accept="image/*"
          class="flex items-center"
        >
          <el-button type="primary" size="small">Upload Image</el-button>
        </el-upload>
      </div>
    </h2>

    <!-- 图片列表区域 -->
    <!-- 使用 Grid 布局，响应式列数 -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4" v-loading="loading">
      <div v-for="banner in banners" :key="banner.sn" class="border border-vscode-border rounded p-2 bg-[#2d2d2d] group relative">
        <!-- 多选框：仅管理员可见 -->
        <div v-if="isAdmin" class="absolute top-2 left-2 z-10">
          <el-checkbox v-model="selectedBanners" :value="banner.sn" size="large" />
        </div>
        <!-- 图片展示 -->
        <!-- 使用 formatUrl 处理图片路径，支持预览 -->
        <el-image 
          :src="formatUrl(banner.path)" 
          fit="cover" 
          class="w-full h-32 rounded bg-black"
          :preview-src-list="[formatUrl(banner.path)]"
        />
        <!-- 图片名称 -->
        <div class="mt-2 text-xs truncate text-gray-400">{{ banner.name }}</div>
        
        <!-- 删除按钮：仅管理员可见，悬停显示 -->
        <div v-if="isAdmin" class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
           <el-button type="danger" circle size="small" icon="Delete" @click="remove(banner.sn)" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * BannerList.vue
 * 
 * @description 后台 Banner 管理页面。展示图片列表，支持上传和删除。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires vue, element-plus, ../api/banner, ../stores/auth, @/utils/url
 */
import { ref, onMounted, computed } from 'vue'
import { getBanners, uploadBanners, deleteBanners } from '../api/banner'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete } from '@element-plus/icons-vue'
import { authStore } from '../stores/auth'
import { formatUrl } from '@/utils/url'

/**
 * 计算属性：是否为管理员
 * @returns {boolean} true if role is 1 (admin)
 */
const isAdmin = computed(() => authStore.role === 1)

/**
 * 响应式变量：Banner 列表数据
 * @type {Ref<Array>}
 */
const banners = ref([])

/**
 * 响应式变量：选中的 Banner SN 列表
 * @type {Ref<Array>}
 */
const selectedBanners = ref([])

/**
 * 响应式变量：加载状态
 * @type {Ref<boolean>}
 */
const loading = ref(false)

/**
 * 获取 Banner 列表
 * @description 从后端获取图片列表数据
 * @async
 */
const fetchData = async () => {
  loading.value = true
  try {
    // 调用 API 获取列表，默认第一页，100条
    const res = await getBanners({ page: 1, size: 100 })
    if (res.data.code === 10000) {
       const d = res.data.data
       // 兼容不同的数据返回格式（数组或分页对象）
       banners.value = Array.isArray(d) ? d : (d.list || [])
    }
  } catch (error) {
    ElMessage.error('获取图片失败')
  } finally {
    loading.value = false
  }
}

/**
 * 上传图片处理
 * @description 处理图片上传请求
 * @param {Object} param - Upload parameter object from el-upload
 * @async
 */
const upload = async (param) => {
  // 构建 FormData 对象
  const formData = new FormData()
  formData.append('images', param.file)
  
  try {
    // 调用上传 API
    const res = await uploadBanners(formData)
    if (res.data.code === 10000) {
      // 检查返回的文件列表状态
      const fileList = res.data.data || []
      const items = Array.isArray(fileList) ? fileList : []
      const firstError = items.find(file => file.msg)

      if (firstError && firstError.msg) {
        ElMessage.warning(firstError.msg || '部分图片上传失败')
      } else {
        ElMessage.success('上传成功')
      }
      
      // 无论成功与否都刷新列表（可能部分成功）
      fetchData()
    } else {
      ElMessage.error(res.data.msg)
    }
  } catch (e) {
    ElMessage.error('上传失败')
  }
}

/**
 * 批量删除 Banner
 * @description 删除选中的图片
 * @async
 */
const handleBatchDelete = () => {
  ElMessageBox.confirm(`确定删除选中的 ${selectedBanners.value.length} 张图片吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteBanners(selectedBanners.value)
      if (res.data.code === 10000) {
        ElMessage.success('批量删除成功')
        selectedBanners.value = [] // 清空选中状态
        fetchData()
      } else {
        ElMessage.error(res.data.msg)
      }
    } catch (e) {
      ElMessage.error('删除失败')
    }
  })
}

/**
 * 删除 Banner
 * @description 删除指定的图片
 * @param {string} sn - Banner serial number
 * @async
 */
const remove = (sn) => {
  ElMessageBox.confirm('确定删除该图片吗?', '提示').then(async () => {
    try {
      // 调用删除 API
      const res = await deleteBanners([sn])
      if (res.data.code === 10000) {
        ElMessage.success('删除成功')
        // 删除成功后刷新列表
        fetchData()
      } else {
        ElMessage.error(res.data.msg)
      }
    } catch (e) {
      ElMessage.error('删除失败')
    }
  })
}

/**
 * 生命周期钩子：组件挂载完成
 * @description 初始化加载数据
 */
onMounted(() => {
  fetchData()
})
</script>
