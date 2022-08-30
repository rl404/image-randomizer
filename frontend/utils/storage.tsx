const accessTokenKey = `access-token`;
const refreshTokenKey = `refresh-token`;
const usernameKey = `username`;

export const deleteStorage = () => {
  localStorage.clear();
};

export const saveUsername = (username: string) => {
  localStorage.setItem(usernameKey, JSON.stringify(username));
};

export const getUsername = (): string => {
  const username = localStorage.getItem(usernameKey);
  if (!username) return '';
  return JSON.parse(username);
};

export const saveAccessToken = (token: string) => {
  localStorage.setItem(accessTokenKey, JSON.stringify(token));
};

export const getAccessToken = (): string => {
  const token = localStorage.getItem(accessTokenKey);
  if (!token) return '';
  return JSON.parse(token);
};

export const saveRefreshToken = (token: string) => {
  localStorage.setItem(refreshTokenKey, JSON.stringify(token));
};

export const getRefreshToken = (): string => {
  const token = localStorage.getItem(refreshTokenKey);
  if (!token) return '';
  return JSON.parse(token);
};
