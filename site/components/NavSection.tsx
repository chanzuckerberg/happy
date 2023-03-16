import {NavItem} from '@/lib/types'
import {useState} from "react";
import {Button, Collapse} from "react-bootstrap";
import SideBarMenu from "react-bootstrap-sidebar-menu";


export type NavigationSubmenuProps = { navItem: NavItem, index: number, currentPage: string | undefined }

export default function NavSection(props: NavigationSubmenuProps) {
    let {navItem, index, currentPage} = props
    let {title} = navItem
    let {subPaths} = navItem.props
    let collapseTarget = `${navItem.title.toLowerCase().split(/\s/,).join('-')}-collapse`
    const [open, setOpen] = useState(true)

    return (
        <>
            <SideBarMenu.Sub>
                <SideBarMenu.Sub>
                    <SideBarMenu.Nav.Icon/>
                    {}
                    <SideBarMenu.Nav.Link href={`/${navItem.props.nextmd.join('/')}`}>
                        {navItem.props.frontMatter.title}
                    </SideBarMenu.Nav.Link>
                </SideBarMenu.Sub>
                <SideBarMenu.Sub>


                    <SideBarMenu.Nav>
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
