import React from 'react';
import { Modal } from 'antd';
import classnames from 'classnames';
import { guid } from '@/utils/utils';
import { useToggle } from 'react-use';
import Table from '@ant-design/pro-table';
import ProCard from '@ant-design/pro-card';
import type { Task } from '@/typings/model';
import type { ProColumns } from '@ant-design/pro-table';
import { humanizeDuration, renderCardTitle } from '../../helper';

import styles from './style.less';

const descriptionsGenerator = (fieldKeys: string[], data: any) => {
  const labels: Record<string, string> = {
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
    dir: '组件子目录',
    keep: '文件占位名',
  };

  // 序列化
  const matrix = fieldKeys.reduce((res, item, index) => {
    if (index % 2 === 0) {
      res.push([item]);
    } else {
      const pos = Math.floor(index / 2);
      res[pos] = res[pos].concat(item);
    }
    return res;
  }, [] as any[][]);

  return matrix.map((keys: any[]) => {
    return (
      <div className={styles.col} key={guid()}>
        {keys.map((key, index) => {
          const isLeftRow = index % 2 === 0;
          return (
            <div className={classnames(styles.row, isLeftRow && styles.left)} key={guid()}>
              <span>{labels[key]}：</span>
              <span>{data[key] || '-'}</span>
            </div>
          );
        })}
      </div>
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
          width="800px"
          title="任务详情"
          visible={visible}
          onOk={() => toggleVisible(false)}
          onCancel={() => toggleVisible(false)}
          cancelButtonProps={{ style: { display: 'none' } }}
        >
          <div className={styles.gutter}>
            <div className={styles.col}>
              <div className={styles.row}>{renderCardTitle('ClearCase')}</div>
              <div className={styles.row}>{renderCardTitle('Git')}</div>
            </div>
            {descriptionsGenerator(
              [
                'pvob',
                'gitURL',
                'component',
                'gitEmail',
                'ccUser',
                'gitUser',
                'ccPassword',
                'gitPassword',
                'dir',
              ],
              data.taskModel,
            )}
            <div className={styles.divider} />
            {data.taskModel.matchInfo.map((matchInfo) =>
              descriptionsGenerator(['stream', 'gitBranch'], matchInfo),
            )}
            <div className={styles.col}>
              <span className={styles.row}>
                <span>是否保留空目录：{data.taskModel.includeEmpty ? '是' : '否'}</span>
                <span style={{ marginLeft: '2em' }}>占位文件名：{data.taskModel.keep}</span>
              </span>
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
