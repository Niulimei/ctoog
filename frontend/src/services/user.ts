import request, { authTokenAction } from '@/utils/request';
import type { User } from '@/typings/model';

/** 登陆操作 */
export const login = async (params: User.Base) => {
  const res = await request.post('/login', {
    data: params,
  });
  if (res.token) {
    authTokenAction.set(res.token);
  }
  return res;
};

export const createUser = (params: User.Base) => {
  return request.post('/users', {
    data: params,
  });
};

/** admin 用户可以获取权限列表 */
export const getUsers = (params: API.PaginationRequestParams) => {
  return request.get('/users', {
    params,
  });
};

export const getCurrentUser = () => {
  return request.get(`/users/self`);
};
