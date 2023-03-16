import Head from 'next/head';
import { HappyNextMarkdownProps } from '../lib/types';

export default function MarkdownPage(props: HappyNextMarkdownProps) {
    const { html, frontMatter } = props;
    return (
        <>
            <Head>
                <title>{frontMatter.title}</title>
            </Head>
            <div>{html && <div dangerouslySetInnerHTML={{ __html: html }} />}</div>
        </>
    );
}
