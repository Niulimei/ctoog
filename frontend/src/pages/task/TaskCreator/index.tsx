import React from 'react';
import { useMount } from 'react-use';
import {useModel} from 'umi';
import {toJS} from 'mobx';
import { observer } from 'mobx-react';
import { useToggle } from 'react-use';
import classnames from 'classnames';
import { guid } from '@/utils/utils';
import type { FormInstance } from 'antd/es/form';
import { task as taskService } from '@/services';
import { useClearCaseSelectEnum } from '@/utils/hooks';
import { Button, message, Form, Checkbox, Input, AutoComplete } from 'antd';
import { MinusOutlined, PlusOutlined } from '@ant-design/icons';
import { ModalForm, ProFormSelect as FormSelect, ProFormText, ProFormTextArea } from '@ant-design/pro-form';

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

// 绑定 getPopupContainer
const ProFormSelect = ({ children, ...props }: React.ComponentProps<any>) => {
  return (
    <div key={props.name} data-field-key={props.name}>
      <FormSelect
        {...props}
        fieldProps={{
          getPopupContainer: () => document.querySelector(`[data-field-key="${props.name}"]`)!,
        }}
      >
        {children}
      </FormSelect>
    </div>
  );
};

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

const FieldAutoComplete = props => {
  const streamObj = toJS(props.options);
  let arr = Object.keys(streamObj).map(item => {
    return {label: item, value: item}
  });

  return (
     <Form.Item {...props}>
       <AutoComplete options={arr || []} />
     </Form.Item>
  )
}

const TaskCreator: React.FC<IModalCreatorProps> = (props) => {
  const { initialState } = useModel('@@initialState');
  const { onSuccess, actionRef } = props;
  const [branchFieldNum, setBranchFieldNum] = React.useState(1);
  const [form] = Form.useForm<IFormFields>();
  const { dispatch: optionDispatch, valueEnum } = useClearCaseSelectEnum();
  const [visible, toggleVisible] = useToggle(false);
  const modalRef = React.useRef<{ taskId: string }>({ taskId: '' });

  const { RouteList = [] } = initialState;

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
          const { pvob, component, matchInfo } = fieldValues;
          optionDispatch('component', { pvob });
          optionDispatch('stream', { component, pvob });

          form.setFieldsValue(fieldValues);
          setBranchFieldNum(matchInfo ? matchInfo.length : 1);
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

  const actionText = isUpdateMode() ? '更新' : '新建';

  const finishHandler = async (values: any) => {
    try {
      if (isUpdateMode()) {
        await taskService.updateTask(modalRef.current.taskId, values);
      } else {
        await taskService.createTask({...values, modelType: 'ClearCase'});
      }
      message.success(`迁移任务${actionText}成功`);
      onSuccess?.();
      return true;
    } catch (err) {
      // message.error(`迁移任务${actionText}出现异常`);
      // eslint-disable-next-line no-console
      console.error(err);
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

  /** matchInfo 字段不重复 */
  const isDuplicateMatchInfoItem = (key: 'stream' | 'gitBranch', inputVal?: string) => {
    if (!inputVal || !inputVal.trim()) return false;
    const matchInfo = form.getFieldValue('matchInfo');
    if (Array.isArray(matchInfo)) {
      const values = matchInfo.map((item) => item[key]);
      const len = values.filter((val) => val === inputVal).length;
      if (len >= 2) return true;
    }
    return false;
  };


  return (
    <ModalForm
      form={form}
      width="850px"
      visible={visible}
      onFinish={finishHandler}
      title={`${actionText}迁移任务`}
      onValuesChange={onFormValuesChange}
      onVisibleChange={(vis) => toggleVisible(vis)}
      modalProps={{ okText: actionText, className: styles.modalForm, centered: true }}
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
              valueEnum: valueEnum.pvob,
              // props
              showSearch: true,
            },
            {
              name: 'gitURL',
              required: true,
              component: ProFormText,
              placeholder: '请输入 Git Repo URL',
              rules: [
                {
                  async validator(_: any, value: string) {
                    if(value.indexOf(' ') !== -1){
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
              name: 'component',
              required: true,
              showSearch: true,
              component: ProFormSelect,
              placeholder: '请选择组件',
              valueEnum: valueEnum.component,
            },
            {
              name: 'gitEmail',
              required: true,
              rules: [
                {
                  pattern: RouteList.includes('jianxin') ? /^[^@]+@ccbft.com$/ : /^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+/,
                  message: RouteList.includes('jianxin') ? '请输入 Git Email，格式：云桌面账号@ccbft.com' : '请输入 Git Email',
                },
              ],
              component: ProFormText,
              placeholder: RouteList.includes('jianxin') ? '请输入 Git Email，格式：云桌面账号@ccbft.com' : '请输入 Git Email',
            },
          ],
          [
            {
              name: 'dir',
              component: ProFormText,
              placeholder: '请输入组件子目录，如果为空则迁移整个组件',
              valueEnum: valueEnum.component,
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
              ...DisablePasswordFieldAutocompleteProps,
            },
            {
              name: 'gitPassword',
              required: true,
              component: ProFormText.Password,
              placeholder: '请输入 Git 密码',
              ...DisablePasswordFieldAutocompleteProps,
            },
          ],
        ])}
        <div className={styles.divider} />
        {formFieldsGenerator(
          Array.from(Array(branchFieldNum), (_, index) => [
            {
              name: ['matchInfo', index, 'stream'],
              component: FieldAutoComplete,
              placeholder: '请选择开发流',
              valueEnum: valueEnum.stream,
              options: valueEnum.stream,
              showSearch: true,
              required: true,
              rules: [
                {
                  // eslint-disable-next-line @typescript-eslint/no-shadow
                  async validator(_: any, value: string) {
                    if (isDuplicateMatchInfoItem('stream', value))
                      throw new Error('不能出现重复的开发流');
                  },
                },
              ],
            },
            {
              name: ['matchInfo', index, 'gitBranch'],
              component: ProFormText,
              placeholder: '请输入Git对应分支',
              required: true,
              valueEnum: valueEnum.stream,
              rules: [
                {
                  // eslint-disable-next-line @typescript-eslint/no-shadow
                  async validator(_: any, value: string) {
                    if (isDuplicateMatchInfoItem('gitBranch', value))
                      throw new Error('不能出现重复的 Git 分支');
                  },
                },
              ],
            },
            <>
              {branchFieldNum !== 1 && (
                <Button
                  key="delete"
                  icon={<MinusOutlined />}
                  className={styles.deleteButton}
                  onClick={() => deleteBranch(index)}
                />
              )}
              {index === branchFieldNum - 1 && (
                <Button key="add" type="primary" icon={<PlusOutlined />} onClick={addBranchField} />
              )}
            </>,
          ]),
        )}
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
                disabled={false}
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
