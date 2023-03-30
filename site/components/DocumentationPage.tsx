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
            <div className="container-fluid d-flex">
                <Navigation nav={nav} currentPageTitle={currentPage?.frontMatter.title}/>
                <main className="flex-column flex-xl-grow-1 mx-1 mx-lg-3 pb-5 py-5 bg-white bg-opacity-25">
                    <div className="flex-shrink-0">
                        <div className="content px-2 px-lg-4">
                            {html &&
                                <div className="content" dangerouslySetInnerHTML={{__html: html}}/>}
                            <div className="d-flex justify-content-evenly pb-3">
                                {previousPage && (
                                    <div
                                        className="col-3 col-xl-2 col-sm-5 p-2 p-xl-1 mx-2 mx-xl-4 float-start bg-light rounded-3 border">
                                        <a href={hrefForNextmd(previousPage.nextmd)}
                                           className="">
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
                                    <div
                                        className="col-3 col-xl-2 col-sm-5 p-2 p-xl-1 mx-2 mx-xl-4 float-end bg-light rounded-3 border">

                                        <a href={hrefForNextmd(nextPage.nextmd)}
                                           className="prev-next-link next dark">

                                            <div className="col-12 text-end p-0">
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

                </main>
            </div>

        </>
    )
        ;
}

// -----
// Utils
// -----

const hrefForNextmd = (nextmd: string[]) => `/${nextmd.join('/')}`;
