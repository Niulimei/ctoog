export default [
  {
    path: '/task',
    name: '任务管理',
    icon: 'dashboard',
    access: 'isLogin',
    routes: [
      {
        name: '迁移计划',
        path: '/task/plan',
        component: './task/Plan',
        exact: true,
        access: 'jianxin'
      },
      {
        name: 'CC迁移任务',
        path: '/task/list',
        component: './task/List',
        exact: true,
        access: 'ccRoute'
      },
      {
        name: 'SVN迁移任务',
        path: '/task/svn',
        component: './task/Svn',
        exact: true,
        access: 'svnRoute'
      },
      {
        hideInMenu: true,
        name: '任务详情',
        path: '/task/detail',
        component: './task/Detail',
      },
      {
        hideInMenu: true,
        name: 'SVN任务详情',
        path: '/task/svnDetail',
        component: './task/Detail/svnList',
        access: 'svnRoute'
      },
      {
        access: 'isAdmin',
        name: '任务执行节点',
        path: '/task/node',
        component: './task/Node',
        exact: true,
      },
      {
        name: 'Gitlab迁移任务',
        path: '/task/gitlab',
        component: './task/gitlab/list',
        exact: true,
      },
      {
        name: 'log',
        path: '/task/log',
        component: './task/gitlab/log',
        exact: true,
      }
    ],
  },
  {
    path: '/admin',
    name: '用户管理',
    icon: 'user',
    routes: [
      {
        access: 'isAdmin',
        name: '用户列表',
        path: '/admin/user',
        component: './admin/UserList',
      },
    ],
  },
  {
    path: '/user',
    layout: false,
    exact: false,
    routes: [
      {
        path: '/user',
        routes: [
          {
            name: 'login',
            path: '/user/login',
            component: './user/Login',
          },
        ],
      },
    ],
  },
  {
    path: '/',
    component: './home',
  },
  {
    component: './404',
  },
];
