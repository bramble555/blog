<template>
  <div class="markdown-renderer bg-bg-secondary/50 rounded-xl overflow-hidden border border-border-primary/50 transition-colors">
    <!-- Content -->
    <div class="p-8 min-h-[200px]">
      <div v-if="loading" class="flex items-center justify-center h-40 text-text-tertiary">
        <svg class="animate-spin h-6 w-6 mr-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        Loading...
      </div>

      <div v-else class="prose max-w-none">
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
/* Global variables are available without explicit import */

/* Deep selectors for rendered markdown content with Semantic Props */
:deep(h1) {
  font-size: 2em;
  font-weight: 800;
  margin-bottom: 0.5em;
  color: var(--text-primary);
  line-height: 1.2;
}
:deep(h2) {
  font-size: 1.5em;
  font-weight: 700;
  margin-top: 1.5em;
  margin-bottom: 0.75em;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-primary);
  padding-bottom: 0.5em;
}
:deep(p) {
  margin-bottom: 1.5em;
  line-height: 1.75;
  color: var(--text-secondary);
}
:deep(pre) {
  background-color: var(--bg-primary); /* Darker contrast block */
  padding: 1.25em;
  border-radius: 0.5em;
  overflow-x: auto;
  border: 1px solid var(--border-primary);
  margin-bottom: 1.5em;
  box-shadow: inset 0 2px 4px rgba(0,0,0,0.05);
}
:deep(code) {
  font-family: 'Fira Code', 'Consolas', monospace;
  background-color: var(--bg-tertiary);
  padding: 0.2em 0.4em;
  border-radius: 4px;
  color: var(--accent-secondary);
  font-size: 0.875em;
}
:deep(pre code) {
  background-color: transparent;
  padding: 0;
  color: var(--text-primary); /* Or appropriate syntax highlighting color if handled by HLJS upstream */
}
:deep(blockquote) {
  border-left: 4px solid var(--accent-primary);
  padding-left: 1em;
  padding-top: 0.5em;
  padding-bottom: 0.5em;
  background-color: var(--bg-tertiary);
  color: var(--text-tertiary);
  margin-bottom: 1.5em;
  font-style: italic;
  border-radius: 0 4px 4px 0;
}
:deep(ul), :deep(ol) {
  padding-left: 1.5em;
  margin-bottom: 1.5em;
  color: var(--text-secondary);
}
:deep(li) {
  margin-bottom: 0.5em;
  position: relative;
}
:deep(a) {
  color: var(--accent-primary);
  text-decoration: none;
  font-weight: 500;
  border-bottom: 1px solid transparent;
  transition: all 0.2s;
}
:deep(a:hover) {
  border-bottom-color: var(--accent-primary);
}
:deep(img) {
  max-width: 100%;
  border-radius: 8px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  margin: 1.5em 0;
}
</style>
