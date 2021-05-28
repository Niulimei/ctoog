import React, {useState} from 'react';
import { observer } from 'mobx-react';
import { useToggle } from 'react-use';
import classnames from 'classnames';
import { guid } from '@/utils/utils';
import type { FormInstance } from 'antd/es/form';
import { task as taskService, svn as svnService } from '@/services';
import { Button, message, Form, Modal, Checkbox, Input, Empty } from 'antd';
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
    // console.log(leftNode);

    return (
      <div className={styles.col} key={key}>
        <div className={classnames(styles.leftTitle)}>{leftNode}</div>
        <div className={classnames(styles.row, styles.left, isRequired(rightNode) ? styles.must : '')}>{rightNode}</div>
        {ExNode ? <div className={classnames(styles.row, styles.right, isRequired(ExNode) ? styles.must : '')}>{ExNode}</div> : null}
      </div>
    );
  });
};

const TaskCreator: React.FC<IModalCreatorProps> = (props) => {
  const { onSuccess, actionRef } = props;
  const [form] = Form.useForm<IFormFields>();
  const [formUser] = Form.useForm<IFormFields>();
  const [svnList, setSvnList] = useState([]);
  const [currentNum, setCurrentNum] = useState(0);
  const [visible, toggleVisible] = useToggle(false);
  const modalRef = React.useRef<{ taskId: string }>({ taskId: '' });

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
        formUser.setFieldsValue({'namePair': preNamesPair});
      } else {
        formUser.setFieldsValue({'namePair': list.map(item => {return {'svnUserName': item, 'gitUserName': item, 'gitEmail': `${item}@example.com`}})})
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

  return (
    <StepsForm
      current={currentNum}
      onFinish={closeModal}
      onCurrentChange={num => setCurrentNum(num)}
      title={`${actionText}迁移任务`}
      stepsFormRender={(dom, submitter) => {
          return (
            <Modal
              title={`${actionText}迁移任务`}
              width={850}
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
          form={formUser}
          title="SVN用户"
          onFinish={finishHandler}
        >
         <div className={styles.gutter}>
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
