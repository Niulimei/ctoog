import React from 'react';
import { Modal, Descriptions } from 'antd';
import {useLocation} from 'umi';
import { useToggle } from 'react-use';
import { plan as planServices } from '@/services';
import styles from './index.less';

interface IProps {
  actionRef?: React.ForwardedRef<{
    open: (id: string, needAutoRefresh?: boolean) => void;
  }>;
}

const TaskWithPlanModal: React.FC<IProps> = ({ actionRef }) => {
  const [content, setContent] = React.useState({});
  const location = useLocation<any>();
  const { id: taskId } = (location as any).query;
  const [visible, toggleVisible] = useToggle(false);

  React.useImperativeHandle(actionRef, () => ({
    open: async () => {
      toggleVisible(true);
      const fieldValues = await planServices.getPlanDetail(taskId, {idType: 'task'});
      setContent(fieldValues);
    },
  }));

  const closeModal = () => {
    toggleVisible(false);
  };

  return (
   <Modal
      width="800px"
      title="计划详情"
      visible={visible}
      onOk={closeModal}
      onCancel={closeModal}
      cancelButtonProps={{ style: { display: 'none' } }}
    >
      <div>
        <Descriptions title="计划信息" bordered>
          <Descriptions.Item label="事业群" span="12">{content?.group}</Descriptions.Item>
          <Descriptions.Item label="项目组" span="12">{content?.team}</Descriptions.Item>
          <Descriptions.Item label="联系人" span="12">{content?.supporter}</Descriptions.Item>
          <Descriptions.Item label="联系电话" span="12">{content?.supporterTel}</Descriptions.Item>
        </Descriptions>
      </div>
    </Modal>
  );
};

export default TaskWithPlanModal;
