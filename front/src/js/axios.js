import axios from "axios";
import { ElMessage } from 'element-plus';
import router from '../router';

const http = axios.create({
	baseURL: '/',
	timeout: 50000
})

http.interceptors.request.use((config) => {
	if (sessionStorage.getItem('token')) {
		config.headers['token'] = sessionStorage.getItem('token')
	  }
	  return config;
}, (error) => {
	return Promise.reject(error);
});

http.interceptors.response.use((res) => {
	if (res.config.url.indexOf('login') > -1) {
		return res;
	}
	if (res.status === 401) {
		// TODO: 退出登录并重新登录
	}

	if (res.status === 200) {
		return res.data;
	}
	}, (error) => {
		if (error.response.status) {
		switch (error.response.status) {
			case 401:
				router.push('/')
				break;
			case 404:
				ElMessage({
					type: 'error',
					message: '请求路径找不到！',
					showClose: true
				});
				break;
			case 500:
			case 502:
				ElMessage({
					type: 'error',
					message: '服务器内部报错！',
					showClose: true
				});
				break;
		default:
			break;
		}
	}
	return Promise.reject(error);
});

export default http;