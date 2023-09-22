<template>
  <el-container class="main-layout">
        <el-aside width="50px">
            <div v-for="menu of menus" class="menu-item" :key="menu.name">
                <div>{{menu.name}}</div>
            </div>
        </el-aside>
        <el-main class="main-body">
            <router-view></router-view>
        </el-main>
    </el-container>
</template>

<script setup>
import { reactive } from 'vue'

const menus = reactive([
    {name: '电影', path: '/sessions'},
    {name: '电视剧', path: '/users'},
    {name: '动漫', path: '/requests'},
    {name: '综艺', path: '/moments'},
])
</script>

<style less>
#app {
    font-family: Avenir, Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    width: 100%;
    height: 100%;
    overflow-y: auto;
}
</style>
