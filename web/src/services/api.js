import axios from 'axios'

export function useApi(host) {
  const instance = axios.create({
    baseURL: host,
    timeout: 10000,
    withCredentials: true,
    headers: {}
  })

  instance.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest'
  instance.defaults.maxRedirects = 0;

  instance.interceptors.response.use(function (response) {
    return response
  }, function (error) {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          window.location.reload()
          break
      }
    }

    return Promise.reject(error)
  })

  return instance
}
