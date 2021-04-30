import React from 'react';
import { Button, message, Form } from 'antd';
import { ModalForm, ProFormGroup, ProFormSelect, ProFormText } from '@ant-design/pro-form';
import ProCard from '@ant-design/pro-card';
import { MinusOutlined, PlusOutlined } from '@ant-design/icons';
import { guid } from '@/utils/utils';
import { renderCardTitle, useSelectOptions } from '../../helper';
import type { FormInstance } from 'antd/es/form';
import { observer } from 'mobx-react';
import { useMount } from 'react-use';
import { task as taskService } from '@/services';

import styles from './style.less';

interface IFormFields {
  pvob: string;
  component: string;
  ccUser: string;
  ccPassword: string;
  gitURL: string;
  gitUser: string;
  gitPassword: string;
  matchInfo: { stream: string; gitBranch: string };
}

/** 空行 */
const EmptyColSpace = <div style={{ marginBottom: 24, height: 32 }} />;
/** 右侧操作空白数 */
const RightButtonTopSpaces = 5;
/** 增减按钮默认样式 */
const ActionButtonDefaultStyles = {
  width: 31,
  height: 31,
  position: 'relative',
  top: -10,
  marginLeft: -30,
} as React.CSSProperties;

type CustomChangeHandlersType = Record<
  keyof IFormFields,
  (form: FormInstance<IFormFields>, value: any, dispatch: any) => void
>;
/** 页面联动交互处理 */
const CustomChangeHandlers: Partial<CustomChangeHandlersType> = {
  pvob(form, value, dispatch) {
    dispatch('component', { pvob: value });
  },
  component(form, value, dispatch) {
    dispatch('stream', { component: value, pvob: form.getFieldValue('pvob') });
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
    if (component === 'EmptyColSpace') {
      return React.cloneElement(EmptyColSpace, { key: guid() });
    }
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

const ModalCreator: React.FC<{ onCreateSuccess?: () => void }> = (props) => {
  const [branchFieldNum, setBranchFieldNum] = React.useState(1);
  const [form] = Form.useForm<IFormFields>();

  const { dispatch: optionDispatch, options } = useSelectOptions();

  useMount(() => {
    optionDispatch('pvob', {});
  });

  const finishHandler = async (values: any) => {
    try {
      await taskService.createTask(values);
      message.success('迁移任务新建成功');
      props.onCreateSuccess?.();
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
    // TODO:
    console.log(pos);
  };

  const onFormValuesChange = (values: Partial<IFormFields>) => {
    Object.entries(values).forEach(([key, val]) => {
      CustomChangeHandlers[key]?.(form, val, optionDispatch);
    });
  };

  return (
    <ModalForm
      form={form}
      width="700px"
      title="新建迁移任务"
      onValuesChange={onFormValuesChange}
      modalProps={{ okText: '新建', className: styles.modalForm }}
      trigger={
        <Button size="small" type="primary">
          新建迁移任务
        </Button>
      }
      onFinish={finishHandler}
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
            },
            {
              name: 'ccPassword',
              component: ProFormText.Password,
              placeholder: '请输入CC密码',
            },
          ])}
          <ProFormGroup>
            {renderElements(branchFieldNum, (index, uid) => (
              <ProFormSelect
                key={uid}
                width="md"
                placeholder="请选择开发流"
                valueEnum={options.stream}
                name={['matchInfo', index, 'stream']}
              />
            ))}
          </ProFormGroup>
        </ProCard>
        <ProCard colSpan="47%" style={{ border: 'none' }}>
          {renderCardTitle('Git')}
          {formFieldsGenerator([
            {
              name: 'gitURL',
              component: ProFormText,
              placeholder: '请输入 Git Repo URL',
              valueEnum: options.pvob,
            },
            {
              component: 'EmptyColSpace',
            },
            {
              name: 'gitUser',
              component: ProFormText,
              placeholder: '请输入 Git 账号',
              valueEnum: options.component,
            },
            {
              name: 'gitPassword',
              component: ProFormText.Password,
              placeholder: '请输入 Git 密码',
            },
          ])}
          <ProFormGroup>
            {renderElements(branchFieldNum, (index, uid) => (
              <ProFormText
                key={uid}
                width="md"
                placeholder="请输入Git对应分支"
                name={['matchInfo', index, 'gitBranch']}
              />
            ))}
          </ProFormGroup>
        </ProCard>
        <ProCard colSpan="6%">
          {renderElements(RightButtonTopSpaces, EmptyColSpace)}

          {renderElements(branchFieldNum - 1, (index, uid) => (
            <Button
              icon={<MinusOutlined />}
              key={uid}
              onClick={() => deleteBranch(index)}
              style={{
                ...ActionButtonDefaultStyles,
                marginBottom: 24,
              }}
            />
          ))}

          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={addBranchField}
            style={ActionButtonDefaultStyles}
          />
        </ProCard>
      </ProCard>
    </ModalForm>
  );
};

export default observer(ModalCreator);
