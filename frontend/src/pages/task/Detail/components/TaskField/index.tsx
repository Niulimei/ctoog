import React from 'react';
import { useToggle } from 'umi';
import classnames from 'classnames';
import { guid } from '@/utils/utils';
import type { Task } from '@/typings/model';

import styles from './style.less';

// const PrivacyPassword: React.FC<{ value: string }> = () => {
//   const [isHidden, toggleHidden] = useToggle(true);
//   return <span>{value}</span>;
// };

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
          const isPasswordField = key.includes('password');

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

const TaskDetail: React.FC<{ data?: Task.Detail['taskModel'] }> = ({ data }) => {
  if (!data) return null;
  return (
    <div className={styles.gutter}>
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
        data,
      )}
      <div className={styles.divider} />
      {data.matchInfo?.map((matchInfo) =>
        descriptionsGenerator(['stream', 'gitBranch'], matchInfo),
      )}
      <div className={styles.col}>
        <span className={styles.row}>
          <span>是否保留空目录：{data.includeEmpty ? '是' : '否'}</span>
          <span style={{ marginLeft: '2em' }}>占位文件名：{data.keep}</span>
        </span>
      </div>
    </div>
  );
};

export default TaskDetail;
