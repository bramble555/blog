<template>
  <div>
    <h2 class="text-xl font-semibold mb-6 flex justify-between items-center">
      Tag Management
      <div class="flex gap-2">
         <el-button v-if="isAdmin" :type="isBatchMode ? 'info' : 'warning'" size="small" @click="toggleBatchMode">
            {{ isBatchMode ? 'Exit Batch' : 'Batch Manage' }}
         </el-button>
         <el-button v-if="isAdmin && selectedTags.length > 0" type="danger" size="small" @click="batchRemove">
            Delete Selected ({{ selectedTags.length }})
         </el-button>
         <el-button type="primary" size="small" @click="openDialog">Add Tag</el-button>
      </div>
    </h2>

    <div class="flex flex-wrap gap-3" v-loading="loading">
      <div 
        v-for="tag in tags" 
        :key="tag.sn"
        class="relative group flex items-center"
      >
        <el-checkbox 
            v-if="isBatchMode" 
            :model-value="selectedTags.includes(tag.sn)" 
            @change="toggleSelect(tag.sn)" 
            class="mr-2"
        />
        <el-tag 
          class="text-base px-4 py-2 cursor-pointer transition-all select-none"
          :class="{ 'ring-2 ring-blue-500 ring-offset-1': selectedTags.includes(tag.sn) }"
          :type="selectedTags.includes(tag.sn) ? '' : 'info'"
          effect="plain"
          round
          @click="handleTagClick(tag.sn)"
        >
          # {{ tag ? tag.title : '' }}
        </el-tag>
        <!-- Individual delete for Admin (only when not in batch mode) -->
        <div v-if="isAdmin && !isBatchMode" class="absolute -top-2 -right-2 hidden group-hover:block z-10">
           <el-button type="danger" circle size="small" :icon="Delete" @click.stop="remove(tag.sn)" />
        </div>
      </div>

      <div v-if="tags.length === 0 && !loading" class="text-gray-500 text-sm italic p-4">
        No tags available.
      </div>
    </div>
    
    <!-- Batch Actions Footer -->
    <div v-if="selectedTags.length > 0 && isAdmin" class="mt-6 p-4 bg-vscode-sidebar rounded-lg flex items-center justify-between border border-vscode-border shadow-lg">
      <span class="text-gray-300">Selected: <span class="font-bold text-blue-400 text-lg mx-1">{{ selectedTags.length }}</span> tags</span>
      <div class="flex gap-3">
        <el-button plain size="default" @click="toggleBatchMode">Cancel</el-button>
        <el-button type="danger" size="default" @click="batchRemove">
          Confirm Delete
        </el-button>
      </div>
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
import { Delete } from '@element-plus/icons-vue'
import { authStore } from '../stores/auth'

const isAdmin = computed(() => authStore.role === 1)
const tags = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const form = reactive({ title: '' })
const selectedTags = ref([])
const isBatchMode = ref(false)

const toggleBatchMode = () => {
  isBatchMode.value = !isBatchMode.value
  if (!isBatchMode.value) {
    selectedTags.value = []
  }
}

const handleTagClick = (sn) => {
  if (isBatchMode.value) {
    toggleSelect(sn)
  }
}

const toggleSelect = (sn) => {
  if (!isAdmin.value) return
  // If not in batch mode, maybe enable it? Or just ignore? 
  // But now toggleSelect is called by checkbox change or tag click in batch mode.
  // Checkbox change event passes value, but here we just need SN to toggle.
  
  const index = selectedTags.value.indexOf(sn)
  if (index > -1) {
    selectedTags.value.splice(index, 1)
  } else {
    selectedTags.value.push(sn)
  }
}

const batchRemove = () => {
  ElMessageBox.confirm(`Delete ${selectedTags.value.length} selected tags?`, 'Warning', {
    confirmButtonText: 'Delete',
    cancelButtonText: 'Cancel',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteTags(selectedTags.value)
      if (res.data.code === 10000) {
        ElMessage.success(res.data.data)
        tags.value = tags.value.filter(t => !selectedTags.value.includes(t.sn))
        selectedTags.value = []
      } else {
        ElMessage.error(res.data.msg)
      }
    } catch (e) {
      ElMessage.error('Batch delete failed')
    }
  })
}


const fetchData = async () => {
  loading.value = true
  try {
    const res = await getTags()
    if (res.data.code === 10000) {
       const d = res.data.data
       if (Array.isArray(d)) {
          tags.value = d
       } else if (d && Array.isArray(d.list)) {
          tags.value = d.list
       }
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
