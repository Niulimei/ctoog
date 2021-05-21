import React from 'react';
import dayjs from 'dayjs';
import { Form, message } from 'antd';
import type { Plan } from '@/typings/model';
import { useToggle, useUpdate } from 'react-use';
import { plan as planServices } from '@/services';
import { ModalForm, ProFormSelect, ProFormDatePicker } from '@ant-design/pro-form';

const StatusOptions = ['未迁移', '已迁移', '已切换'];

interface PlanStatusSwitcherProps {
  actionRef?: React.ForwardedRef<{ openModal: (data: Plan.Item) => void }>;
  onSuccess?: () => void;
}

const PlanStatusSwitcher: React.FC<PlanStatusSwitcherProps> = ({ actionRef, onSuccess }) => {
  const [visible, toggleVisible] = useToggle(false);
  const [plan, setPlan] = React.useState<Plan.Item>();
  const forceUpdate = useUpdate();
  const [form] = Form.useForm();

  React.useImperativeHandle(actionRef, () => ({
    openModal(data) {
      form.resetFields();
      if (data) {
        toggleVisible(true);
        setPlan(data);
      }
    },
  }));

  const isSwitchStatus = () => form.getFieldValue('status') === '已切换';

  const handleFinish = async ({ status, date }: any) => {
    const data: any = {
      status,
    };
    switch (status) {
      case '已迁移':
        data.actual_start_time = date;
        break;
      case '已切换':
        data.actual_switch_time = date;
        break;
      default:
        break;
    }
    await planServices.updatePlan(plan!.id, data);
    message.success('计划变更成功');
    toggleVisible(false);
    onSuccess?.();
  };

  return (
    <ModalForm
      form={form}
      width="400px"
      title="变更计划状态"
      visible={visible}
      layout="horizontal"
      onFinish={handleFinish}
      modalProps={{ okText: '更新' }}
      initialValues={{ date: dayjs().format('YYYY-MM-DD') }}
      onValuesChange={(_, values) => {
        if (values.status) forceUpdate();
      }}
      onVisibleChange={(vis) => toggleVisible(vis)}
    >
      <ProFormSelect
        name="status"
        label="变更状态为"
        placeholder="请选择变更状态"
        options={StatusOptions.filter((status) => status !== plan?.status)}
      />
      {['已迁移', '已切换'].includes(form.getFieldValue('status')) && (
        <ProFormDatePicker
          name="date"
          fieldProps={{ width: '300px' }}
          placeholder="请选择日期"
          label={`实际${isSwitchStatus() ? '切换' : '迁移'}日期`}
        />
      )}
    </ModalForm>
  );
};

export default PlanStatusSwitcher;
