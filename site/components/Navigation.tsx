import {NavItem} from "@/lib/types";
import NavSection from "@/components/NavSection";
import SideBarMenu from 'react-bootstrap-sidebar-menu';
import Link from "next/link";

type NavigationProps = { nav: NavItem[], currentPageTitle: string | undefined }
export default function Navigation(props: NavigationProps) {
    let {nav, currentPageTitle} = props

    return (
        <>
            <SideBarMenu
                className={"flex-xl-shrink-1 min-vh-100 d-md-inline-flex mx-0 mx-lg-1 my-3 p-1 px-lg-3 bg-opacity-75"}
                bg={"light"} variant={"light"} expand="lg"
                hide="sm">
                <SideBarMenu.Body className="">
                    <SideBarMenu.Brand key={"brand"} className="">
                        <Link href="/docs"
                              className="navbar-brand p-0 mb-0">
                            <object type="image/svg+xml" data="/penguin-travel.svg" className="logo">Happy Logo
                            </object>
                        </Link>
                        <div className="p-lg-1 display-5 text-center">Happy Path</div>
                    </SideBarMenu.Brand>
                    <hr className="p-0 m-0 my-2 border-dark"/>
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
