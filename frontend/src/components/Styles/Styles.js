import styled from 'styled-components'

/** Grid **/
export const GridWrapper = styled.div`
    display: grid;
    grid-template-columns: 5em auto;
    grid-template-rows: 6em calc(100vh - 6em);
    grid-template-areas:
        "sidebar header"
        "....... body";
`
/** Sidebar **/
export const SidebarWrapper = styled.div`
    grid-area: sidebar;    
    background-color: #1356AA;
    height: 100vh;
`
export const SidebarHeader = styled.div`
    color: white;
    font-size: 2em;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-top: 1em;
`

export const SidebarItemWrapper = styled.div`
    display: grid;
    grid-template-columns: 1fr;
`

export const SidebarItems = styled.div`
    font-size: 2em;
    display: flex;
    justify-content: center;;
    margin-top: 1em;
    a{
        color: white;
    }
`

/** Header **/
export const HeaderWrapper = styled.div`
    grid-area: header;    
    background-color: #fdfdfd;
`
/** Body **/
export const BodyWrapper = styled.div`
    grid-area: body;  
    background-color: #DFD9E2;
`
