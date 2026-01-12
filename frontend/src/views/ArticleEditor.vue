<template>
  <div class="max-w-4xl mx-auto">
    <div class="mb-6 flex justify-between items-center">
      <h2 class="text-xl font-semibold text-white">{{ isEdit ? 'Edit Article' : 'New Article' }}</h2>
      <button @click="goBack" class="px-3 py-1 text-gray-400 hover:text-white transition-colors">
        Cancel
      </button>
    </div>

    <div class="bg-vscode-sidebar border border-vscode-border rounded-lg p-6 space-y-6">
      
      <!-- Title -->
      <div>
        <label class="block text-sm font-medium text-gray-400 mb-2">Title</label>
        <input v-model="form.title" type="text" placeholder="Enter article title"
          class="w-full bg-vscode-bg border border-vscode-border rounded p-2 text-vscode-text focus:border-vscode-primary focus:outline-none transition-colors"
        >
      </div>

      <!-- Abstract -->
      <div>
        <label class="block text-sm font-medium text-gray-400 mb-2">Abstract</label>
        <textarea v-model="form.abstract" rows="2" placeholder="Short summary..."
          class="w-full bg-vscode-bg border border-vscode-border rounded p-2 text-vscode-text focus:border-vscode-primary focus:outline-none transition-colors"
        ></textarea>
      </div>

      <!-- Category & Tags -->
      <div class="grid grid-cols-2 gap-6">
        <div>
          <label class="block text-sm font-medium text-gray-400 mb-2">Category</label>
          <input v-model="form.category" type="text" placeholder="e.g. Go, Vue, Tech"
            class="w-full bg-vscode-bg border border-vscode-border rounded p-2 text-vscode-text focus:border-vscode-primary focus:outline-none transition-colors"
          >
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-400 mb-2">Tags (comma separated)</label>
          <input v-model="tagsInput" type="text" placeholder="e.g. tutorial, backend"
            class="w-full bg-vscode-bg border border-vscode-border rounded p-2 text-vscode-text focus:border-vscode-primary focus:outline-none transition-colors"
          >
        </div>
      </div>
      
      <!-- Content (Simple Textarea for Markdown) -->
      <div>
         <label class="block text-sm font-medium text-gray-400 mb-2">Content (Markdown)</label>
         <div class="border border-vscode-border rounded bg-vscode-bg h-96 flex flex-col">
            <!-- Toolbar placeholder -->
            <div class="border-b border-vscode-border p-2 bg-[#2d2d2d] text-xs text-gray-500 flex gap-2">
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

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const submitting = ref(false)
const tagsInput = ref('')

const form = reactive({
  title: '',
  abstract: '',
  category: '',
  content: '',
  tags: [] // Array of strings
})

// Initialize
onMounted(async () => {
  if (isEdit.value) {
    // Fetch article detail
    try {
      const res = await getArticle(route.params.id)
      if (res.data.code === 10000) {
        const d = res.data.data
        form.title = d.title
        form.abstract = d.abstract
        form.category = d.category
        form.content = d.content
        form.tags = d.tags || []
        tagsInput.value = (form.tags || []).join(', ')
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
  
  // Process tags
  const tags = tagsInput.value.split(',').map(t => t.trim()).filter(t => t)

  const payload = {
    title: form.title,
    abstract: form.abstract,
    category: form.category,
    content: form.content,
    tags: tags,
    banner_id: "1" // Default placeholder as string because backend uses json:"banner_id,string"
  }

  try {
    let res
    if (isEdit.value) {
      // Update
      // Using generic object map for PUT
      res = await updateArticle(route.params.id, payload)
    } else {
      // Create
      // ParamArticle struct matches
      res = await createArticle(payload)
    }

    if (res.data.code === 10000) {
      router.push('/admin/articles')
    } else {
      alert('Operation failed: ' + res.data.msg)
    }
  } catch (e) {
    console.error(e)
    alert('Network Error')
  } finally {
    submitting.value = false
  }
}
</script>
