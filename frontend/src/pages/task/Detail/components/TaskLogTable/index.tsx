import React from 'react';
import { Button } from 'antd';
import Table from '@ant-design/pro-table';
import type { Task } from '@/typings/model';
import { humanizeDuration } from '@/utils/utils';
import type { ProColumns } from '@ant-design/pro-table';

interface IProps {
  data?: Task.Detail['logList'];
  onDisplayLog?: (id: string) => void;
}

const TaskLogTable: React.FC<IProps> = ({ data, onDisplayLog }) => {
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
    {
      title: '操作',
      align: 'center',
      // @ts-ignore
      render({ logID }: Task.Log) {
        return (
          <Button onClick={() => onDisplayLog?.(logID)} type="link">
            查看日志
          </Button>
        );
      },
    },
  ];
  if (!data) return null;
  return (
    <Table
      pagination={{
        pageSize: 5,
        pageSizeOptions: ['5', '10', '15'],
      }}
      rowKey="logID"
      search={false}
      dataSource={data}
      toolBarRender={false}
      columns={TableColumns}
    />
  );
};

export default TaskLogTable;
