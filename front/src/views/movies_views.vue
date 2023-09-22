<template>
    <div>
        <img :src="object.movie.ImagePath">
        <div>{{ object.movie.Description }}</div>

        <li v-for="href of object.players.Players" v-bind:key="href">
            <div>{{ href.URL }}</div>
        </li>
    </div>
</template>

<script setup>
import { onMounted, watch, reactive } from 'vue';

import { useRoute } from 'vue-router';
const route = useRoute()

import useStore from '@/stores'
const { movieStore } = useStore()

const object = reactive({
    movie: {},
    players: []
})

onMounted(async() => {
    await refresh()
    watch(() => route.params.id, async () => {
        await refresh();
    });
});

async function refresh() {
    object.movie= await movieStore.getMovie(route.params.id)
    object.players = await movieStore.getMoviePlayer(route.params.id)
}
</script>