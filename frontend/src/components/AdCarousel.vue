<template>
  <el-carousel :interval="4000" type="card" height="200px" v-if="adverts.length">
    <el-carousel-item v-for="item in adverts" :key="item.sn">
      <component
        :is="item.href ? 'a' : 'div'"
        :href="item.href || undefined"
        target="_blank"
        class="block w-full h-full relative"
      >
        <img :src="item.images" class="w-full h-full object-cover" :alt="item.title">
        <div class="absolute bottom-0 left-0 w-full bg-black bg-opacity-50 text-white p-2 text-center truncate">
            {{ item.title }}
        </div>
      </component>
    </el-carousel-item>
  </el-carousel>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAdverts } from '@/api/advert'

const adverts = ref([])

onMounted(async () => {
    try {
        const res = await getAdverts({ page: 1, limit: 10, is_show: true })
        if (res.data.code === 10000) {
            const d = res.data.data
            adverts.value = Array.isArray(d) ? d : (d.list || [])
        }
    } catch(e) {
        console.error(e)
    }
})
</script>
