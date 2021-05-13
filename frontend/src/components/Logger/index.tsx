/* eslint-disable react/no-array-index-key */
import React from 'react';
import Log from './Log';

import './style.less';

interface ILoggerProps {
  value?: string;
}

const Logger: React.FC<ILoggerProps> = ({ value }) => {
  const codeRef = React.useRef<any>(null);
  const lexerRef = React.useRef<Log | null>(null);
  React.useEffect(() => {
    lexerRef.current = new Log(codeRef.current);
  }, []);

  React.useEffect(() => {
    lexerRef.current?.write(value || '');
  }, [value]);

  return (
    <div>
      <div ref={codeRef} className="logger__dark" />
    </div>
  );
};

export default Logger;
