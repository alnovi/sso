import axios from 'axios'

const api = axios.create({
    baseURL: '/',
    timeout: 30000,
    withCredentials: true,
    headers: {}
})

export default api
