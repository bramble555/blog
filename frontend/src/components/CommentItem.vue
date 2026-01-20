/**
 * CommentItem.vue
 * 
 * @description 单个评论组件，支持递归显示子评论和折叠功能。
 * 特性：
 * - 显示用户信息、内容和时间
 * - 支持回复评论
 * - 支持删除评论（管理员或本人）
 * - 子评论折叠逻辑：父评论最多显示1条子评论，其余折叠
 * - "展开/收起" 切换功能
 * 
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires element-plus, vue, @element-plus/icons-vue
 */

<template>
  <div class="comment-item">
    <div class="comment-content p-4 border-b border-gray-100/10">
      <div class="user-info flex items-center mb-2">
        <el-avatar :src="formatUrl(comment.user_detail?.avatar)" size="small" class="mr-2"></el-avatar>
        <span class="username font-bold mr-2 text-vscode-text">{{ comment.user_detail?.username }}</span>
        <span class="time text-[#FF6600] text-base">{{ formatDateTime(comment.create_time) }}</span>
      </div>
      <div class="text mb-2 text-vscode-text break-all whitespace-pre-wrap">{{ comment.content }}</div>
      <div class="actions flex items-center gap-4">
        <DiggButton 
          :sn="comment.sn" 
          :count="comment.digg_count" 
          :isDigg="comment.is_digg"
          type="comment"
          @update:count="updateCount"
          @update:isDigg="updateIsDigg"
        />
        <el-button link type="primary" @click="showReply = !showReply">回复</el-button>
        <el-button v-if="canDelete" link type="danger" @click="handleDelete">删除</el-button>
      </div>
      
      <!-- Reply Input -->
      <div v-if="showReply" class="reply-box mt-4 bg-vscode-sidebar border border-vscode-border rounded-xl p-4">
          <el-input 
            v-model="replyContent" 
            type="textarea" 
            :rows="3" 
            placeholder="回复..." 
            class="mb-3"
            resize="vertical"
          />
          <div class="flex justify-end">
            <el-button type="primary" size="small" :loading="replyLoading" @click="submitReply">发送</el-button>
          </div>
      </div>
    </div>
    
    <!-- Recursive Sub Comments -->
    <div v-if="allSubComments.length" class="sub-comments ml-8 border-l-2 border-vscode-border/30 pl-4 mt-2">
      <!-- Fold/Expand Button -->
      <div class="mt-2 mb-2">
        <el-button 
          link 
          type="primary" 
          size="small" 
          @click="isExpanded = !isExpanded"
          class="flex items-center gap-1 font-bold"
        >
          <span class="text-lg mr-1">{{ isExpanded ? '-' : '+' }}</span>
          {{ isExpanded ? '收起回复' : `查看回复 (${allSubComments.length})` }}
        </el-button>
      </div>

      <!-- Hidden comments with transition -->
      <div 
        class="transition-all duration-500 ease-in-out overflow-hidden" 
        :style="{ maxHeight: isExpanded ? '2000px' : '0px', opacity: isExpanded ? 1 : 0 }"
      >
        <CommentItem 
          v-for="sub in allSubComments" 
          :key="sub.sn" 
          :comment="sub"
          @reply-success="$emit('reply-success')"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import DiggButton from './DiggButton.vue'
import { createComment, deleteComment } from '@/api/comment'
import { ElMessage, ElMessageBox } from 'element-plus'
import { authStore } from '@/stores/auth'
import { formatDateTime } from '@/utils/date'
import { formatUrl } from '@/utils/url'

const props = defineProps({
  comment: Object
})

const emit = defineEmits(['reply-success', 'delete-success'])

const showReply = ref(false)
const replyContent = ref('')
const replyLoading = ref(false)
const isExpanded = ref(false)

const allSubComments = computed(() => props.comment.sub_comments || [])

const canDelete = computed(() => {
    // Admin (role=1) or Owner (sn match)
    return authStore.role === 1 || (authStore.sn && authStore.sn == props.comment.user_sn)
})

const updateCount = (val) => {
  props.comment.digg_count = val
}
const updateIsDigg = (val) => {
  props.comment.is_digg = val
}

const submitReply = async () => {
    if(!replyContent.value.trim()) return ElMessage.warning("请输入内容")
    replyLoading.value = true
    try {
        const res = await createComment({
            article_sn: props.comment.article_sn,
            parent_comment_sn: props.comment.sn,
            content: replyContent.value
        })
        if(res.data.code === 10000) {
            ElMessage.success("回复成功")
            showReply.value = false
            replyContent.value = ''
            emit('reply-success')
        } else {
            ElMessage.error(res.data.msg)
        }
    } catch(e) {
        ElMessage.error("回复失败")
    } finally {
        replyLoading.value = false
    }
}

const handleDelete = () => {
    ElMessageBox.confirm('确定删除该评论吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        const res = await deleteComment(props.comment.sn)
        if(res.data.code === 10000) {
            ElMessage.success("删除成功")
            emit('reply-success') // Reuse refresh event
        } else {
            ElMessage.error(res.data.msg)
        }
    })
}
</script>

<style scoped>
.comment-item {
  margin-bottom: 10px;
}
</style>
