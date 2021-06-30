import { useModel, Redirect } from 'umi';
// 根据权限和登陆情况判断用户需要跳转的对应页面
const Home = () => {
  const { initialState } = useModel('@@initialState');
  const { RouteList = [] } = initialState;
  if (RouteList.includes('jianxin')) {
    return <Redirect to='/task/plan' />
  }
  if (RouteList.includes('ccRoute')) {
    return <Redirect to='/task/List'/>
  }
  if (RouteList.includes('svnRoute')) {
    return <Redirect to='/task/svn'/>
  }
  return <Redirect to='/task/node'/>

}

export default Home;
