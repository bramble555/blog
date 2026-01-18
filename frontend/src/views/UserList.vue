<template>
  <div>
    <h2 class="text-xl font-semibold mb-6 flex justify-between">
      User Management
      <div class="flex gap-2">
        <el-button type="primary" size="small" @click="fetchData">Refresh</el-button>
      </div>
    </h2>

    <el-table :data="users" style="width: 100%" v-loading="loading">
      <el-table-column label="SN" width="180">
        <template #default="scope">
          <span v-if="scope.row.sn == authStore.sn" class="text-orange-500 font-bold">我</span>
          <span v-else>{{ scope.row.sn }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Avatar" width="160">
        <template #default="scope">
          <div class="flex items-center gap-2">
            <el-avatar :size="40" :src="formatUrl(scope.row.avatar)" />
            <el-button
              v-if="scope.row.sn == authStore.sn"
              type="primary"
              link
              size="small"
              @click="openAvatarDialog(scope.row)"
            >
              Choose
            </el-button>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="username" label="Username" />
      <el-table-column label="Email" width="250">
        <template #default="scope">
           <span>{{ scope.row.email }}</span>
           <el-button 
             v-if="scope.row.sn == authStore.sn" 
             type="primary" 
             link 
             size="small"
             icon="Edit"
             @click="openBindEmailDialog"
             style="margin-left: 5px"
           >
             Bind
           </el-button>
        </template>
      </el-table-column>
      <el-table-column prop="role" label="Role">
        <template #default="scope">
           <el-tag :type="scope.row.role === 1 ? 'danger' : 'info'">{{ scope.row.role === 1 ? 'Admin' : 'User' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="ip" label="IP" />
      <el-table-column label="Addr" prop="addr" />
      <el-table-column label="Created At" width="180">
        <template #default="scope">
          {{ formatDate(scope.row.create_time) }}
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="200" fixed="right">
        <template #default="scope">
          <el-button 
            type="primary" 
            link 
            @click="editRole(scope.row)"
            v-if="authStore.role === 1"
          >
            Role
          </el-button>
          <el-button 
            type="warning" 
            link 
            @click="openPwdDialog(scope.row)"
            v-if="scope.row.sn == authStore.sn"
          >
            Password
          </el-button>
          <el-button 
            type="danger" 
            link 
            @click="deleteUser(scope.row.sn)"
            v-if="authStore.role === 1"
          >
            Delete
          </el-button>
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

    <!-- Password Update Dialog -->
    <el-dialog v-model="pwdDialogVisible" title="Change My Password" width="400px">
      <el-form 
        ref="pwdFormRef" 
        :model="pwdForm" 
        :rules="pwdRules" 
        label-width="140px"
        status-icon
      >
        <el-form-item label="Current Password" prop="old_pwd">
          <el-input 
            v-model="pwdForm.old_pwd" 
            type="password" 
            show-password 
            placeholder="Current password"
          />
        </el-form-item>
        <el-form-item label="New Password" prop="pwd">
          <el-input 
            v-model="pwdForm.pwd" 
            type="password" 
            show-password 
            placeholder="New password"
          />
        </el-form-item>
        <el-form-item label="Confirm Password" prop="re_pwd">
          <el-input 
            v-model="pwdForm.re_pwd" 
            type="password" 
            show-password 
            placeholder="Confirm password"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="pwdDialogVisible = false">Cancel</el-button>
          <el-button type="primary" :loading="pwdLoading" @click="submitPwd">Update</el-button>
        </span>
      </template>
    </el-dialog>
    <!-- Bind Email Dialog -->
    <el-dialog v-model="bindEmailDialogVisible" title="Bind Email" width="400px">
        <el-form :model="bindEmailForm" :rules="bindEmailRules" ref="bindEmailFormRef" label-position="top">
            <el-form-item label="Email" prop="email">
                <el-input v-model="bindEmailForm.email" placeholder="Enter your email" />
            </el-form-item>
            <el-form-item label="Verification Code" prop="code">
                <div class="flex gap-2 w-full">
                    <el-input v-model="bindEmailForm.code" placeholder="Enter code" />
                    <el-button type="primary" :disabled="codeCountdown > 0" @click="handleSendCode">
                        {{ codeCountdown > 0 ? `${codeCountdown}s` : 'Send Code' }}
                    </el-button>
                </div>
            </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="bindEmailDialogVisible = false">Cancel</el-button>
                <el-button type="primary" @click="handleBindEmail">Confirm</el-button>
            </span>
        </template>
    </el-dialog>

    <el-dialog v-model="avatarDialogVisible" title="Select Avatar Banner" width="720px">
      <div class="space-y-4">
        <div class="flex justify-between items-center">
          <span class="text-sm text-[#FFA500]">选择一张 Banner 作为头像</span>
          <el-button size="small" @click="refreshBanners" :loading="bannerLoading">Refresh</el-button>
        </div>
        <div class="grid grid-cols-3 md:grid-cols-4 gap-4 min-h-[140px]" v-loading="bannerLoading">
          <div
            v-for="banner in bannerList"
            :key="banner.sn"
            class="relative border border-vscode-border rounded-lg overflow-hidden cursor-pointer group"
            :class="banner.sn === selectedBannerSN ? 'ring-2 ring-[#22c55e] border-[#22c55e]' : ''"
            @click="handleSelectBanner(banner.sn)"
          >
            <el-image
              :src="formatUrl(banner.path || ('/uploads/file/' + banner.name))"
              fit="cover"
              class="w-full h-24 bg-black/50"
              lazy
            />
            <div
              class="absolute inset-0 bg-black/40 flex items-center justify-center text-xs text-white opacity-0 group-hover:opacity-100 transition-opacity"
            >
              {{ banner.name }}
            </div>
            <div
              v-if="banner.sn === selectedBannerSN"
              class="absolute top-1 right-1 px-2 py-0.5 rounded-full text-[10px] bg-[#22c55e] text-black font-semibold"
            >
              Selected
            </div>
          </div>
          <div
            v-if="!bannerLoading && bannerList.length === 0"
            class="col-span-full text-center text-sm text-gray-400 py-8"
          >
            No banners available. Please upload images first.
          </div>
        </div>
        <div class="flex justify-between items-center mt-2">
          <el-button
            size="small"
            :disabled="!bannerHasMore || bannerLoading"
            @click="loadMoreBanners"
          >
            {{ bannerHasMore ? 'Load More' : 'No More' }}
          </el-button>
          <div class="text-xs text-gray-400">
            Loaded {{ bannerList.length }} items
          </div>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="avatarDialogVisible = false">Cancel</el-button>
          <el-button type="primary" :disabled="!selectedBannerSN" :loading="selectingAvatar" @click="confirmSelectBanner">
            Save
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
/**
 * UserList.vue
 * 
 * @description 后台用户管理页面。展示注册用户列表，支持角色修改和删除用户。
 * @author GVB Admin
 * @last_modified 2026-01-15
 * @requires vue, element-plus, ../api/user
 */
import { ref, reactive, onMounted } from 'vue'
import { getUsers, deleteUser as deleteUserApi, updateUserRole, updateUserPassword, bindEmail, selectUserBanner } from '../api/user'
import { getBanners } from '../api/banner'
import { formatDate } from '../utils/date'
import { ElMessage, ElMessageBox } from 'element-plus'
import { authStore } from '../stores/auth'
import { formatUrl } from '../utils/url'
import { Edit } from '@element-plus/icons-vue'

// Removed incorrect useAuthStore usage
const loading = ref(false)
const roleDialogVisible = ref(false)
const pwdDialogVisible = ref(false)
const roleForm = reactive({
  user_sn: '',
  role: 2
})
const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const users = ref([])
const pwdLoading = ref(false)
const pwdFormRef = ref(null)
const pwdForm = reactive({
  old_pwd: '',
  pwd: '',
  re_pwd: ''
})

const validatePass = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('Please input the new password'))
  } else if (value === pwdForm.old_pwd) {
    callback(new Error('New password cannot be the same as current password'))
  } else {
    if (pwdForm.re_pwd !== '') {
      if (!pwdFormRef.value) return
      pwdFormRef.value.validateField('re_pwd')
    }
    callback()
  }
}

const validatePass2 = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('Please confirm the new password'))
  } else if (value !== pwdForm.pwd) {
    callback(new Error("Two inputs don't match!"))
  } else {
    callback()
  }
}

const pwdRules = reactive({
  old_pwd: [{ required: true, message: 'Current password is required', trigger: 'blur' }],
  pwd: [{ validator: validatePass, trigger: 'blur' }],
  re_pwd: [{ validator: validatePass2, trigger: 'blur' }]
})

const openPwdDialog = (row) => {
  // Extra safety check
  if (row.sn != authStore.sn) {
    ElMessage.warning('You can only change your own password.')
    return
  }
  pwdForm.old_pwd = ''
  pwdForm.pwd = ''
  pwdForm.re_pwd = ''
  pwdDialogVisible.value = true
}

const submitPwd = () => {
  if (!pwdFormRef.value) return
  pwdFormRef.value.validate(async (valid) => {
    if (valid) {
      pwdLoading.value = true
      try {
        const payload = {
          old_pwd: pwdForm.old_pwd,
          pwd: pwdForm.pwd
        }
        const res = await updateUserPassword(payload)
        if (res.data.code === 10000) {
           ElMessage.success(res.data.msg || 'Password updated successfully. Please login again.')
           pwdDialogVisible.value = false
        } else {
           ElMessage.error(res.data.msg || 'Failed to update password')
        }
      } catch (error) {
        ElMessage.error(error.response?.data?.msg || 'An error occurred')
      } finally {
        pwdLoading.value = false
      }
    }
  })
}

// Bind Email Logic
const bindEmailDialogVisible = ref(false)
const bindEmailForm = ref({
    email: '',
    code: ''
})
const bindEmailFormRef = ref(null)
const codeCountdown = ref(0)
let timer = null

const bindEmailRules = {
    email: [
        { required: true, message: 'Please input email', trigger: 'blur' },
        { type: 'email', message: 'Please input correct email address', trigger: ['blur', 'change'] }
    ],
    code: [
        { required: true, message: 'Please input verification code', trigger: 'blur' }
    ]
}

const avatarDialogVisible = ref(false)
const bannerList = ref([])
const bannerLoading = ref(false)
const bannerPage = ref(1)
const bannerPageSize = 12
const bannerHasMore = ref(true)
const selectedBannerSN = ref(null)
const selectingAvatar = ref(false)

const openAvatarDialog = (row) => {
  if (row.sn != authStore.sn) {
    ElMessage.warning('You can only change your own avatar.')
    return
  }
  avatarDialogVisible.value = true
  selectedBannerSN.value = null
  if (bannerList.value.length === 0) {
    bannerPage.value = 1
    bannerHasMore.value = true
    loadBanners(true)
  }
}

const loadBanners = async (reset = false) => {
  if (bannerLoading.value) return
  bannerLoading.value = true
  try {
    const page = reset ? 1 : bannerPage.value
    const res = await getBanners({ page, size: bannerPageSize })
    if (res.data.code === 10000) {
      const d = res.data.data
      const list = Array.isArray(d) ? d : (d.list || [])
      if (reset) {
        bannerList.value = list
      } else {
        bannerList.value = bannerList.value.concat(list)
      }
      if (list.length < bannerPageSize) {
        bannerHasMore.value = false
      } else {
        bannerHasMore.value = true
        bannerPage.value = page + 1
      }
    } else {
      ElMessage.error(res.data.msg || '获取 Banner 失败')
    }
  } catch (e) {
    ElMessage.error('获取 Banner 失败')
  } finally {
    bannerLoading.value = false
  }
}

const loadMoreBanners = () => {
  if (!bannerHasMore.value) return
  loadBanners(false)
}

const refreshBanners = () => {
  bannerPage.value = 1
  bannerHasMore.value = true
  loadBanners(true)
}

const handleSelectBanner = (sn) => {
  const exists = bannerList.value.some(b => b.sn === sn)
  if (!exists) {
    ElMessage.error('Invalid banner id')
    return
  }
  selectedBannerSN.value = sn
}

const confirmSelectBanner = async () => {
  if (!selectedBannerSN.value) {
    ElMessage.warning('Please choose a banner first.')
    return
  }
  const exists = bannerList.value.some(b => b.sn === selectedBannerSN.value)
  if (!exists) {
    ElMessage.error('Invalid banner id')
    return
  }
  selectingAvatar.value = true
  try {
    const res = await selectUserBanner(selectedBannerSN.value)
    if (res.data.code === 10000) {
      const avatar = res.data.data && res.data.data.avatar
      if (avatar) {
        authStore.avatar = avatar
        localStorage.setItem('avatar', avatar)
        const idx = users.value.findIndex(u => u.sn === authStore.sn)
        if (idx !== -1) {
          users.value[idx].avatar = avatar
        }
      }
      ElMessage.success(res.data.msg || 'Avatar updated.')
      avatarDialogVisible.value = false
    } else {
      ElMessage.error(res.data.msg || 'Failed to update avatar')
    }
  } catch (e) {
    ElMessage.error('Failed to update avatar')
  } finally {
    selectingAvatar.value = false
  }
}

const openBindEmailDialog = () => {
    bindEmailDialogVisible.value = true
    bindEmailForm.value = { email: '', code: '' }
    codeCountdown.value = 0
    if(timer) clearInterval(timer)
}

const handleSendCode = async () => {
    if(!bindEmailForm.value.email) {
        ElMessage.warning('Please input email first')
        return
    }
    
    try {
        const res = await bindEmail({ email: bindEmailForm.value.email })
        if(res.data.code === 10000) {
            ElMessage.success(res.data.msg || 'Verification code sent')
            codeCountdown.value = 60
            timer = setInterval(() => {
                codeCountdown.value--
                if(codeCountdown.value <= 0) clearInterval(timer)
            }, 1000)
        } else {
            ElMessage.error(res.data.msg)
        }
    } catch(e) {
        ElMessage.error('Failed to send code')
    }
}

const handleBindEmail = async () => {
    if(!bindEmailFormRef.value) return
    await bindEmailFormRef.value.validate(async (valid) => {
        if(valid) {
            try {
                const res = await bindEmail({ 
                    email: bindEmailForm.value.email,
                    code: bindEmailForm.value.code
                })
                if(res.data.code === 10000) {
                    ElMessage.success(res.data.msg)
                    bindEmailDialogVisible.value = false
                    fetchData() // Refresh list
                } else {
                    ElMessage.error(res.data.msg)
                }
            } catch(e) {
                ElMessage.error('Failed to bind email')
            }
        }
    })
}

/**
 * 获取用户列表
 */
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getUsers({ page: pagination.page, size: pagination.size })
    if (res.data.code === 10000) {
       const d = res.data.data
       // Robust data handling: d can be array, or object with list, or object with data
       if (Array.isArray(d)) {
         users.value = d
         pagination.total = d.length
       } else if (d && Array.isArray(d.list)) {
         users.value = d.list
         pagination.total = d.count || d.total || d.list.length
       } else if (d && d.data && Array.isArray(d.data)) {
         users.value = d.data
         pagination.total = d.count || d.total || d.data.length
       } else {
         users.value = []
         pagination.total = 0
       }
    }
  } catch (error) {
    ElMessage.error('获取用户失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  pagination.page = page
  fetchData()
}

/**
 * 删除用户
 * @param {string} sn - 用户 SN
 */
const deleteUser = (sn) => {
  ElMessageBox.confirm('确定删除该用户吗?', '提示', {
    type: 'warning'
	}).then(async () => {
		try {
			const res = await deleteUserApi(sn)
      if (res.data.code === 10000) {
        ElMessage.success('删除成功')
        users.value = users.value.filter(u => u.sn !== sn)
      } else {
        ElMessage.error(res.data.msg)
      }
    } catch (e) {
      ElMessage.error('删除失败')
    }
  })
}

/**
 * 打开角色编辑对话框
 * @param {Object} row - 用户对象
 */
const editRole = (row) => {
  roleForm.user_sn = row.sn 
  roleForm.role = row.role
  roleDialogVisible.value = true
}

/**
 * 确认角色更新
 */
const confirmRoleUpdate = async () => {
  try {
    const res = await updateUserRole(roleForm)
    if (res.data.code === 10000) {
      ElMessage.success('更新成功')
      roleDialogVisible.value = false
      fetchData()
    } else {
      ElMessage.error(res.data.msg)
    }
  } catch (e) {
    ElMessage.error('更新失败')
  }
}

onMounted(() => {
  fetchData()
})
</script>
