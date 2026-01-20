<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-xl font-semibold">Messages</h2>
      <el-button type="primary" @click="openSendDialog">Send Message</el-button>
    </div>

    <el-tabs v-model="activeTab" @tab-click="handleTabClick">
      <el-tab-pane label="My Messages" name="my">
         <div v-if="myMessages.length === 0" class="text-gray-500 p-4 text-center">No messages</div>
         <div v-for="msg in myMessages" :key="msg.sn" 
              class="p-4 border-b border-vscode-border hover:bg-[#2d2d2d] cursor-pointer transition-colors duration-200" 
              :class="{'bg-[#252526]': !msg.is_read}"
              @click="viewDetail(msg)">
            <div class="flex justify-between items-start mb-2">
                <div class="font-bold flex items-center gap-2">
                    <span>{{ msg.send_user_name }}</span>
                    <el-tag v-if="!msg.is_read" size="small" type="danger" effect="dark" round>New</el-tag>
                    <el-tag v-else size="small" type="info" effect="plain" round>Read</el-tag>
                </div>
                <div class="text-xs text-gray-500">{{ formatDate(msg.create_time) }}</div>
            </div>
            <div class="text-base text-gray-300 line-clamp-2">{{ msg.content }}</div>
         </div>
         
         <!-- Pagination for My Messages -->
         <div class="mt-4 flex justify-center" v-if="myMessages.length > 0">
            <el-pagination
              background
              layout="prev, pager, next"
              :total="myPagination.total"
              :page-size="myPagination.size"
              :current-page="myPagination.page"
              @current-change="handleMyPageChange"
            />
         </div>
      </el-tab-pane>

      <el-tab-pane label="Sent Messages" name="sent">
          <el-table :data="sentMessages" style="width: 100%">
             <el-table-column prop="rev_user_name" label="Receiver" width="150" />
             <el-table-column prop="content" label="Content" show-overflow-tooltip />
             <el-table-column label="Status" width="100">
                <template #default="scope">
                    <el-tag :type="scope.row.is_read ? 'success' : 'info'">{{ scope.row.is_read ? 'Read' : 'Unread' }}</el-tag>
                </template>
             </el-table-column>
             <el-table-column label="Time" width="180">
               <template #default="scope">
                 {{ formatDate(scope.row.create_time) }}
               </template>
             </el-table-column>
             <el-table-column label="Actions" width="120">
               <template #default="scope">
                 <el-button link type="primary" @click="viewDetail(scope.row, false)">View</el-button>
               </template>
             </el-table-column>
          </el-table>
          
          <div class="mt-4 flex justify-center">
            <el-pagination
              background
              layout="total, prev, pager, next"
              :total="sentPagination.total"
              :page-size="sentPagination.size"
              :current-page="sentPagination.page"
              @current-change="handleSentPageChange"
            />
          </div>
      </el-tab-pane>
      
      <el-tab-pane label="All Messages (Admin)" name="all" v-if="isAdmin">
          <el-table :data="allMessages" style="width: 100%">
             <el-table-column prop="send_user_name" label="Sender" width="150" />
             <el-table-column prop="rev_user_name" label="Receiver" width="150" />
             <el-table-column prop="content" label="Content" show-overflow-tooltip />
             <el-table-column label="Status" width="100">
                <template #default="scope">
                    <el-tag :type="scope.row.is_read ? 'success' : 'info'">{{ scope.row.is_read ? 'Read' : 'Unread' }}</el-tag>
                </template>
             </el-table-column>
             <el-table-column label="Time" width="180">
               <template #default="scope">
                 {{ formatDate(scope.row.create_time) }}
               </template>
             </el-table-column>
          </el-table>
          
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
      </el-tab-pane>
    </el-tabs>

    <!-- Send Message Dialog -->
    <el-dialog v-model="sendDialogVisible" title="Send Message" width="500px" destroy-on-close>
        <el-form :model="sendForm" label-position="top">
            <el-form-item v-if="isAdmin">
                <el-checkbox v-model="sendForm.isBroadcast" label="Broadcast to All Users (Admin Only)" border class="w-full" />
            </el-form-item>
            
            <el-form-item label="Receiver" v-if="!sendForm.isBroadcast" required>
                <el-select 
                    v-model="sendForm.rev_user_sn" 
                    placeholder="Search and select user" 
                    filterable 
                    remote
                    :remote-method="searchUsers"
                    :loading="userLoading"
                    style="width: 100%">
                    <el-option v-for="u in userList" :key="u.sn" :label="u.username" :value="u.sn">
                        <span class="float-left">{{ u.username }}</span>
                        <span class="float-right text-gray-400 text-xs">{{ u.email }}</span>
                    </el-option>
                </el-select>
            </el-form-item>

            <el-form-item label="Content" required>
                <el-input 
                    v-model="sendForm.content" 
                    type="textarea" 
                    rows="5" 
                    placeholder="Type your message here..." 
                    maxlength="500"
                    show-word-limit
                />
            </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="sendDialogVisible = false">Cancel</el-button>
                <el-button type="primary" @click="handleSend" :loading="sending">
                    {{ sendForm.isBroadcast ? 'Broadcast' : 'Send' }}
                </el-button>
            </span>
        </template>
    </el-dialog>

    <!-- Message Detail Dialog -->
    <el-dialog v-model="detailVisible" title="Message Detail" width="500px">
        <div v-if="currentMessage">
            <div class="flex items-center gap-4 mb-6 pb-4 border-b border-vscode-border">
                <el-avatar :size="50" :src="formatUrl(detailAvatar)">{{ detailUserName?.charAt(0) }}</el-avatar>
                <div>
                    <div class="text-lg font-bold">{{ detailUserName }}</div>
                    <div class="text-sm text-gray-400">{{ formatDate(currentMessage.create_time) }}</div>
                </div>
            </div>
            <div class="text-base leading-relaxed whitespace-pre-wrap bg-[#1e1e1e] p-4 rounded border border-vscode-border">
                {{ currentMessage.content }}
            </div>
        </div>
        <template #footer>
            <el-button @click="detailVisible = false">Close</el-button>
        </template>
    </el-dialog>
  </div>
</template>

<script setup>
/**
 * MessageList.vue
 * 
 * @description 用户消息展示页面及管理员消息概览。
 * @author GVB Admin
 * @last_modified 2026-01-19
 * @requires vue, element-plus, ../api/message
 */
import { ref, reactive, onMounted, computed } from 'vue'
import { getMessagesAll, getMyMessages, getSentMessages, sendMessage, broadcastMessage, readMessage } from '../api/message'
import { getUsers } from '../api/user'
import { formatDate } from '../utils/date'
import { formatUrl } from '../utils/url'
import { ElMessage, ElMessageBox } from 'element-plus'
import { authStore } from '../stores/auth'

const isAdmin = computed(() => authStore.role === 1)
const activeTab = ref('my')

// My Messages State
const myMessages = ref([])
const myPagination = reactive({
    page: 1,
    size: 10,
    total: 0
})

const sentMessages = ref([])
const sentPagination = reactive({
    page: 1,
    size: 10,
    total: 0
})

// All Messages State (Admin)
const allMessages = ref([])
const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

// Send Message State
const sendDialogVisible = ref(false)
const sending = ref(false)
const sendForm = reactive({
    isBroadcast: false,
    rev_user_sn: '',
    content: ''
})
const userList = ref([])
const userLoading = ref(false)

// Detail State
const detailVisible = ref(false)
const currentMessage = ref(null)
const detailMode = ref('received')
const detailUserName = computed(() => {
    if (!currentMessage.value) return ''
    return detailMode.value === 'sent' ? currentMessage.value.rev_user_name : currentMessage.value.send_user_name
})
const detailAvatar = computed(() => {
    if (!currentMessage.value) return ''
    return detailMode.value === 'sent' ? currentMessage.value.rev_user_avater : currentMessage.value.send_user_avater
})

/**
 * 获取当前用户的消息
 */
const fetchMy = async () => {
   try {
      const res = await getMyMessages({ page: myPagination.page, size: myPagination.size })
      if (res.data.code === 10000) {
         const d = res.data.data
         if (Array.isArray(d)) {
             myMessages.value = d
             myPagination.total = d.length // Simplistic total if no count returned
         } else {
             myMessages.value = d.list || []
             myPagination.total = d.count || d.total || 0
         }
      }
   } catch(e) {
      console.log('Error fetching my messages', e)
   }
}

/**
 * 获取所有消息（管理员）
 */
const fetchAll = async () => {
   if (!isAdmin.value) return
   try {
      const res = await getMessagesAll({ page: pagination.page, size: pagination.size })
      if (res.data.code === 10000) {
         allMessages.value = res.data.data.list
         pagination.total = res.data.data.count
      }
   } catch(e) {
      console.log('Error fetching all messages')
   }
}

const fetchSent = async () => {
   try {
      const res = await getSentMessages({ page: sentPagination.page, size: sentPagination.size })
      if (res.data.code === 10000) {
         const d = res.data.data
         if (Array.isArray(d)) {
             sentMessages.value = d
             sentPagination.total = d.length
         } else {
             sentMessages.value = d.list || []
             sentPagination.total = d.count || d.total || 0
         }
      }
   } catch(e) {
      console.log('Error fetching sent messages', e)
   }
}

const handlePageChange = (page) => {
  pagination.page = page
  fetchAll()
}

const handleMyPageChange = (page) => {
    myPagination.page = page
    fetchMy()
}

const handleSentPageChange = (page) => {
    sentPagination.page = page
    fetchSent()
}

/**
 * 处理标签页切换
 * @param {Object} tab - 标签页实例
 */
const handleTabClick = (tab) => {
   if (tab.paneName === 'all') fetchAll()
   else if (tab.paneName === 'sent') fetchSent()
   else fetchMy()
}

/**
 * 打开发送对话框
 */
const openSendDialog = () => {
    sendForm.content = ''
    sendForm.rev_user_sn = ''
    sendForm.isBroadcast = false
    sendDialogVisible.value = true
    // Initial fetch of users
    searchUsers('')
}

/**
 * 搜索用户
 */
const searchUsers = async (query) => {
    userLoading.value = true
    try {
        // Fetch users, simplistic approach: fetch page 1 size 50
        // Real implementation should pass query to backend if supported
        const res = await getUsers({ page: 1, size: 50, nickname: query }) 
        if (res.data.code === 10000) {
            userList.value = res.data.data.list || []
        }
    } catch (e) {
        console.error(e)
    } finally {
        userLoading.value = false
    }
}

/**
 * 发送消息逻辑
 */
const handleSend = async () => {
    if (!sendForm.content.trim()) {
        ElMessage.warning('Please enter message content')
        return
    }

    if (sendForm.isBroadcast) {
        // Admin Broadcast
        ElMessageBox.confirm(
            'Are you sure you want to broadcast this message to ALL users?',
            'Confirm Broadcast',
            { confirmButtonText: 'Yes, Broadcast', cancelButtonText: 'Cancel', type: 'warning' }
        ).then(async () => {
            performSend()
        }).catch(() => {})
    } else {
        // Normal Send
        if (!sendForm.rev_user_sn) {
            ElMessage.warning('Please select a receiver')
            return
        }
        performSend()
    }
}

const performSend = async () => {
    sending.value = true
    try {
        let res
        if (sendForm.isBroadcast) {
            res = await broadcastMessage({ content: sendForm.content })
        } else {
            res = await sendMessage({ 
                rev_user_sn: String(sendForm.rev_user_sn), 
                content: sendForm.content 
            })
        }

        if (res.data.code === 10000) {
            ElMessage.success('Message sent successfully')
            sendDialogVisible.value = false
            fetchMy()
            fetchSent()
            if (activeTab.value === 'all') fetchAll()
        } else {
            ElMessage.error(res.data.msg || 'Failed to send')
        }
    } catch (e) {
        ElMessage.error('Network error')
    } finally {
        sending.value = false
    }
}

/**
 * 查看消息详情并更新已读状态
 * @param {Object} msg - 消息对象
 */
const viewDetail = async (msg, markRead = true) => {
   currentMessage.value = msg
   detailVisible.value = true
   detailMode.value = markRead ? 'received' : 'sent'

   // Mark as read if not already
   if (markRead && !msg.is_read) {
       try {
           const res = await readMessage(msg.sn)
           if (res.data.code === 10000) {
               // Update local state immediately for responsiveness
               msg.is_read = true
               // Decrease unread count if we had a global counter
           }
       } catch (e) {
           console.error('Failed to mark as read', e)
           // Keep local state as unread if failed, or retry logic could be added here
       }
   }
}

onMounted(() => {
   fetchMy()
})
</script>

<style scoped>
:deep(.el-tabs__item) {
  font-size: 16px;
}
:deep(.el-table) {
  font-size: 16px;
}
:deep(.el-pagination) {
  font-size: 16px;
}
</style>
