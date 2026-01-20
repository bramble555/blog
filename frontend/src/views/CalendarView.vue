<template>
  <div class="calendar-view">
    <div class="mb-6">
      <h2 class="text-xl font-semibold text-white mb-2">Article Calendar</h2>
      <p class="text-base text-gray-400">Visualize your article publishing frequency and history.</p>
    </div>
    
    <ArticleCalendar />
  </div>
</template>

<script setup>
import ArticleCalendar from '@/components/ArticleCalendar.vue'
</script>
