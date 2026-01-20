<template>
  <div>
    <h2 class="text-xl font-semibold mb-6 flex items-center justify-between">
      <span>Comment Management</span>
      <el-button
        v-if="isAdmin"
        type="danger"
        :disabled="selectedSNList.length === 0"
        @click="removeBatch"
      >
        Delete Selected ({{ selectedSNList.length }})
      </el-button>
    </h2>
    <div class="mb-4">
       <el-input v-model="articleSn" placeholder="Enter Article SN to view comments" style="width: 200px" class="mr-2" />
       <el-button type="primary" @click="fetchData">Query</el-button>
    </div>

    <el-table :data="comments" style="width: 100%" v-loading="loading" @selection-change="handleSelectionChange">
       <el-table-column v-if="isAdmin" type="selection" width="48" />
       <el-table-column prop="sn" label="SN" width="180" />
       <el-table-column prop="user_detail.username" label="User" />
       <el-table-column prop="content" label="Content" />
       <el-table-column label="Time">
          <template #default="scope">
            {{ formatDateTime(scope.row.create_time) }}
          </template>
       </el-table-column>
       <el-table-column label="Actions">
          <template #default="scope">
             <el-button type="danger" link @click="remove(scope.row.sn)">Delete</el-button>
          </template>
       </el-table-column>
    </el-table>

    <!-- Pagination -->
    <div class="mt-4 flex justify-center">
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
import { ref, reactive, onMounted, computed } from 'vue'
import { getComments, deleteComment, removeCommentBatch } from '../api/comment'
import { formatDateTime } from '../utils/date'
import { ElMessage, ElMessageBox } from 'element-plus'
import { authStore } from '../stores/auth'

const comments = ref([])
const loading = ref(false)
const articleSn = ref('')
const selectedSNList = ref([])
const isAdmin = computed(() => authStore.role === 1)
const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const fetchData = async () => {
  loading.value = true
  try {
     const params = {
        page: pagination.page,
        size: pagination.size
     }
     const sn = articleSn.value || '' 
     
     const res = await getComments(sn, params)
     if (res.data.code === 10000) {
        const d = res.data.data
        if (Array.isArray(d)) {
           comments.value = d
           pagination.total = d.length
        } else if (d && (d.list || Array.isArray(d.list))) {
           comments.value = d.list || []
           pagination.total = d.count || d.total || comments.value.length
        }
     } else {
        comments.value = []
        pagination.total = 0
     }
  } catch (e) {
     ElMessage.error('Failed to load comments')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  pagination.page = page
  fetchData()
}

const handleSelectionChange = (rows) => {
  selectedSNList.value = (rows || []).map(r => r.sn)
}

const remove = (sn) => {
   ElMessageBox.confirm('Delete this comment?', 'Warning').then(async () => {
      try {
         const res = await deleteComment(sn)
         if (res.data.code === 10000) {
            ElMessage.success('Deleted')
            selectedSNList.value = selectedSNList.value.filter(v => v !== sn)
            fetchData()
         } else {
            ElMessage.error(res.data.msg)
         }
      } catch (e) {
         ElMessage.error('Delete failed')
      }
   })
}

const removeBatch = () => {
  if (!isAdmin.value || selectedSNList.value.length === 0) return
  ElMessageBox.confirm(`Delete ${selectedSNList.value.length} selected comments?`, 'Warning').then(async () => {
    try {
      const res = await removeCommentBatch(selectedSNList.value)
      if (res.data.code === 10000) {
        ElMessage.success('Deleted')
        selectedSNList.value = []
        fetchData()
      } else {
        ElMessage.error(res.data.msg)
      }
    } catch (e) {
      ElMessage.error('Delete failed')
    }
  })
}

onMounted(() => {
   fetchData()
})
</script>

<style scoped>
:deep(.el-table) {
  font-size: 16px;
}
:deep(.el-input__inner) {
  font-size: 16px;
}
:deep(.el-pagination) {
  font-size: 16px;
}
:deep(.el-button) {
  font-size: 16px;
}
</style>
