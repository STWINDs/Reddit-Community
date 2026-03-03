<template>
  <div class="post-form-container">
    <el-dialog 
      v-model="visible" 
      title="发布新帖" 
      width="600px" 
      @close="resetForm"
    >
      <el-form 
        ref="formRef" 
        :model="form" 
        :rules="rules" 
        label-width="80px"
      >
        <el-form-item label="社区" prop="community_id">
          <el-select v-model="form.community_id" placeholder="选择社区" class="w-full">
            <el-option 
              v-for="item in communityList" 
              :key="item.id" 
              :label="item.name" 
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入标题" maxlength="100" show-word-limit />
        </el-form-item>
        
        <el-form-item label="内容" prop="content">
          <el-input 
            v-model="form.content" 
            type="textarea" 
            placeholder="请输入内容" 
            :rows="6" 
            maxlength="1000" 
            show-word-limit
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="visible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">
            发布
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getCommunityList, createPost } from '@/api/post'
import { ElMessage } from 'element-plus'

const visible = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const communityList = ref([])

const form = reactive({
  community_id: '',
  title: '',
  content: ''
})

const rules = {
  community_id: [
    { required: true, message: '请选择社区', trigger: 'change' }
  ],
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { min: 3, max: 100, message: '标题长度在 3 到 100 个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入内容', trigger: 'blur' },
    { min: 5, message: '内容长度至少 5 个字符', trigger: 'blur' }
  ]
}

const fetchCommunities = async () => {
  try {
    const response = await getCommunityList()
    communityList.value = response.data || []
  } catch (error) {
    ElMessage.error('获取社区列表失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        await createPost({
          ...form,
          // Convert community_id to number if needed
          community_id: Number(form.community_id)
        })
        ElMessage.success('发布成功')
        visible.value = false
        emit('success')
      } catch (error) {
        ElMessage.error(error.message || '发布失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

const open = () => {
  visible.value = true
  fetchCommunities()
}

const emit = defineEmits(['success'])

defineExpose({ open })
</script>

<style scoped>
.w-full {
  width: 100%;
}
</style>
