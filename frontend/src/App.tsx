import React, { useState, useEffect } from 'react';
import MessageBox from './components/MessageBox';
import UserList from './components/UserList';
import MessageInput from './components/MessageInput';
import './App.less';
import Client from './classes/Connection';
function App() {
  const url = process.env.REACT_APP_WS_URL || '';
  const [client] = useState<Client>(() : Client => {
    return new Client(url);
  });
  useEffect(() => {
    return function close() : void {
    };
  }, []);
  function send(text: string) : void {
    client.send(text);
  }
  return (
    <div className="app-content">
      <div className="chat-content">
        <div className="header">
          header
        </div>
        <div className="body">
          <UserList />
          <MessageBox />
        </div>
        <MessageInput onSend={send} />
      </div>
    </div>
  );
}
export default App;
