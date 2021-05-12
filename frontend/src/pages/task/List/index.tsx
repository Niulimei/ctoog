import React from 'react';
import moment from 'moment';
import { useHistory } from 'umi';
import { Task } from '@/typings/model';
import Table from '@ant-design/pro-table';
import { DownOutlined } from '@ant-design/icons';
import { task as taskService } from '@/services';
import TaskCreator from './components/TaskCreator';
import { Button, message, Dropdown, Menu } from 'antd';
import type { ProColumns } from '@ant-design/pro-table';

type Actions = Record<
  'startTask' | 'gotoDetail' | 'updateTask' | 'createTask',
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
      width: 120,
      align: 'center',
      // @ts-ignore
      render(item: Task.Item) {
        return (
          <>
            <Button size="small" type="link" onClick={() => actions.gotoDetail(item.id)}>
              详情
            </Button>
            <Dropdown
              overlay={
                <Menu>
                  <Menu.Item>
                    <Button size="small" type="link" onClick={() => actions.updateTask(item.id)}>
                      修改任务
                    </Button>
                  </Menu.Item>

                  {item.status !== Task.Status.RUNNING && (
                    <Menu.Item>
                      <Button size="small" type="link" onClick={() => actions.startTask(item.id)}>
                        启动任务
                      </Button>
                    </Menu.Item>
                  )}
                </Menu>
              }
            >
              <Button size="small" type="link">
                更多
                <DownOutlined />
              </Button>
            </Dropdown>
          </>
        );
      },
    },
  ];
};

const TaskList: React.FC = () => {
  const tableRef = React.useRef<any>(null);
  const creatorModalRef = React.useRef<any>(null);
  const [pageSize, setPageSize] = React.useState(10);
  const history = useHistory();

  const actions: Actions = {
    /** 查看任务详情 */
    async gotoDetail(id: string) {
      history.push(`/task/detail?id=${id}`);
    },
    /** 启动任务 */
    async startTask(id: string) {
      try {
        await taskService.startTask(id);
        message.success('迁移任务启动成功');
      } catch (err) {
        // message.error('迁移任务启动出现异常');
        // eslint-disable-next-line no-console
        console.error(err);
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
          pageSize,
          onChange(num, size) {
            if (size) {
              setPageSize(size);
            }
          },
        }}
        request={async (params) => {
          const { taskInfo, count } = await taskService.getTasks({
            offset: (params.current! - 1 || 0) * pageSize,
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
      <TaskCreator
        actionRef={creatorModalRef}
        onSuccess={() => {
          tableRef.current.reload();
        }}
      />
    </>
  );
};

export default TaskList;
