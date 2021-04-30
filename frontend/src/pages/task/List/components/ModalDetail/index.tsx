import React from 'react';
import { Modal } from 'antd';
import ProCard from '@ant-design/pro-card';
import { useToggle } from 'react-use';
import type { Task } from '@/typings/model';
import type { ProColumns } from '@ant-design/pro-table';
import Table from '@ant-design/pro-table';
import { guid } from '@/utils/utils';

import styles from './style.less';

const EmptyColSpace = <li key="emptyCol" className={styles.emptyCol} />;

const descriptionsGenerator = (fieldKeys: string[], data: any) => {
  const taskKeyLabel: Record<string, string> = {
    pvob: 'PVOB',
    component: 'Component',
    ccUser: 'CC 用户名',
    ccPassword: 'CC 密码',
    stream: 'CC 开发流',
    gitURL: 'Git Repo URL',
    gitUser: 'Git 用户名',
    gitPassword: 'Git 密码',
    gitBranch: 'Git分支',
  };
  return fieldKeys.map((key) => {
    if (key === 'EmptyColSpace') return EmptyColSpace;
    return (
      <li key={guid()}>
        <span>{taskKeyLabel[key]}：</span>
        {data[key]}
      </li>
    );
  });
};

const TableColumns: ProColumns<Task.Log>[] = [
  {
    title: '任务序号',
    dataIndex: 'id',
  },
  {
    dataIndex: '任务状态',
    renderText(item: Task.Log) {
      return item.status;
    },
  },
  {
    dataIndex: '开始事件',
    renderText(item: Task.Log) {
      return item.startTime;
    },
  },
  {
    dataIndex: '开始事件',
    renderText(item: Task.Log) {
      return item.endTime;
    },
  },
  {
    dataIndex: '历时时长',
    renderText(item: Task.Log) {
      return item.duration;
    },
  },
];

const ModalDetail: React.FC<{ data?: Task.Detail; actionRef: any }> = ({ data, actionRef }) => {
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
          title="任务详情"
          width="700px"
          visible={visible}
          cancelButtonProps={{ style: { display: 'none' } }}
          onOk={() => toggleVisible(false)}
          onCancel={() => toggleVisible(false)}
        >
          <div className={styles.gutter}>
            <div className={styles.row}>
              <h6>ClearCase</h6>
              <ul className={styles.list}>
                {descriptionsGenerator(['pvob', 'component', 'ccUser'], data.taskModel)}
                <div className={styles.pipe} />
                {data.taskModel.matchInfo.map(({ stream }) =>
                  descriptionsGenerator(['stream'], { stream }),
                )}
              </ul>
            </div>
            <div className={styles.row}>
              <h6>Git</h6>
              <ul className={styles.list}>
                {descriptionsGenerator(['gitURL', 'EmptyColSpace', 'gitUser'], data.taskModel)}
                <div className={styles.pipe} />
                {data.taskModel.matchInfo.map(({ gitBranch }) =>
                  descriptionsGenerator(['gitBranch'], { gitBranch }),
                )}
              </ul>
            </div>
          </div>
          {EmptyColSpace}
          <ProCard title="执行历史记录" style={{ marginTop: 22 }}>
            <Table
              search={false}
              toolBarRender={false}
              dataSource={data.logList}
              columns={TableColumns}
            />
          </ProCard>
        </Modal>
      ) : null}
    </div>
  );
};

export default ModalDetail;
