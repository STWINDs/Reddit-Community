<template>
  <div class="home-container">
    <Header />
    <main class="main-content">
      <div class="content-wrapper">
        <div class="home-layout">
          <!-- Left Column: Post List -->
          <div class="post-column">
            <PostList ref="postListRef" @create-post="handleOpenPostForm" />
          </div>
          
          <!-- Right Column: Sidebar (Future) -->
          <div class="sidebar-column">
            <el-card class="about-card">
              <template #header>
                <div class="card-header">
                  <span>关于 Bluebell</span>
                </div>
              </template>
              <p>一个简单而强大的论坛，使用 Go + Vue 3 构建。</p>
            </el-card>
          </div>
        </div>
      </div>
    </main>
    
    <!-- Post Creation Modal -->
    <PostForm ref="postFormRef" @success="handlePostSuccess" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Header from '@/components/common/Header.vue'
import PostList from '@/components/post/PostList.vue'
import PostForm from '@/components/post/PostForm.vue'

const postListRef = ref(null)
const postFormRef = ref(null)

const handleOpenPostForm = () => {
  if (postFormRef.value) {
    postFormRef.value.open()
  }
}

const handlePostSuccess = () => {
  if (postListRef.value) {
    postListRef.value.refresh()
  }
}
</script>

<style scoped lang="scss">
.home-container {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.main-content {
  padding: 20px;
  
  .content-wrapper {
    max-width: 1200px;
    margin: 0 auto;
    
    .home-layout {
      display: flex;
      gap: 20px;
      
      .post-column {
        flex: 1;
        min-width: 0; // Prevent flex item from expanding beyond its container
      }
      
      .sidebar-column {
        width: 300px;
        
        @media (max-width: 992px) {
          display: none;
        }
        
        .about-card {
          margin-bottom: 20px;
          
          .card-header {
            font-weight: bold;
          }
          
          p {
            font-size: 14px;
            color: #606266;
            line-height: 1.6;
          }
        }
      }
    }
  }
}
</style>
