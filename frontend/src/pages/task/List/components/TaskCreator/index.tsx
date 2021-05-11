import React from 'react';
import { useMount } from 'react-use';
import { observer } from 'mobx-react';
import { useToggle } from 'react-use';
import classnames from 'classnames';
import { guid } from '@/utils/utils';
import type { FormInstance } from 'antd/es/form';
import { task as taskService } from '@/services';
import { Button, message, Form, Checkbox, Input } from 'antd';
import { MinusOutlined, PlusOutlined } from '@ant-design/icons';
import { renderCardTitle, useSelectOptions } from '../../helper';
import { ModalForm, ProFormSelect, ProFormText } from '@ant-design/pro-form';

import styles from './style.less';

interface IModalCreatorProps {
  /** 创建成功回调 */
  onSuccess?: () => void;
  actionRef?: React.RefObject<{
    openModal: (mode?: 'create' | 'update', id?: string) => void;
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

type CustomChangeHandlersType = Record<
  keyof IFormFields,
  (form: FormInstance<IFormFields>, value: any, dispatch: any) => void
>;
/** 表单变更处理 */
const CustomChangeHandlers: Partial<CustomChangeHandlersType> = {
  pvob(form, value, dispatch) {
    dispatch('component', { pvob: value });
    const { matchInfo } = form.getFieldsValue(['matchInfo']);
    form.setFieldsValue({
      component: undefined,
      matchInfo: (matchInfo || []).map((info: any) => ({ ...info, stream: '' })),
    });
  },
  component(form, value, dispatch) {
    dispatch('stream', { component: value, pvob: form.getFieldValue('pvob') });
    const { matchInfo } = form.getFieldsValue(['matchInfo']);
    form.setFieldsValue({
      matchInfo: (matchInfo || []).map((info: any) => ({ ...info, stream: '' })),
    });
  },
  ccUser(form, value) {
    if (value !== form.getFieldValue('gitUser')) {
      form.setFieldsValue({
        gitUser: value,
      });
    }
  },
  ccPassword(form, value) {
    if (value !== form.getFieldValue('gitPassword')) {
      form.setFieldsValue({
        gitPassword: value,
      });
    }
  },
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
    return React.createElement(component, { key: name, rules, name, ...restProps });
  };

  return fields.map((nodes: any) => {
    const key = guid();

    const [leftNode, rightNode, actionNode] = nodes.map((node: any) =>
      node.component ? renderFieldComponent(node) : node,
    );

    return (
      <div className={styles.col} key={key}>
        <div className={classnames(styles.row, styles.left)}>{leftNode}</div>
        <div className={classnames(styles.row, styles.right)}>{rightNode}</div>
        {actionNode ? <div className={styles.action}>{actionNode}</div> : null}
      </div>
    );
  });
};

const TaskCreator: React.FC<IModalCreatorProps> = (props) => {
  const { onSuccess, actionRef } = props;
  const [branchFieldNum, setBranchFieldNum] = React.useState(1);
  const [form] = Form.useForm<IFormFields>();
  const { dispatch: optionDispatch, options } = useSelectOptions();
  const [visible, toggleVisible] = useToggle(false);
  const modalRef = React.useRef<{ taskId: string }>({ taskId: '' });

  /** 更新模式
   * 1. 回填表单数据
   * 2. pvob component matchInfo 为可修改配置，其他表单项只读
   */
  const [isUpdateMode, setIsUpdateMode] = useToggle(false);

  React.useImperativeHandle(actionRef, () => {
    return {
      async openModal(mode, id) {
        form.resetFields();
        setIsUpdateMode(mode === 'update');
        if (mode === 'update' && id) {
          modalRef.current.taskId = id;
          const { taskModel: fieldValues } = await taskService.getTaskDetail(id);
          const { pvob, component, matchInfo } = fieldValues;
          optionDispatch('component', { pvob });
          optionDispatch('stream', { component, pvob });

          form.setFieldsValue(fieldValues);
          setBranchFieldNum(matchInfo.length);
          toggleVisible(true);
        } else {
          toggleVisible(true);
        }
      },
    };
  });

  useMount(async () => {
    optionDispatch('pvob', {});
  });

  const actionText = isUpdateMode ? '更新' : '新建';

  const finishHandler = async (values: any) => {
    try {
      if (isUpdateMode) {
        await taskService.updateTask(modalRef.current.taskId, values);
      } else {
        await taskService.createTask(values);
      }
      message.success(`迁移任务${actionText}成功`);
      onSuccess?.();
      return true;
    } catch (err) {
      message.error(`迁移任务${actionText}出现异常`);
      return false;
    }
  };

  const addBranchField = () => {
    setBranchFieldNum((num) => num + 1);
  };

  const deleteBranch = (pos: number) => {
    setBranchFieldNum((num) => num - 1);
    const { matchInfo } = form.getFieldsValue(['matchInfo']);
    if (Array.isArray(matchInfo)) {
      matchInfo.splice(pos, 1);
      form.setFieldsValue({
        matchInfo,
      });
    }
  };

  const onFormValuesChange = (values: Partial<IFormFields>) => {
    Object.entries(values).forEach(([key, val]) => {
      CustomChangeHandlers[key]?.(form, val, optionDispatch);
    });
  };

  return (
    <ModalForm
      form={form}
      width="800px"
      visible={visible}
      onFinish={finishHandler}
      title={`${actionText}迁移任务`}
      onValuesChange={onFormValuesChange}
      onVisibleChange={(vis) => toggleVisible(vis)}
      modalProps={{ okText: actionText, className: styles.modalForm }}
      initialValues={{ keep: '.gitkeep' }}
    >
      <div className={styles.gutter}>
        {formFieldsGenerator([
          [renderCardTitle('ClearCase'), renderCardTitle('Git')],
          [
            {
              name: 'pvob',
              required: true,
              component: ProFormSelect,
              placeholder: '请选择 PVOB',
              valueEnum: options.pvob,
              // props
              showSearch: true,
              // async request(data: any) {
              //   const { keyWords } = data;
              //   return transform(
              //     options.pvob,
              //     (result, val) => {
              //       if (!keyWords || new RegExp(keyWords, 'i').test(val)) {
              //         result.push({ label: val, value: val });
              //       }
              //     },
              //     [] as Record<'label' | 'value', string>[],
              //   );
              // },
            },
            {
              name: 'gitURL',
              required: true,
              component: ProFormText,
              placeholder: '请输入 Git Repo URL',
            },
          ],
          [
            {
              name: 'component',
              required: true,
              component: ProFormSelect,
              placeholder: '请选择组件',
              valueEnum: options.component,
            },
            {
              name: 'gitEmail',
              required: true,
              rules: [{ type: 'email', message: '请正确输入邮箱地址' }],
              component: ProFormText,
              placeholder: '请输入 Git Email，用于提交 Git 代码配置',
            },
          ],
          [
            {
              name: 'dir',
              component: ProFormText,
              placeholder: '请输入组件子目录，如果为空则迁移整个组件',
              valueEnum: options.component,
            },
          ],
          [
            {
              name: 'ccUser',
              required: true,
              component: ProFormText,
              placeholder: '请输入CC用户名',
            },
            {
              name: 'gitUser',
              required: true,
              component: ProFormText,
              placeholder: '请输入 Git 账号',
            },
          ],
          [
            {
              name: 'ccPassword',
              required: true,
              component: ProFormText.Password,
              placeholder: '请输入CC密码',
            },
            {
              name: 'gitPassword',
              required: true,
              component: ProFormText.Password,
              placeholder: '请输入 Git 密码',
            },
          ],
        ])}
        <div className={styles.divider} />
        {formFieldsGenerator(
          Array.from(Array(branchFieldNum), (_, index) => [
            {
              name: ['matchInfo', index, 'stream'],
              component: ProFormSelect,
              placeholder: '请选择开发流',
              valueEnum: options.stream,
              rules: [
                {
                  // eslint-disable-next-line @typescript-eslint/no-shadow
                  async validator(_: any, value: string) {
                    if (value) {
                      const matchInfo = form.getFieldValue('matchInfo');
                      if (Array.isArray(matchInfo)) {
                        const streams = matchInfo.map((item) => item.stream);
                        const len = streams.filter((stream) => stream === value).length;
                        if (len >= 2) {
                          throw new Error('不能出现重复的开发流');
                        }
                      }
                    }
                  },
                },
              ],
            },
            {
              name: ['matchInfo', index, 'gitBranch'],
              component: ProFormText,
              placeholder: '请输入Git对应分支',
              valueEnum: options.stream,
              rules: [
                {
                  // eslint-disable-next-line @typescript-eslint/no-shadow
                  async validator(_: any, value: string) {
                    if (value) {
                      const matchInfo = form.getFieldValue('matchInfo');
                      if (Array.isArray(matchInfo)) {
                        const gitBranchs = matchInfo.map((item) => item.gitBranch);
                        const len = gitBranchs.filter((branch) => branch === value).length;
                        if (len >= 2) {
                          throw new Error('不能出现重复的 Git 分支');
                        }
                      }
                    }
                  },
                },
              ],
            },
            <>
              {index === branchFieldNum - 1 ? (
                <Button key="add" type="primary" icon={<PlusOutlined />} onClick={addBranchField} />
              ) : (
                <Button key="delete" icon={<MinusOutlined />} onClick={() => deleteBranch(index)} />
              )}
            </>,
          ]),
        )}
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
                disabled={isUpdateMode}
                style={{ width: 128, marginLeft: 12 }}
                placeholder="请输入占位文件名"
              />
            </Form.Item>
          </span>
        </div>
      </div>
    </ModalForm>
  );
};

export default observer(TaskCreator);
