import ProTable from '@ant-design/pro-table';
import { Button, Modal } from 'antd';
import { StepsForm, ProFormText, ProFormTextArea, ProFormSelect, ProFormCheckbox } from '@ant-design/pro-form';
import { useRef, useState } from 'react';
import { useSetState } from 'ahooks';
import { PlayCircleOutlined, CodeOutlined, DeleteOutlined } from '@ant-design/icons';
import Log from './log';


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
          <DeleteOutlined onClick={() => { console.log(item) }} />
        </div>
      );
    }
  }
];

const data = [
  {
    taskNo: '001',
    time: '2020-01-01',
    gitlab_group: 'gitlab_group',
    gitlab_project: 'gitlab_project',
    gitee_group: 'gitee_groupooo',
    gitee_repo: 'aa',
    status: 'running',
    id: 0,
  },
  {
    taskNo: '002',
    time: '2022-01-01',
    gitlab_group: 'gitlab_group',
    gitlab_project: 'gitlab_project',
    gitee_group: 'gitee_groupooo',
    gitee_repo: 'aa',
    status: 'running',
    id: 1,
  }
];

const GitlabTaskList = () => {
  const [visible, setVisible] = useState(false);
  const [log, setLog] = useSetState({ id: null, visible: false});
  const [loading, setLoading] = useState(false);
  const [migrationType, setMigrationType] = useState('Group');
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
        columns={getColumns(actions)}
        request={(params, sorter, filter) => {
          // 表单搜索项会从 params 传入，传递给后端接口。
          console.log(params, sorter, filter);
          return Promise.resolve({
            data: data,
            success: true,
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
      <Modal
        title="新建Gitlab迁移任务"
        visible={visible}
        onOk={() => {
          setLoading(true);
          setTimeout(() => {
            setVisible(false);
            setLoading(false);
          }, 3000);
        }}
        confirmLoading={loading}
        onCancel={() => setVisible(false)}
        footer={null}
        width={620}
      >
        <StepsForm
          onFinish={(value) => {
            console.log(value);
          }}
          formRef={formRef}
        >
          <StepsForm.StepForm
            name="base"
            title="Gitlab基础设置"
            layout={"horizontal"}
            labelCol={{
              span: 8,
            }}
            wrapperCol={{
              span: 16,
            }}
          >
            <ProFormText
              name="url"
              width="md"
              label="原平台地址"
              tooltip="最长为 24 位，用于标定的唯一 id"
              placeholder="请输入地址"
              rules={[{ required: true }]}
            />
            <ProFormText
              name="token"
              width="md"
              label="原平台令牌"
              tooltip="最长为 24 位，用于标定的唯一 id"
              placeholder="请输入令牌"
              rules={[{ required: true }]}
            />
          </StepsForm.StepForm>
          <StepsForm.StepForm
            name="move"
            title="Gitlab迁移数据设置"
            onFinish={async () => {
              return true;
            }}
            layout={"horizontal"}
            labelCol={{
              span: 8,
            }}
            wrapperCol={{
              span: 16,
            }}
          >
            <ProFormSelect
              name="type"
              label="迁移类型"
              valueEnum={{
                Group: 'Group',
                Project: 'Project',
              }}
              placeholder="请选择迁移类型"
              rules={[{ required: true, message: '请选择迁移类型' }]}
              onChange={(v) => setMigrationType(v)}
              initialValue={"Group"}
            />
            <ProFormText
              name="path"
              label={migrationType + ' path'}
              placeholder="请输入path"
              rules={[{ required: true }]}
            />
            <ProFormCheckbox.Group
              name="checkbox"
              layout="vertical"
              label="行业分布"
              options={['Issue', 'Milestone', 'Merge Request', 'Wiki', 'Permission']}
            />
          </StepsForm.StepForm>
          <StepsForm.StepForm
            name="gitee"
            title="Gitee基础设置"
            layout={"horizontal"}
            labelCol={{
              span: 8,
            }}
            wrapperCol={{
              span: 16,
            }}
          >
            <ProFormText
              name="target_url"
              label="目标平台地址"
              rules={[{ required: true }]}
            />
            <ProFormText
              name="target_token"
              label="目标平台令牌"
            />
            <ProFormText
              name="target_group"
              label="放置在目标组"
              tooltip="为空则放在企业根目录下"
              rules={[{ required: true }]}
            />
          </StepsForm.StepForm>
        </StepsForm>
      </Modal>
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