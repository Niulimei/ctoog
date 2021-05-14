import React from 'react';
import { Modal, Empty } from 'antd';
import Logger from '@/components/Logger';
import { useToggle, useUnmount } from 'react-use';
import { task as taskService } from '@/services';

interface IProps {
  actionRef?: React.ForwardedRef<{
    open: (id: string, needAutoRefresh?: boolean) => void;
  }>;
}

/** 日志刷新间隔 */
const RefreshLogInterval = 3000;

const TaskLogger: React.FC<IProps> = ({ actionRef }) => {
  const timerRef = React.useRef<any>();
  const [logData, setLogData] = React.useState('');
  const [visible, toggleVisible] = useToggle(false);
  const disposeTimeout = () => {
    clearTimeout(timerRef.current);
  };

  React.useImperativeHandle(actionRef, () => ({
    open: async (id, autoRefresh = false) => {
      const refreshLog = async () => {
        const { content } = await taskService.getLogOutput(id);
        setLogData(content);
        if (autoRefresh) {
          timerRef.current = setTimeout(async () => {
            disposeTimeout();
            await refreshLog();
          }, RefreshLogInterval);
        }
      };
      refreshLog();
      toggleVisible(true);
    },
  }));

  const closeModal = () => {
    setLogData('');
    toggleVisible(false);
    disposeTimeout();
  };

  useUnmount(() => {
    disposeTimeout();
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
      {logData ? <Logger value={logData} /> : <Empty description="暂无日志信息" />}
    </Modal>
  );
};

export default TaskLogger;
