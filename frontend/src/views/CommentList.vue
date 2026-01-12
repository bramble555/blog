<template>
  <div>
    <h2 class="text-xl font-semibold mb-6">Comment Management</h2>
    <!-- Note: Backend API requires article_id to list comments. This view would ideally list all comments or require picking an article. 
         For admin dashboard purposes, usually we want "All Comments" or "Recent Comments".
         If GetArticleCommentsHandler requires article_id, we might be limited.
         Assuming for now we might need to select an Article first or this page is just a placeholder until backend supports "List all comments".
         Let's implement a simple "Enter Article ID" query for now to satisfy component existence.
    -->
    <div class="mb-4">
       <el-input v-model="articleId" placeholder="Enter Article ID to view comments" style="width: 200px" class="mr-2" />
       <el-button type="primary" @click="fetchData">Query</el-button>
    </div>

    <el-table :data="comments" style="width: 100%" v-loading="loading">
       <el-table-column prop="id" label="ID" width="80" />
       <el-table-column prop="username" label="User" />
       <el-table-column prop="content" label="Content" />
       <el-table-column prop="create_time" label="Time" />
       <el-table-column label="Actions">
          <template #default="scope">
             <el-button type="danger" link @click="remove(scope.row.id)">Delete</el-button>
          </template>
       </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { getComments, deleteComment } from '../api/comment'
import { ElMessage, ElMessageBox } from 'element-plus'

const articleId = ref('')
const comments = ref([])
const loading = ref(false)

const fetchData = async () => {
  if (!articleId.value) return
  loading.value = true
  try {
    const res = await getComments({ article_id: articleId.value })
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

const remove = (id) => {
   ElMessageBox.confirm('Delete comment?', 'Warning').then(async () => {
      try {
         const res = await deleteComment(id)
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
