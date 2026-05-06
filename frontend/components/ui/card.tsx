import * as React from 'react';

import { cn } from '@/lib/utils';

/**
 * Card groups content inside a bordered panel.
 *
 * Inputs: standard div props.
 * Outputs: a shadcn-style card container.
 */
export function Card({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        'rounded-2xl border border-border bg-card/80 text-card-foreground shadow-sm backdrop-blur',
        className,
      )}
      {...props}
    />
  );
}

/**
 * CardHeader renders the top section of a card.
 *
 * Inputs: standard div props.
 * Outputs: a card header wrapper.
 */
export function CardHeader({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div className={cn('flex flex-col gap-1.5 p-6', className)} {...props} />
  );
}

/**
 * CardTitle renders a prominent section heading inside a card.
 *
 * Inputs: standard heading props.
 * Outputs: a semantic heading element.
 */
export function CardTitle({
  className,
  ...props
}: React.HTMLAttributes<HTMLHeadingElement>) {
  return (
    <h3
      className={cn('text-lg font-semibold leading-none tracking-tight', className)}
      {...props}
    />
  );
}

/**
 * CardDescription renders supporting text under a title.
 *
 * Inputs: standard paragraph props.
 * Outputs: muted supporting copy.
 */
export function CardDescription({
  className,
  ...props
}: React.HTMLAttributes<HTMLParagraphElement>) {
  return (
    <p className={cn('text-sm text-muted-foreground', className)} {...props} />
  );
}

/**
 * CardContent renders the main body of the card.
 *
 * Inputs: standard div props.
 * Outputs: padded content area.
 */
export function CardContent({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={cn('p-6 pt-0', className)} {...props} />;
}

/**
 * CardFooter renders the bottom action row of the card.
 *
 * Inputs: standard div props.
 * Outputs: padded footer area.
 */
export function CardFooter({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div className={cn('flex items-center p-6 pt-0', className)} {...props} />
  );
}
