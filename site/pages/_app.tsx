import '@/styles/globals.scss'
import type {AppProps} from 'next/app'
import {Col, Container, Row} from 'react-bootstrap'
import {config} from '@fortawesome/fontawesome-svg-core'
import '@fortawesome/fontawesome-svg-core/styles.css'

import {useEffect} from "react";
import Script from "next/script";

config.autoAddCss = false

export default function App({Component, pageProps}: AppProps) {

    return (
        <>
            <Component {...pageProps} />
        </>
    )

}
