import React, {useState, useCallback, useMemo} from 'react';
import moment from 'moment';
import { useHistory, useModel } from 'umi';
import {values, uniq, flatten} from 'lodash';
import { Task } from '@/typings/model';
import ProTable from '@ant-design/pro-table';
/** UploadOutlined */
import { DownOutlined } from '@ant-design/icons';
import { task as taskService } from '@/services';
import TaskCreator from '../TaskCreator';
import { useCacheRequestParams } from '@/utils/hooks';
/** Upload */
import { Button, message, Table, Space, Dropdown, Menu, Tooltip } from 'antd';
import type { ProColumns } from '@ant-design/pro-table';
import classnames from "classnames";
import styles from './style.less';

const StatusOptions = [
  "init",
  "failed",
  "completed",
  "running",
];

type Actions = Record<
  'startTask' | 'gotoDetail' | 'updateTask' | 'createTask',
  (id: number) => void
>;
const getColumns = (actions: Actions, isJianxin): ProColumns<Task.Item>[] => {
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
      hideInSearch: true,
      width: 120,
      // search:
    },
    {
      title: 'CC Component',
      dataIndex: 'component',
      hideInSearch: true,
      width: 120,
      render(component: Task.Item) {
        return (
          <Tooltip placement="topLeft" title={component}>
            <span className={classnames(styles.exlipis)}>{component ? component : '-'}</span>
          </Tooltip>
        )
      }
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
      dataIndex: 'status',
      hideInSearch: false,
      valueEnum: StatusOptions,
      renderText(status: Task.Item) {
        return status;
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
            {
              !isJianxin && (
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
                  <DownOutlined/>
                </Button>
              </Dropdown>)
            }
            {item.status !== Task.Status.RUNNING && isJianxin && (
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
  const { initialState } = useModel('@@initialState');
  const [selectPageRow, setPageSelectRow] = useState({});
  const tableRef = React.useRef<any>(null);
  const creatorModalRef = React.useRef<any>(null);
  const history = useHistory();
  const { params, setParams } = useCacheRequestParams('taskList');

  const { RouteList = [] } = initialState;

  const actions: Actions = {
    /** 查看任务详情 */
    async gotoDetail(id) {
      history.push(`/task/detail?id=${id}`);
    },
    /** 启动任务 */
    async startTask(id) {
      try {
        await taskService.startTask(id);
        message.success('迁移任务启动成功');
        tableRef?.current?.reload();
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

  const setCurrentPageChoosen = useCallback((ids) => {
      const currentPage = params.current;
      setPageSelectRow({
        ...selectPageRow,
        [currentPage]: ids
      });
    }, [params, selectPageRow]);

  const runAll = async () => {
      if (selectedRowKeys.length === 0) {
        message.error('请选择任务');
      } else {
        await actions.startTask(selectedRowKeys);
        setPageSelectRow({});
      }
  }

  const selectedRowKeys = useMemo(() => {
    return uniq(flatten(values(selectPageRow))) || [];
  }, [selectPageRow, params]);

  return (
    <>
      <ProTable
        rowKey="id"
        actionRef={tableRef}
        rowSelection={{
          // 自定义选择项参考: https://ant.design/components/table-cn/#components-table-demo-row-selection-custom
          // 注释该行则默认不显示下拉选项
          selections: [Table.SELECTION_INVERT],
          onChange: setCurrentPageChoosen,
          selectedRowKeys
        }}
        search={{
          defaultCollapsed: false,
          span: 6
        }}
        pagination={{
          pageSize: params.pageSize,
          onChange(num, size = 10) {
            setParams({
              current: num,
              pageSize: size,
            });
          },
        }}
        request={async ({ pageSize = 10, current, status }) => {
          const { taskInfo, count } = await taskService.getTasks({
            offset: (current! - 1 || 0) * pageSize,
            limit: pageSize || 10,
            status: StatusOptions[status],
          });

          return {
            data: taskInfo,
            success: true,
            total: count,
            status: StatusOptions[status],
          };
        }}
        headerTitle="迁移任务"
        columns={getColumns(actions, RouteList.includes('jianxin'))}
        toolBarRender={() => [
          <Button
            size="small"
            type="primary"
            onClick={() => runAll()}
          >批量执行</Button>,
          !RouteList.includes('jianxin') && <Button
            size="small"
            type="primary"
            onClick={() => actions.createTask()}
          >新建CC迁移任务</Button>
        ]}
        tableAlertOptionRender={() => {
        return (
          <Space size={16}>
            <a onClick={() => setPageSelectRow({})}>取消选择</a>
          </Space>
        );
      }}
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
