/**
 * cn joins CSS class names and drops falsy values.
 *
 * Inputs: any number of class name strings or falsy values.
 * Outputs: a single space-delimited class string.
 */
export function cn(...values: Array<string | null | undefined | false>): string {
  return values.filter(Boolean).join(' ');
}
