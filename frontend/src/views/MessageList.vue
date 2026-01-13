<template>
  <div>
    <h2 class="text-xl font-semibold mb-6">Messages</h2>
    <el-tabs v-model="activeTab" @tab-click="handleTabClick">
      <el-tab-pane label="My Messages" name="my">
         <div v-for="grp in myMessages" :key="grp.sn" class="p-3 border-b border-vscode-border hover:bg-[#2d2d2d] cursor-pointer" @click="viewDetail(grp)">
            <div class="font-bold">{{ grp.send_user_name }}</div>
            <div class="text-sm text-gray-400">{{ grp.content }}</div>
         </div>
      </el-tab-pane>
      <el-tab-pane label="All Messages (Admin)" name="all">
          <el-table :data="allMessages" style="width: 100%">
             <el-table-column prop="send_user_name" label="Sender" />
             <el-table-column prop="rev_user_name" label="Receiver" />
             <el-table-column prop="content" label="Content" />
             <el-table-column prop="create_time" label="Time" />
          </el-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getMessagesAll, getMyMessages } from '../api/message'
import { ElMessage } from 'element-plus'

const activeTab = ref('my')
const myMessages = ref([
  // Mock Data as fallback
  { sn: 1, send_user_name: 'System', content: 'Welcome to the system!', create_time: '2023-01-01' }
])
const allMessages = ref([])

const fetchMy = async () => {
   try {
      const res = await getMyMessages()
      if (res.data.code === 10000) {
         // TODO: Connect to Go Backend real logic
         // myMessages.value = res.data.data
      }
   } catch(e) {
      console.log('Backend API not ready for MyMessages')
   }
}

const fetchAll = async () => {
   try {
      const res = await getMessagesAll({ page: 1, size: 20 })
      if (res.data.code === 10000) {
         allMessages.value = res.data.data.list
      }
   } catch(e) {
      console.log('Backend API not ready for AllMessages')
   }
}

const handleTabClick = (tab) => {
   if (tab.paneName === 'all') fetchAll()
   else fetchMy()
}

const viewDetail = (grp) => {
   ElMessage.info('Detail view not implemented yet')
}

onMounted(() => {
   fetchMy()
})
</script>
