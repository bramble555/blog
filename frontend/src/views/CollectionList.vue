<template>
  <div>
    <h2 class="text-xl font-semibold mb-6">My Collections</h2>
    
    <el-table :data="pagedArticles" style="width: 100%" v-loading="loading">
       <el-table-column prop="title" label="Article Title" />
       <el-table-column prop="create_time" label="Collect Time" width="180">
          <template #default="scope">
             {{ formatDate(scope.row.create_time) }}
          </template>
       </el-table-column>
       <el-table-column label="Actions" width="120">
          <template #default="scope">
             <el-button type="primary" link @click="$router.push(`/article/${scope.row.article_sn || scope.row.sn}`)">View</el-button>
             <el-button type="danger" link @click="remove(scope.row)">Uncollect</el-button>
          </template>
       </el-table-column>
    </el-table>

    <div v-if="pagination.total > pagination.size" class="mt-4 flex justify-center">
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="pagination.total"
        :page-size="pagination.size"
        :current-page="pagination.page"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { getCollects, deleteCollectArticle } from '../api/article'
import { collectionStore } from '../stores/collection'
import { formatDate } from '../utils/date'
import { ElMessage, ElMessageBox } from 'element-plus'

const articles = ref([])
const loading = ref(false)
const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const pagedArticles = computed(() => {
  const start = (pagination.page - 1) * pagination.size
  return articles.value.slice(start, start + pagination.size)
})

const handlePageChange = (page) => {
  pagination.page = page
}

const fetchData = async () => {
  loading.value = true
  try {
     const res = await getCollects()
     if (res.data.code === 10000) {
        const d = res.data.data
        let list = []
        if (Array.isArray(d)) {
           list = d
        } else if (d && d.list) {
           list = d.list
        }
        articles.value = list
        pagination.total = list.length
        const totalPage = Math.max(1, Math.ceil(pagination.total / pagination.size))
        if (pagination.page > totalPage) pagination.page = totalPage
        // Update global store
        collectionStore.setCollections(list.map(a => a.article_sn || a.sn))
     }
  } catch (e) {
     ElMessage.error('Failed to load collections')
  } finally {
    loading.value = false
  }
}

const remove = (row) => {
   ElMessageBox.confirm('Remove from collection?', 'Warning').then(async () => {
      const sn = row.article_sn || row.sn
      
      // Optimistic remove from local list and store
      const originalArticles = [...articles.value]
      articles.value = articles.value.filter(a => (a.article_sn || a.sn) !== sn)
      pagination.total = articles.value.length
      const totalPage = Math.max(1, Math.ceil(pagination.total / pagination.size))
      if (pagination.page > totalPage) pagination.page = totalPage
      collectionStore.remove(sn)

      try {
        const res = await deleteCollectArticle([sn])
        if (res.data.code === 10000) {
           ElMessage.success('Removed')
           // fetchData() // No need to re-fetch if we use optimistic update correctly, but safe to do
        } else {
           // Rollback
           articles.value = originalArticles
           collectionStore.add(sn)
           ElMessage.error(res.data.msg)
        }
      } catch (e) {
         // Rollback
         articles.value = originalArticles
         collectionStore.add(sn)
         ElMessage.error('Action failed')
      }
   })
}

onMounted(() => {
   fetchData()
})
</script>
