import type { Metadata } from 'next';
import type { ReactNode } from 'react';

import './globals.css';

export const metadata: Metadata = {
  title: 'Interview Coach',
  description: 'AI interview practice and review workspace',
};

type RootLayoutProps = Readonly<{
  children: ReactNode;
}>;

/**
 * RootLayout renders the application shell for every page.
 *
 * Inputs: `children` is the page subtree that Next.js renders into the body.
 * Outputs: a full HTML document with the global stylesheet applied.
 */
export default function RootLayout({ children }: RootLayoutProps) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-background text-foreground">
        {children}
      </body>
    </html>
  );
}
