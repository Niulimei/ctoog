import React from 'react';
import { Task } from '@/typings/model';
import TaskCreator from '../TaskCreator';
import { useLocation, useHistory } from 'umi';
import TaskField from './components/TaskField';
import TaskLogger from './components/TaskLogger';
import { task as taskService } from '@/services';
import { Empty, Button, Modal, message } from 'antd';
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
  const history = useHistory();
  const location = useLocation<any>();
  const [taskDetail, setTaskDetail] = React.useState<Task.Detail>();
  const { id: taskId } = (location as any).query;
  const taskLoggerRef = React.useRef<any>();
  const taskCreatorRef = React.useRef<any>();
  const [isLoading, setisLoading] = React.useState(false);

  const actions = {
    /** 删除任务 */
    async deleteTask() {
      await taskService.deleteTask(taskId);
      message.success('删除任务成功！');
      history.push('/task/list');
    },
    /** 修改任务 */
    updateTask() {
      if (taskDetail?.taskModel.ccUser) {
        taskCreatorRef.current.openModal('update', taskId);
      } else {
        taskCreatorRef.current.openModal('planUpdate', taskId);
      }
    },
    /** 删除缓存 */
    async clearCache() {
      setisLoading(true);
      try {
        await taskService.deleteCache(taskId);
        message.success('删除缓存成功！');
      } catch (error) {
        // eslint-disable-next-line no-console
        console.error(error);
      } finally {
        setisLoading(false);
      }
    },
    /** 启动任务 */
    async startTask() {},
  };
  React.useEffect(() => {
    taskService.getTaskDetail(taskId).then((data) => {
      if (taskId) {
        if (!data.taskModel.ccUser) {
          Modal.warn({
            width: 480,
            title: '提示',
            afterClose: () => taskCreatorRef.current.openModal('planUpdate', taskId),
            content: '该迁移任务信息不完整，任务信息被补全后才能开始执行',
          });
        }
        setTaskDetail(data);
      }
    });
  }, [taskId]);

  return (
    <>
      <PageContainer
        content={<TaskField data={taskDetail?.taskModel} />}
        tabList={tabList}
        header={{
          title: '任务详情',
          breadcrumb,
        }}
        footer={[
          (taskDetail?.taskModel as any).status !== Task.Status.RUNNING ? (
            <Button key="startTask" onClick={actions.startTask} type="primary">
              启动任务
            </Button>
          ) : null,
          <Button key="updateTask" onClick={actions.updateTask}>
            修改任务
          </Button>,
          <Button key="clearCache" loading={isLoading} onClick={actions.clearCache}>
            删除缓存
          </Button>,
        ]}
      >
        <div style={{ padding: 15, background: '#fff' }}>
          {!taskDetail?.logList ? (
            <Empty description="暂无执行历史记录" />
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
      <TaskCreator
        onSuccess={() => window.location.reload()}
        key="TaskCreator"
        actionRef={taskCreatorRef}
      />
    </>
  );
};

export default TaskDetail;
