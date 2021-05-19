import request from '@/utils/request';
import type { Plan } from '@/typings/model';

/** 创建计划计划 */
export const createPlan = (data: any) => {
  return request.post('/plans', {
    data,
  });
};

/** 获取创建计划列表 */
export const getPlans = (params: API.PaginationRequestParams) => {
  return request.get('/plans', { params });
};

/** 获取创建计划列表 */
export const getPlanDetail = (id: string) => {
  return request.get(`/plans/${id}`);
};

/** 修改创建计划 */
export const updatePlan = (id: string, data: any) => {
  return request.put(`/plans/${id}`, {
    data,
  });
};

/** 删除创建计划 */
export const deletePlan = (id: string) => {
  return request.delete(`/plans/${id}`);
};
