import {NavItem} from "@/lib/types";
import NavSection from "@/components/NavSection";
import SideBarMenu from 'react-bootstrap-sidebar-menu';

type NavigationProps = { nav: NavItem[], currentPageTitle: string | undefined }
export default function Navigation(props: NavigationProps) {
    let {nav, currentPageTitle} = props

    return (
        <>
            <SideBarMenu >
                <SideBarMenu.Header>
                    <SideBarMenu.Brand>
                        <a href="/docs"
                           className="navbar-brand">
                            <object type="image/svg+xml" data="/penguin-app.svg" className="logo">Happy Logo</object>
                            <div className="ml-5 p-3">Happy Path</div>
                        </a>
                    </SideBarMenu.Brand>
                </SideBarMenu.Header>
                <SideBarMenu.Body>
                    <SideBarMenu.Nav>
                        {nav.map((navItem, index) => (
                                <>
                                    <NavSection navItem={navItem} index={index} currentPage={currentPageTitle}/>
                                </>
                            )
                        )}
                    </SideBarMenu.Nav>
                </SideBarMenu.Body>
            </SideBarMenu>
        </>
    )
}
