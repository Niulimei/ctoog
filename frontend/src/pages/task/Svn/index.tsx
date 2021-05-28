import React from 'react';
import moment from 'moment';
import { useHistory } from 'umi';
import { Task } from '@/typings/model';
import Table from '@ant-design/pro-table';
/** UploadOutlined */
import { DownOutlined } from '@ant-design/icons';
import { task as taskService } from '@/services';
import TaskCreator from '../SvnTaskCreator';
import { useCacheRequestParams } from '@/utils/hooks';
/** Upload */
import { Button, message, Dropdown, Menu } from 'antd';
import type { ProColumns } from '@ant-design/pro-table';

type Actions = Record<
  'startTask' | 'gotoDetail' | 'updateTask' | 'createTask',
  (id: number) => void
>;
const getColumns = (actions: Actions): ProColumns<Task.Item>[] => {
  return [
    {
      title: '任务编号',
      dataIndex: 'id',
      width: 80,
      hideInSearch: true,
    },
    {
      title: 'SVN 仓库',
      dataIndex: 'svnUrl',
      ellipsis: true,
      width: 120,
      // search:
    },
    {
      title: 'Git Repo',
      dataIndex: 'gitRepo',
      ellipsis: true,
      width: 180,
      hideInSearch: true,
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
      renderText({ lastCompleteDateTime }: Task.Item) {
        return lastCompleteDateTime ? moment(lastCompleteDateTime).format('MM-DD HH:mm') : '-';
      },
      hideInSearch: true,
    },
    {
      title: '操作',
      width: 120,
      align: 'center',
      hideInSearch: true,
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
  const history = useHistory();
  const { params, setParams } = useCacheRequestParams('taskList');

  const actions: Actions = {
    /** 查看任务详情 */
    async gotoDetail(id) {
      history.push(`/task/svnDetail?id=${id}`);
    },
    /** 启动任务 */
    async startTask(id) {
      try {
        await taskService.startTask(id);
        message.success('迁移任务启动成功');
      } catch (err) {
        // eslint-disable-next-line no-console
        console.error(err);
      }
    },
    /** 更新任务 */
    async updateTask(id) {
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
          pageSize: params.pageSize,
          onChange(num, size = 10) {
            setParams({
              current: num,
              pageSize: size,
            });
          },
        }}
        search={false}
        request={async ({ pageSize = 10, current }) => {
          const { taskInfo, count } = await taskService.getTasks({
            offset: (current! - 1 || 0) * pageSize,
            limit: pageSize || 10,
            modelType: 'svn'
          });

          return {
            data: taskInfo,
            success: true,
            total: count,
          };
        }}
        headerTitle="迁移任务"
        columns={getColumns(actions)}
        toolBarRender={() => [
        ]}
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