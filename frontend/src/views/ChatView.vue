<template>
  <div class="flex h-[calc(100vh-100px)] border border-vscode-border rounded overflow-hidden">
    <!-- Chat List -->
    <div class="w-64 bg-vscode-sidebar border-r border-vscode-border p-4">
       <h3 class="font-bold mb-4">Chat Rooms</h3>
       <div class="text-sm text-gray-400 p-2 hover:bg-vscode-bg rounded cursor-pointer bg-vscode-bg text-white">
          General Group
       </div>
    </div>
    
    <!-- Chat Window -->
    <div class="flex-1 flex flex-col bg-[#1e1e1e]">
       <div class="flex-1 p-4 overflow-y-auto space-y-4">
          <div v-for="(msg, i) in chatHistory" :key="i" class="flex gap-3">
             <el-avatar :size="30" :src="msg.avatar">{{ msg.nick_name?.[0] }}</el-avatar>
             <div>
                <div class="text-xs text-vscode-primary font-bold">{{ msg.nick_name }}</div>
                <div class="text-sm bg-[#3e3e42] p-2 rounded mt-1">{{ msg.content }}</div>
             </div>
          </div>
       </div>
       <div class="p-4 border-t border-vscode-border flex gap-2">
          <el-input v-model="inputText" placeholder="Type a message..." @keyup.enter="send" />
          <el-button type="primary" @click="send">Send</el-button>
       </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { WS_URL } from '../api/chat'

const inputText = ref('')
const chatHistory = ref([
   { nick_name: 'Bot', content: 'Welcome to the chat room!', avatar: '' }
])
let ws = null

// TODO: Connect to Go Backend WebSocket
const connect = () => {
   try {
      ws = new WebSocket(WS_URL)
      ws.onopen = () => {
         chatHistory.value.push({ nick_name: 'Sys', content: 'Connected to server.' })
         // Send join room msg
         ws.send(JSON.stringify({ msg_type: 2, nick_name: 'Admin', content: 'joined' }))
      }
      ws.onmessage = (evt) => {
         try {
            const data = JSON.parse(evt.data)
            // Backend returns ResponseChatGroup wrapped or direct?
            // Assuming standard format
            if (data.active_count) return; // Just heartbeat or count
            chatHistory.value.push(data) 
         } catch(e) {}
      }
      ws.onerror = () => {
         chatHistory.value.push({ nick_name: 'Sys', content: 'Connection error (Mock Mode Active)' })
      }
   } catch(e) {
      console.log('WS Connection failed')
   }
}

const send = () => {
   if(!inputText.value) return
   if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ msg_type: 1, content: inputText.value, nick_name: 'Admin' }))
   } else {
      // Mock echo
      chatHistory.value.push({ nick_name: 'Me', content: inputText.value })
   }
   inputText.value = ''
}

onMounted(() => {
   connect()
})

onUnmounted(() => {
   if(ws) ws.close()
})
</script>
