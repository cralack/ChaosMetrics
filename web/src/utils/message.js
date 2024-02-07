import { ElMessage } from 'element-plus'

export function successMessage(str) {
  ElMessage({
    showClose: true,
    message: str,
    type: 'success',
    duration: 1500,
  })
}

export function errorMessage(str) {
  ElMessage({
    showClose: true,
    message: str,
    type: 'error',
    duration: 1500,
  })
}

