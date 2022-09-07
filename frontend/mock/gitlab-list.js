import Mock from 'mockjs';


export default {
  'GET /api/gitlab/tasks': Mock.mock({
    'list|30-60': [
      {
        'taskNo|1': '@title(1)',
        'time|1': '@datetime()',
        'gitlab_group|1': '@title(2, 4)',
        'gitlab_project|1': '@title(2, 3)',
        'gitee_group|1': '@title(2, 3)',
        'gitee_repo|1': '@title(1, 2)',
        'status': '@title(1)',
        'id|0-99999': 1,
      },
    ],
  }),
}
