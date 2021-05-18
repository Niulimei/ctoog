import React from 'react';
import { Empty, Button } from 'antd';
import { useLocation } from 'umi';
import type { Task } from '@/typings/model';
import TaskField from './components/TaskField';
import TaskLogger from './components/TaskLogger';
import { task as taskService } from '@/services';
import TaskLogTable from './components/TaskLogTable';
import { PageContainer } from '@ant-design/pro-layout';
import { useHistory } from 'react-router-dom';
import { message } from 'antd';





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
  const history = useHistory();


  /** 删除缓存 */
  const deleteCache = async () => {
    await taskService.deleteCache(taskId)
    message.success('删除缓存成功！');
  }
  /** 删除任务 */
  const deleteTask = async () => {
    await taskService.deleteTask(taskId)
    message.success('删除任务成功！');
    history.push('/task/list')
  }
  /** 修改任务 */
  const amendTask = async () => {
    // actions.updateTask(taskId)
  }

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
      footer={[
        <Button onClick={deleteCache} >删除缓存</Button>,
        <Button onClick={amendTask} type="primary">修改任务</Button>,
        <Button onClick={deleteTask} danger type="primary">删除任务</Button>,
      ]}
    >
      <div style={{ padding: 15, background: '#fff' }}>
        {!taskDetail ? (
          <Empty />
        ) : (
          <>
            <TaskLogger actionRef={taskLoggerRef} />
            <TaskLogTable
              onDisplayLog={(task: Task.Log) =>
                taskLoggerRef.current.open(task.logID, task.status !== 'completed')
              }
              data={taskDetail?.logList}
            />
          </>
        )}
      </div>

    </PageContainer>
  );
};

export default TaskDetail;
