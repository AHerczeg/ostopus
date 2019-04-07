import React from 'react';
import ReactDOM from 'react-dom';
import Nodes from './Nodes';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<Nodes />, div);
  ReactDOM.unmountComponentAtNode(div);
});
