import React from 'react';
import { useHistory } from 'umi';
import type { Plan } from '@/typings/model';
import ProTable from '@ant-design/pro-table';
import { DownOutlined } from '@ant-design/icons';
import { plan as planServices } from '@/services';
import PlanCreator from './components/PlanCreator';
import { Button, Menu, Dropdown, message } from 'antd';
import type { ProColumns } from '@ant-design/pro-table';
import PlanStatusSwitcher from './components/PlanStatusSwitcher';

type Actions = Record<
  'updatePlan' | 'deletePlan' | 'execTask' | 'toggleStatus' | 'gotoTaskDetail',
  (payload: Plan.Item) => void
>;

const getColumns = (actions: Actions): ProColumns<Plan.Item>[] => {
  return [
    {
      title: '编号',
      dataIndex: 'id',
      ellipsis: true,
      hideInSearch: true,
      width: 60,
    },
    {
      title: '源仓库类型',
      dataIndex: 'originType',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '迁移状态',
      dataIndex: 'status',
      ellipsis: true,
    },
    {
      title: '目标仓库地址',
      dataIndex: 'targetUrl',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '物理子系统',
      dataIndex: 'subsystem',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '配置库',
      dataIndex: 'configLib',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '事业群',
      dataIndex: 'group',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '项目组',
      dataIndex: 'team',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '计划迁移时间',
      dataIndex: 'plan_start_time',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '实际迁移时间',
      dataIndex: 'actual_start_time',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '计划切换时间',
      dataIndex: 'plan_switch_time',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '实际切换时间',
      dataIndex: 'actual_switch_time',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '操作',
      width: 140,
      hideInSearch: true,
      fixed: 'right',
      // @ts-ignore
      render(item: Plan.Item) {
        return (
          <>
            <Button size="small" type="link" onClick={() => actions.updatePlan(item)}>
              修改
            </Button>
            <Dropdown
              overlay={
                <Menu>
                  <Menu.Item>
                    <Button size="small" type="link" onClick={() => actions.deletePlan(item)}>
                      删除
                    </Button>
                  </Menu.Item>
                  <Menu.Item>
                    <Button size="small" type="link" onClick={() => actions.toggleStatus(item)}>
                      变更状态
                    </Button>
                  </Menu.Item>
                  {item.originType === 'ClearCase' && (
                    <Menu.Item>
                      <Button size="small" type="link" onClick={() => actions.execTask(item)}>
                        执行迁移任务
                      </Button>
                    </Menu.Item>
                  )}
                  {item.originType === 'ClearCase' && item.task_id && (
                    <Menu.Item>
                      <Button size="small" type="link" onClick={() => actions.gotoTaskDetail(item)}>
                        跳转任务详情页
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

const PlanList: React.FC = () => {
  const tableRef = React.useRef<any>();
  const planCreatorRef = React.useRef<any>();
  const planStatusSwitcherRef = React.useRef<any>();
  const history = useHistory();

  const actions: Actions = {
    updatePlan({ id }) {
      planCreatorRef.current.openModal('update', id);
    },
    async deletePlan({ id }) {
      await planServices.deletePlan(id);
      tableRef.current.reload();
      message.success('计划删除成功');
    },
    async execTask({ id }) {
      const { message: taskId } = await planServices.updatePlan(id, { status: '迁移中' });
      history.push(`/task/detail?id=${taskId}`);
    },
    toggleStatus(plan) {
      planStatusSwitcherRef.current.openModal(plan);
    },
    gotoTaskDetail({ id }) {
      history.push(`/task/detail?id=${id}`);
    },
  };
  return (
    <>
      <PlanCreator onSuccess={() => tableRef.current.reload()} actionRef={planCreatorRef} />
      <PlanStatusSwitcher
        onSuccess={() => tableRef.current.reload()}
        actionRef={planStatusSwitcherRef}
      />
      <ProTable
        rowKey="id"
        scroll={{ x: 1500 }}
        actionRef={tableRef}
        request={async ({ pageSize = 10, current }) => {
          const { planInfo, count } = await planServices.getPlans({
            offset: (current! - 1 || 0) * pageSize,
            limit: pageSize || 10,
          });

          return {
            data: planInfo,
            success: true,
            total: count,
          };
        }}
        columns={getColumns(actions)}
        search={false}
        toolBarRender={() => [
          <Button
            size="small"
            type="primary"
            onClick={() => {
              planCreatorRef.current.openModal();
            }}
          >
            新建迁移计划1
          </Button>,
        ]}
      ></ProTable>
    </>
  );
};

export default PlanList;
