<template>
  <div class="markdown-renderer bg-vscode-bg rounded-lg overflow-hidden border border-vscode-border">
    <!-- Content -->
    <div class="p-6 min-h-[200px]">
      <div v-if="loading" class="flex items-center justify-center h-40 text-gray-500">
        <svg class="animate-spin h-6 w-6 mr-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        Loading...
      </div>

      <div v-else class="prose prose-invert max-w-none">
        <div v-html="sanitizedContent"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import DOMPurify from 'dompurify'

const props = defineProps({
  parsed: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const sanitizedContent = computed(() => {
  if (props.parsed) {
    return DOMPurify.sanitize(props.parsed)
  }
  return ''
})
</script>

<style scoped>
/* Deep selectors for rendered markdown content */
:deep(h1) {
  font-size: 2em;
  font-weight: bold;
  margin-bottom: 0.5em;
  color: #569cd6;
}
:deep(h2) {
  font-size: 1.5em;
  font-weight: bold;
  margin-top: 1em;
  margin-bottom: 0.5em;
  color: #569cd6;
}
:deep(p) {
  margin-bottom: 1em;
  line-height: 1.6;
}
:deep(pre) {
  background-color: #1e1e1e;
  padding: 1em;
  border-radius: 4px;
  overflow-x: auto;
  border: 1px solid #3e3e42;
  margin-bottom: 1em;
}
:deep(code) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  background-color: rgba(255, 255, 255, 0.1);
  padding: 0.2em 0.4em;
  border-radius: 3px;
  color: #ce9178;
}
:deep(pre code) {
  background-color: transparent;
  padding: 0;
  color: #d4d4d4;
}
:deep(blockquote) {
  border-left: 4px solid #FFA500;
  padding-left: 1em;
  color: #808080;
  margin-bottom: 1em;
}
:deep(ul), :deep(ol) {
  padding-left: 2em;
  margin-bottom: 1em;
}
:deep(li) {
  margin-bottom: 0.5em;
}
:deep(a) {
  color: #3794ff;
  text-decoration: none;
}
:deep(a:hover) {
  text-decoration: underline;
}
:deep(img) {
  max-width: 100%;
  border-radius: 4px;
}
</style>
