import { useLocalObservable } from 'mobx-react';
import { task } from '@/services';

type OptionItem = Record<string, string>;
type OptionType = 'component' | 'pvob' | 'stream';

/** 初始化状态 */
const initialOptionState = { component: {}, pvob: {}, stream: {} } as Record<
  OptionType,
  OptionItem
>;
/** 获取 options item */
export const useSelectOptions = () => {
  const { set, ...restState } = useLocalObservable(() => ({
    ...initialOptionState,
    set(type: OptionType, options: OptionItem) {
      this[type] = options;
    },
  }));

  const listToOptions = (list: any[]) => {
    // return { mockData: 'mockData' };
    if (!Array.isArray(list)) return {};
    return list.reduce(
      (res, item) => ({
        ...res,
        [item]: item,
      }),
      {},
    ) as OptionItem;
  };

  return {
    async dispatch(type: OptionType, payload: Partial<Record<OptionType, string>>) {
      if (type === 'pvob') {
        const res = await task.getPvobs();
        set('pvob', listToOptions(res));
      } else if (type === 'component') {
        if (!payload.pvob) throw Error('pvob is required');
        const res = await task.getComponents(payload.pvob);
        set('component', listToOptions(res));
      } else if (type === 'stream') {
        if (!payload.pvob || !payload.component) throw Error('pvob is required');
        const res = await task.getStreams(payload.pvob, payload.component);
        set('stream', listToOptions(res));
      }
    },
    options: restState,
  };
};

/**
 * render card title
 * @param {string} title
 * @returns {JSX.Element}
 */
export const renderCardTitle = (title: string) => {
  return <h3 style={{ textAlign: 'center', marginBottom: '20px' }}>{title}</h3>;
};
