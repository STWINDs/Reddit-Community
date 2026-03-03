<template>
  <el-card class="post-card" shadow="hover">
    <div class="post-layout">
      <!-- Vote Control -->
      <div class="vote-control">
        <el-button 
          link 
          :type="voteStatus === 1 ? 'primary' : 'info'" 
          @click.stop="handleVote(1)"
        >
          <el-icon :size="20"><CaretTop /></el-icon>
        </el-button>
        <span class="vote-count" :class="{ 'positive': voteNum > 0, 'negative': voteNum < 0 }">
          {{ voteNum }}
        </span>
        <el-button 
          link 
          :type="voteStatus === -1 ? 'danger' : 'info'" 
          @click.stop="handleVote(-1)"
        >
          <el-icon :size="20"><CaretBottom /></el-icon>
        </el-button>
      </div>

      <!-- Post Content -->
      <div class="post-main">
        <div class="post-header">
          <span class="community-name">c/{{ post.community_name }}</span>
          <span class="dot">·</span>
          <span class="author">Posted by u/{{ post.author_name }}</span>
          <span class="dot">·</span>
          <span class="time">{{ formatTime(post.create_time) }}</span>
        </div>
        <h3 class="post-title" @click="$emit('click', post.id)">{{ post.title }}</h3>
        <p class="post-excerpt">{{ post.content }}</p>
        
        <div class="post-footer">
          <el-button link size="small">
            <el-icon><ChatLineRound /></el-icon>
            <span>评论</span>
          </el-button>
          <el-button link size="small">
            <el-icon><Share /></el-icon>
            <span>分享</span>
          </el-button>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { ref } from 'vue'
import { CaretTop, CaretBottom, ChatLineRound, Share } from '@element-plus/icons-vue'
import { votePost } from '@/api/post'
import { ElMessage } from 'element-plus'

const props = defineProps({
  post: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['click', 'vote-change'])

const voteNum = ref(props.post.vote_num || 0)
const voteStatus = ref(0) // 0: no vote, 1: up, -1: down

const handleVote = async (direction) => {
  // If clicking same direction, cancel vote
  const newDirection = voteStatus.value === direction ? 0 : direction
  
  try {
    await votePost({
      post_id: props.post.id.toString(),
      direction: newDirection.toString() // 将数字转换为字符串
    })
    
    // Update local state (Optimistic update could be better, but let's keep it simple)
    // voteNum logic is handled by backend, but we can reflect it locally
    if (voteStatus.value === 1) voteNum.value--
    if (voteStatus.value === -1) voteNum.value++
    
    if (newDirection === 1) voteNum.value++
    if (newDirection === -1) voteNum.value--
    
    voteStatus.value = newDirection
    emit('vote-change', { id: props.post.id, voteNum: voteNum.value })
  } catch (error) {
    ElMessage.error(error.message || '投票失败')
  }
}

const formatTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString()
}
</script>

<style scoped lang="scss">
.post-card {
  margin-bottom: 12px;
  cursor: pointer;
  :deep(.el-card__body) {
    padding: 0;
  }
}

.post-layout {
  display: flex;
}

.vote-control {
  width: 40px;
  background-color: #f8f9fa;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 0;
  
  .vote-count {
    font-size: 12px;
    font-weight: bold;
    margin: 4px 0;
    color: #1a1a1b;
    
    &.positive { color: #ff4500; }
    &.negative { color: #7193ff; }
  }
}

.post-main {
  flex: 1;
  padding: 8px 12px;
  
  .post-header {
    display: flex;
    align-items: center;
    font-size: 12px;
    color: #787c7e;
    margin-bottom: 8px;
    
    .community-name {
      color: #1a1a1b;
      font-weight: bold;
      margin-right: 4px;
    }
    
    .dot {
      margin: 0 4px;
    }
  }
  
  .post-title {
    font-size: 18px;
    font-weight: 500;
    margin: 0 0 8px 0;
    color: #222;
    
    &:hover {
      text-decoration: underline;
    }
  }
  
  .post-excerpt {
    font-size: 14px;
    color: #4a4a4a;
    line-height: 1.5;
    margin-bottom: 8px;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  
  .post-footer {
    display: flex;
    gap: 12px;
    color: #878a8c;
  }
}
</style>
