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
            <Container fluid={true}>
                <Row className="flex-nowrap">
                    <Col md={3} sm={"auto"}>
                            <Navigation nav={nav} currentPageTitle={currentPage?.frontMatter.title}/>
                    </Col>
                    <Col className="py-3">
                        <div className="col-8" id="headerImage"/>
                        <div className="col-8 flex-shring-0 bg-white">

                            <div className="content">
                                {html && <div dangerouslySetInnerHTML={{__html: html}}/>}
                                <div className="prev-next-container container">

                                    {previousPage && (
                                        <Card className="col-3 p-2 float-start">
                                            <a href={hrefForNextmd(previousPage.nextmd)}
                                               className="prev-next-link prev dark">
                                                {previousPage.navTitle} Previous  <FontAwesomeIcon
                                                className="float-end" size="2x" icon={faCircleArrowLeft}/>
                                            </a>
                                        </Card>
                                    )}

                                    {nextPage && (
                                        <a href={hrefForNextmd(nextPage.nextmd)} className="prev-next-link next dark">
                                            <Card className="col-3 p-2 float-end">

                                                <div className="align-middle p-0">Next {nextPage.navTitle}
                                                    <FontAwesomeIcon
                                                        className="float-end" size="2x" icon={faCircleArrowRight}/>
                                                </div>
                                            </Card></a>
                                    )}
                                </div>
                            </div>
                        </div>
                    </Col>
                </Row>
            </Container>
        </>
    );
}

// -----
// Utils
// -----

const hrefForNextmd = (nextmd: string[]) => `/${nextmd.join('/')}`;
