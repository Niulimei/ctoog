import { extend } from 'umi-request';
import { message } from 'antd';

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
});

request.interceptors.request.use((url, { headers, ...restOpts }) => {
  return {
    url,
    options: {
      ...restOpts,
      headers: {
        ...headers,
        Authorization: authTokenAction.get(),
      },
    },
  };
});

request.interceptors.response.use(async (res) => {
  const responseBody = await res.clone().json();

  if (res.status === 401) {
    // TODO: 登陆信息过期，跳转到 login
  }
  if (res.status === 500 && responseBody) {
    message.error(responseBody.message);
  }
  return res;
});

export default request;
