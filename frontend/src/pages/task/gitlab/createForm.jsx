import { useEffect, useRef, useState } from 'react';
import { Button, Modal } from 'antd';
import { StepsForm, ProFormText, ProFormTextArea, ProFormSelect, ProFormCheckbox } from '@ant-design/pro-form';
import { gitlab as gitlabService } from '@/services';

export default ({ visible, setVisible }) => {
  const [migrationType, setMigrationType] = useState('Group');
  return (
    <Modal
      title="新建Gitlab迁移任务"
      visible={visible}
      onCancel={() => setVisible(false)}
      footer={null}
      width={650}
      destroyOnClose
    >
      <StepsForm
        onFinish={(value) => {
          const params = {
            ...value,
            modelType: "gitlab"
          };
          return gitlabService.createTask(params).then(() => {
            setVisible(false)
          })
        }}
      >
        <StepsForm.StepForm
          name="base"
          title="Gitlab基础设置"
          layout={"horizontal"}
          labelCol={{
            span: 8,
          }}
          wrapperCol={{
            span: 16,
          }}
        >
          <ProFormText
            name="sourceURL"
            width="md"
            label="原平台地址"
            tooltip="最长为 24 位，用于标定的唯一 id"
            placeholder="请输入地址"
            rules={[{ required: true }]}
          />
          <ProFormText
            name="gitlabToken"
            width="md"
            label="原平台令牌"
            tooltip="最长为 24 位，用于标定的唯一 id"
            placeholder="请输入令牌"
            rules={[{ required: true }]}
          />
        </StepsForm.StepForm>
        <StepsForm.StepForm
          name="move"
          title="Gitlab迁移数据设置"
          onFinish={async () => {
            return true;
          }}
          layout={"horizontal"}
          labelCol={{
            span: 8,
          }}
          wrapperCol={{
            span: 16,
          }}
        >
          <ProFormSelect
            name="type"
            label="迁移类型"
            valueEnum={{
              Group: 'Group',
              Project: 'Project',
            }}
            placeholder="请选择迁移类型"
            rules={[{ required: true, message: '请选择迁移类型' }]}
            onChange={(v) => setMigrationType(v)}
            initialValue={"Group"}
          />
          <ProFormText
            name={`gitlab${migrationType}`}
            label={migrationType + ' path'}
            placeholder="请输入path"
            rules={[{ required: true }]}
          />
          <ProFormCheckbox.Group
            name="checkbox"
            layout="vertical"
            label={'包含'+migrationType+'迁移数据'}
            options={['Issue', 'Milestone', 'Merge Request', 'Wiki', 'Permission']}
          />
        </StepsForm.StepForm>
        <StepsForm.StepForm
          name="gitee"
          title="Gitee基础设置"
          layout={"horizontal"}
          labelCol={{
            span: 8,
          }}
          wrapperCol={{
            span: 16,
          }}
        >
          <ProFormText
            name="targetURL"
            label="目标平台地址"
            rules={[{ required: true }]}
          />
          <ProFormText
            name="giteeToken"
            label="目标平台令牌"
          />
          <ProFormText
            name={`gitee${migrationType}`}
            label="放置在目标组"
            tooltip="为空则放在企业根目录下"
            rules={[{ required: true }]}
          />
        </StepsForm.StepForm>
      </StepsForm>
    </Modal>
  );
}
