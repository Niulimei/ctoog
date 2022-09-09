import request from '@/utils/request';

/** 获取迁移任务 */
export const getTasks = (params) => {
  return request.get('/tasks', {
    params,
  });
};
export const createTask = (data) => {
  return request.post('/gitlab/tasks', {
    data,
  });
};