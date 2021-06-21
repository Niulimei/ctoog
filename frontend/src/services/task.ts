import request from '@/utils/request';
import {isArray} from 'lodash';

/** 获取迁移任务 */
export const getTasks = (params: API.PaginationRequestParams) => {
  return request.get('/tasks', {
    params,
  });
};

/** 获取任务详情 */
export const getTaskDetail = (id: string, modelType: string) => {
  return request.get(`/tasks/${id}`, {
    modelType
  });
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
export const startTask = (ids: [] | string) => {
  const idList = isArray(ids) ? ids : [ids];
  return request.post('/tasks/restart', {
    data: {
      id: idList,
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

/** 获取任务节点列表 */
export const getWorkList = (limit: number, offset: number) => {
  return request.get(`/workers`, {
    params: {
      limit,
      offset
    }
  });
};

/** 验证cc和git用户名密码是否正确 */
export const checkCg = (ccInfo: any, gitInfo: any)=>{
  console.log(ccInfo, gitInfo);
  
}
