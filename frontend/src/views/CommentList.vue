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
       <el-table-column prop="user_detail.username" label="User" />
       <el-table-column prop="content" label="Content" />
       <el-table-column label="Time">
          <template #default="scope">
            {{ formatDateTime(scope.row.create_time) }}
          </template>
       </el-table-column>
       <el-table-column label="Actions">
          <template #default="scope">
             <el-button v-if="isAdmin || scope.row.user_sn === uSN" type="danger" link @click="remove(scope.row.sn)">Delete</el-button>
          </template>
       </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
/**
 * CommentList.vue
 * 
 * @description 评论管理页面（管理员）。允许按文章SN查询评论并进行删除。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires element-plus, vue, ../api/comment, ../stores/auth, @/utils/date
 */
import { ref, onMounted, computed } from 'vue'
import { getComments, deleteComment } from '../api/comment'
import { ElMessage, ElMessageBox } from 'element-plus'
import { authStore } from '../stores/auth'
import { formatDateTime } from '@/utils/date'

const comments = ref([])
const loading = ref(false)
const articleSn = ref('')
const pagination = ref({ page: 1, size: 20, total: 0 })

const isAdmin = computed(() => authStore.role === 1)
const uSN = computed(() => authStore.sn)

/**
 * 获取评论数据
 */
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getComments(articleSn.value || undefined, { page: pagination.value.page, size: pagination.value.size })
    if (res.data.code === 10000) {
      const data = res.data.data || {}
      comments.value = data.list || []
      pagination.value.total = data.count || 0
    } else {
      ElMessage.error(res.data.msg)
    }
  } catch (e) {
    ElMessage.error('Failed to load comments')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})

/**
 * 删除单个评论
 * @param {string} sn - Comment serial number
 */
const remove = (sn) => {
   ElMessageBox.confirm('确定删除该评论吗?', '提示', {
      type: 'warning'
   }).then(async () => {
      try {
         const res = await deleteComment(sn)
         if (res.data.code === 10000) {
            ElMessage.success('删除成功')
            fetchData()
         }
      } catch (e) {
         ElMessage.error('Delete failed')
      }
   })
}
</script>
