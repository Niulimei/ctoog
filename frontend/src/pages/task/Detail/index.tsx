import React from 'react';
import { Empty } from 'antd';
import { useLocation } from 'umi';
import type { Task } from '@/typings/model';
import TaskField from './components/TaskField';
import TaskLogger from './components/TaskLogger';
import { task as taskService } from '@/services';
import TaskLogTable from './components/TaskLogTable';
import { PageContainer } from '@ant-design/pro-layout';

/** breadcrumb 配置 */
const breadcrumb = {
  routes: [
    {
      breadcrumbName: '任务列表',
      path: '/task/list',
    },
    {
      breadcrumbName: '任务详情',
      path: '/task/detail',
    },
  ],
};

/** tablist */
const tabList = [
  {
    tab: '执行历史记录',
    key: 'log',
    closable: false,
  },
];

const TaskDetail = () => {
  const location = useLocation();
  const [taskDetail, setTaskDetail] = React.useState<Task.Detail>();
  const { id: taskId } = (location as any).query;
  const taskLoggerRef = React.useRef<any>();

  React.useEffect(() => {
    taskService.getTaskDetail(taskId).then((data) => {
      if (taskId) {
        setTaskDetail(data);
      }
    });
  }, [taskId]);

  return (
    <PageContainer
      content={<TaskField data={taskDetail?.taskModel} />}
      tabList={tabList}
      header={{
        title: '任务详情',
        breadcrumb,
      }}
    >
      <div style={{ padding: 15, background: '#fff' }}>
        {!taskDetail ? (
          <Empty />
        ) : (
          <>
            <TaskLogger actionRef={taskLoggerRef} />
            <TaskLogTable
              onDisplayLog={(id) => taskLoggerRef.current.open(id)}
              data={taskDetail?.logList}
            />
          </>
        )}
      </div>
    </PageContainer>
  );
};

export default TaskDetail;
