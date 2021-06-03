import React from 'react';
import { task } from '@/services';
import { useLocalObservable } from 'mobx-react';

interface RequestParams {
  current?: number;
  pageSize?: number;
  [k: string]: any;
}

export const useCacheRequestParams = (key: string) => {
  const storageKey = `PAGINATION_${key}`;
  const initialParams = React.useRef({
    current: 0,
    pageSize: 10,
  });
  const [params, setParams] = React.useState<RequestParams>(initialParams.current);

  const storage = React.useMemo(
    () => ({
      get() {
        try {
          return JSON.parse(sessionStorage.getItem(storageKey)!) || initialParams.current;
        } catch (e) {
          return initialParams.current;
        }
      },
      set(data: RequestParams) {
        try {
          sessionStorage.setItem(storageKey, JSON.stringify(data));
        } catch (e) {
          console.error(e);
        }
      },
    }),
    [storageKey],
  );

  React.useEffect(() => {
    setParams(storage.get());
  }, [storage]);

  return {
    params,
    setParams(data: RequestParams) {
      setParams(data);
      storage.set(data);
    },
  };
};

type OptionItem = Record<string, string>;
type OptionType = 'component' | 'pvob' | 'stream' | 'clearStream' | 'clearComponent';
const initialOptionState = { component: {}, pvob: {}, stream: {} } as Record<
  OptionType,
  OptionItem
>;
/** 获取 options item */
export const useClearCaseSelectEnum = () => {
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
      let res;
      switch (type) {
        case 'pvob':
          res = await task.getPvobs();
          set('pvob', listToOptions(res));
          break;
        case 'component':
          if (!payload.pvob) throw Error('pvob is required');
          res = await task.getComponents(payload.pvob);
          set('component', listToOptions(res));
          break;
        case 'stream':
          if (!payload.pvob || !payload.component) throw Error('pvob is required');
          res = await task.getStreams(payload.pvob, payload.component);
          set('stream', listToOptions(res));
          break;
        case 'clearStream':
          set('stream', listToOptions([]));
          break;
        case 'clearComponent':
          set('component', listToOptions([]));
          break;
        default:
          break;
      }
    },
    valueEnum: restState,
  };
};
