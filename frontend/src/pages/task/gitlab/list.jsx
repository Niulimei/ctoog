import ProTable from '@ant-design/pro-table';
import { Button, message, Modal, Popconfirm } from 'antd';
import { StepsForm, ProFormText, ProFormTextArea, ProFormSelect, ProFormCheckbox } from '@ant-design/pro-form';
import { useEffect, useRef, useState } from 'react';
import { useBoolean, useSetState } from 'ahooks';
import { PlayCircleOutlined, CodeOutlined, DeleteOutlined } from '@ant-design/icons';
import Log from './log';
import CreateForm from './createForm';
import { task as taskService, gitlab as gitlabService } from '@/services';


const getColumns = (actions) => [
  {
    title: '任务编号',
    dataIndex: 'id',
    hideInSearch: true,
  },
  {
    title: 'Gitlab Group',
    dataIndex: 'gitlabGroup',
    hideInSearch: true,
  },
  {
    title: 'Gitlab Project',
    dataIndex: 'gitlabProject',
    hideInSearch: true,
  },
  {
    title: 'Gitee Group',
    dataIndex: 'giteeGroup',
    hideInSearch: true,
  },
  {
    title: 'Gitee Repo',
    dataIndex: 'gitRepo',
    hideInSearch: true,
  },
  {
    title: '当前状态',
    dataIndex: 'status',
    valueEnum: {pending: 'pending', success: 'success', failed: 'failed', init: 'init'},
  },
  {
    title: '最后一次完成时间',
    dataIndex: 'lastCompleteDateTime',
    hideInSearch: true,
  },
  {
    title: '操作',
    valueType: 'action',
    hideInSearch: true,
    render(item) {
      return (
        <div style={{'display': 'flex', gap: '6px'}}>
          <PlayCircleOutlined onClick={() => {
            const { status, id } = item.props.record;
            if (status === 'pending') {
              message.error('任务pending中');
              return;
            }
            taskService.startTask(id);
          }} />
          <CodeOutlined onClick={() => actions.checkLog(item.props.record.id)} />
          <Popconfirm
            title="删除吗"
            onConfirm={() => {
              taskService.deleteTask(item.props.record.id).then(() => {
                message.success('删除成功');
              });
            }}
            okText="yes"
            cancelText="no"
          >
            <DeleteOutlined />
          </Popconfirm>
        </div>
      );
    }
  }
];

const GitlabTaskList = () => {
  const [visible, setVisible] = useState(false);
  const [loading, { toggle, setTrue, setFalse }] = useBoolean(false);
  const [log, setLog] = useSetState({ id: null, visible: false});
  const formRef = useRef(null);
  const actions = {
    checkLog(id) {
      setLog({
        id,
        visible: true,
      });
    },
  };

  return (
    <>
      <ProTable
        loading={loading}
        columns={getColumns(actions)}
        request={({ current, pageSize, status, ...params }, sorter, filter) => {
          // 表单搜索项会从 params 传入，传递给后端接口。
          console.log(params, sorter, filter);
          setTrue();
          const offset = (current - 1) * pageSize;
          return gitlabService.getTasks({
            limit: pageSize,
            offset,
            status,
            modelType: 'gitlab'
          }).then((data) => {
            setFalse();
              return {
                data: data.taskInfo,
                success: true,
                total: data.total,
              };
          });
        }}
        rowKey="id"
        toolBarRender={() => [
          <Button key="button" type="primary" size="small" onClick={() => setVisible(true)} >
            新建Gitlab迁移任务
          </Button>
        ]}
      >
      </ProTable>
      <CreateForm visible={visible} setVisible={setVisible} />
      <Modal
        title={log.id}
        visible={log.visible}
        onCancel={() => setLog({
          visible: false,
        })}
        width={800}
      >
        <Log id={log.id} />
      </Modal>
    </>
  )
}

export default GitlabTaskList;