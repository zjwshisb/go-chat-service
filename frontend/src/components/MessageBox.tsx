import React, { useState, useEffect } from 'react';
import Message from './Message';
export default function(props : any) {
    return (<div className="chat">
      <Message></Message>
    </div>);
}