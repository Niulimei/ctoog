import { LockOutlined, UserOutlined, EditTwoTone, ProfileTwoTone } from '@ant-design/icons';
import { message, Tabs, Form } from 'antd';
import React, { useState, useEffect } from 'react';
import ProForm, { ProFormText, ProFormSelect } from '@ant-design/pro-form';
import { Link, history, useModel } from 'umi';
import md5 from 'md5';

import { user as UserService } from '@/services';
import type { User } from '@/typings/model';

import styles from './index.less';

/** 此方法会跳转到 redirect 参数所在的位置 */
const goto = () => {
  if (!history) return;
  setTimeout(() => {
    const { query } = history.location;
    const { redirect } = query as { redirect: string };
    history.push(redirect || '/');
  }, 10);
};

const Login: React.FC = () => {
  const [submitting, setSubmitting] = useState(false);
  const [permission, setPermission] = useState([]);
  const [type, setType] = useState<string>('account');
  const { initialState, setInitialState } = useModel('@@initialState');

  const fetchUserInfo = async () => {
    const userInfo = await initialState?.fetchUserInfo?.();
    const routeInfo = await initialState?.fetchRouteInfo?.();
    if (userInfo) {
      setInitialState({
        ...initialState,
        currentUser: userInfo,
        RouteList: routeInfo
      });
    }
  };


  const GroupOptions = [
    '北京事业群',
    '厦门事业群',
    '成都事业群',
    '深圳事业群',
    '上海事业群',
    '广州事业群',
    '广研事业群',
    '武汉事业群',
    '基础技术中心',
    '实施管理中心',
    '大数据中心',
    '产品经营中心',
    '智能云事业部',
    '交付事业部',
  ];

  // const handleSubmit = async ({ password, username }: User.Base) => {


  //   setSubmitting(true);
  //   // 登录
  //   try {
  //     const msg = await UserService.login({
  //       username,
  //       password: md5(password),
  //     });
  //     if (msg.token) {
  //       message.success('登录成功！');
  //       await fetchUserInfo();
  //       goto();
  //     }
  //   } catch (err) {
  //     // eslint-disable-next-line no-console
  //     console.log(err);
  //   } finally {
  //     setSubmitting(false);
  //   }
  // };
  const handleSubmit = async (value: any) => {

    const { username, password, team, group, nickname, bussinessgroup } = value

    if (type === 'account') {
      setSubmitting(true);
      // 登录
      try {
        const msg = await UserService.login({
          username,
          password: md5(password),
        });
        if (msg.token) {
          message.success('登录成功！');
          await fetchUserInfo();
          goto();
        }
      } catch (err) {
        // eslint-disable-next-line no-console
        console.log(err);
      } finally {
        setSubmitting(false);
      }
    } else if (type === 'registery') {
      try {
        await UserService.registerUser({
          username, password: md5(password), team, group, nickname, bussinessgroup
        })
        message.success('注册成功！');
        history.replace('/')
      } catch (error) {
        // eslint-disable-next-line no-console
        console.log(error);
      } finally {
        setSubmitting(false);
      }


    }

  };

  const [form] = Form.useForm();

  useEffect(
    async () => {
      const RouteList = await UserService.getPermission();
      setPermission(RouteList);
    },
    []
  );

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <div className={styles.top}>
          <div className={styles.header}>
            <Link to="/">
              <span className={styles.title}>代码仓库迁移平台</span>
            </Link>
          </div>
          <div className={styles.desc} />
        </div>

        <div className={styles.main}>
          <ProForm
            initialValues={{
              autoLogin: true,
            }}
            submitter={{
              searchConfig: {
                submitText: type === 'account' ? '登录' : '注册',
              },
              render: (_, dom) => dom.pop(),
              submitButtonProps: {
                loading: submitting,
                size: 'large',
                style: {
                  width: '100%',
                },
              },
            }}
            onFinish={async (values) => {
              handleSubmit(values as User.Base);
            }}
            form={form}
          >
            <Tabs activeKey={type} onChange={setType}>
              <Tabs.TabPane key="account" tab="账户密码登录" />
              <Tabs.TabPane key="registery" tab="注册新用户" />
            </Tabs>
            {type === 'account' && (
              <>
                <ProFormText
                  name="username"
                  fieldProps={{
                    size: 'large',
                    prefix: <UserOutlined className={styles.prefixIcon} />,
                  }}
                  placeholder="请输入用户名或手机号"
                  rules={[
                    {
                      required: true,
                      message: '请输入用户名或手机号!',
                    },
                  ]}
                />
                <ProFormText.Password
                  name="password"
                  fieldProps={{
                    size: 'large',
                    prefix: <LockOutlined className={styles.prefixIcon} />,
                  }}
                  placeholder="请输入密码"
                  rules={[
                    {
                      required: true,
                      message: '请输入密码！',
                    },
                  ]}
                />
              </>
            )}
            {
              type === 'registery' && (
                <>
                  {/* 手机号 */}
                  <ProFormText
                    name="username"
                    fieldProps={{
                      size: 'large',
                      prefix: <UserOutlined className={styles.prefixIcon} />,
                    }}
                    placeholder="请输入手机号码"
                    rules={[
                      {
                        required: true,
                        message: '请输入手机号码!',
                      },
                      {
                        pattern: /^1\d{10}$/,
                        message: '不合法的手机号格式!',
                      },
                    ]}
                  />
                  {/* 名字 */}
                  <ProFormText
                    name="nickname"
                    fieldProps={{
                      size: 'large',
                      prefix: <EditTwoTone className={styles.prefixIcon} />,
                    }}
                    placeholder="请输入姓名"
                    rules={[
                      {
                        required: true,
                        message: '请输入姓名!',
                      },
                      {
                        pattern: /[\u4e00-\u9fa5a-z0-9]+$/i,
                        message: '名字不能含有特殊字符!',
                      },

                    ]}
                  />

                  {
                    permission.includes('jianxin') && (
                      <>
                        {/* 事业群 */}
                        <ProFormSelect
                          options={GroupOptions}
                          showSearch
                          placeholder="请选择事业群"
                          name="bussinessgroup"
                        />
                        {/* 项目组 */}
                        <ProFormText
                          name="team"
                          fieldProps={{
                            size: 'large',
                            prefix: <ProfileTwoTone className={styles.prefixIcon} />,
                          }}
                          placeholder="请输入项目组"
                          rules={[
                            {
                              required: true,
                              message: '请输入项目组!',
                            },

                          ]}
                        />
                      </>
                    )
                  }
                  {/* 密码 */}
                  <ProFormText.Password
                    name="password"
                    fieldProps={{
                      size: 'large',
                      prefix: <LockOutlined className={styles.prefixIcon} />,
                    }}
                    placeholder="请输入密码"
                    rules={[
                      {
                        required: true,
                        message: '请输入密码！',
                      },
                      {
                        pattern: /^[a-zA-Z]\w{5,17}$/,
                        message: '密码以字母开头，长度在6~18之间，只能包含字母、数字和下划线',
                      },
                    ]}
                  />
                  <ProFormText.Password
                    placeholder="请再次输入用户密码"
                    fieldProps={{
                      size: 'large',
                      prefix: <LockOutlined className={styles.prefixIcon} />,
                    }}
                    rules={[
                      {
                        async validator(_, value) {
                          if (value !== form.getFieldValue('password')) {
                            throw new Error('两次密码输入不一致');
                          }
                        },
                      },
                    ]}
                    name="retypePassword"
                  />
                </>
              )
            }
          </ProForm>
        </div>
      </div>
    </div>
  );
};

export default Login;
