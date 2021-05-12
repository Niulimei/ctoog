import React from 'react';
import Table from '@ant-design/pro-table';
import type { ProColumns } from '@ant-design/pro-table';
import UserCreator from '@/pages/admin/UserList/components/UserCreator';
import { User } from '@/typings/model';
import { user as UserService } from '@/services';

const PageSize = 10;

const Columns: ProColumns[] = [
  {
    title: '用户名',
    dataIndex: 'username',
  },
  {
    title: '用户权限',
    dataIndex: 'role_id',
    renderText(role) {
      return role === User.Role.ADMIN ? '管理员' : '普通用户';
    },
  },
];

function UserList() {
  const tableRef = React.useRef<any>(null);
  return (
    <Table
      pagination={{
        pageSize: PageSize,
        showSizeChanger: false,
      }}
      rowKey="username"
      actionRef={tableRef}
      request={async (params) => {
        const { userInfo, count } = await UserService.getUsers({
          offset: (params.current! - 1 || 0) * PageSize,
          limit: params.pageSize || 10,
        });
        return {
          data: userInfo,
          success: true,
          total: count,
        };
      }}
      headerTitle="用户列表"
      columns={Columns}
      toolBarRender={() => [
        <UserCreator
          onCreateSuccess={() => {
            tableRef.current.reload();
          }}
        />,
      ]}
      search={false}
    />
  );
}

export default UserList;
