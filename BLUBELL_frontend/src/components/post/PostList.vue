<template>
  <div class="post-list-container">
    <!-- Post Creation -->
    <el-button type="primary" class="create-post-btn" @click="$emit('create-post')">
      <el-icon><Plus /></el-icon>
      <span>发布新帖</span>
    </el-button>

    <!-- Sorting Control -->
    <div class="sort-control">
      <el-button-group>
        <el-button 
          :type="currentOrder === 'time' ? 'primary' : 'default'" 
          @click="changeOrder('time')"
        >
          <el-icon><Clock /></el-icon>
          最新
        </el-button>
        <el-button 
          :type="currentOrder === 'score' ? 'primary' : 'default'" 
          @click="changeOrder('score')"
        >
          <el-icon><Histogram /></el-icon>
          最热
        </el-button>
      </el-button-group>
    </div>

    <!-- Loading state -->
    <div v-if="loading" v-loading="loading" class="loading-state"></div>

    <!-- Post Items -->
    <div v-else-if="posts.length > 0" class="post-items">
      <PostCard 
        v-for="post in posts" 
        :key="post.id" 
        :post="post" 
        @click="handlePostClick"
      />
      
      <!-- Pagination -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          layout="prev, pager, next"
          :total="total"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- Empty state -->
    <el-empty v-else description="暂无帖子" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Plus, Clock, Histogram } from '@element-plus/icons-vue'
import PostCard from './PostCard.vue'
import { getPostList2 } from '@/api/post'
import { ElMessage } from 'element-plus'

const emit = defineEmits(['create-post', 'post-click'])

const posts = ref([])
const loading = ref(false)
const currentOrder = ref('time')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(100) // Backend should ideally provide total count

const fetchPosts = async () => {
  loading.value = true
  try {
    const response = await getPostList2({
      page: currentPage.value,
      size: pageSize.value,
      order: currentOrder.value
    })
    posts.value = response.data || []
  } catch (error) {
    ElMessage.error(error.message || '获取帖子列表失败')
  } finally {
    loading.value = false
  }
}

const changeOrder = (order) => {
  if (currentOrder.value === order) return
  currentOrder.value = order
  currentPage.value = 1
  fetchPosts()
}

const handlePageChange = (page) => {
  currentPage.value = page
  fetchPosts()
}

const handlePostClick = (id) => {
  emit('post-click', id)
}

onMounted(() => {
  fetchPosts()
})

// Expose fetchPosts to refresh from parent
defineExpose({ refresh: fetchPosts })
</script>

<style scoped lang="scss">
.post-list-container {
  max-width: 800px;
  margin: 0 auto;
}

.create-post-btn {
  width: 100%;
  margin-bottom: 20px;
  height: 50px;
  font-size: 16px;
}

.sort-control {
  background: #fff;
  padding: 12px;
  border-radius: 4px;
  margin-bottom: 12px;
  border: 1px solid #ebeef5;
}

.loading-state {
  height: 200px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
