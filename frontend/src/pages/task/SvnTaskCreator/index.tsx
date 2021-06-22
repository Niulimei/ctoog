import React, {useCallback, useState} from 'react';
import { observer } from 'mobx-react';
import { useToggle } from 'react-use';
import {throttle, map, groupBy} from 'lodash';
import classnames from 'classnames';
import { guid } from '@/utils/utils';
import type { FormInstance } from 'antd/es/form';
import { task as taskService, svn as svnService } from '@/services';
import { Button, message, Form, Modal, Checkbox, Input, Empty, Tooltip, Divider } from 'antd';
import { StepsForm, ProFormText, ProFormTextArea } from '@ant-design/pro-form';

import styles from './style.less';

interface IModalCreatorProps {
  /** 创建成功回调 */
  onSuccess?: () => void;
  actionRef?: React.MutableRefObject<{
    openModal: (mode?: 'create' | 'update' | 'planUpdate', id?: string) => void;
  }>;
}

interface IFormFields {
  pvob: string;
  component: string;
  ccUser: string;
  ccPassword: string;
  gitURL: string;
  gitUser: string;
  gitPassword: string;
  matchInfo: { stream: string; gitBranch: string }[];
}

const renderCardTitle = (title: string) => {
  return <h3 style={{ textAlign: 'center', marginBottom: '20px' }}>{title}</h3>;
};

const renderSvnTitle = (title: string) => {
  return <h4 style={{ textAlign: 'left'}}>{title}</h4>;
};

/** pro compoent 禁用提交 */
const DisablePasswordFieldAutocompleteProps = {
  fieldProps: { autoComplete: 'new-password' },
};

const { TextArea } = Input;

/** 生成表单项 */
const formFieldsGenerator = (fields: any) => {
  const renderFieldComponent = ({ component, name, required, rules, ...restProps }: any) => {
    const requiredRules = [
      {
        required: true,
        message: restProps.placeholder ? restProps.placeholder : `${name} 为必填参数`,
      },
    ];
    // eslint-disable-next-line no-param-reassign
    rules = (rules || []).concat(required && requiredRules).filter(Boolean);

    return React.createElement(component, {
      key: name,
      required,
      rules,
      name,
      ...restProps,
    });
  };

  /** 判断当前的节点是不是带有require属性 */
  const isRequired = (node: any): boolean => {
    if(node?.props?.rules){
      if(Array.isArray(node.props.rules) && node.props.rules[0]){
        return node.props.rules[0].hasOwnProperty('required')  || node.props.rules[1].hasOwnProperty('required')
      }
    }
    return false

  }

  return fields.map((nodes: any) => {
    const key = guid();

    const [leftNode, rightNode, actionNode] = nodes.map((node: any) =>
      node.component ? renderFieldComponent(node) : node,
    );
    // console.log(leftNode);

    return (
      <div className={styles.col} key={key}>
        <div className={classnames(styles.row, styles.left, isRequired(leftNode) ? styles.must : '')}>{leftNode}</div>
        <div className={classnames(styles.row, styles.right, isRequired(rightNode) ? styles.must : '')}>{rightNode}</div>
        {actionNode ? <div className={styles.action}>{actionNode}</div> : null}
      </div>
    );
  });
};

const formTreeFieldsGenerator = (fields: any) => {
  const renderFieldComponent = ({ component, name, required, rules, ...restProps }: any) => {
    const requiredRules = [
      {
        required: true,
        message: restProps.placeholder ? restProps.placeholder : `${name} 为必填参数`,
      },
    ];
    // eslint-disable-next-line no-param-reassign
    rules = (rules || []).concat(required && requiredRules).filter(Boolean);

    return React.createElement(component, {
      key: name,
      required,
      rules,
      name,
      ...restProps,
    });
  };

  /** 判断当前的节点是不是带有require属性 */
  const isRequired = (node: any): boolean => {
    if(node?.props?.rules){
      if(Array.isArray(node.props.rules) && node.props?.rules[0]){
        return node.props?.rules[0]?.hasOwnProperty('required')  || node.props?.rules[1]?.hasOwnProperty('required')
      }
    }
    return false

  }

  return fields.map((nodes: any) => {
    const key = guid();

    const [leftNode, rightNode, ExNode] = nodes.map((node: any) =>
      node.component ? renderFieldComponent(node) : node,
    );

    return (
      <div className={styles.col} key={key}>
        <div className={classnames(styles.leftTitles)}>{leftNode}</div>
        <div className={classnames(styles.rightNodes)}>
          <div className={classnames(styles.left, isRequired(rightNode) ? styles.must : '')}>{rightNode}</div>
          {ExNode ? <div className={classnames(styles.right, isRequired(ExNode) ? styles.must : '')}>{ExNode}</div> : null}
        </div>
      </div>
    );
  });
};

const TaskCreator: React.FC<IModalCreatorProps> = (props) => {
  const { onSuccess, actionRef } = props;
  const [form] = Form.useForm<IFormFields>();
  const [formUser] = Form.useForm<IFormFields>();
  const [svnList, setSvnList] = useState([]);
  const [preSvnLIst, setPreSvnLIst] = useState([]);
  const [currentNum, setCurrentNum] = useState(0);
  const [visible, toggleVisible] = useToggle(false);
  const modalRef = React.useRef<{ taskId: string }>({ taskId: '' });

  // 两个脚本域的变量
  const [gitScript, setGitScript] = useState('function main(user) { return user }');
  const [gitEmailScript, setGitEmailScript] = useState('function main(user) { return user + "@example.com" }');

  /** 更新模式
   * 1. 回填表单数据
   * 2. pvob component matchInfo 为可修改配置，其他表单项只读
   */
  const [mode, setMode] = React.useState<'create' | 'update' | 'planUpdate'>();
  const isUpdateMode = (data = mode) => data && ['update', 'planUpdate'].includes(data);

  React.useImperativeHandle(actionRef, () => {
    return {
      async openModal(modalMode, id) {
        form.resetFields();
        setMode(modalMode);
        if (isUpdateMode(modalMode) && id) {
          modalRef.current.taskId = id;
          const { taskModel: fieldValues } = await taskService.getTaskDetail(id);
          form.setFieldsValue(fieldValues);
          if (fieldValues?.status === 'running') {
            toggleVisible(false);
            message.error('任务正在执行不可进行修改');
          } else {
            toggleVisible(true);
          }
        } else {
          toggleVisible(true);
        }
      },
    };
  });

  const actionText = isUpdateMode() ? '更新' : '新建';

  const finishHandler = async (values: any) => {
    try {
      const preFormValue = form.getFieldsValue();
      if (isUpdateMode()) {
        await taskService.updateTask(modalRef.current.taskId, {
          svnUrl: preFormValue?.svn_url,
          ...preFormValue,
          ...values
        });
      } else {
        await taskService.createTask({
          svnUrl: preFormValue?.svn_url,
          modelType: 'svn',
          ...preFormValue,
          ...values
        });
      }
      message.success(`迁移任务${actionText}成功`);
      onSuccess?.();
      return true;
    } catch (err) {
      // eslint-disable-next-line no-console
      console.error(err);
      return false;
    }
  };

  const getSvnUserList = async (values: any) => {
    try {
      const {ccPassword, svn_url, ccUser} = values;
      const list = await svnService.getSvn({
        svn_password: ccPassword,
        svn_url: svn_url,
        svn_user: ccUser,
      });
      setSvnList(list || []);
      const preNamesPair = form.getFieldValue('namePair');
      if (preNamesPair) {
        let nameGroupBy = groupBy(preNamesPair, 'svnUserName');
        let newNamePair = list.map(item => {
          const namePairItem = nameGroupBy[item];
          if (nameGroupBy[item]) {
            return namePairItem[0];
          } else {
            return {'svnUserName': item, 'gitUserName': item, 'gitEmail': `${item}@example.com`};
          }
        })
        formUser.setFieldsValue({'namePair': newNamePair});
        setPreSvnLIst(newNamePair);
      } else {
        formUser.setFieldsValue({'namePair': list.map(item => {return {'svnUserName': item, 'gitUserName': item, 'gitEmail': `${item}@example.com`}})})
        setPreSvnLIst(list.map(item => {return {'svnUserName': item, 'gitUserName': item, 'gitEmail': `${item}@example.com`}}));
      }
      setCurrentNum(1);
      return true;
    } catch (err) {
      console.error(err);
      setCurrentNum(0);
      return false;
    }
  };

  const closeModal = () => {
    toggleVisible(false);
    setCurrentNum(0);
  }

  const runScript = useCallback(
    type => {
      const { namePair } = formUser.getFieldsValue(['namePair']);
      try {
        if (type === 'git') {
          const gitOutPut = map(preSvnLIst, item => {
            return {
              gitEmail: item?.gitEmail,
              gitUserName: eval(gitScript + 'main("'+item?.gitUserName+'")'),
              svnUserName: item?.svnUserName
            }
          })
          formUser.setFieldsValue({'namePair': gitOutPut});
        } else {
          const gitEmailOutPut = map(namePair, item => {
            return {
              gitEmail: eval(gitEmailScript + 'main("'+item?.gitUserName+'")'),
              gitUserName: item?.gitUserName,
              svnUserName: item?.svnUserName
            }
          })
          formUser.setFieldsValue({'namePair': gitEmailOutPut});
        }
      } catch (e) {
        message.error(e);
      }
    }, [gitScript, gitEmailScript, preSvnLIst]
  );

  return (
    <StepsForm
      current={currentNum}
      containerStyle={{
        maxHeight: "calc(100vh - 190px)",
        overflowY: 'scroll',
        overflowX: 'hidden',
        width: 'max-content',
        padding: "0 20px"
      }}
      onFinish={closeModal}
      onCurrentChange={num => setCurrentNum(num)}
      title={`${actionText}迁移任务`}
      stepsFormRender={(dom, submitter) => {
          return (
            <Modal
              title={`${actionText}迁移任务`}
              width={850}
              centered
              wrapClassName={styles.svnModal}
              onCancel={closeModal}
              visible={visible}
              footer={submitter}
              destroyOnClose
            >
              {dom}
            </Modal>
          );
        }}
    >
      <StepsForm.StepForm
          name="base"
          form={form}
          style={{overflow: 'hidden'}}
          title="SVN基础设置"
          onFinish={getSvnUserList}
        >
      <div className={styles.gutter}>
        {formFieldsGenerator([
          [renderCardTitle('SVN'), renderCardTitle('Git')],
          [
            {
              name: 'svn_url',
              required: true,
              width: 'md',
              component: ProFormText,
              placeholder: '请输入 SVN 地址',
              showSearch: true,
            },
            {
              name: 'gitURL',
              required: true,
              width: 'md',
              component: ProFormText,
              placeholder: '请输入 Git Repo URL',
              rules: [
                {
                  async validator(_: any, value: string) {
                    if(value?.indexOf(' ') !== -1){
                      // value.replace(/\s+/g,"");
                      throw new Error('路径中不能出现空格')
                    }
                  },
                },
              ],
            },
          ],
          [
            {
              name: 'ccUser',
              required: true,
              width: 'md',
              component: ProFormText,
              placeholder: '请输入 SVN 用户名',
            },
            {
              name: 'gitUser',
              required: true,
              width: 'md',
              component: ProFormText,
              placeholder: '请输入 Git 用户名',
            },
          ],
          [
            {
              name: 'ccPassword',
              required: true,
              width: 'md',
              component: ProFormText.Password,
              placeholder: '请输入 SVN 密码',
              ...DisablePasswordFieldAutocompleteProps,
            },
            {
              name: 'gitPassword',
              required: true,
              width: 'md',
              component: ProFormText.Password,
              placeholder: '请输入 Git 密码',
              ...DisablePasswordFieldAutocompleteProps,
            },
          ],
        ])}
        <div className={classnames(styles.ignore)}>
          <ProFormTextArea
            name="gitignore"
            placeholder="请输入 gitignore信息"
          />
        </div>
         <div className={classnames(styles.col, styles.keep)}>
          <span>
            <Form.Item valuePropName="checked" noStyle name="includeEmpty">
              <Checkbox />
            </Form.Item>
            <span className={styles.label}>是否保留空目录</span>
          </span>
          <span className={styles.keep}>
            <span className={classnames(styles.label, styles.keepLabel)}>占位文件名</span>
            <Form.Item
              name={['keep']}
              rules={[
                {
                  async validator(_, value) {
                    const { includeEmpty } = form.getFieldsValue(['includeEmpty']);
                    if (includeEmpty && !value) {
                      throw new Error('文件名称不能为空');
                    }
                    const invalidated = /[^a-z0-9-_.]+/.test(value);
                    if (invalidated) {
                      throw new Error(
                        '文件名称字符只能包括：字母、数字、"."(点)、"_"(下划线)和"-"(连字符)',
                      );
                    }
                  },
                },
              ]}
            >
              <Input
                size="small"
                style={{ width: 128, marginLeft: 12 }}
                placeholder="请输入占位文件名"
              />
            </Form.Item>
          </span>
        </div>
      </div>
        </StepsForm.StepForm>
      <StepsForm.StepForm
          name="user"
          style={{overflow: 'hidden'}}
          form={formUser}
          title="用户设置"
          onFinish={finishHandler}
        >
         <div className={classnames([styles.gutter, styles.gutterEx])}>
           <div className={styles.svnTitle}>
             <p className={styles.firstNode}>SVN</p>
             {/* <Divider type="vertical" orientation='left' style={{height: '100%'}}/> */}
             <p className={styles.secondNode}>Git</p>
           </div>
           <div className={styles.svnTitle}>
             <div className={styles.firstNode}></div>
             <div className={styles.secondNodes}>
               <span className={styles.gitTitle}>Git用户名</span>
               <span className={styles.gitTitle}>Git邮箱</span>
             </div>
           </div>
            <div className={styles.svnTitle}>
             <div className={styles.firstNode}>
               <Tooltip title="脚本执行main函数">
                 {/* <span>js脚本</span> */}
               </Tooltip>
             </div>
             <div className={styles.secondNodes}>
                <div className={classnames([styles.innerRun, styles.innerRunLeft])}>
                  <ProFormTextArea fieldProps={{value: gitScript, onChange: e => throttle(() => setGitScript(e?.target?.value), 1000)()}} />
                  <Button type="primary" size="small" onClick={() => runScript('git')}>执行</Button>
                </div>
                <div className={styles.innerRun}>
                  <ProFormTextArea fieldProps={{value: gitEmailScript, onChange: e => setGitEmailScript(e?.target?.value)}} />
                  <Button type="primary" size="small" onClick={() => runScript('email')}>执行</Button>
                </div>
             </div>
           </div>
        {formTreeFieldsGenerator(
          svnList.map((item, index) => [
            {
              name: ['namePair', index, 'svnUserName'],
              component: ProFormText,
              width: 'ld',
              readonly: true,
            },
            {
              name: ['namePair', index, 'gitUserName'],
              component: ProFormText,
              placeholder: '请输入Git用户名',
              width: 'ld',
              required: true,
              rules: [
                {
                  // eslint-disable-next-line @typescript-eslint/no-shadow
                  async validator(_: any, value: string) {

                  },
                },
              ],
            },
            {
              name: ['namePair', index, 'gitEmail'],
              component: ProFormText,
              placeholder: '请输入Git 邮箱',
              width: 'ld',
              required: true,
              rules: [
                {
                  // eslint-disable-next-line @typescript-eslint/no-shadow
                  async validator(_: any, value: string) {

                  },
                },
              ],
            },
          ]),
        )}
           {
             svnList.length === 0 && (
               <Empty />
             )
           }
           </div>
      </StepsForm.StepForm>
    </StepsForm>
  );
};

export default observer(TaskCreator);
