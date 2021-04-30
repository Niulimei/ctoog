import React from 'react';
import moment from 'moment';
import { Button, message } from 'antd';
import Table from '@ant-design/pro-table';
import type { Task } from '@/typings/model';
import { task as taskService } from '@/services';
import ModalDetail from './components/ModalDetail';
import ModalCreator from './components/ModalCreator';
import type { ProColumns } from '@ant-design/pro-table';

type Actions = Record<'refreshTask' | 'displayDetail', (id: number) => void>;
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
      width: 100,
      // @ts-ignore
      render(item: Task.Item) {
        return (
          <>
            <Button size="small" type="link" onClick={() => actions.displayDetail(item.id)}>
              查看详情
            </Button>
            <Button size="small" type="link" onClick={() => actions.refreshTask(item.id)}>
              刷新任务
            </Button>
          </>
        );
      },
    },
  ];
};

const TaskList: React.FC = () => {
  const tableRef = React.useRef<any>(null);
  const detailModalRef = React.useRef<any>(null);
  const [taskDetail, setTaskDetail] = React.useState<Task.Detail>();

  const refreshTask = async (id: number) => {
    try {
      await taskService.refreshTask(id);
      message.success('迁移任务刷新成功');
    } catch (err) {
      message.error('迁移任务刷新出现异常');
    }
  };
  const displayDetail = async (id: number) => {
    const res = await taskService.getTaskDetail(id);
    setTaskDetail(res);
    detailModalRef.current.openModal();
  };

  const actions = { displayDetail, refreshTask };

  return (
    <>
      <Table
        rowKey="id"
        actionRef={tableRef}
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
          <ModalCreator
            onCreateSuccess={() => {
              tableRef.current.reload();
            }}
          />,
        ]}
        search={false}
      />
      <ModalDetail actionRef={detailModalRef} data={taskDetail} />
    </>
  );
};

export default TaskList;
