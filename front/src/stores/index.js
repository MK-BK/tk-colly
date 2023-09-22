import userMovieStore from './movie'

export default function useStore() {
    return {
        movieStore: userMovieStore(),
    }
}