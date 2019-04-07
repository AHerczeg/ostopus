import React from 'react';
import ReactDOM from 'react-dom';
import Querys from './Querys';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<Querys />, div);
  ReactDOM.unmountComponentAtNode(div);
});
