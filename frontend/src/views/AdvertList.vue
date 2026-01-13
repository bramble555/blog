<template>
  <div>
    <h2 class="text-xl font-semibold mb-6 flex justify-between">
      Advert Management
      <el-button type="primary" size="small" @click="openDialog">Add Advert</el-button>
    </h2>

    <el-table :data="adverts" style="width: 100%" v-loading="loading">
      <el-table-column prop="title" label="Title" />
      <el-table-column label="Image" width="100">
        <template #default="scope">
          <el-image :src="scope.row.images" class="w-16 h-10 rounded" fit="cover" />
        </template>
      </el-table-column>
      <el-table-column prop="href" label="Link" />
      <el-table-column label="Visible" width="80">
         <template #default="scope">
            <el-switch v-model="scope.row.is_show" disabled />
         </template>
      </el-table-column>
      <el-table-column label="Actions" width="100">
        <template #default="scope">
          <el-button type="danger" link @click="remove(scope.row.sn)">Delete</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="New Advert" width="30%">
      <el-form :model="form" label-width="80px">
        <el-form-item label="Title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="Link">
          <el-input v-model="form.href" />
        </el-form-item>
        <el-form-item label="Image URL">
          <el-input v-model="form.images" />
        </el-form-item>
        <el-form-item label="Show">
          <el-switch v-model="form.is_show" />
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
import { ref, reactive, onMounted } from 'vue'
import { getAdverts, createAdvert, deleteAdverts } from '../api/advert'
import { ElMessage, ElMessageBox } from 'element-plus'

const adverts = ref([])
const loading = ref(false)
const dialogVisible = ref(false)

const form = reactive({
  title: '',
  href: '',
  images: '',
  is_show: true
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getAdverts({ page: 1, size: 100 })
    if (res.data.code === 10000) {
       const d = res.data.data
       adverts.value = Array.isArray(d) ? d : (d.list || [])
    }
  } catch (error) {
    ElMessage.error('Failed to load adverts')
  } finally {
    loading.value = false
  }
}

const openDialog = () => {
  form.title = ''
  form.href = ''
  form.images = ''
  form.is_show = true
  dialogVisible.value = true
}

const submit = async () => {
  try {
    const res = await createAdvert(form)
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
  ElMessageBox.confirm('Delete this advert?', 'Warning').then(async () => {
    try {
      const res = await deleteAdverts([sn])
      if (res.data.code === 10000) {
        ElMessage.success('Deleted')
        fetchData()
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
