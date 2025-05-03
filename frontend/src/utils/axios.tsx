import axios, { InternalAxiosRequestConfig } from 'axios';
import { deleteStorage, getAccessToken, getRefreshToken, saveAccessToken } from './storage';

export const axios2 = axios.create();

axios2.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    if (config.headers) config.headers.authorization = `Bearer ${getAccessToken()}`;
    return config;
  },
  (error) => {
    Promise.reject(error);
  },
);

axios2.interceptors.response.use(
  (resp) => {
    return resp;
  },
  (error) => {
    const originalRequest = error.config;
    if (error.response.status !== 401 || originalRequest._retry) {
      return Promise.reject(error);
    }

    originalRequest._retry = true;

    return axios
      .post('/api/token/refresh', null, {
        headers: {
          Authorization: `Bearer ${getRefreshToken()}`,
        },
      })
      .then((resp) => {
        saveAccessToken(resp.data.data.access_token);
        return axios2(originalRequest);
      })
      .catch((err) => {
        deleteStorage();
        if (err.response && err.response.status === 401) {
          window.location.href = '/';
          return Promise.reject(error);
        }
        return Promise.reject(error);
      });
  },
);
