<template>
  <div class="max-w-4xl mx-auto">
    <div class="mb-6 flex justify-between items-center">
      <h2 class="text-xl font-semibold text-white">{{ isEdit ? 'Edit Article' : 'New Article' }}</h2>
      <button @click="goBack" class="px-3 py-1 text-[#FFA500] hover:text-white transition-colors">
        Cancel
      </button>
    </div>

    <div class="bg-vscode-sidebar border border-vscode-border rounded-lg p-6 space-y-6">
      
      <!-- Title -->
      <div>
        <label class="block text-sm font-medium text-[#FFA500] mb-2">Title</label>
        <input v-model="form.title" type="text" placeholder="Enter article title"
          class="w-full bg-vscode-bg border border-vscode-border rounded p-2 text-vscode-text focus:border-vscode-primary focus:outline-none transition-colors"
        >
      </div>

      <!-- Abstract -->
      <div>
        <label class="block text-sm font-medium text-[#FFA500] mb-2">Abstract</label>
        <textarea v-model="form.abstract" rows="2" placeholder="Short summary..."
          class="w-full bg-vscode-bg border border-vscode-border rounded p-2 text-vscode-text focus:border-vscode-primary focus:outline-none transition-colors"
        ></textarea>
      </div>

      <!-- Tags -->
      <div>
        <label class="block text-sm font-medium text-[#FFA500] mb-2">Tags (comma separated)</label>
        <input v-model="tagsInput" type="text" placeholder="e.g. tutorial, backend"
          class="w-full bg-vscode-bg border border-vscode-border rounded p-2 text-vscode-text focus:border-vscode-primary focus:outline-none transition-colors"
        >
      </div>
      
      <!-- Content (Simple Textarea for Markdown) -->
      <div>
         <label class="block text-sm font-medium text-[#FFA500] mb-2">Content (Markdown)</label>
         <div class="border border-vscode-border rounded bg-vscode-bg h-96 flex flex-col">
            <!-- Toolbar placeholder -->
            <div class="border-b border-vscode-border p-2 bg-[#2d2d2d] text-xs text-[#FFA500] flex gap-2">
               <span>Markdown Supported</span>
            </div>
            <textarea v-model="form.content" 
              class="flex-1 w-full bg-transparent p-4 text-vscode-text focus:outline-none font-mono text-sm resize-none"
              placeholder="# Write your content here..."
            ></textarea>
         </div>
      </div>

      <!-- Actions -->
      <div class="flex justify-end pt-4 border-t border-vscode-border">
         <button @click="submit" :disabled="submitting" 
           class="px-6 py-2 bg-vscode-primary text-white rounded font-medium hover:bg-opacity-90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
         >
           {{ submitting ? 'Saving...' : (isEdit ? 'Update Article' : 'Publish Article') }}
         </button>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticle, createArticle, updateArticle } from '../api/article'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.sn)
const submitting = ref(false)
const tagsInput = ref('')

const form = reactive({
  title: '',
  abstract: '',
  content: '',
  tags: []
})

// Store original data for change detection
const originalForm = reactive({
  title: '',
  abstract: '',
  content: '',
  tags: []
})

// Initialize
onMounted(async () => {
  if (isEdit.value) {
    // Fetch article detail
    try {
      const res = await getArticle(route.params.sn)
      if (res.data.code === 10000) {
        const d = res.data.data
        form.title = d.title
        form.abstract = d.abstract
        form.content = d.content
        let serverTags = d.tags
        let tagArr = []
        if (Array.isArray(serverTags)) {
          tagArr = serverTags
        } else if (typeof serverTags === 'string') {
          try {
            const parsed = JSON.parse(serverTags)
            if (Array.isArray(parsed)) {
              tagArr = parsed
            } else if (typeof parsed === 'string') {
              tagArr = parsed.split(/[,，]/).map(t => t.trim()).filter(t => t)
            }
          } catch {
            tagArr = serverTags.split(/[,，]/).map(t => t.trim()).filter(t => t)
          }
        }
        form.tags = tagArr
        tagsInput.value = tagArr.join(', ')
        
        // Save original state
        originalForm.title = form.title
        originalForm.abstract = form.abstract
        originalForm.content = form.content
        originalForm.tags = [...tagArr]
      } else {
        alert('Failed to load article')
      }
    } catch (e) {
      console.error(e)
      alert('Network Error')
    }
  }
})

const goBack = () => router.back()

const submit = async () => {
  if (!form.title || !form.content) {
    alert('Title and Content are required')
    return
  }

  submitting.value = true

  try {
    let res
    if (isEdit.value) {
      // Update: Only send changed fields
      const changes = {}
      if (form.title !== originalForm.title) changes.title = form.title
      if (form.abstract !== originalForm.abstract) changes.abstract = form.abstract
      if (form.content !== originalForm.content) changes.content = form.content
      
      // Check tags change
      const currentTags = tagsInput.value.split(/[,，]/).map(t => t.trim()).filter(t => t)
      const tagsChanged = JSON.stringify(currentTags) !== JSON.stringify(originalForm.tags)
      
      if (tagsChanged) {
        changes.tags = tagsInput.value
      }

      if (Object.keys(changes).length === 0) {
        ElMessage.info('No changes detected')
        submitting.value = false
        return
      }

      // Using generic object map for PUT
      res = await updateArticle(route.params.sn, changes)
    } else {
      // Create: Send full payload
      const payload = {
        title: form.title,
        abstract: form.abstract,
        content: form.content,
        tags: tagsInput.value,
        banner_sn: "1" // Default placeholder
      }
      res = await createArticle(payload)
    }

    if (res.data.code === 10000) {
      if (isEdit.value) {
         ElMessage.success('Updated successfully')
         router.push('/admin/articles')
      } else {
         try {
            // Check if data is string (new JSON response) or object (old behavior fallback)
            const responseData = typeof res.data.data === 'string' 
               ? JSON.parse(res.data.data) 
               : res.data.data
            
            // Check custom status field if present
            if (responseData && responseData.status && responseData.status !== 'ok') {
               ElMessage.warning('Upload succeeded but returned status: ' + responseData.status)
            } else {
               ElMessage.success('Published successfully')
            }
            router.push('/admin/articles')
         } catch(e) {
            console.warn('Response parsing error:', e)
            // Fallback success
            ElMessage.success('Published successfully')
            router.push('/admin/articles')
         }
      }
    } else {
      ElMessage.error('Operation failed: ' + res.data.msg)
    }
  } catch (e) {
    console.error(e)
    ElMessage.error('Network Error')
  } finally {
    submitting.value = false
  }
}
</script>
