import http  from '@/js/axios'
import { defineStore } from 'pinia'

const userMovieStore = defineStore('movie', {
	state: () => ({ 
		token: '',
		user: {}
	}),
	actions: {
		async listMovie(number) {
			return await http.get(`/api/movies?offset=${number}`)
		},

		async getMovie(id) {
			return await http.get(`/api/movies/${id}`)
		},

		async getMoviePlayer(id) {
			return await http.get(`/api/movies_players/${id}`)
		}
	},
})

export default userMovieStore