import { ModalForm, ProFormText } from '@ant-design/pro-form';
import { Button, Form, message } from 'antd';
import { user as userService } from '@/services';
import md5 from 'md5';

function UserCreator(props: React.PropsWithChildren<{ onCreateSuccess?: () => void }>) {
  const [form] = Form.useForm();
  return (
    <ModalForm
      onFinish={async ({ password, username }) => {
        try {
          await userService.createUser({
            username,
            password: md5(password),
          });
          message.success('新建用户成功');
          props.onCreateSuccess?.();
          return true;
        } catch (err) {
          message.error('新建用户出现异常');
          return false;
        }
      }}
      form={form}
      width="500px"
      title="新建用户"
      modalProps={{ okText: '新建' }}
      trigger={
        <Button size="small" type="primary">
          新建用户
        </Button>
      }
    >
      <ProFormText placeholder="请输入用户名" name="username" />
      <ProFormText.Password placeholder="请输入用户密码" name="password" />
      <ProFormText.Password
        placeholder="请再次输入用户密码"
        rules={[
          {
            async validator(_, value) {
              if (value !== form.getFieldValue('password')) {
                throw new Error('两次密码输入不一致');
              }
            },
          },
        ]}
        name="retypePassword"
      />
    </ModalForm>
  );
}

export default UserCreator;
