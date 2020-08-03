import React, { useState } from 'react';
import './App.less';
import { Input, Button } from 'antd';
const { TextArea } = Input;

function App() {
  const [text, setText] = useState('');
  const [socket] = useState<WebSocket>(() : WebSocket => {
      const conn : WebSocket = new WebSocket(process.env.REACT_APP_WS_URL ? process.env.REACT_APP_WS_URL : '');
      conn.onopen = () : void => {
        console.log('open');
      };
      return conn;
    }
  );
  function send() {
    socket.send(text);
    setText('');
  }
  return (
    <div className="app-content">
      <div className="chat-content">
        <div className="header">
          header
        </div>
        <div className="body">
          <div className="user">
          </div>
          <div className="chat">
          </div>
        </div>
        <div className="footer">
          <div className="content">
            <TextArea autoSize={{minRows: 8, maxRows: 8}}
              onChange={e => setText(e.target.value)}
              value={text}
            />
          </div>
          <div className="action">
            <Button onClick={send}
              type="primary"
            >发送</Button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
