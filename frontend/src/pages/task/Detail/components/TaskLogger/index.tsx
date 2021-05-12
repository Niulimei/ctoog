import React from 'react';
import { Modal } from 'antd';
import Logger from '@/components/Logger';
import { useToggle, useUnmount } from 'react-use';
import { task as taskService } from '@/services';

interface IProps {
  actionRef?: React.ForwardedRef<{
    open: (id: string) => void;
  }>;
}

/** 日志刷新间隔 */
const RefreshLogInterval = 1500;

const TaskLogger: React.FC<IProps> = ({ actionRef }) => {
  const timerRef = React.useRef<any>();
  const [logData, setLogData] = React.useState('');
  const [visible, toggleVisible] = useToggle(false);

  React.useImperativeHandle(actionRef, () => ({
    open: async (id: string) => {
      const setLogOutput = async () => {
        const { content } = await taskService.getLogOutput(id);
        setLogData(content);
      };

      await setLogOutput();
      toggleVisible(true);
      timerRef.current = setInterval(() => {
        setLogOutput();
      }, RefreshLogInterval);
    },
  }));

  const closeModal = () => {
    setLogData('');
    toggleVisible(false);
    clearInterval(timerRef.current);
  };

  useUnmount(() => {
    clearInterval(timerRef.current);
  });

  return (
    <Modal
      width="800px"
      title="任务执行日志"
      visible={visible}
      onOk={closeModal}
      onCancel={closeModal}
      cancelButtonProps={{ style: { display: 'none' } }}
    >
      <Logger value={logData} />
    </Modal>
  );
};

export default TaskLogger;
