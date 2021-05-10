import React from 'react';
import { Modal } from 'antd';
import { guid } from '@/utils/utils';
import { useToggle } from 'react-use';
import Table from '@ant-design/pro-table';
import ProCard from '@ant-design/pro-card';
import type { Task } from '@/typings/model';
import { humanizeDuration } from '../../helper';
import type { ProColumns } from '@ant-design/pro-table';

import styles from './style.less';

const descriptionsGenerator = (fieldKeys: string[], data: any) => {
  const taskKeyLabel: Record<string, string> = {
    pvob: 'PVOB',
    component: 'Component',
    ccUser: 'CC 用户名',
    gitBranch: 'Git分支',
    ccPassword: 'CC 密码',
    stream: 'CC 开发流',
    gitURL: 'Git Repo URL',
    gitUser: 'Git 用户名',
    gitEmail: 'Git Email',
    gitPassword: 'Git 密码',
    includeEmpty: '是否保留空目录',
  };
  const valueFormatter: any = {
    includeEmpty(val: boolean) {
      return val ? '是' : '否';
    },
  };
  return fieldKeys.map((key) => {
    const formatter = valueFormatter[key];

    return (
      <li key={guid()}>
        <span>{taskKeyLabel[key]}：</span>
        {formatter ? formatter(data[key]) : data[key]}
      </li>
    );
  });
};

const TableColumns: ProColumns<Task.Log>[] = [
  {
    title: '任务序号',
    dataIndex: 'logID',
  },
  {
    title: '任务状态',
    renderText(item: Task.Log) {
      return item.status;
    },
  },
  {
    title: '开始时间',
    renderText(item: Task.Log) {
      return item.startTime;
    },
  },
  {
    title: '结束时间',
    renderText(item: Task.Log) {
      return item.endTime;
    },
  },
  {
    title: '历时时长',
    renderText(item: Task.Log) {
      return item.duration ? humanizeDuration(Number(item.duration)) : '-';
    },
  },
];

const TaskDetail: React.FC<{ data?: Task.Detail; actionRef: any }> = ({ data, actionRef }) => {
  const [visible, toggleVisible] = useToggle(false);

  React.useImperativeHandle(actionRef, () => ({
    openModal() {
      toggleVisible(true);
    },
  }));

  return (
    <div>
      {data ? (
        <Modal
          width="850px"
          title="任务详情"
          visible={visible}
          onOk={() => toggleVisible(false)}
          onCancel={() => toggleVisible(false)}
          cancelButtonProps={{ style: { display: 'none' } }}
        >
          <div className={styles.gutter}>
            <div className={styles.row}>
              <h6>ClearCase</h6>
              <ul className={styles.list}>
                {descriptionsGenerator(['pvob', 'component', 'ccUser'], data.taskModel)}
                <div className={styles.divider} />
                {data.taskModel.matchInfo.map(({ stream }) =>
                  descriptionsGenerator(['stream'], { stream }),
                )}
                {descriptionsGenerator(['includeEmpty'], data.taskModel)}
              </ul>
            </div>
            <div className={styles.row}>
              <h6>Git</h6>
              <ul className={styles.list}>
                {descriptionsGenerator(['gitURL', 'gitEmail', 'gitUser'], data.taskModel)}
                <div className={styles.divider} />
                {data.taskModel.matchInfo.map(({ gitBranch }) =>
                  descriptionsGenerator(['gitBranch'], { gitBranch }),
                )}
              </ul>
            </div>
          </div>
          <ProCard title="执行历史记录" style={{ marginTop: 22 }}>
            <Table
              pagination={{
                pageSize: 5,
              }}
              rowKey="logID"
              search={false}
              toolBarRender={false}
              columns={TableColumns}
              dataSource={data.logList}
            />
          </ProCard>
        </Modal>
      ) : null}
    </div>
  );
};

export default TaskDetail;
