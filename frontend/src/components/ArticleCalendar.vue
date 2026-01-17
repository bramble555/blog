<template>
  <div class="calendar-container w-full bg-vscode-sidebar border border-vscode-border rounded-lg p-4">
    <div class="flex justify-between items-center mb-4">
      <h3 class="text-lg font-semibold text-white">Article Calendar</h3>
      <div v-if="loading" class="text-sm text-gray-400">Loading...</div>
      <div v-else-if="error" class="text-sm text-red-400">{{ error }}</div>
    </div>
    
    <div ref="chartRef" class="w-full h-48 md:h-64"></div>

    <div v-if="selectedDate" class="mt-4 border-t border-vscode-border pt-4">
      <h4 class="text-md font-medium text-[#FFA500] mb-2">
        {{ selectedDate }} ({{ selectedCount }} articles)
      </h4>
      <div v-if="articlesLoading" class="text-gray-400 text-sm">Loading...</div>
      <div v-else-if="selectedCount === 0" class="text-gray-400 text-sm">
        No articles published on this day.
      </div>
      <div v-else class="space-y-2">
        <div 
          v-for="article in selectedArticles" 
          :key="article.sn"
          @click="goToArticle(article.sn)"
          class="p-2 bg-[#2d2d2d] rounded cursor-pointer hover:bg-[#3d3d3d] transition-colors"
        >
          <div class="text-sm font-medium text-white">{{ article.title }}</div>
          <div class="text-xs text-gray-400">{{ article.abstract }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'
import { getArticlesCalendar, getArticles } from '../api/article'
import { useRouter } from 'vue-router'

const router = useRouter()
const chartRef = ref(null)
let chartInstance = null
const loading = ref(false)
const error = ref(null)
const calendarData = ref({})
const selectedDate = ref(null)
const selectedCount = ref(0)
const selectedArticles = ref([])
const articlesLoading = ref(false)

// Cache key
const CACHE_KEY = 'article_calendar_cache'
const CACHE_EXPIRY = 5 * 60 * 1000 // 5 minutes

const initChart = () => {
  if (!chartRef.value) return
  
  chartInstance = echarts.init(chartRef.value)
  
  const dateList = Object.keys(calendarData.value).map(date => [date, calendarData.value[date]])
  
  // Get range from data
  const years = [...new Set(Object.keys(calendarData.value).map(d => d.split('-')[0]))]
  const range = years.length > 0 ? years : [new Date().getFullYear()]

  // Calculate visual map max
  const maxCount = Math.max(...Object.values(calendarData.value), 5)

  const option = {
    // 提示框组件配置
    tooltip: {
      position: 'top', // 提示框显示位置：顶部
      formatter: function (p) {
        // 格式化提示内容：日期 + 文章数量
        const format = echarts.format.formatTime('yyyy-MM-dd', p.data[0]);
        return format + ': ' + p.data[1] + ' articles';
      }
    },
    // 视觉映射组件配置 (右下角的图例)
    visualMap: {
      min: 0, // 最小值
      max: maxCount, // 最大值
      calculable: false, // 是否显示拖拽手柄：否
      orient: 'horizontal', // 布局方向：水平
      left: 'right', // 水平位置：靠右
      bottom: 0, // 垂直位置：底部
      itemWidth: 10, // 图例标记的宽度
      itemHeight: 10, // 图例标记的高度
      text: ['More', 'Less'], // 两端的文本说明
      inRange: {
        // 颜色映射范围
        // 第一个颜色为文章数为0时的颜色 (浅灰色)
        // 后续颜色为文章数增加时的渐变色 (绿色系)
        color: ['#404040', '#0e4429', '#006d32', '#26a641', '#39d353']
      },
      textStyle: {
        color: '#ccc', // 文字颜色
        fontSize: 12 // 文字大小
      }
    },
    // 日历坐标系组件配置
    calendar: {
      top: 30, // 距离容器顶部距离
      left: 30, // 距离容器左侧距离
      right: 30, // 距离容器右侧距离
      cellSize: ['auto', 13], // 单元格尺寸 [宽, 高]
      range: new Date().getFullYear(), // 日历范围：当前年份
      itemStyle: {
        color: '#404040', // 单元格默认背景色 (无数据时) - 浅灰色
        borderWidth: 3, // 单元格边框宽度 (用于模拟间距)
        borderColor: '#252526' // 单元格边框颜色 (与容器背景色一致，形成间距效果)
      },
      yearLabel: { show: false }, // 年份标签：不显示
      dayLabel: { 
        color: '#ccc', // 星期标签颜色
        nameMap: 'en', // 星期显示语言：英文
        firstDay: 1 // 一周开始于：周一
      },
      monthLabel: { color: '#ccc', nameMap: 'en' }, // 月份标签配置
      splitLine: { show: false } // 分隔线：不显示
    },
    // 系列列表配置
    series: [{
      type: 'heatmap', // 图表类型：热力图
      coordinateSystem: 'calendar', // 坐标系：日历坐标系
      data: dateList, // 数据源
      itemStyle: {
        borderRadius: 2, // 数据块圆角
        borderWidth: 3, // 数据块边框宽度
        borderColor: '#252526' // 数据块边框颜色 (与背景一致)
      }
    }]
  }

  // Handle cross-year (show last 12 months)
  // ECharts calendar range can be an array ['2025-01-01', '2026-01-01']
  const dates = Object.keys(calendarData.value).sort()
  if (dates.length > 0) {
      option.calendar.range = [dates[0], dates[dates.length - 1]]
  }

  chartInstance.setOption(option)
  
  chartInstance.on('click', function (params) {
    selectedDate.value = params.data[0]
    selectedCount.value = params.data[1]
    fetchArticlesByDate(params.data[0])
  })
}

const resizeHandler = () => {
  chartInstance && chartInstance.resize()
}

const fetchArticlesByDate = async (date) => {
  if (!date) return
  articlesLoading.value = true
  selectedArticles.value = []
  try {
    const res = await getArticles({ page: 1, size: 100, date: date })
    if (res.data.code === 10000) {
      selectedArticles.value = res.data.data.list || []
    }
  } catch (e) {
    console.error(e)
  } finally {
    articlesLoading.value = false
  }
}

const goToArticle = (sn) => {
    router.push(`/articles/${sn}`)
}

const fetchData = async () => {
  // Check cache
  const cached = localStorage.getItem(CACHE_KEY)
  if (cached) {
    try {
      const { data, timestamp } = JSON.parse(cached)
      if (Date.now() - timestamp < CACHE_EXPIRY) {
        calendarData.value = data
        initChart()
        return
      }
    } catch (e) {
      localStorage.removeItem(CACHE_KEY)
    }
  }

  loading.value = true
  try {
    const res = await getArticlesCalendar()
    if (res.data.code === 10000) {
      calendarData.value = res.data.data
      // Save to cache
      localStorage.setItem(CACHE_KEY, JSON.stringify({
        data: calendarData.value,
        timestamp: Date.now()
      }))
      initChart()
    } else {
      error.value = res.data.msg || 'Failed to load calendar'
    }
  } catch (e) {
    error.value = 'Network error'
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
  window.addEventListener('resize', resizeHandler)
})

onUnmounted(() => {
  window.removeEventListener('resize', resizeHandler)
  chartInstance && chartInstance.dispose()
})
</script>

<style scoped>
/* Ensure tooltip text is readable */
:deep(.echarts-tooltip) {
  color: #333 !important;
}
</style>
