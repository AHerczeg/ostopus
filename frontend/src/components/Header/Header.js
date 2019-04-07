import * as React from "react";
import {
    Alignment,
    Button,
    Classes,
    Navbar,
    NavbarDivider,
    NavbarGroup,
    NavbarHeading,
} from "@blueprintjs/core";
import { Link } from "react-router-dom";

export default class Header extends React.PureComponent {

    render() {
        return (
                <Navbar>
                    <NavbarGroup align={Alignment.LEFT}>
                        <NavbarHeading>OSTopus</NavbarHeading>
                        <NavbarDivider />
                        <Link to="/">
                            <Button className={Classes.MINIMAL} icon="home" text="Dashboard" />
                        </Link>
                        <Link to="/nodes">
                            <Button className={Classes.MINIMAL} icon="diagram-tree" text="Nodes" />
                        </Link>
                        <Link to="/querys">
                            <Button className={Classes.MINIMAL} icon="database" text="Querys" />
                        </Link>
                        
                    </NavbarGroup>
                    <NavbarGroup align={Alignment.RIGHT}>
                        <NavbarDivider />
                        <Button className="bp3-minimal" icon="user" />
                        
                        <Link to="/settings">
                            <Button className="bp3-minimal" icon="settings" />
                        </Link>
                    </NavbarGroup>
                </Navbar>
        );
    }
}