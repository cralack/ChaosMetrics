<template>
  <el-form @submit.native.prevent="login">
    <el-form-item label="Username">
      <el-input v-model="loginForm.username" placeholder="Username"></el-input>
    </el-form-item>
    <el-form-item label="Password">
      <el-input type="password" v-model="loginForm.password" placeholder="Password"></el-input>
    </el-form-item>
    <el-form-item label="Emai">
      <el-input type="text" v-model="loginForm.email" placeholder="Email"></el-input>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="login">Login</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import axios from 'axios';
export default {
  data() {
    return {
      loginForm: {
        username: '',
        password: '',
        email:'',
      }
    };
  },
  methods: {
    login() {
      axios.post('http://localhost:8080/user/login', this.loginForm)
          .then(response => {
            console.log('Login success:', response);
            this.$message.success('Login successful');
            // Handle login success, e.g., redirect or store token
          })
          .catch(error => {
            console.error('Login error:', error);
            this.$message.error('Login failed');
          });
    }
  }
}
</script>
