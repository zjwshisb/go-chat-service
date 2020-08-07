import React, { useState } from 'react';
import { Input, Button } from 'antd';
const { TextArea } = Input;
export default function(props : any) {
  const [text, setText] = useState('');
  return (<div className="footer">
      <div className="content">
        <TextArea autoSize={{minRows: 8, maxRows: 8}}
          onChange={e => setText(e.target.value)}
          value={text}
        />
      </div>
      <div className="action">
        <Button onClick={() => props.onSend(text)}
          type="primary"
        >发送</Button>
      </div>
    </div>
  );
}