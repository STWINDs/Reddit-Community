<template>
  <el-form
    ref="formRef"
    :model="form"
    :rules="rules"
    label-width="0px"
    class="register-form-content"
  >
    <el-form-item prop="username">
      <el-input
        v-model="form.username"
        placeholder="用户名"
        prefix-icon="User"
        size="large"
      />
    </el-form-item>
    <el-form-item prop="email">
      <el-input
        v-model="form.email"
        placeholder="邮箱"
        prefix-icon="Message"
        size="large"
      />
    </el-form-item>
    <el-form-item prop="password">
      <el-input
        v-model="form.password"
        type="password"
        placeholder="密码"
        prefix-icon="Lock"
        size="large"
        show-password
      />
    </el-form-item>
    <el-form-item prop="confirmPassword">
      <el-input
        v-model="form.confirmPassword"
        type="password"
        placeholder="确认密码"
        prefix-icon="Lock"
        size="large"
        show-password
      />
    </el-form-item>
    <el-form-item>
      <el-button
        type="primary"
        :loading="loading"
        class="submit-btn"
        size="large"
        @click="handleSubmit"
      >
        注册
      </el-button>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { User, Lock, Message } from '@element-plus/icons-vue'
import { validateUsername, validateEmail, validatePassword } from '@/utils/validate'

const props = defineProps({
  loading: Boolean
})

const emit = defineEmits(['submit'])

const formRef = ref(null)
const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const validatePass = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请输入密码'))
  } else {
    const { valid, message } = validatePassword(value)
    if (!valid) {
      callback(new Error(message))
    } else {
      if (form.confirmPassword !== '') {
        formRef.value.validateField('confirmPassword')
      }
      callback()
    }
  }
}

const validatePass2 = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== form.password) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { validator: (rule, value, cb) => {
        const { valid, message } = validateUsername(value)
        valid ? cb() : cb(new Error(message))
      }, trigger: 'blur' 
    }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }
  ],
  password: [
    { required: true, validator: validatePass, trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validatePass2, trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate((valid) => {
    if (valid) {
      // Exclude confirmPassword from submission
      const { confirmPassword, ...submitData } = form
      emit('submit', submitData)
    }
  })
}
</script>

<style scoped>
.submit-btn {
  width: 100%;
}
</style>
