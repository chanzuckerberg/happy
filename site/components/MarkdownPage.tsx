import Head from 'next/head';
import {HappyNextMarkdownProps} from '../lib/types';
import {Nav, Navbar} from "react-bootstrap";

export default function MarkdownPage(props: HappyNextMarkdownProps) {
    const {html, frontMatter} = props;
    return (
        <>
            <Head>
                <title>{frontMatter.title}</title>
            </Head>
            <div className={"container bg-white bg-opacity-50"}>
                <Nav className={"border-bottom border-dark"}>
                    <Navbar.Brand>
                    </Navbar.Brand>
                    <object style={{height: "4em"}} type="image/svg+xml" data="/penguin-travel.svg"
                            className="logo">Happy Logo
                    </object>
                    <Navbar.Text>
                        Happy Path
                    </Navbar.Text>

                </Nav>
                <div className={"mx-lg-5 px-xl-5"}>
                    {html && <div dangerouslySetInnerHTML={{__html: html}}/>}
                </div>
            </div>
        </>
    )
        ;
}
