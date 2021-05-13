/* eslint-disable react/no-array-index-key */
import React from 'react';
import Log from './Log';

import './style.less';

interface ILoggerProps {
  value?: string;
  actionRef?: React.ForwardedRef<{
    clear: () => void;
  }>;
}

const Logger: React.FC<ILoggerProps> = ({ value, actionRef }) => {
  const codeRef = React.useRef<any>(null);
  const lexerRef = React.useRef<Log | null>(null);

  React.useImperativeHandle(actionRef, () => ({
    clear: () => {
      lexerRef.current?.clearOutput();
    },
  }));

  React.useEffect(() => {
    lexerRef.current = new Log(codeRef.current);
  }, []);

  React.useEffect(() => {
    lexerRef.current?.clearOutput();
    lexerRef.current?.write(value || '');
  }, [value]);

  return (
    <div>
      <div ref={codeRef} className="logger__dark" />
    </div>
  );
};

export default Logger;
