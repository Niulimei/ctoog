import request from '@/utils/request';

/** 获取迁移任务 */
export const getTasks = (params: API.PaginationRequestParams) => {
  return request.get('/tasks', {
    params,
  });
};

/** 获取任务详情 */
export const getTaskDetail = (id: string) => {
  return request.get(`/tasks/${id}`);
};

/** 创建迁移任务 */
export const createTask = (data: any) => {
  return request.post('/tasks', {
    data,
  });
};

/** 更新迁移任务 */
export const updateTask = (id: string, data: any) => {
  return request.put(`/tasks/${id}`, {
    data,
  });
};

/** 启动迁移任务 */
export const startTask = (id: number) => {
  return request.post('/tasks/restart', {
    data: {
      id,
    },
  });
};

/** 获取 pvob 列表 */
export const getPvobs = () => {
  return request.get('/pvobs');
};

/** 获取 component 列表 */
export const getComponents = (pvobId: string) => {
  return request.get(`/pvobs/${encodeURIComponent(pvobId)}/components`);
};

/** 获取 stream 列表 */
export const getStreams = (pvobId: string, componentId: string) => {
  return request.get(
    `/pvobs/${encodeURIComponent(pvobId)}/components/${encodeURIComponent(componentId)}/streams`,
  );
};

/** 获取 log output */
export const getLogOutput = (logId: string) => {
  return request.get(`/tasks/cmdout/${logId}`);
};

/** 按钮删除任务 */
export const deleteTask = (taskId: string) => {
  return request.delete(`/tasks/${taskId}`);
};

/** 按钮删除缓存 */
export const deleteCache = (taskId: string) => {
  return request.delete(`/tasks/cache/${taskId}`);
};
