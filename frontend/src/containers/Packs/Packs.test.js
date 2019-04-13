import React from 'react';
import ReactDOM from 'react-dom';
import Packs from './Packs';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<Packs />, div);
  ReactDOM.unmountComponentAtNode(div);
});
