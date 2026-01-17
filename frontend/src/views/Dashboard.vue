<template>
  <div class="space-y-6">
    <!-- Stat Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-vscode-sidebar p-4 rounded border border-vscode-border flex flex-col items-center justify-center hover:border-vscode-primary transition-colors cursor-default">
        <div class="text-3xl font-bold text-white mb-2">{{ sumData.article_count }}</div>
        <div class="text-xs text-gray-400 uppercase tracking-wider">Articles</div>
      </div>
      <div class="bg-vscode-sidebar p-4 rounded border border-vscode-border flex flex-col items-center justify-center hover:border-vscode-primary transition-colors cursor-default">
        <div class="text-3xl font-bold text-white mb-2">{{ sumData.user_count }}</div>
        <div class="text-xs text-gray-400 uppercase tracking-wider">Users</div>
      </div>
      <div class="bg-vscode-sidebar p-4 rounded border border-vscode-border flex flex-col items-center justify-center hover:border-vscode-primary transition-colors cursor-default">
        <div class="text-3xl font-bold text-white mb-2">{{ sumData.message_count }}</div>
        <div class="text-xs text-gray-400 uppercase tracking-wider">Messages</div>
      </div>
      <div class="bg-vscode-sidebar p-4 rounded border border-vscode-border flex flex-col items-center justify-center hover:border-vscode-primary transition-colors cursor-default">
        <div class="text-3xl font-bold text-white mb-2">{{ sumData.chat_group_count }}</div>
        <div class="text-xs text-gray-400 uppercase tracking-wider">Chat Groups</div>
      </div>
    </div>

    <!-- Charts -->
    <div class="bg-vscode-sidebar p-6 rounded border border-vscode-border">
      <h3 class="text-lg font-medium text-white mb-6">User Login Trend (Last 7 Days)</h3>
      <div ref="chartRef" class="w-full h-[400px]"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, reactive } from 'vue'
import * as echarts from 'echarts'
import { getDataSum, getUserLoginData } from '@/api/statistic'
import { ElMessage } from 'element-plus'

const sumData = reactive({
  user_count: 0,
  article_count: 0,
  message_count: 0,
  chat_group_count: 0
})

const chartRef = ref(null)
let chartInstance = null

const initChart = (data) => {
  if (!chartRef.value) return
  
  chartInstance = echarts.init(chartRef.value)
  
  const dates = data.map(item => {
    const date = new Date(item.login_date)
    return `${date.getMonth() + 1}-${date.getDate()}`
  })
  const counts = data.map(item => item.login_count)

  const option = {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'line'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dates,
      axisLine: {
        lineStyle: {
          color: '#4B5563'
        }
      },
      axisLabel: {
        color: '#9CA3AF'
      }
    },
    yAxis: {
      type: 'value',
      axisLine: {
        show: false
      },
      splitLine: {
        lineStyle: {
          color: '#374151'
        }
      },
      axisLabel: {
        color: '#9CA3AF'
      }
    },
    series: [
      {
        name: 'Logins',
        type: 'line',
        smooth: true,
        data: counts,
        symbol: 'circle',
        symbolSize: 8,
        itemStyle: {
          color: '#FFA500'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            {
              offset: 0,
              color: 'rgba(255, 165, 0, 0.5)'
            },
            {
              offset: 1,
              color: 'rgba(255, 165, 0, 0.05)'
            }
          ])
        }
      }
    ]
  }

  chartInstance.setOption(option)
}

const fetchData = async () => {
  try {
    // Parallel requests but handle errors individually
    const sumPromise = getDataSum().catch(e => {
        console.error('Failed to get data sum', e)
        return { data: { code: -1 } }
    })
    const loginPromise = getUserLoginData().catch(e => {
        console.error('Failed to get user login data', e)
        return { data: { code: -1 } }
    })

    const [sumRes, loginRes] = await Promise.all([sumPromise, loginPromise])

    if (sumRes.data && sumRes.data.code === 10000) {
      Object.assign(sumData, sumRes.data.data)
    }
    
    if (loginRes.data && loginRes.data.code === 10000) {
      initChart(loginRes.data.data)
    }
  } catch (error) {
    ElMessage.error('Failed to load dashboard data')
    console.error(error)
  }
}

const handleResize = () => {
  chartInstance?.resize()
}

onMounted(() => {
  fetchData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance?.dispose()
})
</script>
