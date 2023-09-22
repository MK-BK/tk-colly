<template>
   <el-card class="card">
        <el-image style="width: 300px; height: 300px" :src="movie.ImagePath" fit="cover" />
        <div @click="viewMovie(movie)">{{movie.Name}}</div>
        <div>{{ movie.Description }}</div>
    </el-card>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
const props = defineProps(['movie'])

let movie = props.movie

function viewMovie(movie) {
    router.push('/movie/' + movie.ID)
}
</script>

<style lang="less" scoped>
.card {
    height: 400px;
    margin-bottom: 20px;
    margin: 20px 20px;
}
</style>