/**
 * @see https://umijs.org/zh-CN/plugins/plugin-access
 * */

import { User } from '@/typings/model';

export default function access(initialState: { currentUser?: User.Base | undefined }) {
  const { currentUser, RouteList = []  } = initialState || {};
  return {
    isLogin: currentUser,
    svnRoute: () => RouteList.includes('svnRoute'),
    jianxin: () => RouteList.includes('jianxin'),
    isAdmin: currentUser && currentUser.role_id === User.Role.ADMIN,
  };
}
