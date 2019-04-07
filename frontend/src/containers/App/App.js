import React, { Component, Fragment } from 'react';
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import Dashboard from '../../containers/Dashboard/Dashboard'
import Nodes from '../../containers/Nodes/Nodes'
import Querys from '../../containers/Querys/Querys'
import Settings from '../../containers/Settings/Settings'
import Header from '../../components/Header/Header'

export default class App extends Component {

  render() {
      return(
        <Fragment>
          <Router>
              <Header />
              <Switch>
                <Route path="/" exact component={Dashboard} />
                <Route path="/nodes" component={Nodes} />
                <Route path="/querys" component={Querys} />
                <Route path="/settings" component={Settings} />
              </Switch>
            </Router>
        </Fragment>
      );
  }
}