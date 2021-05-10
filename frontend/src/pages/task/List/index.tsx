import React from 'react';
import moment from 'moment';
import { Button, message } from 'antd';
import { Task } from '@/typings/model';
import Table from '@ant-design/pro-table';
import { task as taskService } from '@/services';
import TaskDetail from './components/TaskDetail';
import TaskCreator from './components/TaskCreator';
import type { ProColumns } from '@ant-design/pro-table';

type Actions = Record<
  'startTask' | 'displayDetail' | 'updateTask' | 'createTask',
  (id: string) => void
>;
const getColumns = (actions: Actions): ProColumns<Task.Item>[] => {
  return [
    {
      title: '任务编号',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: 'CC PVOB',
      dataIndex: 'pvob',
      ellipsis: true,
      width: 120,
    },
    {
      title: 'CC Component',
      dataIndex: 'component',
      ellipsis: true,
      width: 120,
    },
    {
      title: 'Git Repo',
      dataIndex: 'gitRepo',
      ellipsis: true,
      width: 180,
    },
    {
      title: '当前状态',
      width: 100,
      renderText(item: Task.Item) {
        return item.status;
      },
    },
    {
      title: '最后一次完成时间',
      ellipsis: true,
      width: 130,
      renderText(item: Task.Item) {
        return moment(item.lastCompleteDateTime).format('MM-DD HH:mm');
      },
    },
    {
      title: '操作',
      width: 80,
      align: 'center',
      // @ts-ignore
      render(item: Task.Item) {
        return (
          <>
            <Button size="small" type="link" onClick={() => actions.displayDetail(item.id)}>
              查看任务
            </Button>
            <Button size="small" type="link" onClick={() => actions.updateTask(item.id)}>
              修改任务
            </Button>
            {item.status !== Task.Status.RUNNING && (
              <Button size="small" type="link" onClick={() => actions.startTask(item.id)}>
                启动任务
              </Button>
            )}
          </>
        );
      },
    },
  ];
};

const TaskList: React.FC = () => {
  const tableRef = React.useRef<any>(null);
  const detailModalRef = React.useRef<any>(null);
  const creatorModalRef = React.useRef<any>(null);

  const [taskDetail, setTaskDetail] = React.useState<Task.Detail>();

  const actions: Actions = {
    /** 查看任务详情 */
    async displayDetail(id: string) {
      const res = await taskService.getTaskDetail(id);
      setTaskDetail(res);
      detailModalRef.current.openModal();
    },
    /** 启动任务 */
    async startTask(id: string) {
      try {
        await taskService.startTask(id);
        message.success('迁移任务启动成功');
      } catch (err) {
        message.error('迁移任务启动出现异常');
      }
    },
    /** 更新任务 */
    async updateTask(id: string) {
      creatorModalRef.current.openModal('update', id);
    },
    /** 创建任务 */
    async createTask() {
      creatorModalRef.current.openModal('create');
    },
  };

  return (
    <>
      <Table
        rowKey="id"
        actionRef={tableRef}
        pagination={{
          pageSize: 10,
        }}
        request={async (params) => {
          const { taskInfo, count } = await taskService.getTasks({
            offset: params.current! - 1 || 0,
            limit: params.pageSize || 10,
          });
          return {
            data: taskInfo,
            success: true,
            total: count,
          };
        }}
        headerTitle="任务列表"
        columns={getColumns(actions)}
        toolBarRender={() => [
          <Button
            size="small"
            type="primary"
            onClick={() => {
              creatorModalRef.current.openModal();
            }}
          >
            新建迁移任务
          </Button>,
        ]}
        search={false}
      />
      <TaskDetail actionRef={detailModalRef} data={taskDetail} />
      <TaskCreator
        actionRef={creatorModalRef}
        onSuccess={() => {
          tableRef.current.reload();
        }}
      />
      ,
    </>
  );
};

export default TaskList;
