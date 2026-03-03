// 邮箱验证
export const validateEmail = (email) => {
  const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return regex.test(email)
}

// 密码强度验证
export const validatePassword = (password) => {
  if (password.length < 6) {
    return { valid: false, message: '密码长度至少6位' }
  }
  if (!/[A-Za-z]/.test(password) || !/[0-9]/.test(password)) {
    return { valid: false, message: '密码必须包含字母和数字' }
  }
  return { valid: true }
}

// 用户名验证
export const validateUsername = (username) => {
  if (username.length < 3 || username.length > 20) {
    return { valid: false, message: '用户名长度应在3-20位之间' }
  }
  if (!/^[A-Za-z0-9_]+$/.test(username)) {
    return { valid: false, message: '用户名只能包含字母、数字和下划线' }
  }
  return { valid: true }
}
