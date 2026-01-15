<template>
  <div>
    <h2 class="text-xl font-semibold mb-6 flex justify-between">
      Advert Management
      <el-button type="primary" size="small" @click="handleAdd">Add Advert</el-button>
    </h2>

    <el-table :data="adverts" style="width: 100%" v-loading="loading">
      <el-table-column prop="title" label="Title" />
      <el-table-column label="Image" width="100">
        <template #default="scope">
          <el-image :src="formatUrl(scope.row.images)" class="w-16 h-10 rounded" fit="cover" />
        </template>
      </el-table-column>
      <el-table-column prop="href" label="Link" />
      <el-table-column label="Visible" width="80">
         <template #default="scope">
            <el-switch v-model="scope.row.is_show" @change="toggleShow(scope.row)" />
         </template>
      </el-table-column>
      <el-table-column label="Actions" width="100">
        <template #default="scope">
          <el-button type="danger" link @click="handleDelete(scope.row.sn)">Delete</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="New Advert" width="30%">
      <el-form :model="form" label-width="80px">
        <el-form-item label="Title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="Link">
          <el-input v-model="form.href" placeholder="可选，留空则仅展示图片" />
        </el-form-item>
        <el-form-item label="Image">
          <el-upload
            action=""
            :http-request="handleImageUpload"
            :show-file-list="false"
            accept="image/*"
          >
            <el-button type="primary" size="small">Upload Image</el-button>
          </el-upload>
          <div v-if="form.images" class="mt-3">
            <el-image :src="formatUrl(form.images)" class="w-24 h-16 rounded" fit="cover" />
          </div>
        </el-form-item>
        <el-form-item label="Show">
          <el-switch v-model="form.is_show" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">Cancel</el-button>
          <el-button type="primary" @click="handleSubmit">Confirm</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
/**
 * AdvertList.vue
 * 
 * @description 后台广告管理页面。允许上传、编辑和删除广告。
 * @author GVB Admin
 * @last_modified 2026-01-14
 * @requires element-plus, vue, ../api/advert, @element-plus/icons-vue
 */
import { ref, onMounted } from 'vue'
import { getAdverts, deleteAdverts, createAdvert, updateAdvert, updateAdvertShow, uploadAdvertImage } from '../api/advert'
import { formatUrl } from '@/utils/url'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'

const adverts = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const form = ref({
    title: '',
    href: '',
    images: '',
    is_show: true
})
const isEdit = ref(false)
const editId = ref('')

/**
 * 获取广告列表
 */
const fetchData = async () => {
    loading.value = true
    try {
        const res = await getAdverts({ page: 1, size: 100 })
        if (res.data.code === 10000) {
            const d = res.data.data
            adverts.value = Array.isArray(d) ? d : (d.list || [])
        }
    } catch (e) {
        ElMessage.error("获取广告失败")
    } finally {
        loading.value = false
    }
}

/**
 * 打开添加对话框
 */
const handleAdd = () => {
    isEdit.value = false
    form.value = { title: '', href: '', images: '', is_show: true }
    dialogVisible.value = true
}

/**
 * 打开编辑对话框
 * @param {Object} row - 广告对象
 */
const handleEdit = (row) => {
    isEdit.value = true
    editId.value = row.id
    form.value = { ...row }
    dialogVisible.value = true
}

/**
 * 提交表单
 */
const handleSubmit = async () => {
	if (!form.value.title) {
		ElMessage.error('Title 不能为空')
		return
	}
	if (!form.value.images) {
		ElMessage.error('请先上传广告图片')
		return
	}
    try {
        let res
        if (isEdit.value) {
            res = await updateAdvert(editId.value, form.value)
        } else {
            res = await createAdvert(form.value)
        }
        if (res.data.code === 10000) {
            ElMessage.success(isEdit.value ? "更新成功" : "创建成功")
            dialogVisible.value = false
            fetchData()
        }
    } catch (e) {
        ElMessage.error("操作失败")
    }
}

const handleImageUpload = async (param) => {
	const fd = new FormData()
	fd.append('images', param.file)
	try {
		const res = await uploadAdvertImage(fd)
		if (res.data.code === 10000) {
			const list = res.data.data || []
			const items = Array.isArray(list) ? list : []
			const okItem = items.find((i) => !i.msg && i.file_name)
			const firstError = items.find((i) => i.msg)
			if (okItem && okItem.file_name) {
				form.value.images = okItem.file_name
			}
			if (firstError && firstError.msg) {
				ElMessage.warning(firstError.msg || '图片上传存在问题')
			} else {
				ElMessage.success('图片上传成功')
			}
		} else {
			ElMessage.error(res.data.msg || '图片上传失败')
		}
	} catch (e) {
		ElMessage.error('图片上传失败')
	}
}

const toggleShow = async (row) => {
    const prev = !row.is_show
    try {
        const res = await updateAdvertShow(row.sn, row.is_show)
        if (res.data.code !== 10000) {
            row.is_show = prev
            ElMessage.error(res.data.msg || 'Update failed')
        }
    } catch (e) {
        row.is_show = prev
        ElMessage.error('Network error')
    }
}

/**
 * 删除广告
 * @param {string} id - 广告 ID
 */
const handleDelete = (id) => {
    ElMessageBox.confirm('确定删除该广告吗?', '提示', {
        type: 'warning'
    }).then(async () => {
        const res = await deleteAdverts([id])
        if (res.data.code === 10000) {
            ElMessage.success("删除成功")
            fetchData()
        }
    })
}

onMounted(() => {
    fetchData()
})
</script>
