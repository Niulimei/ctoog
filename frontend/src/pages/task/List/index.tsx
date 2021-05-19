import React from 'react';
import moment from 'moment';
import { useHistory } from 'umi';
import { Task } from '@/typings/model';
import Table from '@ant-design/pro-table';
/** UploadOutlined */
import { DownOutlined } from '@ant-design/icons';
import { task as taskService } from '@/services';
import TaskCreator from '../TaskCreator';
import { useCacheRequestParams } from '@/utils/hooks';
/** Upload */
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
      hideInSearch: true,
    },
    {
      title: 'CC PVOB',
      dataIndex: 'pvob',
      ellipsis: true,
      width: 120,
      // search:
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
    async gotoDetail(id: string) {
      history.push(`/task/detail?id=${id}`);
    },
    /** 启动任务 */
    async startTask(id: string) {
      try {
        await taskService.startTask(id);
        message.success('迁移任务启动成功');
      } catch (err) {
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

  /** 批量上传的组件 */
  // const props = {
  //   name: 'file',
  //   action: 'https://www.mocky.io/v2/5cc8019d300000980a055e76',
  //   /** 上传后的钩子 */
  //   onChange(info: any) {
  //     if (info.file.status !== 'uploading') {
  //       console.log(info.file, info.fileList);
  //     }
  //     if (info.file.status === 'done') {
  //       message.success(`${info.file.name} 上传成功！`);
  //       // 调用接口
  //     } else if (info.file.status === 'error') {
  //       message.error(`${info.file.name} 上传失败！.`);
  //     }
  //   },
  // };

  /** 限制上传excel的判断字符串 */
  // const excelTypeStr = `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,application/vnd.ms-excel `

  /** 上传前的钩子函数 */
  // const validateExcel = (file: any) => {

  //   return new Promise<boolean>((resolve, reject) => {
  //     if (file.name.slice(-5).toLowerCase() === '.xlsx' || file.name.slice(-4).toLowerCase() === '.xls') {
  //       resolve(true)
  //     } else {
  //       reject()
  //       message.error('只能上传.XLS.XLSX格式的文件!')
  //     }
  //   })
  // }

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
          /** 批量上传 */
          // <Upload accept={excelTypeStr} beforeUpload={validateExcel} showUploadList={false} {...props}>
          //   <Button icon={<UploadOutlined />}>批量上传</Button>
          // </Upload>,
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
