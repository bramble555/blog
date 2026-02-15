<script setup>
import { onMounted } from 'vue'

onMounted(() => {
  // Ensure the document has the 'dark' class just in case some legacy styles or Element Plus components rely on it
  document.documentElement.classList.add('dark')
})
</script>

<template>
  <el-config-provider namespace="el">
    <router-view></router-view>
  </el-config-provider>
</template>

<style>
/* Global CSS is in style.css */
</style>
