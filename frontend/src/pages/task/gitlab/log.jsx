import { task as taskService } from '@/services';
import { useEffect, useState } from 'react';
import './log.css';

export default ({ id }) => {
  const [ logContent, setLogContent ] = useState('');
  useEffect(() => {
    taskService.getLogOutput(id).then(({ content }) => {
      setLogContent(content);
    })
  }, [id]);
  return (
    <pre className={"terminal"}>
      <code>
      {logContent}
      </code>
    </pre>
  );
}