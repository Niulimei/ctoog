import React from 'react';

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
