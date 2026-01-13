<template>
  <div>
    <h2 class="text-xl font-semibold mb-6 flex justify-between">
      Banner Management
      <el-upload
        v-if="isAdmin"
        action=""
        :http-request="upload"
        :show-file-list="false"
        accept="image/*"
      >
        <el-button type="primary" size="small">Upload Image</el-button>
      </el-upload>
    </h2>

    <div class="grid grid-cols-2 md:grid-cols-4 gap-4" v-loading="loading">
      <div v-for="banner in banners" :key="banner.sn" class="border border-vscode-border rounded p-2 bg-[#2d2d2d] group relative">
        <el-image 
          :src="banner.path" 
          fit="cover" 
          class="w-full h-32 rounded bg-black"
          :preview-src-list="[banner.path]"
        />
        <div class="mt-2 text-xs truncate text-gray-400">{{ banner.name }}</div>
        
        <div v-if="isAdmin" class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
           <el-button type="danger" circle size="small" icon="Delete" @click="remove(banner.sn)" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { getBanners, uploadBanners, deleteBanners } from '../api/banner'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete } from '@element-plus/icons-vue'
import { authStore } from '../stores/auth'

const isAdmin = computed(() => authStore.role === 1)
const banners = ref([])
const loading = ref(false)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getBanners({ page: 1, size: 100 })
    if (res.data.code === 10000) {
       const d = res.data.data
       banners.value = Array.isArray(d) ? d : (d.list || [])
    }
  } catch (error) {
    ElMessage.error('Failed to load banners')
  } finally {
    loading.value = false
  }
}

const upload = async (param) => {
  const formData = new FormData()
  formData.append('images', param.file)
  try {
    const res = await uploadBanners(formData)
    if (res.data.code === 10000) {
      ElMessage.success('Uploaded')
      fetchData()
    } else {
      ElMessage.error(res.data.msg)
    }
  } catch (e) {
    ElMessage.error('Upload failed')
  }
}

const remove = (sn) => {
  ElMessageBox.confirm('Delete this image?', 'Warning').then(async () => {
    try {
      const res = await deleteBanners([sn])
      if (res.data.code === 10000) {
        ElMessage.success('Deleted')
        banners.value = banners.value.filter(b => b.sn !== sn)
      } else {
        ElMessage.error(res.data.msg)
      }
    } catch(e) {
      ElMessage.error('Delete failed')
    }
  })
}

onMounted(() => {
  fetchData()
})
</script>
