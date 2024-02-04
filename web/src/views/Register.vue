<template>
  <el-form @submit.native.prevent="register">
    <el-form-item label="Username">
      <el-input v-model="registerForm.username" placeholder="Username"></el-input>
    </el-form-item>
    <el-form-item label="Email">
      <el-input v-model="registerForm.email" placeholder="Email"></el-input>
    </el-form-item>
    <el-form-item label="Password">
      <el-input type="password" v-model="registerForm.password" placeholder="Password"></el-input>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="register">Register</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import axios from 'axios';
export default {
  data() {
    return {
      registerForm: {
        username: '',
        email: '',
        password: ''
      }
    };
  },
  methods: {
    register() {
      axios.post('http://localhost:8080/user/register', this.registerForm)
        .then(response => {
          console.log('Registration success:', response);
          this.$message.success('Registration successful');
          // Handle registration success, e.g., redirect
        })
        .catch(error => {
          console.error('Registration error:', error);
          this.$message.error('Registration failed');
        });
    }
  }
}
</script>
