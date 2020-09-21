export const Host = process.env.REACT_APP_API_HOST

export function getUserURL(username) {
  return `${Host}/user/${username}/image.jpg`
}

export async function login(data) {
  const result = await fetch(`${Host}/user/login`, {
    method: 'POST',
    body: JSON.stringify(data)
  })
  return result.json()
}

export async function register(data) {
  const result = await fetch(`${Host}/user/register`, {
    method: 'POST',
    body: JSON.stringify(data)
  })
  return result.json()
}

export async function getList(username, token) {
  const result = await fetch(`${Host}/user/${username}`, {
    method: 'GET',
    headers: {
      'token': token,
    }
  })
  return result.json()
}

export async function updateList(token, data) {
  const result = await fetch(`${Host}/user/update`, {
    method: 'POST',
    body: JSON.stringify(data),
    headers: {
      'token': token,
    }
  })
  return result.json()
}