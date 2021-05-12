import React from 'react';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import type { IUnControlledCodeMirror } from 'react-codemirror2';

import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/material.css';

require('codemirror/mode/shell/shell');

const defaultCodeMirrorOptions = {
  lineNumbers: true,
  readOnly: true,
  theme: 'material',
  mode: 'shell',
  cursorHeight: 0,
  workDelay: 100,
};

const Logger: React.FC<IUnControlledCodeMirror> = (props) => {
  return (
    <CodeMirror
      options={{
        ...defaultCodeMirrorOptions,
      }}
      {...props}
    />
  );
};

export default Logger;
