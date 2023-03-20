import {NavItem} from "@/lib/types";
import NavSection from "@/components/NavSection";
import SideBarMenu from 'react-bootstrap-sidebar-menu';
import Link from "next/link";

type NavigationProps = { nav: NavItem[], currentPageTitle: string | undefined }
export default function Navigation(props: NavigationProps) {
    let {nav, currentPageTitle} = props

    return (
        <>
            <SideBarMenu className="pt-2 pt-lg-3 px-lg-1 navigation fixed sidebar-menu-scroll" bg={"dark"} variant={"dark"} >
                <SideBarMenu.Body>
                    <SideBarMenu.Brand key={"brand"} className="px-2 mb-0 pb-0">
                        <Link href="/docs"
                              className="navbar-brand align-content-center">
                            <object type="image/svg+xml" data="/penguin-travel.svg" className="logo">Happy Logo</object>
                            <div className="p-lg-1">Happy Path</div>
                        </Link>
                    </SideBarMenu.Brand>
                    <hr className="my-4 border-light p-0 m-0"/>
                    <SideBarMenu.Nav key={`top-nav`}>
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
