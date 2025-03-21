const duration = 5000

export const notifyInfo = (msg) => {
  return {title: 'Информация', content: msg, duration: duration, keepAliveOnHover: true}
}

export const notifyWarn = (msg) => {
  return {title: 'Внимание', content: msg, duration: duration, keepAliveOnHover: true}
}

export const notifyError = (msg) => {
  return {title: 'Ошибка', content: msg, duration: duration, keepAliveOnHover: true}
}
