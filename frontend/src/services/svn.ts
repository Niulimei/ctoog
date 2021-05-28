import request from '@/utils/request';

/** SVN */
export const getSvn = (params: API.PaginationRequestParams) => {
  return request.get('/svn_username_pairs', {
    params,
  });
};
