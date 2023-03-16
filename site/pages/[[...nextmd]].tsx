import NextMarkdown from 'next-markdown';
import DocumentationPage from '../components/DocumentationPage';
import MarkdownPage from '../components/MarkdownPage';
import remarkPrism from 'remark-prism';
import 'prismjs/themes/prism-tomorrow.css';
import {GetStaticPropsContext} from 'next';
import {DocumentationPageProps, HappyNextMarkdownProps} from '../lib/types';

const nextmd = NextMarkdown({
    pathToContent: './pages-markdown',
    remarkPlugins: [remarkPrism],
});

export const getStaticProps = async (context: GetStaticPropsContext<{ nextmd: string[] }>) => {
    if (isDocumentation(context.params?.nextmd)) {
        console.log("Trying to get document")
        return {
            props: {
                ...(await nextmd.getStaticProps(context)).props,
                nav: [
                    {title: 'Getting Started', ...(await nextmd.getStaticPropsForNextmd(['docs', 'quickstart']))},
                    {title: 'API', ...(await nextmd.getStaticPropsForNextmd(['docs', 'api']))},
                    {title: 'EKS', ...(await nextmd.getStaticPropsForNextmd(['docs', 'eks']))}
                ],
            },
        };
    } else {
        return nextmd.getStaticProps(context);
    }
};

export const getStaticPaths = nextmd.getStaticPaths;

export default function MyMarkdownPage(props: DocumentationPageProps | HappyNextMarkdownProps) {
    if (isDocumentation(props.nextmd)) {
        return <DocumentationPage {...(props as DocumentationPageProps)} />;
    } else {
        return <MarkdownPage {...(props as HappyNextMarkdownProps)} />;
    }
}

// ----------
// Utils
// ----------

const isDocumentation = (nextmd: string[] | undefined) => nextmd?.includes('docs'); // tslint:disable-line:no-shadowed-variable
