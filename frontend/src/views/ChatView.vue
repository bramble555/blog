<template>
  <div class="flex h-[calc(100vh-100px)] border border-vscode-border rounded overflow-hidden relative select-none">
    <!-- Chat List -->
    <div 
      class="bg-vscode-sidebar border-r border-vscode-border p-4 hidden md:block flex-shrink-0 relative group"
      :style="{ width: sidebarWidth + 'px' }"
    >
       <h3 class="font-bold mb-4">Chat Rooms</h3>
       <div class="text-base text-gray-400 p-2 hover:bg-vscode-bg rounded cursor-pointer bg-vscode-bg text-white">
          General Group
       </div>
       <!-- Sidebar Resizer Handle -->
       <div 
          class="absolute right-0 top-0 w-1 h-full cursor-col-resize hover:bg-blue-500 transition-colors z-10 opacity-0 group-hover:opacity-100"
          :class="{ 'opacity-100 bg-blue-500': resizeState.target === 'sidebar' }"
          @mousedown.prevent="startResizeSidebar"
       ></div>
    </div>
    
    <!-- Chat Window -->
    <div class="flex-1 flex flex-col bg-[#1e1e1e] min-w-0">
       <div class="flex items-center justify-between px-4 py-3 border-b border-vscode-border bg-[#1f1f1f]">
          <div class="text-base text-gray-300 font-medium">General Group</div>
          <div class="text-base text-gray-400">在线人数：{{ onlineCount }}</div>
       </div>
       <div class="flex-1 p-4 overflow-y-auto space-y-4" ref="chatContainer">
          <div v-for="(msg, i) in chatHistory" :key="i">
             <!-- System Message -->
             <div v-if="isSystemMsg(msg)" class="text-center text-base text-gray-500 my-2">
                {{ msg.content }}
             </div>
             
             <!-- Chat Message -->
             <div v-else class="flex gap-3" :class="{'flex-row-reverse': isMyMsg(msg)}">
                <el-avatar :size="40" :src="avatarSrc(msg)" class="flex-shrink-0">{{ msg.nick_name?.[0]?.toUpperCase() }}</el-avatar>
                 <div class="max-w-[70%]" :class="{'items-end': isMyMsg(msg), 'items-start': !isMyMsg(msg)}">
                    <div class="flex gap-2 mb-1" :class="{'flex-row-reverse': isMyMsg(msg)}">
                       <span class="text-base text-gray-400">{{ msg.nick_name }}</span>
                       <span class="text-base text-gray-500">{{ formatTime(msg.date) }}</span>
                    </div>
                   <div v-if="msg.msg_type === 4" class="rounded overflow-hidden max-w-[260px]">
                      <el-image :src="imageSrc(msg)" fit="cover" />
                   </div>
                   <div v-else class="p-2 rounded text-base break-words" 
                        :class="isMyMsg(msg) ? 'bg-[#0e639c] text-white' : 'bg-[#3e3e42] text-gray-200'">
                      {{ msg.content }}
                   </div>
                </div>
             </div>
          </div>
       </div>
       
       <!-- Input Area Resizer Handle -->
       <div 
          class="h-1 w-full cursor-row-resize hover:bg-blue-500 transition-colors flex-shrink-0 bg-[#2d2d2d] relative z-10"
          :class="{ 'bg-blue-500': resizeState.target === 'input' }"
          @mousedown.prevent="startResizeInput"
       ></div>

       <div class="border-t border-vscode-border flex flex-col" :style="{ height: inputHeight + 'px' }">
          <div class="flex h-full p-4 bg-[#1e1e1e]">
              <el-input 
                v-model="inputText" 
                type="textarea" 
                class="h-full custom-textarea flex-1 mr-6"
                placeholder="Type a message... (Enter to send, Shift+Enter for new line)" 
                resize="none"
                @keydown.enter.exact.prevent="send"
              />
              <div class="flex gap-2 h-full shrink-0 items-end">
                <el-button type="primary" class="!w-10 !h-10 !px-0" @click="send" title="Send">
                    <el-icon><Promotion /></el-icon>
                </el-button>
                <el-upload action="" :http-request="uploadImage" :show-file-list="false" accept="image/*">
                   <el-button class="!w-10 !h-10 !px-0 !ml-0" title="Image">
                      <el-icon><Picture /></el-icon>
                   </el-button>
                </el-upload>
                <el-button class="!w-10 !h-10 !px-0 !ml-0" @click="openHistory" title="History">
                    <el-icon><Clock /></el-icon>
                </el-button>
              </div>
          </div>
       </div>
    </div>

    <!-- History Drawer -->
    <el-drawer
        v-model="historyDrawerVisible"
        title="Chat History"
        direction="rtl"
        size="400px"
    >
        <div class="flex flex-col h-full">
            <div class="flex-1 overflow-y-auto space-y-4 p-2" v-loading="historyLoading">
                 <div v-for="(msg, i) in historyList" :key="i" class="border-b border-gray-700 pb-2">
                     <div class="flex justify-between text-base text-gray-400 mb-1">
                         <span>{{ msg.nick_name }}</span>
                         <span>{{ formatTime(msg.create_time) }}</span>
                     </div>
                     <div v-if="msg.msg_type === 4" class="text-base text-gray-200">
                         <el-image :src="imageSrc(msg)" fit="cover" class="max-w-[240px]" />
                     </div>
                     <div v-else class="text-base text-gray-200 break-words">{{ msg.content }}</div>
                 </div>
                 <div v-if="historyList.length === 0 && !historyLoading" class="text-center text-gray-500 mt-4">
                     No history records
                 </div>
            </div>
            <div class="pt-4 border-t border-gray-700 flex justify-center">
                 <el-pagination
                    background
                    layout="prev, pager, next"
                    :total="historyTotal"
                    :page-size="historyPageSize"
                    v-model:current-page="historyPage"
                    @current-change="fetchHistory"
                    small
                 />
            </div>
        </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch, reactive } from 'vue'
import { WS_URL, getChatHistory, uploadChatImage } from '../api/chat'
import { authStore } from '../stores/auth'
import { formatUrl } from '@/utils/url'
import { Clock, Promotion, Picture } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// State
const inputText = ref('')
const chatHistory = ref([])
const onlineCount = ref(0)
const chatContainer = ref(null)
let ws = null

// History State
const historyDrawerVisible = ref(false)
const historyList = ref([])
const historyLoading = ref(false)
const historyPage = ref(1)
const historyPageSize = 20
const historyTotal = ref(0)

// Resize State
const sidebarWidth = ref(256)
const inputHeight = ref(100)
const resizeState = reactive({
    isResizing: false,
    target: null, // 'sidebar' | 'input'
    startPos: 0,
    startSize: 0
})

const isMyMsg = (msg) => {
   return String(msg?.user_sn ?? '') !== '' && String(authStore.sn) !== '' && parseInt(msg.user_sn) === parseInt(authStore.sn)
}

const isSystemMsg = (msg) => {
   return [1, 2, 5].includes(Number(msg?.msg_type))
}

const avatarSrc = (msg) => {
   const raw = String(msg?.avatar ?? '').trim()
   if (raw) return formatUrl(raw)
   if (isMyMsg(msg) && authStore.avatar) return formatUrl(authStore.avatar)
   return ''
}

const imageSrc = (msg) => {
   return formatUrl(msg?.content)
}

const formatTime = (dateStr) => {
   if (!dateStr) return ''
   const date = new Date(dateStr)
   return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const scrollToBottom = async () => {
   await nextTick()
   if (chatContainer.value) {
      setTimeout(() => {
          chatContainer.value.scrollTop = chatContainer.value.scrollHeight
      }, 50)
   }
}

// History Logic
const openHistory = () => {
    historyDrawerVisible.value = true
    historyPage.value = 1
    fetchHistory()
}

const fetchHistory = async () => {
    historyLoading.value = true
    try {
        const res = await getChatHistory({ page: historyPage.value, size: historyPageSize })
        if (res.data.code === 10000) {
            historyList.value = res.data.data.list
            historyTotal.value = res.data.data.count
        } else {
            ElMessage.error(res.data.msg || 'Failed to fetch history')
        }
    } catch (e) {
        ElMessage.error('Error fetching history')
    } finally {
        historyLoading.value = false
    }
}

// Resizing Logic
const startResizeSidebar = (e) => {
    resizeState.isResizing = true
    resizeState.target = 'sidebar'
    resizeState.startPos = e.clientX
    resizeState.startSize = sidebarWidth.value
    document.addEventListener('mousemove', handleResize)
    document.addEventListener('mouseup', stopResize)
    document.body.style.cursor = 'col-resize'
}

const startResizeInput = (e) => {
    resizeState.isResizing = true
    resizeState.target = 'input'
    resizeState.startPos = e.clientY
    resizeState.startSize = inputHeight.value
    document.addEventListener('mousemove', handleResize)
    document.addEventListener('mouseup', stopResize)
    document.body.style.cursor = 'row-resize'
}

const handleResize = (e) => {
    if (!resizeState.isResizing) return
    
    if (resizeState.target === 'sidebar') {
        const delta = e.clientX - resizeState.startPos
        const newWidth = resizeState.startSize + delta
        sidebarWidth.value = Math.max(150, Math.min(500, newWidth))
    } else if (resizeState.target === 'input') {
        const delta = resizeState.startPos - e.clientY // Moving up increases height
        const newHeight = resizeState.startSize + delta
        inputHeight.value = Math.max(60, Math.min(500, newHeight))
    }
}

const stopResize = () => {
    resizeState.isResizing = false
    resizeState.target = null
    document.removeEventListener('mousemove', handleResize)
    document.removeEventListener('mouseup', stopResize)
    document.body.style.cursor = ''
}

const connect = () => {
   try {
      if (ws) {
         ws.close()
         ws = null
      }
      if (!authStore.token) {
         chatHistory.value.push({ msg_type: 5, content: '请先登录后再进入群聊' })
         return
      }
      const url = `${WS_URL}?token=${authStore.token}`
      ws = new WebSocket(url)
      
      ws.onopen = () => {
          // Connection established
      }
      
      ws.onmessage = (evt) => {
         try {
            const data = JSON.parse(evt.data)
            if (typeof data?.online_count === 'number') {
               onlineCount.value = data.online_count
            }
            chatHistory.value.push(data)
            scrollToBottom()
         } catch(e) {}
      }
      
      ws.onerror = () => {
         chatHistory.value.push({ 
             msg_type: 5, 
             content: 'Connection failed. Please refresh.' 
         })
      }
   } catch(e) {
      console.error('WS Connection failed', e)
   }
}

const send = () => {
   if(!inputText.value.trim()) return
   
   if (ws && ws.readyState === WebSocket.OPEN) {
      const msg = { 
          msg_type: 3,
          content: inputText.value, 
      }
      ws.send(JSON.stringify(msg))
      inputText.value = ''
   } else {
      chatHistory.value.push({ msg_type: 5, content: 'WebSocket 未连接，请刷新页面重试' })
   }
}

const uploadImage = async (param) => {
   const formData = new FormData()
   formData.append('image', param.file)
   try {
      const res = await uploadChatImage(formData)
      if (res.data.code === 10000) {
         const url = res.data.data
         if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ msg_type: 4, content: url }))
         } else {
            chatHistory.value.push({ msg_type: 5, content: 'WebSocket 未连接，请刷新页面重试' })
         }
      } else {
         ElMessage.error(res.data.msg)
      }
   } catch (e) {
      ElMessage.error('上传失败')
   }
}

watch(() => authStore.token, (newVal) => {
   chatHistory.value = [] 
   connect()
})

onMounted(() => {
   connect()
})

onUnmounted(() => {
   if(ws) ws.close()
   stopResize()
})
</script>

<style scoped>
:deep(.el-textarea__inner) {
    height: 100% !important;
    resize: none;
    font-size: 16px !important;
}
:deep(.el-input__wrapper) {
    height: 100%;
}
</style>
