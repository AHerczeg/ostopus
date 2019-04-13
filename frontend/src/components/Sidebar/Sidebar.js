import * as React from "react";
import { faHome, faLayerGroup, faNetworkWired, faStream, faCogs, faBoxOpen, faPollH } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Link } from "react-router-dom";
import { SidebarWrapper, SidebarHeader, SidebarItems, SidebarItemWrapper } from '../Styles/Styles'


export default class Sidebar extends React.Component {
    render() {
        return (
            <SidebarWrapper>
                <SidebarHeader>
                <FontAwesomeIcon icon={faLayerGroup} size="sm"></FontAwesomeIcon>
                </SidebarHeader>
                <SidebarItemWrapper>
                    <SidebarItems>
                        <Link to="/"><FontAwesomeIcon icon={faHome} size="sm"/></Link>
                    </SidebarItems>
                    <SidebarItems>
                        <Link to="/nodes"><FontAwesomeIcon icon={faNetworkWired} size="sm"/></Link>
                    </SidebarItems>
                    <SidebarItems>
                        <Link to="/querys"><FontAwesomeIcon icon={faStream} size="sm"/></Link>
                    </SidebarItems>
                    <SidebarItems>
                        <Link to="/packs"><FontAwesomeIcon icon={faBoxOpen} size="sm"/></Link>
                    </SidebarItems>
                    <SidebarItems>
                        <Link to="/logs"><FontAwesomeIcon icon={faPollH} size="sm"/></Link>
                    </SidebarItems>
                    <SidebarItems>
                        <Link to="/settings"><FontAwesomeIcon icon={faCogs} size="sm"/></Link>
                    </SidebarItems>
                </SidebarItemWrapper>
            </SidebarWrapper>
        );
    }
}