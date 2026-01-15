<template>
  <div>
    <h2 class="text-xl font-semibold mb-6 flex justify-between">
      Tag Management
      <el-button type="primary" size="small" @click="openDialog">Add Tag</el-button>
    </h2>

    <div class="flex flex-wrap gap-2" v-loading="loading">
      <el-tag 
        v-for="tag in tags" 
        :key="tag.sn" 
        class="text-lg px-4 py-1"
        :closable="isAdmin" @close="remove(tag.sn)"
      >
        {{ tag.title }}
      </el-tag>
    </div>

    <el-dialog v-model="dialogVisible" title="New Tag" width="30%">
      <el-form :model="form">
        <el-form-item label="Tag Title">
          <el-input v-model="form.title" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">Cancel</el-button>
          <el-button type="primary" @click="submit">Confirm</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { getTags, createTag, deleteTags } from '../api/tag'
import { ElMessage, ElMessageBox } from 'element-plus'
import { authStore } from '../stores/auth'

const isAdmin = computed(() => authStore.role === 1)
const tags = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const form = reactive({ title: '' })

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getTags({ page: 1, size: 100 })
    if (res.data.code === 10000) {
       const d = res.data.data
       tags.value = Array.isArray(d) ? d : (d.list || [])
    }
  } catch (error) {
    ElMessage.error('Failed to load tags')
  } finally {
    loading.value = false
  }
}

const openDialog = () => {
  form.title = ''
  dialogVisible.value = true
}

const submit = async () => {
  try {
    const res = await createTag(form)
    if (res.data.code === 10000) {
      ElMessage.success('Created')
      dialogVisible.value = false
      fetchData()
    } else {
      ElMessage.error(res.data.msg)
    }
  } catch (e) {
    ElMessage.error('Create failed')
  }
}

const remove = (sn) => {
  ElMessageBox.confirm('Delete this tag?', 'Warning').then(async () => {
    try {
      const res = await deleteTags([sn])
      if (res.data.code === 10000) {
        ElMessage.success('Deleted')
        tags.value = tags.value.filter(t => t.sn !== sn)
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
