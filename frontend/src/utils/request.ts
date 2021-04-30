import { extend } from 'umi-request';

/** token 存储操作 */
export const authTokenAction = {
  key: 'U_TOKEN',
  set(token: string) {
    localStorage.setItem(this.key, token);
  },
  get() {
    return localStorage.getItem(this.key) || '';
  },
  clear() {
    localStorage.removeItem(this.key);
  },
};

const request = extend({
  prefix: '/api',
  timeout: 10000,
  headers: {
    Authorization: authTokenAction.get(),
  },
});

request.interceptors.response.use((res, options) => {
  if (res.status === 401) {
    // TODO: 登陆信息过期，跳转到 login
  }
  if (res.status === 500 && res.body) {
    // TODO: 登陆失败
  }
  return res;
});

export default request;
