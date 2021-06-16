import request from '@/utils/request';

/** 创建计划计划 */
export const createPlan = (data: any) => {
  return request.post('/plans', {
    data,
  });
};

/** 获取创建迁移计划 */
export const getPlans = (params: API.PaginationRequestParams) => {
  return request.get('/plans', { params });
};

/** 获取创建迁移计划 */
export const getPlanDetail = (id: string, params: any) => {
  return request.get(`/plans/${id}`, {params});
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

/** 批量导入计划 */
export const importPlan = data => {
  return request.post(`/import/plan`, data);
};
