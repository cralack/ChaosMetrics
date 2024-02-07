<template>
  <div>
    <el-result
      icon="warning"
      title="404 Not Found"
      sub-title="页面走丢了"
    >
      <template #extra>
        <el-button
          type="primary"
          @click="$router.push('/')"
        >Back</el-button>
      </template>
    </el-result>
  </div>
</template>
