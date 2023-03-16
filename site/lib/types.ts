import { NextMarkdownProps } from 'next-markdown';

export type FMatter = { title: string, slug: string };

export type HappyNextMarkdownProps = NextMarkdownProps<FMatter, FMatter>;

export type NavItem = { title: string; props: HappyNextMarkdownProps };
export type DocumentationPageProps = HappyNextMarkdownProps & { nav: NavItem[] };

