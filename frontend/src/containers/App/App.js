import React, { Component } from 'react';
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import Dashboard from '../../containers/Dashboard/Dashboard'
import Nodes from '../../containers/Nodes/Nodes'
import Querys from '../../containers/Querys/Querys'
import Packs from '../../containers/Packs/Packs'
import Logs from '../../containers/Logs/Logs'
import Settings from '../../containers/Settings/Settings'
import Sidebar from '../../components/Sidebar/Sidebar'
import {GridWrapper, HeaderWrapper as Header} from '../../components/Styles/Styles'


export default class App extends Component {

  render() {
      return(
        <GridWrapper>
          <Router>
              <Sidebar />
              <Header />
              <Switch>
                <Route path="/" exact component={Dashboard} />
                <Route path="/nodes" component={Nodes} />
                <Route path="/querys" component={Querys} />
                <Route path="/packs" component={Packs} />
                <Route path="/logs" component={Logs} />
                <Route path="/settings" component={Settings} />
              </Switch>
            </Router>
        </GridWrapper>
      );
  }
}