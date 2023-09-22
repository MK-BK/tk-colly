import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
    {
		path: '/',
		component: () => import(/* webpackChunkName: "about" */ '@/views/movies.vue'),
	},
	{
		path: '/movie/:id',
		component: () => import(/* webpackChunkName: "system" */ '@/views/movies_views.vue'),
	},
]

const router = createRouter({
	history: createWebHashHistory(),
	routes
})
export default router;
