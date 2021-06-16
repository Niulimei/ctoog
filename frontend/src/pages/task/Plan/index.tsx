import React, {useState, useMemo} from 'react';
import { useHistory, useModel } from 'umi';
import type { Plan } from '@/typings/model';
import ProTable from '@ant-design/pro-table';
import { DownOutlined, SmileOutlined } from '@ant-design/icons';
import { plan as planServices } from '@/services';
import {authTokenAction} from '@/utils/request';
import PlanCreator from './components/PlanCreator';
import { Button, Menu, Dropdown, message, notification, Upload, Modal } from 'antd';
import type { ProColumns } from '@ant-design/pro-table';
import PlanStatusSwitcher from './components/PlanStatusSwitcher';
import styles from './index.less';
import classnames from "classnames";

type Actions = Record<
  'updatePlan' | 'deletePlan' | 'execTask' | 'execSvnTask' | 'toggleStatus' | 'gotoTaskDetail' | 'gotoSvnTaskDetail',
  (payload: Plan.Item) => void
>;

const StatusOptions = [
  "未迁移",
  "已迁移",
  "已切换",
  "迁移中",
];

const RepoSvnOptions = [
  'ClearCase',
  'ICDP(Gerrit)',
  '私服',
  'svn',
];

const NoRepoOptions = [
  'ClearCase',
  'ICDP(Gerrit)',
  '私服',
];

const GroupOptions = [
  '北京事业群',
  '厦门事业群',
  '成都事业群',
  '深圳事业群',
  '上海事业群',
  '广州事业群',
  '广研事业群',
  '武汉事业群',
  '基础技术中心',
  '实施管理中心',
  '大数据中心',
  '产品经营中心',
  '智能云事业部',
  '交付事业部',
];

const getColumns = (actions: Actions, hasSvn: boolean): ProColumns<Plan.Item>[] => {
  const handleMenuClick = (item: any, key: any) => {
    actions[key]?.(item);
  };
  return [

    {
      title: '源仓库类型',
      dataIndex: 'originType',
      ellipsis: true,
      valueEnum: hasSvn ? RepoSvnOptions : NoRepoOptions,
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
      title: '迁移状态',
      dataIndex: 'status',
      ellipsis: true,
      valueEnum: StatusOptions,
    },
    {
      title: '源仓库信息',
      ellipsis: true,
      width: 250,
      hideInSearch: true,
      // @ts-ignore
      render(_, item: Plan.Item) {
        return item.pvob ? `${item.pvob}/${item.component}` : item.originUrl;
      },
    },
    {
      title: '事业群',
      dataIndex: 'group',
      ellipsis: true,
      hideInSearch: false,
      valueEnum: GroupOptions,
    },
    {
      title: '项目组',
      dataIndex: 'team',
      ellipsis: true,
      hideInSearch: false,
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
      title: '任务编号',
      dataIndex: 'task_id',
      ellipsis: true,
      hideInSearch: true,
      width: 80,
    },
    {
      title: '对接人姓名',
      dataIndex: 'supporter',
      ellipsis: true,
      hideInSearch: false,
      width: 80,
      hideInTable: true
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
                <Menu onClick={(e) => handleMenuClick(item, e.key)}>
                  <Menu.Item key="deletePlan">
                    <Button size="small" type="link">
                      删除
                    </Button>
                  </Menu.Item>
                  <Menu.Item key="toggleStatus">
                    <Button size="small" type="link">
                      变更状态
                    </Button>
                  </Menu.Item>
                  {item.originType === 'ClearCase' && (
                    <Menu.Item key={item.status === '未迁移' ? 'execTask' : 'gotoTaskDetail'}>
                      <Button size="small" type="link">
                        执行迁移任务
                      </Button>
                    </Menu.Item>
                  )}
                  {item.originType === 'svn' && (
                    <Menu.Item key={item.status === '未迁移' ? 'execSvnTask' : 'gotoSvnTaskDetail'}>
                      <Button size="small" type="link">
                        执行迁移任务
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
  const { initialState } = useModel('@@initialState');
  const [importLoading, setImportLoading] = useState<boolean>(false);
  const [visible, setVisible] = useState<boolean>(false);
  const tableRef = React.useRef<any>();
  const planCreatorRef = React.useRef<any>();
  const planStatusSwitcherRef = React.useRef<any>();
  const history = useHistory();

  const { RouteList = [] } = initialState;
  const hasSvn = useMemo(() => RouteList.includes('svnRoute'), [RouteList]);

  const beforeonChange = (props) => {
    const {file = {}} = props;
    let status = file?.status;
    setImportLoading(true);
    if (status === "uploading") {
      setImportLoading(true);
    } else if (status === 'done') {
      message.success('上传成功!');
      tableRef?.current?.reload();
      setImportLoading(false);
    } else if (status === 'error') {
      const status = file?.error?.status;
      const errorMsg = file?.response;
      if (status === 400 || status === 500) {
        message.error(errorMsg || '上传失败!');
      } else {
        message.error('上传失败!');
      }
      setImportLoading(false);
    } else {
      setImportLoading(false);
    }
  };

  const actions: Actions = {
    updatePlan({ id }) {
      planCreatorRef.current.openModal('update', id);
    },
    async deletePlan({ id }) {
      await planServices.deletePlan(id);
      tableRef.current.reload();
      message.success('计划删除成功');
    },
    async execSvnTask({ id }) {
      const { message: msg } = await planServices.updatePlan(id, { status: '迁移中' });
      const taskId = id && msg;
      history.push(`/task/svnDetail?id=${taskId}`);
    },
    async execTask({ id }) {
      const { message: msg } = await planServices.updatePlan(id, { status: '迁移中' });
      const taskId = id && msg;
      history.push(`/task/detail?id=${taskId}`);
    },
    toggleStatus(plan) {
      planStatusSwitcherRef.current.openModal(plan);
    },
    gotoTaskDetail({ task_id }) {
      history.push(`/task/detail?id=${task_id}`);
    },
     gotoSvnTaskDetail({ task_id }) {
      history.push(`/task/svnDetail?id=${task_id}`);
    },
  };
  /** 提示信息 */
  const description = () => {
    return (
      <>
        <div className={classnames(styles.modalTxt)}>
          <span>目前有事业群反映，物理子系统过多，逐个填报过于繁琐,如果想批量录入，工作组提供了临时的批量导入方案。
          现提供仓库迁移范围模板，事业群同事可以编辑迁移范围模板excel，由工作组每日进行导入</span>
          <span>(1）模板存放地址：云上，\128.194.1.13\全生命周期it管理\工作目录\仓库迁移信息\仓库迁移范围信息-模板.xlsx
          事业群可以复制模板，修改名称为【仓库迁移范围信息-事业群.xlsx】进行填写,填写前请先阅读同级文件【readme.txt】</span>
          <span>(2）如果无法访问共享的，请使用sftp填报地址：128.194.225.15 用户名密码：repinf/inf0525
          存放位置：/home/ap/repinf</span>
        </div>
        <Upload
          action="/import/plan"
          name="uploadFile"
          className={classnames(styles.uploadBtn)}
          headers={{
            Authorization: authTokenAction.get()
          }}
          withCredentials={true}
          showUploadList={false}
          onChange={beforeonChange}
          beforeUpload={() => setVisible(false)}
        >
          <Button
            size="large"
            type="primary"
            loading={importLoading}
          >上传文件
          </Button>
        </Upload>
      </>
    )
  }
  return (
    <>
      <PlanCreator onSuccess={() => tableRef.current.reload()} actionRef={planCreatorRef} />
      <PlanStatusSwitcher
        onSuccess={() => tableRef.current.reload()}
        actionRef={planStatusSwitcherRef}
      />
      <ProTable
        headerTitle="迁移计划"
        rowKey="id"
        scroll={{ x: 1500 }}
        actionRef={tableRef}
        request={async ({ pageSize = 10, current, group: groupIndex, team, supporter, status, originType}) => {
          const group = GroupOptions[groupIndex];
          const { planInfo, count } = await planServices.getPlans({
            offset: (current! - 1 || 0) * pageSize,
            limit: pageSize || 10,
            team,
            supporter,
            group,
            status: StatusOptions[status],
            originType: hasSvn ? RepoSvnOptions[originType] : NoRepoOptions[originType],
          });
          return {
            data: planInfo,
            success: true,
            total: count,
            team,
            supporter,
            group,
            status: StatusOptions[status],
            originType: hasSvn ? RepoSvnOptions[originType] : NoRepoOptions[originType],
          };
        }}
        columns={getColumns(actions, hasSvn)}
        search={{
          defaultCollapsed: false,
          span: 6
        }}
        toolBarRender={() => [
          <Button size="small"
            type="primary"
            onClick={() => {
              setVisible(true);
            }}
          >批量计划导入
          </Button>,
          <Button
            size="small"
            type="primary"
            onClick={() => {
              planCreatorRef.current.openModal();
            }}
          >
            新建迁移计划
          </Button>,
        ]}
      ></ProTable>
      <Modal
        title="批量计划导入"
        visible={visible}
        onCancel={() => setVisible(false)}
        okButtonProps={{ style: { display: 'none' } }}
      >
        {description()}
      </Modal>
    </>
  );
};

export default PlanList;
