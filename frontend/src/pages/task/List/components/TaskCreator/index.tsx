import React from 'react';
import { useMount } from 'react-use';
import { guid } from '@/utils/utils';
import { observer } from 'mobx-react';
import { useToggle } from 'react-use';
import ProCard from '@ant-design/pro-card';
import type { FormInstance } from 'antd/es/form';
import { task as taskService } from '@/services';
import { Button, message, Form, Checkbox } from 'antd';
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

/** 空行 */
const EmptyColSpace = <div className={styles.emptyCol} />;
/** 右侧操作空白数 */
const RightButtonTopSpaceNum = 5;

type CustomChangeHandlersType = Record<
  keyof IFormFields,
  (form: FormInstance<IFormFields>, value: any, dispatch: any) => void
>;
/** 页面联动交互处理 */
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
const formFieldsGenerator = (fields: any[]) => {
  return fields.map(({ component, name, ...restProps }) => {
    const rules = [
      {
        required: true,
        message: restProps.placeholder ? restProps.placeholder : `${name} 为必填参数`,
      },
    ];
    return React.createElement(component, { key: name, rules, name, ...restProps });
  });
};

type ElementGetter = React.ReactElement | ((index: number, uid: string) => React.ReactElement);
/** 渲染指定个数元素 */
const renderElements = (num: number, nodeGetter: ElementGetter) => {
  return Array.from(Array(num), (_, index) => {
    if (typeof nodeGetter === 'function') return nodeGetter(index, guid());
    return React.cloneElement(nodeGetter, { key: guid() });
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
      message.success('迁移任务新建成功');
      onSuccess?.();
      return true;
    } catch (err) {
      message.error('迁移任务新建出现异常');
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
      title={`${actionText}迁移任务`}
      onFinish={finishHandler}
      onValuesChange={onFormValuesChange}
      onVisibleChange={(vis) => toggleVisible(vis)}
      modalProps={{ okText: actionText, className: styles.modalForm }}
    >
      <ProCard split="vertical" ghost>
        <ProCard colSpan="47%">
          {renderCardTitle('ClearCase')}
          {formFieldsGenerator([
            {
              name: 'pvob',
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
              name: 'component',
              component: ProFormSelect,
              placeholder: '请选择组件',
              valueEnum: options.component,
            },
            {
              name: 'ccUser',
              component: ProFormText,
              placeholder: '请输入CC用户名',
              readonly: isUpdateMode,
            },
            {
              name: 'ccPassword',
              component: ProFormText.Password,
              placeholder: '请输入CC密码',
              readonly: isUpdateMode,
            },
          ])}
          <div className={styles.dynamicFields}>
            {renderElements(branchFieldNum, (index, uid) => (
              <ProFormSelect
                key={uid}
                width="md"
                placeholder="请选择开发流"
                valueEnum={options.stream}
                name={['matchInfo', index, 'stream']}
              />
            ))}
          </div>
          <Form.Item valuePropName="checked" noStyle name="includeEmpty">
            <Checkbox />
          </Form.Item>
          <span className={styles.checkboxLabel}>是否保留空目录</span>
        </ProCard>
        <ProCard colSpan="47%" style={{ border: 'none' }}>
          {renderCardTitle('Git')}
          {formFieldsGenerator([
            {
              name: 'gitURL',
              component: ProFormText,
              placeholder: '请输入 Git Repo URL',
              readonly: isUpdateMode,
            },
            {
              name: 'gitEmail',
              component: ProFormText,
              placeholder: '请输入 Git Email，用于提交 Git 代码配置',
              readonly: isUpdateMode,
            },
            {
              name: 'gitUser',
              component: ProFormText,
              placeholder: '请输入 Git 账号',
              readonly: isUpdateMode,
            },
            {
              name: 'gitPassword',
              component: ProFormText.Password,
              placeholder: '请输入 Git 密码',
              readonly: isUpdateMode,
            },
          ])}
          <div className={styles.dynamicFields}>
            {renderElements(branchFieldNum, (index, uid) => (
              <ProFormText
                key={uid}
                width="md"
                placeholder="请输入Git对应分支"
                name={['matchInfo', index, 'gitBranch']}
              />
            ))}
          </div>
        </ProCard>
        <ProCard colSpan="6%">
          {renderElements(RightButtonTopSpaceNum, EmptyColSpace)}

          <div className={styles.dynamicFields}>
            {renderElements(branchFieldNum - 1, (index, uid) => (
              <Button
                icon={<MinusOutlined />}
                key={uid}
                onClick={() => deleteBranch(index)}
                className={styles.actionButton}
              />
            ))}

            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={addBranchField}
              className={styles.actionButton}
            />
          </div>
        </ProCard>
      </ProCard>
    </ModalForm>
  );
};

export default observer(TaskCreator);
