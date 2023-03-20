import Head from 'next/head';
import {DocumentationPageProps} from '@/lib/types';
import Navigation from "@/components/Navigation";
import React from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCircleArrowLeft, faCircleArrowRight} from "@fortawesome/free-solid-svg-icons"
import {Card, Col, Container, Row} from "react-bootstrap";

export default function DocumentationPage(props: DocumentationPageProps) {
    const {html, frontMatter, nav, nextmd} = props;

    const pages = nav
        .map((e) => e.props.subPaths?.map((sp) => ({...sp, navTitle: e.title})))
        .flatMap((e) => (e ? e : []));

    const currentIndex = pages.findIndex((v) => JSON.stringify(v.nextmd) === JSON.stringify(nextmd));
    const currentPage = currentIndex !== -1 ? pages[currentIndex] : null;
    const previousPage = currentIndex < 0 ? null : pages[currentIndex - 1];
    const nextPage = currentIndex === pages.length - 1 ? null : pages[currentIndex + 1];

    return (
        <>
            <Head>
                <title>{frontMatter.title}</title>
            </Head>
            <div className="d-flex flex-row flex-nowrap p-0">

                <Col lg={"auto"} md={"3"} className="bg-dark p-0">
                    <Navigation nav={nav} currentPageTitle={currentPage?.frontMatter.title}/>
                </Col>

                <Col lg={8} md={7} className="main float-end">
                    <Container fluid={true} className="pt-1 pt-lg-3 ps-md-5-left ms-md-5">
                        <Row>
                            <div className="pt-5 flex-shrink-0">
                                <div className="content">
                                    {html && <div className="content" dangerouslySetInnerHTML={{__html: html}}/>}
                                    <div className="prev-next-container container">
                                        {previousPage && (
                                            <div className="col-3 col-sm-5 p-2 float-start bg-light rounded-3 border">
                                                <a href={hrefForNextmd(previousPage.nextmd)}
                                                   className="prev-next-link prev dark">
                                                    <div className="col-12 text-start float-start p-0">
                                                        <FontAwesomeIcon
                                                            className="px-1 float-start" size="2x"
                                                            icon={faCircleArrowLeft}/>
                                                        {previousPage.navTitle}
                                                    </div>
                                                </a>
                                            </div>
                                        )}

                                        {nextPage && (
                                            <div className="col-3 col-sm-5 p-2 float-end bg-light rounded-3 border">

                                                <a href={hrefForNextmd(nextPage.nextmd)}
                                                   className="prev-next-link next dark">

                                                    <div className="col-12 text-end float-end p-0">
                                                        {nextPage.navTitle}
                                                        <FontAwesomeIcon
                                                            className="px-1 float-end" size="2x"
                                                            icon={faCircleArrowRight}/>
                                                    </div>

                                                </a>
                                            </div>
                                        )}
                                    </div>
                                </div>
                            </div>

                        </Row>

                    </Container>
                </Col>
            </div>
        </>
    );
}

// -----
// Utils
// -----

const hrefForNextmd = (nextmd: string[]) => `/${nextmd.join('/')}`;
