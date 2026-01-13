<template>
  <div>
    <h2 class="text-xl font-semibold mb-6">Comment Management</h2>
    <!-- Note: Backend API requires article_sn to list comments. This view would ideally list all comments or require picking an article. 
         For admin dashboard purposes, usually we want "All Comments" or "Recent Comments".
         If GetArticleCommentsHandler requires article_sn, we might be limited.
         Assuming for now we might need to select an Article first or this page is just a placeholder until backend supports "List all comments".
         Let's implement a simple "Enter Article SN" query for now to satisfy component existence.
    -->
    <div class="mb-4">
       <el-input v-model="articleSn" placeholder="Enter Article SN to view comments" style="width: 200px" class="mr-2" />
       <el-button type="primary" @click="fetchData">Query</el-button>
    </div>

    <el-table :data="comments" style="width: 100%" v-loading="loading">
       <el-table-column prop="sn" label="SN" width="180" />
       <el-table-column prop="username" label="User" />
       <el-table-column prop="content" label="Content" />
       <el-table-column prop="create_time" label="Time" />
       <el-table-column label="Actions">
          <template #default="scope">
             <el-button v-if="isAdmin || scope.row.user_sn === uSN" type="danger" link @click="remove(scope.row.sn)">Delete</el-button>
          </template>
       </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { getComments, deleteComment } from '../api/comment'
import { ElMessage, ElMessageBox } from 'element-plus'
import { authStore } from '../stores/auth'

const isAdmin = computed(() => authStore.role === 1)
const uSN = computed(() => authStore.sn)

const articleSn = ref('')
const comments = ref([])
const loading = ref(false)

const fetchData = async () => {
  if (!articleSn.value) return
  loading.value = true
  try {
    const res = await getComments({ article_sn: articleSn.value })
    if (res.data.code === 10000) {
      comments.value = res.data.data // Depending on structure, might be list or tree
    } else {
      ElMessage.error(res.data.msg)
    }
  } catch (e) {
    ElMessage.error('Failed to load comments')
  } finally {
    loading.value = false
  }
}

const remove = (sn) => {
   ElMessageBox.confirm('Delete comment?', 'Warning').then(async () => {
      try {
         const res = await deleteComment(sn)
         if (res.data.code === 10000) {
            ElMessage.success('Deleted')
            fetchData()
         }
      } catch (e) {
         ElMessage.error('Delete failed')
      }
   })
}
</script>
