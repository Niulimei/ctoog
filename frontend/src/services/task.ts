import request from '@/utils/request';

/** 获取任务列表 */
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

/** 刷新迁移任务 */
export const refreshTask = (id: string) => {
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
  return request.get(`/pvobs/${encodeURIComponent(pvobId)}/components/${encodeURIComponent(componentId)}/streams`);
};
