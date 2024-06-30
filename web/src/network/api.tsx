import axios from 'axios';
import config from './config';
import Cookies from 'js-cookie';
const env = process.env.NODE_ENV || 'production';

// 创建 Axios 实例
export const api = axios.create({
    baseURL: config[env].apiUrl,
    withCredentials: true,
    headers: {
        "Content-Type": "application/json; charset=UTF-8",
    }
});


// 响应拦截器
api.interceptors.response.use(
    (response) => {
        return response;
    },
    (error) => {
        if (error.response.status === 401) {
            // window.location.href = '/noauth'
        }
        if (error.response.status === 403) {
            window.location.href = '/forbidden';
        }
        if (error.response.status === 410) {
            Object.keys(Cookies.get()).forEach(function (cookieName) {
                Cookies.remove(cookieName);
                console.log(cookieName);
            });
            window.location.href = '/noauth';
        }
        if (error.response.status === 500) {
            alert("程序出错了！快联系开发者看看！\n" + error.response.data);
        }
        if (error.response.status === 502) {
            if (error.response.data === "") {
                alert("服务器出故障了！快联系开发者看看！");
            }
        }
        return Promise.reject(error);
    }
);
