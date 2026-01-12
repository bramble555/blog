<template>
  <el-config-provider namespace="el">
    <router-view></router-view>
  </el-config-provider>
</template>

<script setup>
import { onMounted } from 'vue'
// Ensure dark mode is on
onMounted(() => {
  document.documentElement.classList.add('dark')
})
</script>

<style>
/* Global CSS is in style.css */
</style>
