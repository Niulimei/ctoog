import ProTable from '@ant-design/pro-table';
import { Button, Modal, Popconfirm } from 'antd';
import { StepsForm, ProFormText, ProFormTextArea, ProFormSelect, ProFormCheckbox } from '@ant-design/pro-form';
import { useEffect, useRef, useState } from 'react';
import { useBoolean, useSetState } from 'ahooks';
import { PlayCircleOutlined, CodeOutlined, DeleteOutlined } from '@ant-design/icons';
import Log from './log';
import CreateForm from './createForm';
import { gitlab as gitlabService } from '@/services';


const getColumns = (actions) => [
  {
    title: '任务编号',
    dataIndex: 'taskNo',
    hideInSearch: true,
  },
  {
    title: 'Gitlab Group',
    dataIndex: 'gitlab_group',
    hideInSearch: true,
  },
  {
    title: 'Gitlab Project',
    dataIndex: 'gitlab_project',
    hideInSearch: true,
  },
  {
    title: 'Gitee Group',
    dataIndex: 'gitee_group',
    hideInSearch: true,
  },
  {
    title: 'Gitee Repo',
    dataIndex: 'gitee_repo',
    hideInSearch: true,
  },
  {
    title: '当前状态',
    dataIndex: 'status',
  },
  {
    title: '最后一次完成时间',
    dataIndex: 'time',
    hideInSearch: true,
  },
  {
    title: '操作',
    valueType: 'action',
    hideInSearch: true,
    render(item) {
      return (
        <div style={{'display': 'flex', gap: '6px'}}>
          <PlayCircleOutlined onClick={() => { console.log(item) }} />
          <CodeOutlined onClick={() => actions.checkLog(item.props.record.id)} />
          <Popconfirm
            title="删除吗"
            onConfirm={() => console.log('delete')}
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
        request={(params, sorter, filter) => {
          // 表单搜索项会从 params 传入，传递给后端接口。
          console.log(params, sorter, filter);
          setTrue();
          return gitlabService.getTasks({limit: 20, offset: 0, modelType: 'gitlab'}).then((data) => {
            console.log(data);
            setFalse();
              return {
                data: data.list,
                success: true,
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