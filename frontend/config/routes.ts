export default [
  {
    path: '/task',
    name: '任务管理',
    icon: 'dashboard',
    access: 'isLogin',
    routes: [
      {
        name: '任务列表',
        icon: 'user',
        path: '/task/list',
        component: './task/List',
        exact: true,
      },
    ],
  },
  {
    path: '/admin',
    name: '用户管理',
    icon: 'table',
    routes: [
      {
        access: 'isAdmin',
        name: '用户列表',
        icon: 'admin',
        path: '/admin/user',
        component: './admin/UserList',
      },
    ],
  },
  {
    path: '/user',
    layout: false,
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
    redirect: '/task/list',
  },
  {
    component: './404',
  },
];
