<template>
  <el-button 
    :type="isDigg ? 'primary' : 'default'" 
    :class="{ 'is-digged': isDigg }"
    @click="handleDigg"
    :loading="loading"
    size="small"
    round
  >
    <el-icon class="mr-1"><Pointer /></el-icon>
    <span>{{ count }}</span>
  </el-button>
</template>

<script setup>
import { ref } from 'vue'
import { postArticleDigg, postCommentDigg } from '@/api/digg'
import { ElMessage } from 'element-plus'
import { Pointer } from '@element-plus/icons-vue'

const props = defineProps({
  count: { type: Number, default: 0 },
  isDigg: { type: Boolean, default: false },
  sn: { type: String, required: true },
  type: { type: String, default: 'article' } // 'article' or 'comment'
})

const emit = defineEmits(['update:count', 'update:isDigg'])

const loading = ref(false)

const handleDigg = async () => {
  loading.value = true
  try {
    const api = props.type === 'article' ? postArticleDigg : postCommentDigg
    const res = await api(props.sn)
    if (res.data.code === 10000) {
      const newIsDigg = !props.isDigg
      const newCount = newIsDigg ? props.count + 1 : props.count - 1
      
      emit('update:isDigg', newIsDigg)
      emit('update:count', newCount)
      ElMessage.success(res.data.data)
    } else {
      ElMessage.error(res.data.msg)
    }
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.is-digged {
  transform: scale(1.05);
  transition: transform 0.1s;
}
.mr-1 {
  margin-right: 4px;
}
</style>
