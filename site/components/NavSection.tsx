import {NavItem} from '@/lib/types'
import {useState} from "react";
import {Button, Collapse} from "react-bootstrap";
import SideBarMenu from "react-bootstrap-sidebar-menu";


export type NavigationSubmenuProps = { navItem: NavItem, index: number, currentPage: string | undefined }

export default function NavSection(props: NavigationSubmenuProps) {
    let {navItem, index, currentPage} = props
    let {title} = navItem
    let {subPaths} = navItem.props

    return (
        <>
            <SideBarMenu.Sub id={`${title}-sub`}>
                <SideBarMenu.Sub>
                    <SideBarMenu.Nav key={`sidebar-${title}`}>
                        <SideBarMenu.Text>
                            {navItem.props.frontMatter.slug != undefined ?
                                <>
                                    <a href={`/${navItem.props.nextmd.join('/')}`}>
                                        {title}
                                    </a>
                                </> : <>{title}</>}
                        </SideBarMenu.Text>
                        {subPaths?.map((nextDoc, docIndex) => (
                            <>
                                <SideBarMenu.Nav.Link href={`/${nextDoc.nextmd.join('/')}`}>
                                    {nextDoc.frontMatter.title}
                                </SideBarMenu.Nav.Link>
                            </>
                        ))}
                    </SideBarMenu.Nav>
                </SideBarMenu.Sub>

            </SideBarMenu.Sub>

        </>
    )
}
