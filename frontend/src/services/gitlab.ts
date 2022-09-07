import request from '@/utils/request';

/** 获取迁移任务 */
export const getTasks = (params: API.PaginationRequestParams) => {
  return request.get('/gitlab/tasks', {
    params,
  });
};
export const getTasks1 = (params: API.PaginationRequestParams) => {
  return request.get('/gitlab/tasks1', {
    params,
  });
};