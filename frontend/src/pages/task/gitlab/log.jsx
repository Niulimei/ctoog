import { task as taskService } from '@/services';
import { useEffect, useState } from 'react';
import { LoadingOutlined } from '@ant-design/icons';
import './log.css';

export default ({ id }) => {
  const [ logContent, setLogContent ] = useState('');
  const [load, setLoad] = useState(false);
  useEffect(() => {
    setLoad(true)
    taskService.getLogOutput(id).then(({ content }) => {
      setLogContent(content);
      setLoad(false)
    })
  }, [id]);
  return (
    <>
    {
        load ? <div style={{ textAlign: 'center' }}><LoadingOutlined /></div> : <pre className={"terminal"}>
        <code>
          {logContent}
        </code>
      </pre>
    }
    </>
  );
}