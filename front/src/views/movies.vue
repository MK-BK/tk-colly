<template>
    <el-main class="main-body">
        <div v-if="!movies.length" class="empty-content">movies列表为空</div>
        <el-row>
            <el-col v-for="movie in movies" :key="movie.ID" :span="6">
                <card :movie="movie"></card>
            </el-col>
        </el-row>
        <el-pagination layout="prev, pager, next" :total="50" class="mt-10" @current-change="change" />
    </el-main>
</template>

<script setup>
import card from '@/views/component/card.vue'
import { ref, onMounted } from 'vue'

import useStore from '@/stores'
const { movieStore } = useStore()

const movies = ref([])

onMounted(async () => {
    await refresh()
})

async function refresh() {
    movies.value = await movieStore.listMovie(1)
}

async function change(number) {
    movies.value = await movieStore.listMovie(number)   
}
</script>