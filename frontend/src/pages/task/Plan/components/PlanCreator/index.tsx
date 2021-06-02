import React, {useCallback, useState} from 'react';
import {
  ModalForm,
  ProFormText,
  ProFormRadio,
  ProFormSelect,
  ProFormTextArea,
  ProFormDatePicker,
} from '@ant-design/pro-form';
import { useModel } from 'umi';
import dayjs from 'dayjs';
import { guid } from '@/utils/utils';
import classnames from 'classnames';
import {toJS} from 'mobx';
import { observer } from 'mobx-react';
import type { Plan } from '@/typings/model';
import { Row, Col, Form, message, Button, Checkbox, Input, AutoComplete } from 'antd';
import { MinusOutlined, PlusOutlined } from '@ant-design/icons';
import type { FormInstance } from 'antd/es/form';
import { plan as planServices, task as taskService } from '@/services';
import { useClearCaseSelectEnum } from '@/utils/hooks';
import { useToggle, useMount, useUpdate } from 'react-use';

import styles from './style.less';

interface FormSectionProps {
  title?: string;
  left?: React.ReactNode;
  right?: React.ReactNode;
  wholeLine?: React.ReactNode;
}
const FormSection: React.FC<FormSectionProps> = ({ left, title, right, wholeLine }) => (
  <div className={styles.section}>
    {title && <h4 className={styles.sectionHeader}>{title}</h4>}
    <Row gutter={24}>
      {left && <Col span={11}>{left}</Col>}
      {right && (
        <Col span={11} offset={2}>
          {right}
        </Col>
      )}
    </Row>
    {wholeLine && <div>{wholeLine}</div>}
  </div>
);

const FieldAutoComplete = props => {
  const streamObj = toJS(props.options);
  let arr = Object.keys(streamObj).map(item => {
    return {label: item, value: item}
  });
  const [streamOptions, setStreamOptions] = useState([]);

  const onSearch = useCallback(value => {
    if (value) {
      setStreamOptions([{label: value, value}]);
      arr.map(item => {
        if (item.value === value) {
          setStreamOptions([]);
        }
      })
    } else {
      setStreamOptions([]);
    }

  }, [setStreamOptions]);

  return (
     <Form.Item {...props}>
       <AutoComplete
         options={streamOptions.concat(arr) || []}
         filterOption={true}
         onSearch={onSearch}
         getPopupContainer={triggerNode => triggerNode.parentElement}
       />
     </Form.Item>
  )
}

/** 生成表单项 */
const formFieldsGenerator = (fields: any) => {
  const renderFieldComponent = ({ component, name, required, rules, ...restProps }: any) => {
    const requiredRules = [
      {
        required: true,
        message: restProps.placeholder ? restProps.placeholder : `${name} 为必填参数`,
      },
    ];
    // eslint-disable-next-line no-param-reassign
    rules = (rules || []).concat(required && requiredRules).filter(Boolean);

    return React.createElement(component, {
      key: name,
      required,
      rules,
      name,
      ...restProps,
    });
  };

  /** 判断当前的节点是不是带有require属性 */
  const isRequired = (node: any): boolean => {
    if(node?.props?.rules){
      if(Array.isArray(node.props.rules) && node.props.rules[0]){
        return node.props.rules[0].hasOwnProperty('required')  || node.props.rules[1].hasOwnProperty('required')
      }
    }
    return false

  }

  return fields.map((nodes: any) => {
    const key = guid();

    const [leftNode, rightNode, actionNode] = nodes.map((node: any) =>
      node.component ? renderFieldComponent(node) : node,
    );
    // console.log(leftNode);

    return (
      <div className={styles.col} key={key}>
        <div className={classnames(styles.row, styles.left, isRequired(leftNode) ? styles.must : '')}>{leftNode}</div>
        <div className={classnames(styles.row, styles.right, isRequired(rightNode) ? styles.must : '')}>{rightNode}</div>
        {actionNode ? <div className={styles.action}>{actionNode}</div> : null}
      </div>
    );
  });
};

const OriginTypeOptions = ['ClearCase', 'ICDP(Gerrit)', '私服'];

const ExType = ['svn'];

const TranslateCCTypeOptions = ['项目组自己迁移', '工作组帮迁移'];

const TranslateTypeOptions = ['项目组自己迁移'];

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

const ProjectTypeOptions = ['Java', 'Python', 'C/C++', 'JavaScript', '其他'];

const InitialValues = {
  originType: 'ClearCase',
  translateType: '项目组自己迁移',
};
const CustomChangeHandlers: Partial<
  Record<keyof Plan.Base, (form: FormInstance<Plan.Base>, value: any, dispatch: any) => void>
> = {
  originType(form, value, dispatch) {
    dispatch({ type: 'forceUpdate' });
  },
  pvob(form, value, dispatch) {
    dispatch({ type: 'getComponentValueEnum', payload: value });
    form.setFieldsValue({ configLib: value });
  },
  component(form, value, dispatch) {
    dispatch({type: 'stream', payload: { component: value, pvob: form.getFieldValue('pvob') }});
    const { matchInfo } = form.getFieldsValue(['matchInfo']);
    form.setFieldsValue({
      matchInfo: (matchInfo || []).map((info: any) => ({ ...info, stream: '' })),
    });
  },
  plan_start_time(form, value) {
    const switchTime = form.getFieldValue('plan_switch_time');
    if (dayjs(switchTime).isAfter(dayjs(value)))
      form.setFields([{ name: ['plan_switch_time'], errors: [] }])
    else
      form.setFields([{ name: ['plan_switch_time'], errors: ['计划切换日期应在计划迁移日期之后'] }])
  }
};

interface IPlanCreatorProps {
  actionRef?: React.ForwardedRef<{ openModal: (mode?: 'create' | 'update', id?: string) => void }>;
  onSuccess?: () => void;
}

const PlanCreator: React.FC<IPlanCreatorProps> = ({ actionRef, onSuccess }) => {
  const { initialState } = useModel('@@initialState');
  const [form] = Form.useForm();
  const [branchFieldNum, setBranchFieldNum] = React.useState(1);
  const { dispatch: clearCaseEnumDispatch, valueEnum } = useClearCaseSelectEnum();
  const [visible, toggleVisible] = useToggle(false);
  const forceUpdate = useUpdate();
  const modalRef = React.useRef<{ planId: string }>({ planId: '' });
  const { RouteList = [] } = initialState;
  /** 更新模式
   * 1. 回填表单数据
   * 2. pvob component matchInfo 为可修改配置，其他表单项只读
   */
  const [isUpdateMode, setIsUpdateMode] = useToggle(false);

  React.useImperativeHandle(actionRef, () => ({
    async openModal(mode = 'create', id) {
      form.resetFields();
      setIsUpdateMode(mode === 'update');
      if (mode === 'update' && id) {
        modalRef.current.planId = id;
        const fieldValues = await planServices.getPlanDetail(id);
        let taskModels = {};
        if (fieldValues?.task_id) {
          const {taskModel: taskFields} = await taskService.getTaskDetail(fieldValues?.task_id);
          taskModels = taskFields;
          modalRef.current.task_id = fieldValues?.task_id;
        }
        if (fieldValues.originType === 'ClearCase') {
          clearCaseEnumDispatch('pvob', {});
          clearCaseEnumDispatch('component', { pvob: fieldValues.pvob });
          clearCaseEnumDispatch('stream', { component: fieldValues.component, pvob: fieldValues.pvob });
          form.setFieldsValue({...fieldValues, ...taskModels});
        } else {
          form.setFieldsValue({...fieldValues, gitignore:taskModels?.gitignore });
        }

        if (taskModels?.status === 'running') {
          toggleVisible(false);
          message.error('任务正在执行不可进行修改');
        } else {
          toggleVisible(true);
        }

      } else {
        toggleVisible(true);
      }
    },
  }));

  const actionText = isUpdateMode ? '更新' : '新建';

  const handleFinish = async (values: Plan.Base) => {
    try {
      if (isUpdateMode) {
        if (values?.originType === 'svn') {
          await planServices.updatePlan(modalRef.current.planId, values);
          await taskService.updateTask(modalRef.current.task_id, {
            svn_url: values?.originUrl,
            modelType: values?.originType,
            gitURL: values?.targetUrl,
            ...values
          });
        } else if (values?.originType === 'ClearCase') {
          let newTaskId = null;
          if (modalRef?.current?.task_id) {
            await taskService.updateTask(modalRef.current.task_id, {
              gitURL: values?.targetUrl,
              ...values
            });
            newTaskId = modalRef.current.task_id;
          } else {
            const {message: otherTaskId} = await taskService.createTask({
               gitURL: values?.targetUrl,
               ...values,
             });
            newTaskId = otherTaskId;
          }

          await planServices.updatePlan(modalRef.current.planId, {...values, task_id: Number(newTaskId)});
        } else {
          await planServices.updatePlan(modalRef.current.planId, values);
        }
      } else {
        let task_id = ''
         if (values?.originType === 'svn') {
           const {message: taskId} = await taskService.createTask({
             ...values,
             gitURL: values?.targetUrl,
             svn_url: values?.originUrl,
             modelType: values?.originType,
           });
           task_id = taskId;
        }  else if (values?.originType === 'ClearCase')  {
           const {message: otherTaskId} = await taskService.createTask({
             gitURL: values?.targetUrl,
             ...values,
           });
           task_id = otherTaskId;
        }
        await planServices.createPlan({...values, task_id: Number(task_id)});
      }
      message.success(`迁移计划${actionText}成功`);
      onSuccess?.();
      return true;
    } catch (err) {
      // message.error(`迁移任务${actionText}出现异常`);
      // eslint-disable-next-line no-console
      console.error(err);
      return false;
    }
  };

  useMount(() => {
    clearCaseEnumDispatch('pvob', {});
  });

  const dispatch = ({ type, payload }: any) => {
    switch (type) {
      case 'forceUpdate':
        forceUpdate();
        break;
      case 'getComponentValueEnum':
        clearCaseEnumDispatch('component', { pvob: payload });
        break;
      case 'stream':
        clearCaseEnumDispatch('stream', { pvob: payload?.pvob, component: payload?.component });
        break;
      default:
        break;
    }
  };

  const handleFormValuesChange = (changedValues: Partial<Plan.Base>) => {
    Object.entries(changedValues).forEach(([key, val]) => {
      CustomChangeHandlers[key]?.(form, val, dispatch);
    });
  };

   const addBranchField = () => {
    setBranchFieldNum((num) => num + 1);
  };

   const deleteBranch = (pos: number) => {
    setBranchFieldNum((num) => num - 1);
    const { matchInfo } = form.getFieldsValue(['matchInfo']);
    if (Array.isArray(matchInfo)) {
      matchInfo.splice(pos, 1);
      form.setFieldsValue({
        matchInfo,
      });
    }
  };

     /** matchInfo 字段不重复 */
  const isDuplicateMatchInfoItem = (key: 'stream' | 'gitBranch', inputVal?: string) => {
    if (!inputVal || !inputVal.trim()) return false;
    const matchInfo = form.getFieldValue('matchInfo');
    if (Array.isArray(matchInfo)) {
      const values = matchInfo.map((item) => item[key]);
      const len = values.filter((val) => val === inputVal).length;
      if (len >= 2) return true;
    }
    return false;
  };

  return (
    <ModalForm<Plan.Base>
      form={form}
      width="1000px"
      visible={visible}
      className={classnames(styles.formContainer)}
      layout="horizontal"
      onFinish={handleFinish}
      initialValues={InitialValues}
      title={`${actionText}迁移计划`}
      modalProps={{
        okText: actionText,
        bodyStyle: {maxHeight: 'calc(100vh - 108px)', overflowY: "scroll", overflowX: "hidden"}, style: {top: 0}
      }}
      onValuesChange={handleFormValuesChange}
      onVisibleChange={(vis) => toggleVisible(vis)}
    >
      <div className={classnames(styles.container)}>
      <FormSection
        left={
          <>
            <h6 className={styles.colTitle}>源仓库</h6>
            <ProFormRadio.Group
              name="originType"
              radioType="button"
              label="迁移任务类型"
              options={RouteList.includes('svnRoute') ? OriginTypeOptions.concat(ExType) : OriginTypeOptions}
              rules={[{ required: true, message: '请选择迁移任务类型' }]}
            />

            {form.getFieldValue('originType') === 'ClearCase' ? (
              <>
                <ProFormSelect
                  name="pvob"
                  label="PVOB"
                  fieldProps={{
                    getPopupContainer: triggerNode => triggerNode.parentElement
                  }}
                  placeholder="请选择 PVOB"
                  valueEnum={valueEnum.pvob}
                  rules={[{ required: true, message: '请选择 PVOB' }]}
                  showSearch
                />
                <ProFormSelect
                  name="component"
                  label="组件"
                   fieldProps={{
                    getPopupContainer: triggerNode => triggerNode.parentElement
                  }}
                  placeholder="请选择组件"
                  valueEnum={valueEnum.component}
                  rules={[{ required: true, message: '请选择组件' }]}
                  showSearch
                />
                <ProFormText
                  name="dir"
                  label="子目录"
                  placeholder="请输入组件子目录，如果为空则将迁移整个组件"
                />
                <ProFormText
                  name="ccUser"
                  label="CC用户名"
                  rules={[{ required: true, message: '请填写 CC 用户名' }]}
                  placeholder="请输入 CC 用户名"
                />
                <ProFormText.Password
                  name="ccPassword"
                  label="CC密码"
                  rules={[{ required: true, message: '请填写 CC 密码' }]}
                  placeholder="请输入 CC 密码"
                />
              </>
            ) : (
              <ProFormText
                name="originUrl"
                label="仓库地址"
                placeholder="请输入仓库地址"
                rules={[{ required: true, message: '请输入仓库地址' }]}
              />
            )}
          </>
        }
        right={
          <>
            <h6 className={styles.colTitle}>目标仓库及计划</h6>
            <ProFormRadio.Group
              name="translateType"
              radioType="button"
              label="迁移方式"
              options={form.getFieldValue('originType') === 'ClearCase' ? TranslateCCTypeOptions : TranslateTypeOptions}
              rules={[{ required: true, message: '请选择迁移方式' }]}
            />
            <ProFormText
              name="targetUrl"
              label="目标仓库地址"
              placeholder="请填写目标仓库地址"
              rules={[{ required: true, message: '请填写目标仓库地址' }]}
            />
            <ProFormDatePicker
              name="plan_start_time"
              label="计划迁移日期"
              fieldProps={{
                getPopupContainer: triggerNode => triggerNode.parentElement
              }}
              placeholder="请选择计划迁移日期"
              rules={[{ required: true, message: '请选择计划迁移日期' }]}
            />
            <ProFormDatePicker
              name="plan_switch_time"
              label="计划切换日期"
               fieldProps={{
                 getPopupContainer: triggerNode => triggerNode.parentElement
               }}
              placeholder="请选择计划切换日期"
              rules={[
                { required: true, message: '请选择计划切换日期' },
                {
                  async validator(rule, val) {
                    const stratTime = form.getFieldValue('plan_start_time');
                    if (dayjs(stratTime).isAfter(dayjs(val)))
                      throw new Error('计划切换日期应在计划迁移日期之后');
                  },
                },
              ]}
            />
            <ProFormTextArea
              name="gitignore"
              label="gitignore"
              placeholder="请输入 gitignore"
            />
            {
              form.getFieldValue('originType') === 'ClearCase' && (
                <>
                  <ProFormText
                    name="gitEmail"
                    label="Git Email"
                    placeholder="请填写Git Email"
                    rules={[{ required: true, message: '请输入 Git Email，格式：云桌面账号@ccbft.com', pattern: /^[^@]+@ccbft.com$/, }]}
                  />
                  <ProFormText
                    name="gitUser"
                    label="git 用户名"
                    rules={[{ required: true, message: '请填写 git 用户名' }]}
                    placeholder="请输入 git 用户名"
                  />
                  <ProFormText.Password
                    name="gitPassword"
                    label="git 密码"
                    rules={[{ required: true, message: '请填写 git 密码' }]}
                    placeholder="请输入 git 密码"
                  />
                </>
              )
            }
          </>
        }
      />
      {
        form.getFieldValue('originType') === 'ClearCase' && <div className={styles.gutter}>
          <h4 className={styles.sectionHeader}>CC开发流对应分支关系</h4>
          {formFieldsGenerator(
            Array.from(Array(branchFieldNum), (_, index) => [
              {
                name: ['matchInfo', index, 'stream'],
                component: FieldAutoComplete,
                label: '开发流',
                placeholder: '请选择开发流',
                valueEnum: valueEnum.stream,
                options: valueEnum.stream,
                showSearch: true,
                required: true,
                rules: [
                  {
                    // eslint-disable-next-line @typescript-eslint/no-shadow
                    async validator(_: any, value: string) {
                      if (isDuplicateMatchInfoItem('stream', value))
                        throw new Error('不能出现重复的开发流');
                    },
                  },
                ],
              },
              {
                name: ['matchInfo', index, 'gitBranch'],
                component: ProFormText,
                placeholder: '请输入Git对应分支',
                label: '分支',
                required: true,
                valueEnum: valueEnum.stream,
                rules: [
                  {
                    // eslint-disable-next-line @typescript-eslint/no-shadow
                    async validator(_: any, value: string) {
                      if (isDuplicateMatchInfoItem('gitBranch', value))
                        throw new Error('不能出现重复的 Git 分支');
                    },
                  },
                ],
              },
              <>
                {branchFieldNum !== 1 && (
                  <Button
                    key="delete"
                    icon={<MinusOutlined/>}
                    className={styles.deleteButton}
                    onClick={() => deleteBranch(index)}
                  />
                )}
                {index === branchFieldNum - 1 && (
                  <Button key="add" type="primary" icon={<PlusOutlined/>} onClick={addBranchField}/>
                )}
              </>,
            ]),
          )}
           <div className={classnames(styles.col, styles.keep)}>
          <span>
            <Form.Item valuePropName="checked" noStyle name="includeEmpty">
              <Checkbox />
            </Form.Item>
            <span className={styles.label}>是否保留空目录</span>
          </span>
          <span className={styles.keep}>
            <span className={classnames(styles.label, styles.keepLabel)}>占位文件名</span>
            <Form.Item
              name={['keep']}
              rules={[
                {
                  async validator(_, value) {
                    const { includeEmpty } = form.getFieldsValue(['includeEmpty']);
                    if (includeEmpty && !value) {
                      throw new Error('文件名称不能为空');
                    }
                    const invalidated = /[^a-z0-9-_.]+/.test(value);
                    if (invalidated) {
                      throw new Error(
                        '文件名称字符只能包括：字母、数字、"."(点)、"_"(下划线)和"-"(连字符)',
                      );
                    }
                  },
                },
              ]}
            >
              <Input
                size="small"
                style={{ width: 128, marginLeft: 12, marginTop: 10 }}
                placeholder="请输入占位文件名"
              />
            </Form.Item>
          </span>
        </div>
        </div>
      }
      <FormSection
        title="系统管理信息"
        left={
          <div className={classnames(styles.expecialField)}>
            <ProFormText
              name="subsystem"
              label="物理子系统英文简称"
              placeholder="请输入物理子系统英文简称"
              rules={[{ required: true, message: '请输入物理子系统英文简称' }]}
            />
          </div>
        }
        right={
          <ProFormText
            name="configLib"
            label="配置库"
            placeholder="请输入配置库名称"
            rules={[{ required: true, message: '请输入配置库名称' }]}
          />
        }
      />
      <FormSection
        title="执行组织信息"
        left={
          <>
            <ProFormSelect
              name="group"
              label="事业群"
              fieldProps={{
                getPopupContainer: triggerNode => triggerNode.parentElement
              }}
              placeholder="请选择事业群"
              options={GroupOptions}
              rules={[{ required: true, message: '请选择 事业群' }]}
            />
            <ProFormText
              name="team"
              label="项目组"
              placeholder="请填写项目组名称"
              rules={[{ required: true, message: '请填写项目组名称' }]}
            />
          </>
        }
        right={
          <>
            <ProFormText
              name="supporter"
              label="对接人姓名"
              placeholder="请填写对接人姓名"
              rules={[{ required: true, message: '请填写对接人姓名' }]}
            />
            <ProFormText
              label="联系人电话"
              name="supporterTel"
              placeholder="请填写联系人电话"
              rules={[{ required: true, message: '请填写联系人电话' }]}
            />
          </>
        }
        wholeLine={<ProFormTextArea name="tip" label="备注" placeholder="请填写备注信息" />}
      />
      <FormSection
        title="源业务信息"
        left={
          <ProFormSelect
            label="工程类型"
            name="projectType"
            placeholder="请选择工程类型"
            options={ProjectTypeOptions}
          />
        }
        wholeLine={
          <>
            <ProFormTextArea name="purpose" label="业务用途" placeholder="请填写业务用途" />
            <ProFormTextArea
              name="effect"
              label="影响范围"
              placeholder="请填写本次迁移对业务、技术方面的影响范围大小"
            />
          </>
        }
      />
        </div>
    </ModalForm>
  );
};

export default observer(PlanCreator);
