<template>
  <div class="not-found">
    <h1>404</h1>
    <p>页面未找到</p>
    <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
  </div>
</template>

<style scoped>
.not-found {
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #f0f2f5;
}

h1 {
  font-size: 72px;
  color: #409EFF;
  margin-bottom: 20px;
}

p {
  font-size: 24px;
  color: #606266;
  margin-bottom: 30px;
}
</style>
