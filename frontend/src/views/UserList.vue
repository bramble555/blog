<template>
  <div>
    <h2 class="text-xl font-semibold mb-6 flex justify-between">
      User Management
      <el-button type="primary" size="small" @click="fetchData">Refresh</el-button>
    </h2>

    <el-table :data="users" style="width: 100%" v-loading="loading">
      <el-table-column prop="sn" label="SN" width="180" />
      <el-table-column label="Avatar" width="80">
        <template #default="scope">
          <el-avatar :size="40" :src="scope.row.avatar" />
        </template>
      </el-table-column>
      <el-table-column prop="username" label="Username" />
      <el-table-column prop="email" label="Email" />
      <el-table-column prop="role" label="Role">
        <template #default="scope">
           <el-tag :type="scope.row.role === 1 ? 'danger' : 'info'">{{ scope.row.role === 1 ? 'Admin' : 'User' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="ip" label="IP" />
      <el-table-column label="Addr" prop="addr" />
      <el-table-column label="Created At" prop="create_time" width="180" />
      <el-table-column label="Actions" width="150" fixed="right">
        <template #default="scope">
          <el-button type="primary" link @click="editRole(scope.row)">Role</el-button>
          <el-button type="danger" link @click="deleteUser(scope.row.sn)">Delete</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- Role Edit Dialog -->
    <el-dialog v-model="roleDialogVisible" title="Update User Role" width="30%">
      <el-form :model="roleForm">
        <el-form-item label="Role">
          <el-select v-model="roleForm.role" placeholder="Select role">
            <el-option label="Admin" :value="1" />
            <el-option label="Normal User" :value="2" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="roleDialogVisible = false">Cancel</el-button>
          <el-button type="primary" @click="confirmRoleUpdate">Confirm</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getUsers, deleteUsers, updateUserRole } from '../api/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const users = ref([])
const loading = ref(false)

const roleDialogVisible = ref(false)
const roleForm = reactive({ user_sn: 0, role: 1 })

const fetchData = async () => {
  loading.value = true
  try {
    // Assuming backend returns { data: { list: [], ... }, code: 10000 } or just list
    const res = await getUsers({ page: 1, size: 100 })
    if (res.data.code === 10000) {
       const d = res.data.data
       users.value = Array.isArray(d) ? d : (d.list || [])
    }
  } catch (error) {
    ElMessage.error('Failed to fetch users')
  } finally {
    loading.value = false
  }
}

const deleteUser = (sn) => {
  ElMessageBox.confirm('Are you sure to delete this user?', 'Warning', {
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteUsers([sn])
      if (res.data.code === 10000) {
        ElMessage.success('Deleted')
        users.value = users.value.filter(u => u.sn !== sn)
      } else {
        ElMessage.error(res.data.msg)
      }
    } catch (e) {
      ElMessage.error('Delete failed')
    }
  })
}

const editRole = (row) => {
  roleForm.user_sn = row.sn // Note: row.sn is usually a number in JS but user_sn,string in JSON. API expects int64.
  roleForm.role = row.role
  roleDialogVisible.value = true
}

const confirmRoleUpdate = async () => {
  try {
    const res = await updateUserRole(roleForm)
    if (res.data.code === 10000) {
      ElMessage.success('Updated')
      roleDialogVisible.value = false
      fetchData()
    } else {
      ElMessage.error(res.data.msg)
    }
  } catch (e) {
    ElMessage.error('Update failed')
  }
}

onMounted(() => {
  fetchData()
})
</script>
