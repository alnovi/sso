export const config = (par, def) => {
  if (def === undefined) {
    def = null
  }

  let val = import.meta.env[par]

  return val ? val : def
}

export const meta = (par, def) => {
  if (def === undefined) {
    def = null
  }

  let metaObj = document.head.querySelector(`meta[name=${par}]`)
  if (!metaObj) {
    return def
  }

  let value = metaObj.getAttribute("content")
  if (!value) {
    return def
  }

  if (value === "<nil>" || value.startsWith("{{")) {
    return def
  }

  return value
}

export const validMsg = (message, field, name) => {
  if (message && field) {
    message = message.replace(field, name)
  }

  if (message && message.length > 0) {
    message = message.charAt(0).toUpperCase() + message.slice(1)
  }

  return message || ""
}

export const validStatus = (message) => {
  return !!message ? 'error' : ''
}
