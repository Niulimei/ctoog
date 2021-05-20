﻿export default [
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
      },
      {
        name: '迁移任务',
        path: '/task/list',
        component: './task/List',
        exact: true,
      },
      {
        hideInMenu: true,
        name: '任务详情',
        path: '/task/detail',
        component: './task/Detail',
      },
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
    redirect: '/task/plan',
  },
  {
    component: './404',
  },
];
