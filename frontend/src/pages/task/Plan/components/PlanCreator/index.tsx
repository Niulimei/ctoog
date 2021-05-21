import React from 'react';
import {
  ModalForm,
  ProFormText,
  ProFormRadio,
  ProFormSelect,
  ProFormTextArea,
  ProFormDatePicker,
} from '@ant-design/pro-form';
import dayjs from 'dayjs';
import { observer } from 'mobx-react';
import type { Plan } from '@/typings/model';
import { Row, Col, Form, message } from 'antd';
import type { FormInstance } from 'antd/es/form';
import { plan as planServices } from '@/services';
import { useClearCaseSelectEnum } from '@/utils/hooks';
import { useToggle, useMount, useUpdate } from 'react-use';

import styles from './style.less';

interface FormSectionProps {
  title?: string;
  left?: React.ReactNode;
  right?: React.ReactNode;
  wholeLine?: React.ReactNode;
}
const FormSection: React.FC<FormSectionProps> = ({ left, title, right, wholeLine }) => (
  <div className={styles.section}>
    {title && <h4 className={styles.sectionHeader}>{title}</h4>}
    <Row gutter={24}>
      {left && <Col span={11}>{left}</Col>}
      {right && (
        <Col span={11} offset={2}>
          {right}
        </Col>
      )}
    </Row>
    {wholeLine && <div>{wholeLine}</div>}
  </div>
);

const OriginTypeOptions = ['ClearCase', 'ICDP(Gerrit)', '私服'];

const TranslateTypeOptions = ['项目组自己迁移', '工作组帮迁移'];

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

const ProjectTypeOptions = ['Java', 'Python', 'C/C++', 'JavaScript', '其他'];

const InitialValues = {
  originType: 'ClearCase',
  translateType: '项目组自己迁移',
};
const CustomChangeHandlers: Partial<
  Record<keyof Plan.Base, (form: FormInstance<Plan.Base>, value: any, dispatch: any) => void>
> = {
  originType(form, value, dispatch) {
    dispatch({ type: 'forceUpdate' });
  },
  pvob(form, value, dispatch) {
    dispatch({ type: 'getComponentValueEnum', payload: value });
    form.setFieldsValue({ configLib: value });
  },
};

interface IPlanCreatorProps {
  actionRef?: React.ForwardedRef<{ openModal: (mode?: 'create' | 'update', id?: string) => void }>;
  onSuccess?: () => void;
}

const PlanCreator: React.FC<IPlanCreatorProps> = ({ actionRef, onSuccess }) => {
  const [form] = Form.useForm();
  const { dispatch: clearCaseEnumDispatch, valueEnum } = useClearCaseSelectEnum();
  const [visible, toggleVisible] = useToggle(false);
  const forceUpdate = useUpdate();
  const modalRef = React.useRef<{ planId: string }>({ planId: '' });

  /** 更新模式
   * 1. 回填表单数据
   * 2. pvob component matchInfo 为可修改配置，其他表单项只读
   */
  const [isUpdateMode, setIsUpdateMode] = useToggle(false);

  React.useImperativeHandle(actionRef, () => ({
    async openModal(mode = 'create', id) {
      form.resetFields();
      setIsUpdateMode(mode === 'update');
      if (mode === 'update' && id) {
        modalRef.current.planId = id;
        const fieldValues = await planServices.getPlanDetail(id);
        if (fieldValues.originType === 'ClearCase') {
          clearCaseEnumDispatch('pvob', {});
          clearCaseEnumDispatch('component', { pvob: fieldValues.pvob });
        }
        form.setFieldsValue(fieldValues);
        toggleVisible(true);
      } else {
        toggleVisible(true);
      }
    },
  }));

  const actionText = isUpdateMode ? '更新' : '新建';

  const handleFinish = async (values: Plan.Base) => {
    try {
      if (isUpdateMode) {
        await planServices.updatePlan(modalRef.current.planId, values);
      } else {
        await planServices.createPlan(values);
      }
      message.success(`迁移计划${actionText}成功`);
      onSuccess?.();
      return true;
    } catch (err) {
      // message.error(`迁移任务${actionText}出现异常`);
      // eslint-disable-next-line no-console
      console.error(err);
      return false;
    }
  };

  useMount(() => {
    clearCaseEnumDispatch('pvob', {});
  });

  const dispatch = ({ type, payload }: any) => {
    switch (type) {
      case 'forceUpdate':
        forceUpdate();
        break;
      case 'getComponentValueEnum':
        clearCaseEnumDispatch('component', { pvob: payload });
        break;
      default:
        break;
    }
  };

  const handleFormValuesChange = (changedValues: Partial<Plan.Base>) => {
    Object.entries(changedValues).forEach(([key, val]) => {
      CustomChangeHandlers[key]?.(form, val, dispatch);
    });
  };

  return (
    <ModalForm<Plan.Base>
      form={form}
      width="1000px"
      visible={visible}
      layout="horizontal"
      onFinish={handleFinish}
      initialValues={InitialValues}
      title={`${actionText}迁移计划`}
      modalProps={{ okText: actionText }}
      onValuesChange={handleFormValuesChange}
      onVisibleChange={(vis) => toggleVisible(vis)}
    >
      <FormSection
        left={
          <>
            <h6 className={styles.colTitle}>源仓库</h6>
            <ProFormRadio.Group
              name="originType"
              radioType="button"
              label="迁移任务类型"
              options={OriginTypeOptions}
              rules={[{ required: true, message: '请选择迁移任务类型' }]}
            />

            {form.getFieldValue('originType') === 'ClearCase' ? (
              <>
                <ProFormSelect
                  name="pvob"
                  label="PVOB"
                  placeholder="请选择 PVOB"
                  valueEnum={valueEnum.pvob}
                  rules={[{ required: true, message: '请选择 PVOB' }]}
                  showSearch
                />
                <ProFormSelect
                  name="component"
                  label="组件"
                  placeholder="请选择组件"
                  valueEnum={valueEnum.component}
                  rules={[{ required: true, message: '请选择组件' }]}
                  showSearch
                />
                <ProFormText
                  name="dir"
                  label="子目录"
                  placeholder="请输入组件子目录，如果为空则将迁移整个组件"
                />
              </>
            ) : (
              <ProFormText
                name="originUrl"
                label="仓库地址"
                placeholder="请输入仓库地址"
                rules={[{ required: true, message: '请输入仓库地址' }]}
              />
            )}
          </>
        }
        right={
          <>
            <h6 className={styles.colTitle}>目标仓库及计划</h6>
            <ProFormRadio.Group
              name="translateType"
              radioType="button"
              label="迁移方式"
              options={TranslateTypeOptions}
              rules={[{ required: true, message: '请选择迁移方式' }]}
            />
            <ProFormText name="targetUrl" label="目标仓库地址" placeholder="请填写目标仓库地址" />
            <ProFormDatePicker
              name="plan_start_time"
              label="计划迁移日期"
              placeholder="请选择计划迁移日期"
              rules={[{ required: true, message: '请选择计划迁移日期' }]}
            />
            <ProFormDatePicker
              name="plan_switch_time"
              label="计划切换日期"
              placeholder="请选择计划切换日期"
              rules={[
                { required: true, message: '请选择计划切换日期' },
                {
                  async validator(rule, val) {
                    const stratTime = form.getFieldValue('plan_start_time');
                    if (dayjs(stratTime).isAfter(dayjs(val)))
                      throw new Error('计划切换日期应在计划迁移日期之后');
                  },
                },
              ]}
            />
          </>
        }
      />
      <FormSection
        title="系统管理信息"
        left={
          <ProFormText
            name="subsystem"
            label="物理子系统英文简称"
            placeholder="请输入物理子系统英文简称"
            rules={[{ required: true, message: '请输入物理子系统英文简称' }]}
          />
        }
        right={
          <ProFormText
            name="configLib"
            label="配置库"
            placeholder="请输入配置库名称"
            rules={[{ required: true, message: '请输入配置库名称' }]}
          />
        }
      />
      <FormSection
        title="执行组织信息"
        left={
          <>
            <ProFormSelect
              name="group"
              label="事业群"
              placeholder="请选择事业群"
              options={GroupOptions}
              rules={[{ required: true, message: '请选择 PVOB' }]}
            />
            <ProFormText
              name="team"
              label="项目组"
              placeholder="请填写项目组名称"
              rules={[{ required: true, message: '请填写项目组名称' }]}
            />
          </>
        }
        right={
          <>
            <ProFormText
              name="supporter"
              label="对接人姓名"
              placeholder="请填写对接人姓名"
              rules={[{ required: true, message: '请填写对接人姓名' }]}
            />
            <ProFormText
              label="联系人电话"
              name="supporterTel"
              placeholder="请填写联系人电话"
              rules={[{ required: true, message: '请填写联系人电话' }]}
            />
          </>
        }
        wholeLine={<ProFormTextArea name="tip" label="备注" placeholder="请填写备注信息" />}
      />
      <FormSection
        title="源业务信息"
        left={
          <ProFormSelect
            label="工程类型"
            name="projectType"
            placeholder="请选择工程类型"
            options={ProjectTypeOptions}
          />
        }
        wholeLine={
          <>
            <ProFormTextArea name="purpose" label="业务用途" placeholder="请填写业务用途" />
            <ProFormTextArea
              name="effect"
              label="影响范围"
              placeholder="请填写本次迁移对业务、技术方面的影响范围大小"
            />
          </>
        }
      />
    </ModalForm>
  );
};

export default observer(PlanCreator);
